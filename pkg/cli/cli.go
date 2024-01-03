package cli

import (
	"context"

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
}

// NewCLIBase creates new base struct
func NewCLIBase(apiKey, output string) *Base {
	base := new(Base)
	base.configurePrinter()
	base.configureClient(apiKey)
	return base
}

func (b *Base) configureClient(apiKey string) {
	config := &oauth2.Config{}
	ts := config.TokenSource(context.Background(), &oauth2.Token{AccessToken: apiKey})
	b.Client = govultr.NewClient(oauth2.NewClient(context.Background(), ts))
}

func (b *Base) configurePrinter() {
	b.Printer = &printer.Output{}
}
