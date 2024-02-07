package kubernetes

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// ClustersSummaryPrinter ...
type ClustersSummaryPrinter struct {
	Clusters []govultr.Cluster `json:"vke_clusters"`
	Meta     *govultr.Meta     `json:"meta"`
}

// JSON ...
func (c *ClustersSummaryPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ClustersSummaryPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ClustersSummaryPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"LABEL",
		"STATUS",
		"REGION",
		"VERSION",
		"NODEPOOL#",
		"NODE#",
	}}
}

// Data ...
func (c *ClustersSummaryPrinter) Data() [][]string {
	if len(c.Clusters) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range c.Clusters {
		nodePoolCount := len(c.Clusters[i].NodePools)
		var nodeCount int = 0

		for j := range c.Clusters[i].NodePools {
			nodeCount += len(c.Clusters[i].NodePools[j].Nodes)
		}

		data = append(data, []string{
			c.Clusters[i].ID,
			c.Clusters[i].Label,
			c.Clusters[i].Status,
			c.Clusters[i].Region,
			c.Clusters[i].Version,
			strconv.Itoa(nodePoolCount),
			strconv.Itoa(nodeCount),
		})
	}

	return data
}

// Paging ...
func (c *ClustersSummaryPrinter) Paging() [][]string {
	return printer.NewPaging(c.Meta.Total, &c.Meta.Links.Next, &c.Meta.Links.Prev).Compose()
}

// ======================================

// ClustersPrinter ...
type ClustersPrinter struct {
	Clusters []govultr.Cluster `json:"vke_clusters"`
	Meta     *govultr.Meta     `json:"meta"`
}

// JSON ...
func (c *ClustersPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ClustersPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ClustersPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (c *ClustersPrinter) Data() [][]string {
	if len(c.Clusters) == 0 {
		return [][]string{0: {"No active kubernetes clusters"}}
	}

	var data [][]string
	for i := range c.Clusters {
		data = append(data,
			[]string{"---------------------------"},
			[]string{"ID", c.Clusters[i].ID},
			[]string{"LABEL", c.Clusters[i].Label},
			[]string{"DATE CREATED", c.Clusters[i].DateCreated},
			[]string{"CLUSTER SUBNET", c.Clusters[i].ClusterSubnet},
			[]string{"SERVICE SUBNET", c.Clusters[i].ServiceSubnet},
			[]string{"IP", c.Clusters[i].IP},
			[]string{"ENDPOINT", c.Clusters[i].Endpoint},
			[]string{"HIGH AVAIL", strconv.FormatBool(c.Clusters[i].HAControlPlanes)},
			[]string{"VERSION", c.Clusters[i].Version},
			[]string{"REGION", c.Clusters[i].Region},
			[]string{"STATUS", c.Clusters[i].Status},
			[]string{" "},
			[]string{"NODE POOLS"},
		)

		for j := range c.Clusters[i].NodePools {
			data = append(data,
				[]string{"ID", c.Clusters[i].NodePools[j].ID},
				[]string{"DATE CREATED", c.Clusters[i].NodePools[j].DateCreated},
				[]string{"DATE UPDATED", c.Clusters[i].NodePools[j].DateUpdated},
				[]string{"LABEL", c.Clusters[i].NodePools[j].Label},
				[]string{"TAG", c.Clusters[i].NodePools[j].Tag},
				[]string{"PLAN", c.Clusters[i].NodePools[j].Plan},
				[]string{"STATUS", c.Clusters[i].NodePools[j].Status},
				[]string{"NODE QUANTITY", strconv.Itoa(c.Clusters[i].NodePools[j].NodeQuantity)},
				[]string{"AUTO SCALER", strconv.FormatBool(c.Clusters[i].NodePools[j].AutoScaler)},
				[]string{"MIN NODES", strconv.Itoa(c.Clusters[i].NodePools[j].MinNodes)},
				[]string{"MAX NODES", strconv.Itoa(c.Clusters[i].NodePools[j].MaxNodes)},
				[]string{" "},
				[]string{"NODES"},
			)

			if len(c.Clusters[i].NodePools[j].Nodes) != 0 {
				// Shouldn't ever be zero
				data = append(data, []string{"ID", "DATE CREATED", "LABEL", "STATUS"})
			}

			for k := range c.Clusters[i].NodePools[j].Nodes {
				data = append(data,
					[]string{
						c.Clusters[i].NodePools[j].Nodes[k].ID,
						c.Clusters[i].NodePools[j].Nodes[k].DateCreated,
						c.Clusters[i].NodePools[j].Nodes[k].Label,
						c.Clusters[i].NodePools[j].Nodes[k].Status,
					},
				)
			}

			data = append(data, []string{" "})
		}
	}

	return data
}

// Paging ...
func (c *ClustersPrinter) Paging() [][]string {
	return printer.NewPaging(c.Meta.Total, &c.Meta.Links.Next, &c.Meta.Links.Prev).Compose()
}

// ======================================

// ClusterPrinter ...
type ClusterPrinter struct {
	Cluster *govultr.Cluster `json:"vke_cluster"`
}

// JSON ...
func (c *ClusterPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ClusterPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ClusterPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (c *ClusterPrinter) Data() [][]string {

	var data [][]string
	data = append(data,
		[]string{"ID", c.Cluster.ID},
		[]string{"LABEL", c.Cluster.Label},
		[]string{"DATE CREATED", c.Cluster.DateCreated},
		[]string{"CLUSTER SUBNET", c.Cluster.ClusterSubnet},
		[]string{"SERVICE SUBNET", c.Cluster.ServiceSubnet},
		[]string{"IP", c.Cluster.IP},
		[]string{"ENDPOINT", c.Cluster.Endpoint},
		[]string{"HIGH AVAIL", strconv.FormatBool(c.Cluster.HAControlPlanes)},
		[]string{"VERSION", c.Cluster.Version},
		[]string{"REGION", c.Cluster.Region},
		[]string{"STATUS", c.Cluster.Status},
		[]string{" "},
		[]string{"NODE POOLS"},
	)

	for i := range c.Cluster.NodePools {
		data = append(data,
			[]string{"ID", c.Cluster.NodePools[i].ID},
			[]string{"DATE CREATED", c.Cluster.NodePools[i].DateCreated},
			[]string{"DATE UPDATED", c.Cluster.NodePools[i].DateUpdated},
			[]string{"LABEL", c.Cluster.NodePools[i].Label},
			[]string{"TAG", c.Cluster.NodePools[i].Tag},
			[]string{"PLAN", c.Cluster.NodePools[i].Plan},
			[]string{"STATUS", c.Cluster.NodePools[i].Status},
			[]string{"NODE QUANTITY", strconv.Itoa(c.Cluster.NodePools[i].NodeQuantity)},
			[]string{"AUTO SCALER", strconv.FormatBool(c.Cluster.NodePools[i].AutoScaler)},
			[]string{"MIN NODES", strconv.Itoa(c.Cluster.NodePools[i].MinNodes)},
			[]string{"MAX NODES", strconv.Itoa(c.Cluster.NodePools[i].MaxNodes)},
			[]string{" "},
			[]string{"NODES"},
		)

		if len(c.Cluster.NodePools[i].Nodes) != 0 {
			// Shouldn't ever be zero
			data = append(data, []string{"ID", "DATE CREATED", "LABEL", "STATUS"})
		}

		for j := range c.Cluster.NodePools[i].Nodes {
			data = append(data,
				[]string{
					c.Cluster.NodePools[i].Nodes[j].ID,
					c.Cluster.NodePools[i].Nodes[j].DateCreated,
					c.Cluster.NodePools[i].Nodes[j].Label,
					c.Cluster.NodePools[i].Nodes[j].Status,
				},
			)
		}

		data = append(data, []string{" "})
	}

	return data
}

// Paging ...
func (c *ClusterPrinter) Paging() [][]string {
	return nil
}

// ======================================

// NodePoolsPrinter ...
type NodePoolsPrinter struct {
	NodePools []govultr.NodePool `json:"node_pools"`
	Meta      *govultr.Meta      `json:"meta"`
}

// JSON ...
func (n *NodePoolsPrinter) JSON() []byte {
	return printer.MarshalObject(n, "json")
}

// YAML ...
func (n *NodePoolsPrinter) YAML() []byte {
	return printer.MarshalObject(n, "yaml")
}

// Columns ...
func (n *NodePoolsPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (n *NodePoolsPrinter) Data() [][]string {
	if len(n.NodePools) == 0 {
		// this shouldn't be possible since at least one nodepool is required
		return [][]string{0: {"No active nodepools on cluster"}}
	}

	var data [][]string
	for i := range n.NodePools {
		data = append(data,
			[]string{"---------------------------"},
			[]string{"ID", n.NodePools[i].ID},
			[]string{"DATE CREATED", n.NodePools[i].DateCreated},
			[]string{"DATE UPDATED", n.NodePools[i].DateUpdated},
			[]string{"LABEL", n.NodePools[i].Label},
			[]string{"TAG", n.NodePools[i].Tag},
			[]string{"PLAN", n.NodePools[i].Plan},
			[]string{"STATUS", n.NodePools[i].Status},
			[]string{"NODE QUANTITY", strconv.Itoa(n.NodePools[i].NodeQuantity)},
			[]string{"AUTO SCALER", strconv.FormatBool(n.NodePools[i].AutoScaler)},
			[]string{"MIN NODES", strconv.Itoa(n.NodePools[i].MinNodes)},
			[]string{"MAX NODES", strconv.Itoa(n.NodePools[i].MaxNodes)},
			[]string{" "},
			[]string{"NODES"},
		)

		if len(n.NodePools[i].Nodes) != 0 {
			// Shouldn't ever be zero
			data = append(data, []string{"ID", "DATE CREATED", "LABEL", "STATUS"})
		}

		for j := range n.NodePools[i].Nodes {
			data = append(data,
				[]string{
					n.NodePools[i].Nodes[j].ID,
					n.NodePools[i].Nodes[j].DateCreated,
					n.NodePools[i].Nodes[j].Label,
					n.NodePools[i].Nodes[j].Status,
				},
			)
		}
	}

	return data
}

// Paging ...
func (n *NodePoolsPrinter) Paging() [][]string {
	return printer.NewPaging(n.Meta.Total, &n.Meta.Links.Next, &n.Meta.Links.Prev).Compose()
}

// ======================================

// NodePoolPrinter ...
type NodePoolPrinter struct {
	NodePool *govultr.NodePool `json:"node_pool"`
}

// JSON ...
func (n *NodePoolPrinter) JSON() []byte {
	return printer.MarshalObject(n, "json")
}

// YAML ...
func (n *NodePoolPrinter) YAML() []byte {
	return printer.MarshalObject(n, "yaml")
}

// Columns ...
func (n *NodePoolPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (n *NodePoolPrinter) Data() [][]string {
	var data [][]string
	data = append(data,
		[]string{"ID", n.NodePool.ID},
		[]string{"DATE CREATED", n.NodePool.DateCreated},
		[]string{"DATE UPDATED", n.NodePool.DateUpdated},
		[]string{"LABEL", n.NodePool.Label},
		[]string{"TAG", n.NodePool.Tag},
		[]string{"PLAN", n.NodePool.Plan},
		[]string{"STATUS", n.NodePool.Status},
		[]string{"NODE QUANTITY", strconv.Itoa(n.NodePool.NodeQuantity)},
		[]string{"AUTO SCALER", strconv.FormatBool(n.NodePool.AutoScaler)},
		[]string{"MIN NODES", strconv.Itoa(n.NodePool.MinNodes)},
		[]string{"MAX NODES", strconv.Itoa(n.NodePool.MaxNodes)},
		[]string{" "},
		[]string{"NODES"},
	)

	if len(n.NodePool.Nodes) != 0 {
		// Shouldn't ever be zero
		data = append(data, []string{"ID", "DATE CREATED", "LABEL", "STATUS"})
	}

	for i := range n.NodePool.Nodes {
		data = append(data,
			[]string{
				n.NodePool.Nodes[i].ID,
				n.NodePool.Nodes[i].DateCreated,
				n.NodePool.Nodes[i].Label,
				n.NodePool.Nodes[i].Status,
			},
		)
	}

	return data
}

// Paging ...
func (n *NodePoolPrinter) Paging() [][]string {
	return nil
}

// ======================================

// VersionsPrinter ...
type VersionsPrinter struct {
	Versions []string `json:"versions"`
}

// JSON ...
func (v *VersionsPrinter) JSON() []byte {
	return printer.MarshalObject(v, "json")
}

// YAML ...
func (v *VersionsPrinter) YAML() []byte {
	return printer.MarshalObject(v, "yaml")
}

// Columns ...
func (v *VersionsPrinter) Columns() [][]string {
	return [][]string{0: {"VERSIONS"}}
}

// Data ...
func (v *VersionsPrinter) Data() [][]string {
	if len(v.Versions) == 0 {
		return [][]string{0: {"---"}}
	}

	var data [][]string

	for i := range v.Versions {
		data = append(data, []string{v.Versions[i]})
	}

	return data
}

// Paging ...
func (v *VersionsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// UpgradesPrinter ...
type UpgradesPrinter struct {
	Upgrades []string `json:"available_upgrades"`
}

// JSON ...
func (u *UpgradesPrinter) JSON() []byte {
	return printer.MarshalObject(u, "json")
}

// YAML ...
func (u *UpgradesPrinter) YAML() []byte {
	return printer.MarshalObject(u, "yaml")
}

// Columns ...
func (u *UpgradesPrinter) Columns() [][]string {
	return [][]string{0: {"UPGRADES"}}
}

// Data ...
func (u *UpgradesPrinter) Data() [][]string {
	if len(u.Upgrades) == 0 {
		return [][]string{0: {"---"}}
	}

	var data [][]string

	for i := range u.Upgrades {
		data = append(data, []string{u.Upgrades[i]})
	}

	return data
}

// Paging ...
func (u *UpgradesPrinter) Paging() [][]string {
	return nil
}

// ======================================

// ConfigPrinter ...
type ConfigPrinter struct {
	Config *govultr.KubeConfig
}

// JSON ...
func (c *ConfigPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ConfigPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ConfigPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (c *ConfigPrinter) Data() [][]string {
	return [][]string{0: {c.Config.KubeConfig}}
}

// Paging ...
func (c *ConfigPrinter) Paging() [][]string {
	return nil
}
