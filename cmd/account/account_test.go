package account

import (
	"context"
	"reflect"
	"testing"

	"github.com/vultr/vultr-cli/pkg/cli"

	"github.com/vultr/govultr/v2"
)

type mockVultrAccount struct {
	client *govultr.Client
}

func (m mockVultrAccount) Get(ctx context.Context) (*govultr.Account, error) {
	return &govultr.Account{
		Balance:        10,
		PendingCharges: 100,
		Name:           "John Smith",
		Email:          "john@example.com",
		ACL:            []string{"manage_users", "subscriptions", "billing"},
	}, nil
}

func TestNewAccountOptions(t *testing.T) {
	accountOption := NewAccountOptions(&cli.Base{Client: &govultr.Client{Account: mockVultrAccount{nil}}})

	ref := reflect.TypeOf(accountOption)
	if _, ok := ref.MethodByName("Get"); !ok {
		t.Errorf("Missing get function")
	}

	if _, ok := ref.MethodByName("validate"); ok {
		t.Errorf("validate isn't exported shouldn't be accessible")
	}

	aInterface := reflect.TypeOf(new(AccountInterface)).Elem()
	if !ref.Implements(aInterface) {
		t.Errorf("Options does not implement AccountInterface")
	}
}

func TestNewCmdAccount(t *testing.T) {
	cmd := NewCmdAccount(&cli.Base{Client: &govultr.Client{Account: mockVultrAccount{nil}}})

	if cmd.Short != "get account information" {
		t.Errorf("invalid short")
	}

	if cmd.Use != "account" {
		t.Errorf("invalid account")
	}

}

func TestOptions_Get(t *testing.T) {
	a := NewAccountOptions(&cli.Base{Client: &govultr.Client{Account: mockVultrAccount{nil}}})

	expectedAccount := &govultr.Account{
		Balance:        10,
		PendingCharges: 100,
		Name:           "John Smith",
		Email:          "john@example.com",
		ACL:            []string{"manage_users", "subscriptions", "billing"},
	}

	account, _ := a.Get()

	if !reflect.DeepEqual(account, expectedAccount) {
		t.Errorf("OSOptions.list returned %v expected %v", account, expectedAccount)
	}

}
