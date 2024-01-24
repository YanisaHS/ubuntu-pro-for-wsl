package cloudinit_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/canonical/ubuntu-pro-for-wsl/common/golden"
	"github.com/canonical/ubuntu-pro-for-wsl/windows-agent/internal/cloudinit"
	"github.com/canonical/ubuntu-pro-for-wsl/windows-agent/internal/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		breakWriteAgentData bool
		wantErr             bool
	}{
		"Success": {},
		"Error when cloud-init agent file cannot be written": {breakWriteAgentData: true, wantErr: true},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			dataDir := t.TempDir()

			conf := &mockConfig{
				subcriptionErr: tc.breakWriteAgentData,
			}

			_, err := cloudinit.New(ctx, conf, dataDir)
			if tc.wantErr {
				require.Error(t, err, "Cloud-init creation should have returned an error")
				return
			}

			require.NoError(t, err, "Cloud-init creation should have returned no error")
			require.Len(t, conf.notify, 1, "Cloud-init should have attached a callback to the config")

			// Assert that the subscribed function works
			path := filepath.Join(dataDir, "agent.yaml")
			require.NoErrorf(t, os.RemoveAll(path), "Removing the agent cloud-init should not fail")

			conf.triggerNotify()

			// We don't assert on specifics, as they are tested in WriteAgentData tests.
			require.FileExists(t, path, "agent data file was not created when updating the config")
		})
	}
}

func TestWriteAgentData(t *testing.T) {
	t.Parallel()

	// All error cases share a golden file so we need to protect it during updates
	var sharedGolden goldenMutex

	const landscapeConfigOld string = `[irrelevant]
info=this section should have been omitted

[client]
data=This is an old data field
info=This is the old configuration
`

	const landscapeConfigNew string = `[irrelevant]
info=this section should have been omitted

[client]
info = This is the new configuration
url = www.example.com/new/rickroll
`

	testCases := map[string]struct {
		// Contents
		skipProToken      bool
		skipLandscapeConf bool

		// Break marshalling
		breakSubscription bool
		breakLandscape    bool

		// Landcape parsing
		landscapeNoClientSection bool
		badLandscape             bool

		// Break writing to file
		breakDir      bool
		breakTempFile bool
		breakFile     bool

		wantErr bool
	}{
		"Success":                                    {},
		"Success without pro token":                  {skipProToken: true},
		"Success without Landscape":                  {skipLandscapeConf: true},
		"Success without Landscape [client] section": {landscapeNoClientSection: true},
		"Success with empty contents":                {skipProToken: true, skipLandscapeConf: true},

		"Error obtaining pro token":             {breakSubscription: true, wantErr: true},
		"Error obtaining Landscape config":      {breakLandscape: true, wantErr: true},
		"Error with erroneous Landscape config": {badLandscape: true, wantErr: true},

		"Error when the datadir cannot be created":   {breakDir: true, wantErr: true},
		"Error when the temp file cannot be written": {breakTempFile: true, wantErr: true},
		"Error when the temp file cannot be renamed": {breakFile: true, wantErr: true},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			root := t.TempDir()
			dir := filepath.Join(root, "cloud-init")
			path := filepath.Join(dir, "agent.yaml")

			conf := &mockConfig{
				proToken:      "OLD_PRO_TOKEN",
				landscapeConf: landscapeConfigOld,
			}

			// Test a clean filesystem (New calls WriteAgentData internally)
			ci, err := cloudinit.New(ctx, conf, dir)
			require.NoError(t, err, "cloudinit.New should return no error")
			require.FileExists(t, path, "New() should have created an agent cloud-init file")

			// Test overriding the file: New() created the agent.yaml file
			conf.subcriptionErr = tc.breakSubscription
			conf.landscapeErr = tc.breakLandscape

			conf.proToken = "NEW_PRO_TOKEN"
			if tc.skipProToken {
				conf.proToken = ""
			}

			conf.landscapeConf = landscapeConfigNew
			if tc.badLandscape {
				conf.landscapeConf = "This is not valid ini"
			}
			if tc.landscapeNoClientSection {
				conf.landscapeConf = "[irrelevant]\ninfo=This section should be ignored"
			}
			if tc.skipLandscapeConf {
				conf.landscapeConf = ""
			}

			if tc.breakTempFile {
				require.NoError(t, os.RemoveAll(path+".tmp"), "Setup: Agent cloud-init file should not fail to delete")
				require.NoError(t, os.MkdirAll(path+".tmp", 0600), "Setup: could not create directory to mess with cloud-init")
			}

			if tc.breakFile {
				require.NoError(t, os.RemoveAll(path), "Setup: Agent cloud-init file should not fail to delete")
				require.NoError(t, os.MkdirAll(path, 0600), "Setup: could not create directory to mess with cloud-init")
			}

			if tc.breakDir {
				require.NoError(t, os.RemoveAll(dir), "Setup: Agent cloud-init file should not fail to delete")
				require.NoError(t, os.WriteFile(dir, nil, 0600), "Setup: could not create file to mess with cloud-init directory")
			}

			err = ci.WriteAgentData(ctx)
			var opts []golden.Option
			if tc.wantErr {
				require.Error(t, err, "WriteAgentData should have returned an error")
				errorGolden := filepath.Join(golden.TestFamilyPath(t), "golden", "error-cases")
				opts = append(opts, golden.WithGoldenPath(errorGolden))
			} else {
				require.NoError(t, err, "WriteAgentData should return no errors")
			}

			// Assert that the file was updated (success case) or that the old one remains (error case)
			if tc.breakFile || tc.breakDir {
				// Cannot really assert on anything: we removed the old file
				return
			}

			got, err := os.ReadFile(path)
			require.NoError(t, err, "There should be no error reading the cloud-init agent file")

			sharedGolden.Lock()
			defer sharedGolden.Unlock()

			want := golden.LoadWithUpdateFromGolden(t, string(got), opts...)
			require.Equal(t, want, string(got), "Agent cloud-init file does not match the golden file")
		})
	}
}

// goldenMutex is a mutex that only works when golden update is enabled.
type goldenMutex struct {
	sync.Mutex
}

func (mu *goldenMutex) Lock() {
	if !golden.UpdateEnabled() {
		return
	}
	mu.Mutex.Lock()
}

func (mu *goldenMutex) Unlock() {
	if !golden.UpdateEnabled() {
		return
	}
	mu.Mutex.Unlock()
}

func TestWriteDistroData(t *testing.T) {
	t.Parallel()

	const oldCloudInit = `# cloud-init
# I'm an old piece of user data
data:
	is_this_data: Yes, it is
	new: false
`

	const newCloudInit = `# cloud-init
# I'm a shiny new piece of user data
data:
	new: true
`

	testCases := map[string]struct {
		// Break marshalling
		emptyData bool
		noOldData bool

		// Break writing to file
		breakDir      bool
		breakTempFile bool
		breakFile     bool

		wantErr bool
	}{
		"Success":                  {},
		"Success with no old data": {noOldData: true},
		"Success with empty data":  {emptyData: true},

		"Error when the datadir cannot be created":   {breakDir: true, wantErr: true},
		"Error when the temp file cannot be written": {breakTempFile: true, wantErr: true},
		"Error when the temp file cannot be renamed": {breakFile: true, wantErr: true},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			distroName := "CoolDistro"

			root := t.TempDir()
			dir := filepath.Join(root, "cloud-init")
			path := filepath.Join(dir, "landscape", distroName+".user-data")

			conf := &mockConfig{}

			// Test a clean filesystem (New calls WriteAgentData internally)
			ci, err := cloudinit.New(ctx, conf, dir)
			require.NoError(t, err, "Setup: cloud-init New should return no errors")

			if !tc.noOldData {
				require.NoError(t, os.MkdirAll(filepath.Dir(path), 0600), "Setup: could not write old distro data directory")
				require.NoError(t, os.WriteFile(path, []byte(oldCloudInit), 0600), "Setup: could not write old distro data")
			}

			if tc.breakTempFile {
				require.NoError(t, os.RemoveAll(path+".tmp"), "Setup: Distro cloud-init file should not fail to delete")
				require.NoError(t, os.MkdirAll(path+".tmp", 0600), "Setup: could not create directory to mess with cloud-init")
			}

			if tc.breakFile {
				require.NoError(t, os.RemoveAll(path), "Setup: Distro cloud-init file should not fail to delete")
				require.NoError(t, os.MkdirAll(path, 0600), "Setup: could not create directory to mess with cloud-init")
			}

			if tc.breakDir {
				require.NoError(t, os.RemoveAll(dir), "Setup: Distro cloud-init file should not fail to delete")
				require.NoError(t, os.WriteFile(dir, nil, 0600), "Setup: could not create file to mess with cloud-init directory")
			}

			var input string
			if !tc.emptyData {
				input = newCloudInit
			}

			err = ci.WriteDistroData(distroName, input)
			var want string
			if tc.wantErr {
				require.Error(t, err, "WriteAgentData should have returned an error")
				want = oldCloudInit
			} else {
				require.NoError(t, err, "WriteAgentData should return no errors")
				want = input
			}

			// Assert that the file was updated (success case) or that the old one remains (error case)
			if tc.breakFile || tc.breakDir {
				// Cannot really assert on anything: we removed the old file
				return
			}

			got, err := os.ReadFile(path)
			require.NoError(t, err, "There should be no error reading the distro's cloud-init file")
			require.Equal(t, want, string(got), "Agent cloud-init file does not match the golden file")
		})
	}
}

func TestRemoveDistroData(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		fileDoesNotExist bool
		dirDoesNotExist  bool

		wantErr bool
	}{
		"Success":                                  {},
		"Success when the file did not exist":      {fileDoesNotExist: true},
		"Success when the directory did not exist": {dirDoesNotExist: true},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			distroName := "CoolDistro"

			root := t.TempDir()
			distroDir := filepath.Join(root, "landscape")
			path := filepath.Join(distroDir, distroName+".user-data")

			ci, err := cloudinit.New(ctx, &mockConfig{}, root)
			require.NoError(t, err, "Setup: cloud-init New should return no errors")

			if !tc.dirDoesNotExist {
				require.NoError(t, os.MkdirAll(distroDir, 0700), "Setup: could not set up directory")
				if !tc.fileDoesNotExist {
					require.NoError(t, os.WriteFile(path, []byte("hello, world!"), 0600), "Setup: could not set up directory")
				}
			}

			err = ci.RemoveDistroData(distroName)
			require.NoError(t, err, "RemoveDistroData should return no errors")
			require.NoFileExists(t, path, "RemoveDistroData should remove the distro cloud-init data file")
		})
	}
}

type mockConfig struct {
	proToken       string
	subcriptionErr bool

	landscapeConf string
	landscapeErr  bool

	notify []func()
}

func (c mockConfig) triggerNotify() {
	for _, f := range c.notify {
		f()
	}
}

func (c mockConfig) Subscription(ctx context.Context) (string, config.Source, error) {
	if c.subcriptionErr {
		return "", config.SourceNone, errors.New("culd not get subscription: mock error")
	}

	if c.proToken == "" {
		return "", config.SourceNone, nil
	}

	return c.proToken, config.SourceUser, nil
}

func (c mockConfig) LandscapeClientConfig(context.Context) (string, config.Source, error) {
	if c.landscapeErr {
		return "", config.SourceNone, errors.New("could not get landscape configuration: mock error")
	}

	if c.landscapeConf == "" {
		return "", config.SourceNone, nil
	}

	return c.landscapeConf, config.SourceUser, nil
}

func (c *mockConfig) Notify(f func()) {
	c.notify = append(c.notify, f)
}