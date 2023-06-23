// Package contracts manages Microsoft-Store-entitled subscriptions.
package contracts

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/canonical/ubuntu-pro-for-windows/windows-agent/internal/contracts/client"
	"github.com/canonical/ubuntu-pro-for-windows/windows-agent/internal/contracts/microsoftstore"
	"github.com/ubuntu/decorate"
)

const defaultProURL = "https://contracts.canonical.com"

type options struct {
	proURL *url.URL
}

// Option is an optional argument for ProToken.
type Option func(*options)

// WithProURL overrides the Ubuntu Pro contract server URL.
func WithProURL(proURL *url.URL) Option {
	return func(o *options) {
		o.proURL = proURL
	}
}

// ProToken directs the dance between the Microsoft Store and the Ubuntu Pro contract server to
// validate a store entitlement and obtain its associated pro token. If there is no entitlement,
// the token is returned as an empty string.
func ProToken(ctx context.Context, args ...Option) (token string, err error) {
	var opts options

	for _, f := range args {
		f(&opts)
	}

	if opts.proURL == nil {
		url, err := url.Parse(defaultProURL)
		if err != nil {
			return token, fmt.Errorf("could not parse default contract server URL %q: %v", defaultProURL, err)
		}
		opts.proURL = url
	}

	contractClient := client.New(opts.proURL, &http.Client{Timeout: 30 * time.Second})
	msftStore := microsoftstore.New()

	token, err = proToken(ctx, contractClient, msftStore)
	if err != nil {
		return "", err
	}

	return token, nil
}

func proToken(ctx context.Context, serverClient *client.Client, msftStore *microsoftstore.Store) (proToken string, err error) {
	defer decorate.OnError(&err, "could not obtain pro token")

	adToken, err := serverClient.GetServerAccessToken(ctx)
	if err != nil {
		return "", err
	}

	storeToken, err := msftStore.SubscriptionToken(ctx, adToken)
	if err != nil {
		return "", err
	}

	proToken, err = serverClient.GetProToken(ctx, storeToken)
	if err != nil {
		return "", err
	}

	return proToken, nil
}
