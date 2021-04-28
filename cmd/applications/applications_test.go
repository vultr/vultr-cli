package applications

import (
	"context"
	"github.com/vultr/govultr/v2"
	"reflect"
	"testing"
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
	client := &govultr.Client{Application: mockVultrApplications{nil}}
	appOption := NewApplicationOptions(client)

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
		t.Errorf("PlanOptions.list returned %v expected %v", app, expected)
	}

	if !reflect.DeepEqual(expectedMeta, meta) {
		t.Errorf("PlanOptions.list returned %v expected %v", meta, expectedMeta)
	}
}

func TestNewApplicationOptions(t *testing.T) {
	client := &govultr.Client{Application: mockVultrApplications{nil}}
	appOption := NewApplicationOptions(client)

	ref := reflect.TypeOf(appOption)
	if _, ok := ref.MethodByName("List"); !ok {
		t.Errorf("Missing list function")
	}


	if _, ok := ref.MethodByName("validate"); ok {
		t.Errorf("validate isn't exported shouldn't be accessible")
	}

	pInterface := reflect.TypeOf(new(Interface)).Elem()
	if !ref.Implements(pInterface){
		t.Errorf("Options does not implement Interface")
	}
}

func TestNewCmdApplications(t *testing.T) {
	client := &govultr.Client{Application: mockVultrApplications{nil}}
	cmd := NewCmdApplications(client)


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