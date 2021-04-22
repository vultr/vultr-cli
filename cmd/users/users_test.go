package users

import (
	"context"
	"testing"

	"github.com/vultr/govultr/v2"
)

type mockVultrUser struct {
	client *govultr.Client
}

func (m mockVultrUser) Create(ctx context.Context, userCreate *govultr.UserReq) (*govultr.User, error) {
	panic("implement me")
}

func (m mockVultrUser) Get(ctx context.Context, userID string) (*govultr.User, error) {
	return &govultr.User{
		ID:     "1234",
		Name:   "CLI Tests",
		Email:  "cli@vultr.com",
		APIKey: "1234",
		ACL:    []string{"test"},
	}, nil
}

func (m mockVultrUser) Update(ctx context.Context, userID string, userReq *govultr.UserReq) error {
	panic("implement me")
}

func (m mockVultrUser) Delete(ctx context.Context, userID string) error {
	return nil
}

func (m mockVultrUser) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.User, *govultr.Meta, error) {
	panic("implement me")
}

func TestUserOptions_Create(t *testing.T) {
	panic("implement me")
}

func TestUserOptions_Get(t *testing.T) {
	//client := &govultr.Client{User: mockVultrUser{nil}}
	//
	//userExpected := &govultr.User{
	//	ID:     "1234",
	//	Name:   "CLI Tests",
	//	Email:  "cli@vultr.com",
	//	APIKey: "1234",
	//	ACL:    []string{"test"},
	//}
	//
	//u := UserOptions{Client: client, Args: []string{"1234"}}
	//user := u.Get(context.Background(), u)
	//if !reflect.DeepEqual(userExpected, user) {
	//	t.Errorf("UserOptions.get returned %+v, expected %+v", user, userExpected)
	//}
}

func TestUserOptions_Delete(t *testing.T) {
	client := &govultr.Client{User: mockVultrUser{nil}}

	u := UserOptions{Client: client, Args: []string{"1234"}}
	u.Delete()
}
