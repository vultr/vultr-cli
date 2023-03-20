package govultr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

const vkePath = "/v2/kubernetes/clusters"

// KubernetesService is the interface to interact with kubernetes endpoint on the Vultr API
// Link : https://www.vultr.com/api/#tag/kubernetes
type KubernetesService interface {
	CreateCluster(ctx context.Context, createReq *ClusterReq) (*Cluster, *http.Response, error)
	GetCluster(ctx context.Context, id string) (*Cluster, *http.Response, error)
	ListClusters(ctx context.Context, options *ListOptions) ([]Cluster, *Meta, *http.Response, error)
	UpdateCluster(ctx context.Context, vkeID string, updateReq *ClusterReqUpdate) error
	DeleteCluster(ctx context.Context, id string) error
	DeleteClusterWithResources(ctx context.Context, id string) error

	CreateNodePool(ctx context.Context, vkeID string, nodePoolReq *NodePoolReq) (*NodePool, *http.Response, error)
	ListNodePools(ctx context.Context, vkeID string, options *ListOptions) ([]NodePool, *Meta, *http.Response, error)
	GetNodePool(ctx context.Context, vkeID, nodePoolID string) (*NodePool, *http.Response, error)
	UpdateNodePool(ctx context.Context, vkeID, nodePoolID string, updateReq *NodePoolReqUpdate) (*NodePool, *http.Response, error)
	DeleteNodePool(ctx context.Context, vkeID, nodePoolID string) error

	DeleteNodePoolInstance(ctx context.Context, vkeID, nodePoolID, nodeID string) error
	RecycleNodePoolInstance(ctx context.Context, vkeID, nodePoolID, nodeID string) error

	GetKubeConfig(ctx context.Context, vkeID string) (*KubeConfig, *http.Response, error)
	GetVersions(ctx context.Context) (*Versions, *http.Response, error)

	GetUpgrades(ctx context.Context, vkeID string) ([]string, *http.Response, error)
	Upgrade(ctx context.Context, vkeID string, body *ClusterUpgradeReq) error
}

// KubernetesHandler handles interaction with the kubernetes methods for the Vultr API
type KubernetesHandler struct {
	client *Client
}

// Cluster represents a full VKE cluster
type Cluster struct {
	ID            string     `json:"id"`
	Label         string     `json:"label"`
	DateCreated   string     `json:"date_created"`
	ClusterSubnet string     `json:"cluster_subnet"`
	ServiceSubnet string     `json:"service_subnet"`
	IP            string     `json:"ip"`
	Endpoint      string     `json:"endpoint"`
	Version       string     `json:"version"`
	Region        string     `json:"region"`
	Status        string     `json:"status"`
	NodePools     []NodePool `json:"node_pools"`
}

// NodePool represents a pool of nodes that are grouped by their label and plan type
type NodePool struct {
	ID           string `json:"id"`
	DateCreated  string `json:"date_created"`
	DateUpdated  string `json:"date_updated"`
	Label        string `json:"label"`
	Plan         string `json:"plan"`
	Status       string `json:"status"`
	NodeQuantity int    `json:"node_quantity"`
	MinNodes     int    `json:"min_nodes"`
	MaxNodes     int    `json:"max_nodes"`
	AutoScaler   bool   `json:"auto_scaler"`
	Tag          string `json:"tag"`
	Nodes        []Node `json:"nodes"`
}

// Node represents a node that will live within a nodepool
type Node struct {
	ID          string `json:"id"`
	DateCreated string `json:"date_created"`
	Label       string `json:"label"`
	Status      string `json:"status"`
}

// KubeConfig will contain the kubeconfig b64 encoded
type KubeConfig struct {
	KubeConfig string `json:"kube_config"`
}

// ClusterReq struct used to create a cluster
type ClusterReq struct {
	Label     string        `json:"label"`
	Region    string        `json:"region"`
	Version   string        `json:"version"`
	NodePools []NodePoolReq `json:"node_pools"`
}

// ClusterReqUpdate struct used to update update a cluster
type ClusterReqUpdate struct {
	Label string `json:"label"`
}

// NodePoolReq struct used to create a node pool
type NodePoolReq struct {
	NodeQuantity int    `json:"node_quantity"`
	Label        string `json:"label"`
	Plan         string `json:"plan"`
	Tag          string `json:"tag"`
	MinNodes     int    `json:"min_nodes,omitempty"`
	MaxNodes     int    `json:"max_nodes,omitempty"`
	AutoScaler   *bool  `json:"auto_scaler"`
}

// NodePoolReqUpdate struct used to update a node pool
type NodePoolReqUpdate struct {
	NodeQuantity int     `json:"node_quantity,omitempty"`
	Tag          *string `json:"tag,omitempty"`
	MinNodes     int     `json:"min_nodes,omitempty"`
	MaxNodes     int     `json:"max_nodes,omitempty"`
	AutoScaler   *bool   `json:"auto_scaler,omitempty"`
}

type vkeClustersBase struct {
	VKEClusters []Cluster `json:"vke_clusters"`
	Meta        *Meta     `json:"meta"`
}

type vkeClusterBase struct {
	VKECluster *Cluster `json:"vke_cluster"`
}

type vkeNodePoolsBase struct {
	NodePools []NodePool `json:"node_pools"`
	Meta      *Meta      `json:"meta"`
}

type vkeNodePoolBase struct {
	NodePool *NodePool `json:"node_pool"`
}

// Versions that are supported for VKE
type Versions struct {
	Versions []string `json:"versions"`
}

// AvailableUpgrades for a given VKE cluster
type availableUpgrades struct {
	AvailableUpgrades []string `json:"available_upgrades"`
}

// ClusterUpgradeReq struct for vke upgradse
type ClusterUpgradeReq struct {
	UpgradeVersion string `json:"upgrade_version,omitempty"`
}

// CreateCluster will create a Kubernetes cluster.
func (k *KubernetesHandler) CreateCluster(ctx context.Context, createReq *ClusterReq) (*Cluster, *http.Response, error) {
	req, err := k.client.NewRequest(ctx, http.MethodPost, vkePath, createReq)
	if err != nil {
		return nil, nil, err
	}

	var k8 = new(vkeClusterBase)
	resp, err := k.client.DoWithContext(ctx, req, &k8)
	if err != nil {
		return nil, resp, err
	}

	return k8.VKECluster, resp, nil
}

// GetCluster will return a Kubernetes cluster.
func (k *KubernetesHandler) GetCluster(ctx context.Context, id string) (*Cluster, *http.Response, error) {
	req, err := k.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", vkePath, id), nil)
	if err != nil {
		return nil, nil, err
	}

	k8 := new(vkeClusterBase)
	resp, err := k.client.DoWithContext(ctx, req, &k8)
	if err != nil {
		return nil, resp, err
	}

	return k8.VKECluster, resp, nil
}

// ListClusters will return all kubernetes clusters.
func (k *KubernetesHandler) ListClusters(ctx context.Context, options *ListOptions) ([]Cluster, *Meta, *http.Response, error) {
	req, err := k.client.NewRequest(ctx, http.MethodGet, vkePath, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	newValues, err := query.Values(options)
	if err != nil {
		return nil, nil, nil, err
	}

	req.URL.RawQuery = newValues.Encode()

	k8s := new(vkeClustersBase)
	resp, err := k.client.DoWithContext(ctx, req, &k8s)
	if err != nil {
		return nil, nil, resp, err
	}

	return k8s.VKEClusters, k8s.Meta, resp, nil
}

// UpdateCluster updates label on VKE cluster
func (k *KubernetesHandler) UpdateCluster(ctx context.Context, vkeID string, updateReq *ClusterReqUpdate) error {
	req, err := k.client.NewRequest(ctx, http.MethodPut, fmt.Sprintf("%s/%s", vkePath, vkeID), updateReq)
	if err != nil {
		return err
	}

	_, err = k.client.DoWithContext(ctx, req, nil)
	return err
}

// DeleteCluster will delete a Kubernetes cluster.
func (k *KubernetesHandler) DeleteCluster(ctx context.Context, id string) error {
	req, err := k.client.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", vkePath, id), nil)
	if err != nil {
		return err
	}
	_, err = k.client.DoWithContext(ctx, req, nil)
	return err
}

// DeleteClusterWithResources will delete a Kubernetes cluster and all related resources.
func (k *KubernetesHandler) DeleteClusterWithResources(ctx context.Context, id string) error {
	req, err := k.client.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s/delete-with-linked-resources", vkePath, id), nil)
	if err != nil {
		return err
	}
	_, err = k.client.DoWithContext(ctx, req, nil)
	return err
}

// CreateNodePool creates a nodepool on a VKE cluster
func (k *KubernetesHandler) CreateNodePool(ctx context.Context, vkeID string, nodePoolReq *NodePoolReq) (*NodePool, *http.Response, error) {
	req, err := k.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/%s/node-pools", vkePath, vkeID), nodePoolReq)
	if err != nil {
		return nil, nil, err
	}

	n := new(vkeNodePoolBase)
	resp, err := k.client.DoWithContext(ctx, req, n)
	if err != nil {
		return nil, resp, err
	}

	return n.NodePool, resp, nil
}

// ListNodePools will return all nodepools for a given VKE cluster
func (k *KubernetesHandler) ListNodePools(ctx context.Context, vkeID string, options *ListOptions) ([]NodePool, *Meta, *http.Response, error) {
	req, err := k.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s/node-pools", vkePath, vkeID), nil)
	if err != nil {
		return nil, nil, nil, err
	}

	newValues, err := query.Values(options)
	if err != nil {
		return nil, nil, nil, err
	}

	req.URL.RawQuery = newValues.Encode()

	n := new(vkeNodePoolsBase)
	resp, err := k.client.DoWithContext(ctx, req, &n)
	if err != nil {
		return nil, nil, resp, err
	}

	return n.NodePools, n.Meta, resp, nil
}

// GetNodePool will return a single nodepool
func (k *KubernetesHandler) GetNodePool(ctx context.Context, vkeID, nodePoolID string) (*NodePool, *http.Response, error) {
	req, err := k.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s/node-pools/%s", vkePath, vkeID, nodePoolID), nil)
	if err != nil {
		return nil, nil, err
	}

	n := new(vkeNodePoolBase)
	resp, err := k.client.DoWithContext(ctx, req, &n)
	if err != nil {
		return nil, resp, err
	}

	return n.NodePool, resp, nil
}

// UpdateNodePool will allow you change the quantity of nodes within a nodepool
func (k *KubernetesHandler) UpdateNodePool(ctx context.Context, vkeID, nodePoolID string, updateReq *NodePoolReqUpdate) (*NodePool, *http.Response, error) {
	req, err := k.client.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("%s/%s/node-pools/%s", vkePath, vkeID, nodePoolID), updateReq)
	if err != nil {
		return nil, nil, err
	}

	np := new(vkeNodePoolBase)
	resp, err := k.client.DoWithContext(ctx, req, np)
	if err != nil {
		return nil, resp, err
	}

	return np.NodePool, resp, nil
}

// DeleteNodePool will remove a nodepool from a VKE cluster
func (k *KubernetesHandler) DeleteNodePool(ctx context.Context, vkeID, nodePoolID string) error {
	req, err := k.client.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s/node-pools/%s", vkePath, vkeID, nodePoolID), nil)
	if err != nil {
		return err
	}

	_, err = k.client.DoWithContext(ctx, req, nil)
	return err
}

// DeleteNodePoolInstance will remove a specified node from a nodepool
func (k *KubernetesHandler) DeleteNodePoolInstance(ctx context.Context, vkeID, nodePoolID, nodeID string) error {
	req, err := k.client.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s/node-pools/%s/nodes/%s", vkePath, vkeID, nodePoolID, nodeID), nil)
	if err != nil {
		return err
	}

	_, err = k.client.DoWithContext(ctx, req, nil)
	return err
}

// RecycleNodePoolInstance will recycle (destroy + redeploy) a given node on a nodepool
func (k *KubernetesHandler) RecycleNodePoolInstance(ctx context.Context, vkeID, nodePoolID, nodeID string) error {
	req, err := k.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/%s/node-pools/%s/nodes/%s/recycle", vkePath, vkeID, nodePoolID, nodeID), nil)
	if err != nil {
		return err
	}

	_, err = k.client.DoWithContext(ctx, req, nil)
	return err
}

// GetKubeConfig returns the kubeconfig for the specified VKE cluster
func (k *KubernetesHandler) GetKubeConfig(ctx context.Context, vkeID string) (*KubeConfig, *http.Response, error) {
	req, err := k.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s/config", vkePath, vkeID), nil)
	if err != nil {
		return nil, nil, err
	}

	kc := new(KubeConfig)
	resp, err := k.client.DoWithContext(ctx, req, &kc)
	if err != nil {
		return nil, resp, err
	}

	return kc, resp, nil
}

// GetVersions returns the supported kubernetes versions
func (k *KubernetesHandler) GetVersions(ctx context.Context) (*Versions, *http.Response, error) {
	uri := "/v2/kubernetes/versions"
	req, err := k.client.NewRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}

	versions := new(Versions)
	resp, err := k.client.DoWithContext(ctx, req, &versions)
	if err != nil {
		return nil, resp, err
	}

	return versions, resp, nil
}

// GetUpgrades returns all version a VKE cluster can upgrade to
func (k *KubernetesHandler) GetUpgrades(ctx context.Context, vkeID string) ([]string, *http.Response, error) {
	req, err := k.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s/available-upgrades", vkePath, vkeID), nil)
	if err != nil {
		return nil, nil, err
	}

	upgrades := new(availableUpgrades)
	resp, err := k.client.DoWithContext(ctx, req, &upgrades)
	if err != nil {
		return nil, resp, err
	}

	return upgrades.AvailableUpgrades, resp, nil
}

// Upgrade beings a VKE cluster upgrade
func (k *KubernetesHandler) Upgrade(ctx context.Context, vkeID string, body *ClusterUpgradeReq) error {

	req, err := k.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/%s/upgrades", vkePath, vkeID), body)
	if err != nil {
		return err
	}

	_, err = k.client.DoWithContext(ctx, req, nil)
	return err
}
