package applications

import (
	"context"
	"reflect"
	"testing"

	"github.com/vultr/vultr-cli/pkg/cli"

	"github.com/vultr/govultr/v2"
)

type mockVultrApplications struct {
	client *govultr.Client
}

func (m mockVultrApplications) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.Application, *govultr.Meta, error) {
	return []govultr.Application{
			{
				ID:         1,
				Name:       "LEMP",
				ShortName:  "LEMP",
				DeployName: "Lemp on CentOS 6",
			},
		}, &govultr.Meta{
			Total: 1,
			Links: nil,
		}, nil
}

func TestOptions_List(t *testing.T) {
	appOption := NewApplicationOptions(&cli.Base{Client: &govultr.Client{Application: mockVultrApplications{nil}}})

	expected := []govultr.Application{
		{
			ID:         1,
			Name:       "LEMP",
			ShortName:  "LEMP",
			DeployName: "Lemp on CentOS 6",
		},
	}
	expectedMeta := &govultr.Meta{
		Total: 1,
		Links: nil,
	}

	app, meta, _ := appOption.List()

	if !reflect.DeepEqual(expected, app) {
		t.Errorf("AppOptions.list returned %v expected %v", app, expected)
	}

	if !reflect.DeepEqual(expectedMeta, meta) {
		t.Errorf("AppOptions.list returned %v expected %v", meta, expectedMeta)
	}
}

func TestNewApplicationOptions(t *testing.T) {
	appOption := NewApplicationOptions(&cli.Base{Client: &govultr.Client{Application: mockVultrApplications{nil}}})

	ref := reflect.TypeOf(appOption)
	if _, ok := ref.MethodByName("List"); !ok {
		t.Errorf("Missing list function")
	}

	if _, ok := ref.MethodByName("validate"); ok {
		t.Errorf("validate isn't exported shouldn't be accessible")
	}

	pInterface := reflect.TypeOf(new(Interface)).Elem()
	if !ref.Implements(pInterface) {
		t.Errorf("Options does not implement Interface")
	}
}

func TestNewCmdApplications(t *testing.T) {
	cmd := NewCmdApplications(&cli.Base{Client: &govultr.Client{Application: mockVultrApplications{nil}}})

	if cmd.Short != "display applications" {
		t.Errorf("invalid short")
	}

	if cmd.Use != "apps" {
		t.Errorf("invalid apps")
	}

	alias := []string{"a", "application", "applications", "app"}
	if !reflect.DeepEqual(cmd.Aliases, alias) {
		t.Errorf("expected alias %v got %v", alias, cmd.Aliases)
	}
}
