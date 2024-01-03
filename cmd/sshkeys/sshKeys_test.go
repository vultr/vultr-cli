package sshkeys

import (
	"context"
	"reflect"
	"testing"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/pkg/cli"
)

type MockVultrSSHKeys struct {
	client *govultr.Client
}

func (m MockVultrSSHKeys) Create(ctx context.Context, sshKeyReq *govultr.SSHKeyReq) (*govultr.SSHKey, error) {
	return &govultr.SSHKey{
		ID:          "71023f44-ac87-4073-a3d0-c05677096a4d",
		Name:        "test",
		SSHKey:      "ssh-rsa test",
		DateCreated: "2021-05-01T18:50:06+00:00",
	}, nil
}

func (m MockVultrSSHKeys) Get(ctx context.Context, sshKeyID string) (*govultr.SSHKey, error) {
	return &govultr.SSHKey{
		ID:          "71023f44-ac87-4073-a3d0-c05677096a4d",
		Name:        "test",
		SSHKey:      "ssh-rsa test",
		DateCreated: "2021-05-01T18:50:06+00:00",
	}, nil
}

func (m MockVultrSSHKeys) Update(ctx context.Context, sshKeyID string, sshKeyReq *govultr.SSHKeyReq) error {
	return nil
}

func (m MockVultrSSHKeys) Delete(ctx context.Context, sshKeyID string) error {
	return nil
}

func (m MockVultrSSHKeys) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.SSHKey, *govultr.Meta, error) {
	return []govultr.SSHKey{
			{
				ID:          "71023f44-ac87-4073-a3d0-c05677096a4d",
				Name:        "test",
				SSHKey:      "ssh-rsa test",
				DateCreated: "2021-05-01T18:50:06+00:00",
			},
		}, &govultr.Meta{
			Total: 0,
			Links: nil,
		},
		nil
}

func setup() *cli.Base {
	return &cli.Base{Client: &govultr.Client{SSHKey: MockVultrSSHKeys{nil}}}
}

func TestOptions_Create(t *testing.T) {
	ssh := NewSSHKeyOptions(setup())

	expected := &govultr.SSHKey{
		ID:          "71023f44-ac87-4073-a3d0-c05677096a4d",
		Name:        "test",
		SSHKey:      "ssh-rsa test",
		DateCreated: "2021-05-01T18:50:06+00:00",
	}

	sshKey, err := ssh.Create()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(expected, sshKey) {
		t.Errorf("Options.create returned %v expected %v", sshKey, expected)
	}
}

func TestOptions_Get(t *testing.T) {
	ssh := NewSSHKeyOptions(setup())
	ssh.Base.Args = []string{"test"}

	expected := &govultr.SSHKey{
		ID:          "71023f44-ac87-4073-a3d0-c05677096a4d",
		Name:        "test",
		SSHKey:      "ssh-rsa test",
		DateCreated: "2021-05-01T18:50:06+00:00",
	}

	sshKey, err := ssh.Get()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(expected, sshKey) {
		t.Errorf("Options.get returned %v expected %v", sshKey, expected)
	}
}

func TestOptions_List(t *testing.T) {
	ssh := NewSSHKeyOptions(setup())
	ssh.Base.Args = []string{"test"}

	expected := []govultr.SSHKey{{
		ID:          "71023f44-ac87-4073-a3d0-c05677096a4d",
		Name:        "test",
		SSHKey:      "ssh-rsa test",
		DateCreated: "2021-05-01T18:50:06+00:00",
	}}
	expectedMeta := &govultr.Meta{
		Total: 0,
		Links: nil,
	}

	sshKey, meta, err := ssh.List()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(expected, sshKey) {
		t.Errorf("Options.list returned %v expected %v", sshKey, expected)
	}
	if !reflect.DeepEqual(expectedMeta, meta) {
		t.Errorf("Options.list returned %v expected %v", meta, expectedMeta)
	}
}

func TestOptions_Update(t *testing.T) {
	ssh := NewSSHKeyOptions(setup())
	ssh.Base.Args = []string{"test"}

	err := ssh.Update()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestOptions_Delete(t *testing.T) {
	ssh := NewSSHKeyOptions(setup())
	ssh.Base.Args = []string{"test"}

	err := ssh.Delete()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestNewCmdSSHKey(t *testing.T) {
	cmd := NewCmdSSHKey(setup())

	if cmd.Short != "ssh-key commands" {
		t.Errorf("invalid short")
	}

	if cmd.Use != "ssh-key" {
		t.Errorf("invalid ssh-key")
	}

	alias := []string{"ssh", "ssh-keys", "sshkeys"}
	if !reflect.DeepEqual(cmd.Aliases, alias) {
		t.Errorf("expected alias %v got %v", alias, cmd.Aliases)
	}
}

func TestNewSSHKeyOptions(t *testing.T) {
	options := NewSSHKeyOptions(setup())

	ref := reflect.TypeOf(options)
	if _, ok := ref.MethodByName("List"); !ok {
		t.Errorf("Missing list function")
	}

	if _, ok := ref.MethodByName("Create"); !ok {
		t.Errorf("Missing list function")
	}

	if _, ok := ref.MethodByName("Update"); !ok {
		t.Errorf("Missing Update function")
	}
	if _, ok := ref.MethodByName("Get"); !ok {
		t.Errorf("Missing Get function")
	}
	if _, ok := ref.MethodByName("Delete"); !ok {
		t.Errorf("Missing delete function")
	}
	if _, ok := ref.MethodByName("validate"); ok {
		t.Errorf("validate isn't exported shouldn't be accessible")
	}

	rInterface := reflect.TypeOf(new(Interface)).Elem()
	if !ref.Implements(rInterface) {
		t.Errorf("Options does not implement Interface")
	}
}
