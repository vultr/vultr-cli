package govultr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

const vpcPath = "/v2/vpcs"

// VPCService is the interface to interact with the VPC endpoints on the Vultr API
// Link : https://www.vultr.com/api/#tag/vpcs
type VPCService interface {
	Create(ctx context.Context, createReq *VPCReq) (*VPC, *http.Response, error)
	Get(ctx context.Context, vpcID string) (*VPC, *http.Response, error)
	Update(ctx context.Context, vpcID string, description string) error
	Delete(ctx context.Context, vpcID string) error
	List(ctx context.Context, options *ListOptions) ([]VPC, *Meta, *http.Response, error)
}

// VPCServiceHandler handles interaction with the VPC methods for the Vultr API
type VPCServiceHandler struct {
	client *Client
}

// VPC represents a Vultr VPC
type VPC struct {
	ID           string `json:"id"`
	Region       string `json:"region"`
	Description  string `json:"description"`
	V4Subnet     string `json:"v4_subnet"`
	V4SubnetMask int    `json:"v4_subnet_mask"`
	DateCreated  string `json:"date_created"`
}

// VPCReq represents parameters to create or update a VPC resource
type VPCReq struct {
	Region       string `json:"region"`
	Description  string `json:"description"`
	V4Subnet     string `json:"v4_subnet"`
	V4SubnetMask int    `json:"v4_subnet_mask"`
}

type vpcsBase struct {
	VPCs []VPC `json:"vpcs"`
	Meta *Meta `json:"meta"`
}

type vpcBase struct {
	VPC *VPC `json:"vpc"`
}

// Create creates a new VPC. A VPC can only be used at the location for which it was created.
func (n *VPCServiceHandler) Create(ctx context.Context, createReq *VPCReq) (*VPC, *http.Response, error) {
	req, err := n.client.NewRequest(ctx, http.MethodPost, vpcPath, createReq)
	if err != nil {
		return nil, nil, err
	}

	vpc := new(vpcBase)
	resp, err := n.client.DoWithContext(ctx, req, vpc)
	if err != nil {
		return nil, resp, err
	}

	return vpc.VPC, resp, nil
}

// Get gets the VPC of the requested ID
func (n *VPCServiceHandler) Get(ctx context.Context, vpcID string) (*VPC, *http.Response, error) {
	uri := fmt.Sprintf("%s/%s", vpcPath, vpcID)
	req, err := n.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	vpc := new(vpcBase)
	resp, err := n.client.DoWithContext(ctx, req, vpc)
	if err != nil {
		return nil, resp, err
	}

	return vpc.VPC, resp, nil
}

// Update updates a VPC
func (n *VPCServiceHandler) Update(ctx context.Context, vpcID string, description string) error {
	uri := fmt.Sprintf("%s/%s", vpcPath, vpcID)

	vpcReq := RequestBody{"description": description}
	req, err := n.client.NewRequest(ctx, http.MethodPut, uri, vpcReq)
	if err != nil {
		return err
	}

	_, err = n.client.DoWithContext(ctx, req, nil)
	return err
}

// Delete deletes a VPC. Before deleting, a VPC must be disabled from all instances
func (n *VPCServiceHandler) Delete(ctx context.Context, vpcID string) error {
	uri := fmt.Sprintf("%s/%s", vpcPath, vpcID)
	req, err := n.client.NewRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}
	_, err = n.client.DoWithContext(ctx, req, nil)
	return err
}

// List lists all VPCs on the current account
func (n *VPCServiceHandler) List(ctx context.Context, options *ListOptions) ([]VPC, *Meta, *http.Response, error) {
	req, err := n.client.NewRequest(ctx, http.MethodGet, vpcPath, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	newValues, err := query.Values(options)
	if err != nil {
		return nil, nil, nil, err
	}

	req.URL.RawQuery = newValues.Encode()

	vpcs := new(vpcsBase)
	resp, err := n.client.DoWithContext(ctx, req, vpcs)
	if err != nil {
		return nil, nil, resp, err
	}

	return vpcs.VPCs, vpcs.Meta, resp, nil
}
