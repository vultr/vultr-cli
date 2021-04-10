package users

import (
	"context"
	"fmt"
	"github.com/vultr/govultr/v2"
	"testing"
)

type mockVultrUser struct {
	client *govultr.Client
}

func (m mockVultrUser) Create(ctx context.Context, userCreate *govultr.UserReq) (*govultr.User, error) {
	panic("implement me")
}

func (m mockVultrUser) Get(ctx context.Context, userID string) (*govultr.User, error) {
	return &govultr.User{
		ID:         "1234",
		Name:       "CLI Tests",
		Email:      "cli@vultr.com",
		APIEnabled: govultr.BoolToBoolPtr(true),
		APIKey:     "1234",
		ACL:        []string{"test"},
	}, nil
}

func (m mockVultrUser) Update(ctx context.Context, userID string, userReq *govultr.UserReq) error {
	panic("implement me")
}

func (m mockVultrUser) Delete(ctx context.Context, userID string) error {
	panic("implement me")
}

func (m mockVultrUser) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.User, *govultr.Meta, error) {
	panic("implement me")
}

func TestUserOptions_Get(t *testing.T) {
	client := &govultr.Client{User: mockVultrUser{nil}}
	u := UserOptions{Client: client, Args: []string{"1234"}}
	user := u.Get()
	fmt.Println(user)
}
