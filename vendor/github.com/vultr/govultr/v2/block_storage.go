package govultr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

type RequestBody map[string]interface{}

// BlockStorageService is the interface to interact with Block-Storage endpoint on the Vultr API
type BlockStorageService interface {
	Create(ctx context.Context, blockReq *BlockStorageCreate) (*BlockStorage, error)
	Get(ctx context.Context, blockID string) (*BlockStorage, error)
	Update(ctx context.Context, blockID string, blockReq *BlockStorageUpdate) error
	Delete(ctx context.Context, blockID string) error
	List(ctx context.Context, options *ListOptions) ([]BlockStorage, *Meta, error)

	Attach(ctx context.Context, blockID string, attach *BlockStorageAttach) error
	Detach(ctx context.Context, blockID string, detach *BlockStorageDetach) error
	Resize(ctx context.Context, blockID string, sizeGB int) error
}

// BlockStorageServiceHandler handles interaction with the block-storage methods for the Vultr API
type BlockStorageServiceHandler struct {
	client *Client
}

// BlockStorage represents Vultr Block-Storage
type BlockStorage struct {
	ID                 string  `json:"id"`
	Cost               float32 `json:"cost"`
	Status             string  `json:"status"`
	SizeGB             int     `json:"size_gb"`
	Region             string  `json:"region"`
	DateCreated        string  `json:"date_created"`
	AttachedToInstance string  `json:"attached_to_instance"`
	Label              string  `json:"label"`
}

// BlockStorageReq
type BlockStorageCreate struct {
	Region string `json:"region"`
	SizeGB int    `json:"size_gb"`
	Label  string `json:"label,omitempty"`
}

// BlockStorageUpdate
type BlockStorageUpdate struct {
	SizeGB int    `json:"size_gb,omitempty"`
	Label  string `json:"label,omitempty"`
}

// BlockStorageAttach
type BlockStorageAttach struct {
	InstanceID string `json:"instance_id"`
	Live       bool   `json:"live,omitempty"`
}

// BlockStorageDetach
type BlockStorageDetach struct {
	Live bool `json:"live,omitempty"`
}

type blockStoragesBase struct {
	Blocks []BlockStorage `json:"blocks"`
	Meta   *Meta          `json:"meta"`
}

type blockStorageBase struct {
	Block *BlockStorage `json:"block"`
}

// Create builds out a block storage
func (b *BlockStorageServiceHandler) Create(ctx context.Context, blockReq *BlockStorageCreate) (*BlockStorage, error) {
	uri := "/v2/blocks"

	req, err := b.client.NewRequest(ctx, http.MethodPost, uri, blockReq)
	if err != nil {
		return nil, err
	}

	block := new(blockStorageBase)
	if err = b.client.DoWithContext(ctx, req, block); err != nil {
		return nil, err
	}

	return block.Block, nil
}

// Get returns a single block storage instance based ony our blockID you provide from your Vultr Account
func (b *BlockStorageServiceHandler) Get(ctx context.Context, blockID string) (*BlockStorage, error) {
	uri := fmt.Sprintf("/v2/blocks/%s", blockID)

	req, err := b.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	block := new(blockStorageBase)
	if err = b.client.DoWithContext(ctx, req, block); err != nil {
		return nil, err
	}

	return block.Block, nil
}

// SetLabel allows you to set/update the label on your Vultr Block storage
func (b *BlockStorageServiceHandler) Update(ctx context.Context, blockID string, blockReq *BlockStorageUpdate) error {
	uri := fmt.Sprintf("/v2/blocks/%s", blockID)

	req, err := b.client.NewRequest(ctx, http.MethodPatch, uri, blockReq)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// Delete will remove block storage instance from your Vultr account
func (b *BlockStorageServiceHandler) Delete(ctx context.Context, blockID string) error {
	uri := fmt.Sprintf("/v2/blocks/%s", blockID)

	req, err := b.client.NewRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// List returns a list of all block storage instances on your Vultr Account
func (b *BlockStorageServiceHandler) List(ctx context.Context, options *ListOptions) ([]BlockStorage, *Meta, error) {
	uri := "/v2/blocks"

	req, err := b.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	newValues, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = newValues.Encode()

	blocks := new(blockStoragesBase)
	if err = b.client.DoWithContext(ctx, req, blocks); err != nil {
		return nil, nil, err
	}

	return blocks.Blocks, blocks.Meta, nil
}

// Attach will link a given block storage to a given Vultr instance
// If Live is set to true the block storage will be attached without reloading the instance
func (b *BlockStorageServiceHandler) Attach(ctx context.Context, blockID string, attach *BlockStorageAttach) error {
	uri := fmt.Sprintf("/v2/blocks/%s/attach", blockID)

	req, err := b.client.NewRequest(ctx, http.MethodPost, uri, attach)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// Detach will de-link a given block storage to the Vultr instance it is attached to
// If Live is set to true the block storage will be detached without reloading the instance
func (b *BlockStorageServiceHandler) Detach(ctx context.Context, blockID string, detach *BlockStorageDetach) error {
	uri := fmt.Sprintf("/v2/blocks/%s/detach", blockID)

	req, err := b.client.NewRequest(ctx, http.MethodPost, uri, detach)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// Resize allows you to resize your Vultr block storage
func (b *BlockStorageServiceHandler) Resize(ctx context.Context, blockID string, sizeGB int) error {
	uri := fmt.Sprintf("/v2/blocks/%s/resize", blockID)
	body := &RequestBody{"size_gb": sizeGB}

	req, err := b.client.NewRequest(ctx, http.MethodPost, uri, body)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}
