package operatingSystems

import (
	"context"
	"reflect"
	"testing"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/pkg/cli"
)

type mockVultrOS struct {
	client *govultr.Client
}

func (m mockVultrOS) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.OS, *govultr.Meta, error) {
	return []govultr.OS{{
			ID:     127,
			Name:   "CentOS 6 x64",
			Arch:   "x64",
			Family: "centos",
		}}, &govultr.Meta{
			Total: 0,
			Links: nil,
		}, nil
}

func TestOptions_List(t *testing.T) {
	os := NewOSOptions(&cli.Base{Client: &govultr.Client{OS: mockVultrOS{nil}}})

	expectedOS := []govultr.OS{{
		ID:     127,
		Name:   "CentOS 6 x64",
		Arch:   "x64",
		Family: "centos",
	}}

	expectedMeta := &govultr.Meta{
		Total: 0,
		Links: nil,
	}

	osList, meta, _ := os.List()

	if !reflect.DeepEqual(osList, expectedOS) {
		t.Errorf("OSOptions.list returned %v expected %v", osList, expectedMeta)
	}

	if !reflect.DeepEqual(expectedMeta, meta) {
		t.Errorf("PlanOptions.metal returned %v expected %v", meta, expectedMeta)
	}
}

func TestNewCmdOS(t *testing.T) {
	cmd := NewCmdOS(&cli.Base{Client: &govultr.Client{OS: mockVultrOS{nil}}})

	if cmd.Short != "list available operating systems" {
		t.Errorf("invalid short")
	}

	if cmd.Use != "os" {
		t.Errorf("invalid os")
	}

	alias := []string{"o"}
	if !reflect.DeepEqual(cmd.Aliases, alias) {
		t.Errorf("expected alias %v got %v", alias, cmd.Aliases)
	}
}

func TestNewOSOptions(t *testing.T) {
	options := NewOSOptions(&cli.Base{Client: &govultr.Client{OS: mockVultrOS{nil}}})

	ref := reflect.TypeOf(options)
	if _, ok := ref.MethodByName("List"); !ok {
		t.Errorf("Missing list function")
	}

	if _, ok := ref.MethodByName("validate"); ok {
		t.Errorf("validate isn't exported shouldn't be accessible")
	}

	oInterface := reflect.TypeOf(new(Interface)).Elem()
	if !ref.Implements(oInterface) {
		t.Errorf("Options does not implement Interface")
	}
}
