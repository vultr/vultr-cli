package backups

import (
	"context"
	"reflect"
	"testing"

	"github.com/vultr/govultr/v2"
	"github.com/vultr/vultr-cli/pkg/cli"
)

type mockVultrBackup struct {
	client *govultr.Client
}

func setup() *cli.Base {
	return &cli.Base{Client: &govultr.Client{Backup: mockVultrBackup{nil}}}
}

func (m mockVultrBackup) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.Backup, *govultr.Meta, error) {
	return []govultr.Backup{{
			ID:          "041954ba-6c1c-4bd8-bd12-6948c4709b5c",
			DateCreated: "2021-05-09T19:02:06+00:00",
			Description: "my-backup",
			Size:        100000,
			Status:      "complete",
		}}, &govultr.Meta{
			Total: 0,
			Links: nil,
		}, nil
}

func (m mockVultrBackup) Get(ctx context.Context, backupID string) (*govultr.Backup, error) {
	return &govultr.Backup{
		ID:          "041954ba-6c1c-4bd8-bd12-6948c4709b5c",
		DateCreated: "2021-05-09T19:02:06+00:00",
		Description: "my-backup",
		Size:        100000,
		Status:      "complete",
	}, nil
}

func TestOptions_List(t *testing.T) {
	backup := NewBackupOptions(&cli.Base{Client: &govultr.Client{Backup: mockVultrBackup{nil}}})

	expectedBackup := []govultr.Backup{{
		ID:          "041954ba-6c1c-4bd8-bd12-6948c4709b5c",
		DateCreated: "2021-05-09T19:02:06+00:00",
		Description: "my-backup",
		Size:        100000,
		Status:      "complete",
	}}

	expectedMeta := &govultr.Meta{
		Total: 0,
		Links: nil,
	}

	backupList, meta, _ := backup.List()

	if !reflect.DeepEqual(backupList, expectedBackup) {
		t.Errorf("BackupOptions.list returned %v expected %v", backupList, expectedMeta)
	}

	if !reflect.DeepEqual(expectedMeta, meta) {
		t.Errorf("BackupOptions.meta returned %v expected %v", meta, expectedMeta)
	}
}

func TestOptions_Get(t *testing.T) {
	backup := NewBackupOptions(setup())
	backup.Base.Args = []string{"my-backup"}

	expected := &govultr.Backup{
		ID:          "041954ba-6c1c-4bd8-bd12-6948c4709b5c",
		DateCreated: "2021-05-09T19:02:06+00:00",
		Description: "my-backup",
		Size:        100000,
		Status:      "complete",
	}

	backups, err := backup.Get()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(expected, backups) {
		t.Errorf("Options.get returned %v expected %v", backup, expected)
	}
}

func TestNewCmdBackup(t *testing.T) {
	cmd := NewCmdBackup(&cli.Base{Client: &govultr.Client{Backup: mockVultrBackup{nil}}})

	if cmd.Short != "backup commands" {
		t.Errorf("invalid short")
	}

	if cmd.Use != "backup" {
		t.Errorf("invalid backup")
	}

	alias := []string{"backups", "backup", "b"}
	if !reflect.DeepEqual(cmd.Aliases, alias) {
		t.Errorf("expected alias %v got %v", alias, cmd.Aliases)
	}
}

func TestNewBackupOptions(t *testing.T) {
	options := NewBackupOptions(&cli.Base{Client: &govultr.Client{Backup: mockVultrBackup{nil}}})

	ref := reflect.TypeOf(options)
	if _, ok := ref.MethodByName("List"); !ok {
		t.Errorf("Missing list function")
	}

	if _, ok := ref.MethodByName("validate"); ok {
		t.Errorf("validate isn't exported and shouldn't be accessible")
	}

	bInterface := reflect.TypeOf(new(BackupOptionsInterface)).Elem()
	if !ref.Implements(bInterface) {
		t.Errorf("Options does not implement BackupOptionsInterface")
	}
}
