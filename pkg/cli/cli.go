package cli

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"golang.org/x/oauth2"
)

// BaseInterface that is required for any struct that is used as a base
type BaseInterface interface {
	configureClient(apiKey string)
	configurePrinter()
}

// Base contains the basic needs for all CLI commands
type Base struct {
	Args    []string
	Client  *govultr.Client
	Options *govultr.ListOptions
	Printer *printer.Output
	Context context.Context
	HasAuth bool
}

// NewCLIBase creates new base struct
func NewCLIBase(apiKey, userAgent, output string) *Base {
	base := new(Base)
	base.configurePrinter()
	base.configureClient(apiKey, userAgent)
	base.configureContext()
	return base
}

func (b *Base) configureClient(apiKey, userAgent string) {
	var token string
	b.HasAuth = false

	token = viper.GetString("api-key")
	if token == "" {
		token = apiKey
	}

	if token == "" {
		b.Client = govultr.NewClient(nil)
	} else {
		config := &oauth2.Config{}
		ts := config.TokenSource(context.Background(), &oauth2.Token{AccessToken: apiKey})
		b.Client = govultr.NewClient(oauth2.NewClient(context.Background(), ts))
		b.HasAuth = true
	}

	b.Client.SetRateLimit(1 * time.Second)
	b.Client.SetUserAgent(userAgent)
}

func (b *Base) configurePrinter() {
	b.Printer = &printer.Output{}
}

func (b *Base) configureContext() {
	b.Context = context.Background()
}
