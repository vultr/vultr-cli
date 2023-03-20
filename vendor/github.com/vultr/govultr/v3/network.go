package govultr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

const netPath = "/v2/private-networks"

// NetworkService is the interface to interact with the network endpoints on the Vultr API
// Link : https://www.vultr.com/api/#tag/private-Networks
// Deprecated: NetworkService should no longer be used. Instead, use VPCService.
type NetworkService interface {
	// Deprecated: NetworkService Create should no longer be used. Instead, use VPCService Create.
	Create(ctx context.Context, createReq *NetworkReq) (*Network, *http.Response, error)
	// Deprecated: NetworkService Get should no longer be used. Instead, use VPCService Get.
	Get(ctx context.Context, networkID string) (*Network, *http.Response, error)
	// Deprecated: NetworkService Update should no longer be used. Instead, use VPCService Update.
	Update(ctx context.Context, networkID string, description string) error
	// Deprecated: NetworkService Delete should no longer be used. Instead, use VPCService Delete.
	Delete(ctx context.Context, networkID string) error
	// Deprecated: NetworkService List should no longer be used. Instead, use VPCService List.
	List(ctx context.Context, options *ListOptions) ([]Network, *Meta, *http.Response, error)
}

// NetworkServiceHandler handles interaction with the network methods for the Vultr API
// Deprecated: NetworkServiceHandler should no longer be used. Instead, use VPCServiceHandler.
type NetworkServiceHandler struct {
	client *Client
}

// Network represents a Vultr private network
// Deprecated: Network should no longer be used. Instead, use VPC.
type Network struct {
	NetworkID    string `json:"id"`
	Region       string `json:"region"`
	Description  string `json:"description"`
	V4Subnet     string `json:"v4_subnet"`
	V4SubnetMask int    `json:"v4_subnet_mask"`
	DateCreated  string `json:"date_created"`
}

// NetworkReq represents parameters to create or update Network resource
// Deprecated: NetworkReq should no longer be used. Instead, use VPCReq.
type NetworkReq struct {
	Region       string `json:"region"`
	Description  string `json:"description"`
	V4Subnet     string `json:"v4_subnet"`
	V4SubnetMask int    `json:"v4_subnet_mask"`
}

type networksBase struct {
	Networks []Network `json:"networks"`
	Meta     *Meta     `json:"meta"`
}

type networkBase struct {
	Network *Network `json:"network"`
}

// Create a new private network. A private network can only be used at the location for which it was created.
// Deprecated: NetworkServiceHandler Create should no longer be used. Instead, use VPCServiceHandler Create.
func (n *NetworkServiceHandler) Create(ctx context.Context, createReq *NetworkReq) (*Network, *http.Response, error) {
	req, err := n.client.NewRequest(ctx, http.MethodPost, netPath, createReq)
	if err != nil {
		return nil, nil, err
	}

	network := new(networkBase)
	resp, err := n.client.DoWithContext(ctx, req, network)
	if err != nil {
		return nil, resp, err
	}

	return network.Network, resp, nil
}

// Get gets the private networks of the requested ID
// Deprecated: NetworkServiceHandler Get should no longer be used.  Instead use VPCServiceHandler Create.
func (n *NetworkServiceHandler) Get(ctx context.Context, networkID string) (*Network, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s", netPath, networkID)
	req, err := n.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	network := new(networkBase)
	resp, err := n.client.DoWithContext(ctx, req, network)
	if err != nil {
		return nil, resp, err
	}

	return network.Network, resp, nil
}

// Update updates a private network
// Deprecated: NetworkServiceHandler Update should no longer be used. Instead, use VPCServiceHandler Update.
func (n *NetworkServiceHandler) Update(ctx context.Context, networkID string, description string) error {
	uri := fmt.Sprintf("%s/%s", netPath, networkID)

	netReq := RequestBody{"description": description}
	req, err := n.client.NewRequest(ctx, http.MethodPut, uri, netReq)
	if err != nil {
		return err
	}
	_, err = n.client.DoWithContext(ctx, req, nil)
	return err
}

// Delete a private network. Before deleting, a network must be disabled from all instances
// Deprecated: NetworkServiceHandler Delete should no longer be used. Instead, use VPCServiceHandler Delete.
func (n *NetworkServiceHandler) Delete(ctx context.Context, networkID string) error {
	uri := fmt.Sprintf("%s/%s", netPath, networkID)
	req, err := n.client.NewRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	_, err = n.client.DoWithContext(ctx, req, nil)
	return err
}

// List lists all private networks on the current account
// Deprecated: NetworkServiceHandler List should no longer be used. Instead, use VPCServiceHandler List.
func (n *NetworkServiceHandler) List(ctx context.Context, options *ListOptions) ([]Network, *Meta, *http.Response, error) {
	req, err := n.client.NewRequest(ctx, http.MethodGet, netPath, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	newValues, err := query.Values(options)
	if err != nil {
		return nil, nil, nil, err
	}

	req.URL.RawQuery = newValues.Encode()

	networks := new(networksBase)
	resp, err := n.client.DoWithContext(ctx, req, networks)
	if err != nil {
		return nil, nil, resp, err
	}

	return networks.Networks, networks.Meta, resp, nil
}
