package landscape

import (
	"context"
	"errors"
	"fmt"
	"os/user"

	landscapeapi "github.com/canonical/landscape-hostagent-api"
	log "github.com/canonical/ubuntu-pro-for-wsl/windows-agent/internal/grpc/logstreamer"
	"github.com/canonical/ubuntu-pro-for-wsl/windows-agent/internal/proservices/landscape/distroinstall"
	"github.com/ubuntu/decorate"
	"github.com/ubuntu/gowsl"
)

// executor is in charge of executing commands received from the Landscape server.
type executor struct {
	serviceData
}

func (e executor) exec(ctx context.Context, command *landscapeapi.Command) (err error) {
	defer decorate.OnError(&err, "could not execute command %s", commandString(command))

	switch cmd := command.GetCmd().(type) {
	case *landscapeapi.Command_AssignHost_:
		return e.assignHost(ctx, cmd.AssignHost)
	case *landscapeapi.Command_Start_:
		return e.start(ctx, cmd.Start)
	case *landscapeapi.Command_Stop_:
		return e.stop(ctx, cmd.Stop)
	case *landscapeapi.Command_Install_:
		return e.install(ctx, cmd.Install)
	case *landscapeapi.Command_Uninstall_:
		return e.uninstall(ctx, cmd.Uninstall)
	case *landscapeapi.Command_SetDefault_:
		return e.setDefault(ctx, cmd.SetDefault)
	case *landscapeapi.Command_ShutdownHost_:
		return e.shutdownHost(ctx, cmd.ShutdownHost)
	default:
		return fmt.Errorf("unknown command type %T: %v", command.GetCmd(), command.GetCmd())
	}
}

func commandString(command *landscapeapi.Command) string {
	switch cmd := command.GetCmd().(type) {
	case *landscapeapi.Command_AssignHost_:
		return fmt.Sprintf("Assign host (uid: %q)", cmd.AssignHost.GetUid())
	case *landscapeapi.Command_Start_:
		return fmt.Sprintf("Start (id: %q)", cmd.Start.GetId())
	case *landscapeapi.Command_Stop_:
		return fmt.Sprintf("Stop (id: %q)", cmd.Stop.GetId())
	case *landscapeapi.Command_Install_:
		return fmt.Sprintf("Install (id: %q)", cmd.Install.GetId())
	case *landscapeapi.Command_Uninstall_:
		return fmt.Sprintf("Uninstall (id: %q)", cmd.Uninstall.GetId())
	case *landscapeapi.Command_SetDefault_:
		return fmt.Sprintf("SetDefault (id: %q)", cmd.SetDefault.GetId())
	case *landscapeapi.Command_ShutdownHost_:
		return "ShutdownHost"
	default:
		return "Unknown"
	}
}

func (e executor) assignHost(ctx context.Context, cmd *landscapeapi.Command_AssignHost) error {
	conf := e.config()

	if uid, err := conf.LandscapeAgentUID(ctx); err != nil {
		log.Warningf(ctx, "Possibly overriding current landscape client UID: could not read current Landscape UID: %v", err)
	} else if uid != "" {
		log.Warning(ctx, "Overriding current landscape client UID")
	}

	if err := conf.SetLandscapeAgentUID(ctx, cmd.GetUid()); err != nil {
		return err
	}

	return nil
}

func (e executor) start(ctx context.Context, cmd *landscapeapi.Command_Start) (err error) {
	log.Debugf(ctx, "Landscape: received command Start. Target: %s", cmd.GetId())
	d, ok := e.database().Get(cmd.GetId())
	if !ok {
		return fmt.Errorf("distro %q not in database", cmd.GetId())
	}

	return d.LockAwake()
}

func (e executor) stop(ctx context.Context, cmd *landscapeapi.Command_Stop) (err error) {
	log.Debugf(ctx, "Landscape: received command Stop. Target: %s", cmd.GetId())
	d, ok := e.database().Get(cmd.GetId())
	if !ok {
		return fmt.Errorf("distro %q not in database", cmd.GetId())
	}

	return d.ReleaseAwake()
}

func (e executor) install(ctx context.Context, cmd *landscapeapi.Command_Install) (err error) {
	log.Debugf(ctx, "Landscape: received command Install. Target: %s", cmd.GetId())

	if cmd.GetId() == "" {
		return errors.New("Landscape install: empty distro name")
	}

	distro := gowsl.NewDistro(ctx, cmd.GetId())
	if registered, err := distro.IsRegistered(); err != nil {
		return err
	} else if registered {
		return errors.New("Landscape install: already installed")
	}

	if err := e.cloudInit().WriteDistroData(cmd.GetId(), cmd.GetCloudinit()); err != nil {
		return fmt.Errorf("Landscape install: skipped installation: %v", err)
	}

	if err := gowsl.Install(ctx, distro.Name()); err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}
		// Avoid error states by cleaning up on error
		err := distro.Uninstall(ctx)
		if err != nil {
			log.Infof(ctx, "Landscape Install: distro %q: failed to uninstall after failed Install: %v", distro.Name(), err)
		}
	}()

	if err := distroinstall.InstallFromExecutable(ctx, distro); err != nil {
		return err
	}

	// TODO: The rest of this function will need to be rethought once cloud-init support exists.
	windowsUser, err := user.Current()
	if err != nil {
		return err
	}

	userName := windowsUser.Username
	if !distroinstall.UsernameIsValid(userName) {
		userName = "ubuntu"
	}

	uid, err := distroinstall.CreateUser(ctx, distro, userName, windowsUser.Name)
	if err != nil {
		return err
	}

	if err := distro.DefaultUID(uid); err != nil {
		return fmt.Errorf("could not set user as default: %v", err)
	}

	return nil
}

func (e executor) uninstall(ctx context.Context, cmd *landscapeapi.Command_Uninstall) (err error) {
	log.Debugf(ctx, "Landscape: received command Uninstall. Target: %s", cmd.GetId())
	d, ok := e.database().Get(cmd.GetId())
	if !ok {
		return fmt.Errorf("Landscape uninstall: distro %q not in database", cmd.GetId())
	}

	if err := d.Uninstall(ctx); err != nil {
		// Uninstall's error message already includes the distro name.
		return fmt.Errorf("Landscape uninstall: %v", err)
	}

	if err := e.cloudInit().RemoveDistroData(d.Name()); err != nil {
		log.Warningf(ctx, "Landscape uninstall: distro %q: %v", d.Name(), err)
	}

	return nil
}

func (e executor) setDefault(ctx context.Context, cmd *landscapeapi.Command_SetDefault) error {
	log.Debugf(ctx, "Landscape: received command SetDefault. Target: %s", cmd.GetId())
	d := gowsl.NewDistro(ctx, cmd.GetId())
	return d.SetAsDefault()
}

//nolint:unparam // cmd is not used, but kep here for consistency with other commands.
func (e executor) shutdownHost(ctx context.Context, cmd *landscapeapi.Command_ShutdownHost) error {
	log.Debug(ctx, "Landscape: received command ShutdownHost")
	return gowsl.Shutdown(ctx)
}
