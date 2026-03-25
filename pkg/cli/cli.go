package cli

import (
	"context"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"golang.org/x/oauth2"
)

// Base contains the basic needs for all CLI commands
type Base struct {
	Args      []string
	Client    *govultr.Client
	Options   *govultr.ListOptions
	Printer   *printer.Output
	Context   context.Context
	UserAgent string
}

// NewCLIBase creates new base struct
func NewCLIBase(userAgent string) *Base {
	base := Base{}
	base.UserAgent = userAgent
	base.configurePrinter()
	base.configureClient(nil)
	base.configureContext()
	return &base
}

func (b *Base) configureClient(oauthClient *http.Client) {
	b.Client = govultr.NewClient(oauthClient)
	b.Client.SetRateLimit(1 * time.Second)
	b.Client.SetUserAgent(b.UserAgent)
}

func (b *Base) configurePrinter() {
	b.Printer = &printer.Output{}
}

func (b *Base) configureContext() {
	b.Context = context.Background()
}

func (b *Base) HasAuth() bool {
	var token string

	if viper.IsSet("api-key") {
		token = viper.GetString("api-key")

		if token == "" {
			return false
		}
	}

	config := &oauth2.Config{}
	ts := config.TokenSource(context.Background(), &oauth2.Token{AccessToken: token})
	b.configureClient(oauth2.NewClient(context.Background(), ts))

	return true
}
