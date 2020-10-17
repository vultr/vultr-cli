package govultr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

const bmPath = "/v2/bare-metals"

// BareMetalServerService is the interface to interact with the bare metal endpoints on the Vultr API
type BareMetalServerService interface {
	Create(ctx context.Context, bmCreate *BareMetalReq) (*BareMetalServer, error)
	Get(ctx context.Context, serverID string) (*BareMetalServer, error)
	Update(ctx context.Context, serverID string, bmReq *BareMetalReq) error
	Delete(ctx context.Context, serverID string) error
	List(ctx context.Context, options *ListOptions) ([]BareMetalServer, *Meta, error)

	GetBandwidth(ctx context.Context, serverID string) (*Bandwidth, error)

	ListIPv4s(ctx context.Context, serverID string, options *ListOptions) ([]IPv4, *Meta, error)
	ListIPv6s(ctx context.Context, serverID string, options *ListOptions) ([]IPv6, *Meta, error)

	Halt(ctx context.Context, serverID string) error
	Reboot(ctx context.Context, serverID string) error
	Reinstall(ctx context.Context, serverID string) error

	MassStart(ctx context.Context, serverList []string) error
	MassHalt(ctx context.Context, serverList []string) error
	MassReboot(ctx context.Context, serverList []string) error
}

// BareMetalServerServiceHandler handles interaction with the bare metal methods for the Vultr API
type BareMetalServerServiceHandler struct {
	client *Client
}

// BareMetalServer represents a bare metal server on Vultr
type BareMetalServer struct {
	ID              string   `json:"id"`
	Os              string   `json:"os"`
	RAM             string   `json:"ram"`
	Disk            string   `json:"disk"`
	MainIP          string   `json:"main_ip"`
	CPUCount        int      `json:"cpu_count"`
	Region          string   `json:"region"`
	DefaultPassword string   `json:"default_password"`
	DateCreated     string   `json:"date_created"`
	Status          string   `json:"status"`
	NetmaskV4       string   `json:"netmask_v4"`
	GatewayV4       string   `json:"gateway_v4"`
	Plan            string   `json:"plan"`
	V6Network       string   `json:"v6_network"`
	V6MainIP        string   `json:"v6_main_ip"`
	V6NetworkSize   int      `json:"v6_network_size"`
	Label           string   `json:"label"`
	Tag             string   `json:"tag"`
	OsID            int      `json:"os_id"`
	AppID           int      `json:"app_id"`
	Features        []string `json:"features"`
}

// BareMetalReq represents the optional parameters that can be set when creating or updating a bare metal server
type BareMetalReq struct {
	Region          string   `json:"region,omitempty"`
	Plan            string   `json:"plan,omitempty"`
	OsID            int      `json:"os_id,omitempty"`
	StartupScriptID string   `json:"script_id,omitempty"`
	SnapshotID      string   `json:"snapshot_id,omitempty"`
	EnableIPv6      bool     `json:"enable_ipv6,omitempty"`
	Label           string   `json:"label,omitempty"`
	SSHKeyIDs       []string `json:"sshkey_id,omitempty"`
	AppID           int      `json:"app_id,omitempty"`
	UserData        string   `json:"user_data,omitempty"`
	ActivationEmail bool     `json:"activation_email,omitempty"`
	Hostname        string   `json:"hostname,omitempty"`
	Tag             string   `json:"tag,omitempty"`
	ReservedIPv4    string   `json:"reserved_ipv4,omitempty"`
}

// BareMetalServerBandwidth represents bandwidth information for a bare metal server
type BareMetalServerBandwidth struct {
	IncomingBytes int `json:"incoming_bytes"`
	OutgoingBytes int `json:"outgoing_bytes"`
}

type bareMetalsBase struct {
	BareMetals []BareMetalServer `json:"bare_metals"`
	Meta       *Meta             `json:"meta"`
}

type bareMetalBase struct {
	BareMetal *BareMetalServer `json:"bare_metal"`
}

type BMBareMetalBase struct {
	BareMetalBandwidth map[string]BareMetalServerBandwidth `json:"bandwidth"`
}

// Create a new bare metal server.
func (b *BareMetalServerServiceHandler) Create(ctx context.Context, bmCreate *BareMetalReq) (*BareMetalServer, error) {
	req, err := b.client.NewRequest(ctx, http.MethodPost, bmPath, bmCreate)
	if err != nil {
		return nil, err
	}

	bm := new(bareMetalBase)
	if err = b.client.DoWithContext(ctx, req, bm); err != nil {
		return nil, err
	}

	return bm.BareMetal, nil
}

// Get gets the server with the given ID
func (b *BareMetalServerServiceHandler) Get(ctx context.Context, serverID string) (*BareMetalServer, error) {
	uri := fmt.Sprintf("%s/%s", bmPath, serverID)
	req, err := b.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	bms := new(bareMetalBase)
	if err = b.client.DoWithContext(ctx, req, bms); err != nil {
		return nil, err
	}

	return bms.BareMetal, nil
}

// Update will update the given bare metal. Empty values are ignored
func (b *BareMetalServerServiceHandler) Update(ctx context.Context, serverID string, bmReq *BareMetalReq) error {
	uri := fmt.Sprintf("%s/%s", bmPath, serverID)
	req, err := b.client.NewRequest(ctx, http.MethodPatch, uri, bmReq)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// Delete a bare metal server.
// All data will be permanently lost, and the IP address will be released. There is no going back from this call.
func (b *BareMetalServerServiceHandler) Delete(ctx context.Context, serverID string) error {
	uri := fmt.Sprintf("%s/%s", bmPath, serverID)
	req, err := b.client.NewRequest(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// List lists all bare metal servers on the current account. This includes both pending and active servers.
func (b *BareMetalServerServiceHandler) List(ctx context.Context, options *ListOptions) ([]BareMetalServer, *Meta, error) {
	req, err := b.client.NewRequest(ctx, http.MethodGet, bmPath, nil)
	if err != nil {
		return nil, nil, err
	}

	newValues, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = newValues.Encode()

	bms := new(bareMetalsBase)
	if err = b.client.DoWithContext(ctx, req, bms); err != nil {
		return nil, nil, err
	}

	return bms.BareMetals, bms.Meta, nil
}

// Bandwidth will get the bandwidth used by a bare metal server
func (b *BareMetalServerServiceHandler) GetBandwidth(ctx context.Context, serverID string) (*Bandwidth, error) {
	uri := fmt.Sprintf("%s/%s/bandwidth", bmPath, serverID)
	req, err := b.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	bms := new(Bandwidth)
	if err = b.client.DoWithContext(ctx, req, &bms); err != nil {
		return nil, err
	}

	return bms, nil
}

// ListIPv4s will List the IPv4 information of a bare metal server.
// IP information is only available for bare metal servers in the "active" state.
func (b *BareMetalServerServiceHandler) ListIPv4s(ctx context.Context, serverID string, options *ListOptions) ([]IPv4, *Meta, error) {
	uri := fmt.Sprintf("%s/%s/ipv4", bmPath, serverID)
	req, err := b.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	newValues, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = newValues.Encode()

	ipv4 := new(ipBase)
	if err = b.client.DoWithContext(ctx, req, ipv4); err != nil {
		return nil, nil, err
	}

	return ipv4.IPv4S, ipv4.Meta, nil
}

// ListIPv6s lists the IPv6 information of a bare metal server.
// IP information is only available for bare metal servers in the "active" state.
// If the bare metal server does not have IPv6 enabled, then an empty array is returned.
func (b *BareMetalServerServiceHandler) ListIPv6s(ctx context.Context, serverID string, options *ListOptions) ([]IPv6, *Meta, error) {
	uri := fmt.Sprintf("%s/%s/ipv6", bmPath, serverID)
	req, err := b.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	newValues, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = newValues.Encode()

	ipv6 := new(ipBase)
	if err = b.client.DoWithContext(ctx, req, ipv6); err != nil {
		return nil, nil, err
	}

	return ipv6.IPv6S, ipv6.Meta, nil
}

// Halt a bare metal server.
// This is a hard power off, meaning that the power to the machine is severed.
// The data on the machine will not be modified, and you will still be billed for the machine.
func (b *BareMetalServerServiceHandler) Halt(ctx context.Context, serverID string) error {
	uri := fmt.Sprintf("%s/%s/halt", bmPath, serverID)
	req, err := b.client.NewRequest(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// Reboot a bare metal server. This is a hard reboot, which means that the server is powered off, then back on.
func (b *BareMetalServerServiceHandler) Reboot(ctx context.Context, serverID string) error {
	uri := fmt.Sprintf("%s/%s/reboot", bmPath, serverID)

	req, err := b.client.NewRequest(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// Reinstall the operating system on a bare metal server.
// All data will be permanently lost, but the IP address will remain the same. There is no going back from this call.
func (b *BareMetalServerServiceHandler) Reinstall(ctx context.Context, serverID string) error {
	uri := fmt.Sprintf("%s/%s/reinstall", bmPath, serverID)
	req, err := b.client.NewRequest(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// Start will start a list of bare metal servers the machine is already running, it will be restarted.
func (b *BareMetalServerServiceHandler) MassStart(ctx context.Context, serverList []string) error {
	uri := fmt.Sprintf("%s/start", instancePath)

	reqBody := RequestBody{"baremetal_ids": serverList}
	req, err := b.client.NewRequest(ctx, http.MethodPost, uri, reqBody)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// Halt will pause a list of bare metals.
func (b *BareMetalServerServiceHandler) MassHalt(ctx context.Context, serverList []string) error {
	uri := fmt.Sprintf("%s/halt", instancePath)

	reqBody := RequestBody{"baremetal_ids": serverList}
	req, err := b.client.NewRequest(ctx, http.MethodPost, uri, reqBody)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// MassReboot reboots a list of instances.
func (b *BareMetalServerServiceHandler) MassReboot(ctx context.Context, serverList []string) error {
	uri := fmt.Sprintf("%s/reboot", instancePath)

	reqBody := RequestBody{"baremetal_ids": serverList}
	req, err := b.client.NewRequest(ctx, http.MethodPost, uri, reqBody)
	if err != nil {
		return err
	}

	if err = b.client.DoWithContext(ctx, req, nil); err != nil {
		return err
	}

	return nil
}
