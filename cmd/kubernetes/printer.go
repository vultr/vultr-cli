package kubernetes

import (
	"fmt"
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
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
		var nodeCount = 0

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
	return printer.NewPagingFromMeta(c.Meta).Compose()
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
			[]string{"FIREWALL GROUP ID", c.Clusters[i].FirewallGroupID},
			[]string{"VERSION", c.Clusters[i].Version},
			[]string{"REGION", c.Clusters[i].Region},
			[]string{"STATUS", c.Clusters[i].Status},
			[]string{" "},
			[]string{"NODE POOLS"},
		)

		for j := range c.Clusters[i].NodePools {
			data = append(data,
				[]string{"ID", c.Clusters[i].NodePools[j].ID},
				[]string{"LABEL", c.Clusters[i].NodePools[j].Label},
				[]string{"TAG", c.Clusters[i].NodePools[j].Tag},
				[]string{"PLAN", c.Clusters[i].NodePools[j].Plan},
				[]string{"STATUS", c.Clusters[i].NodePools[j].Status},
				[]string{"NODE QUANTITY", strconv.Itoa(c.Clusters[i].NodePools[j].NodeQuantity)},
				[]string{"AUTO SCALER", strconv.FormatBool(c.Clusters[i].NodePools[j].AutoScaler)},
				[]string{"MIN NODES", strconv.Itoa(c.Clusters[i].NodePools[j].MinNodes)},
				[]string{"MAX NODES", strconv.Itoa(c.Clusters[i].NodePools[j].MaxNodes)},
				[]string{"VPC ONLY", strconv.FormatBool(c.Clusters[i].NodePools[j].VPCOnly)},
				[]string{"DATE CREATED", c.Clusters[i].NodePools[j].DateCreated},
				[]string{"DATE UPDATED", c.Clusters[i].NodePools[j].DateUpdated},
			)

			data = append(data,
				[]string{" "},
				[]string{"NODES"},
			)

			if len(c.Clusters[i].NodePools[j].Nodes) != 0 {
				// Shouldn't ever be zero
				data = append(data, []string{"ID", "LABEL", "STATUS", "PRIVATE IP", "DATE CREATED"})
			}

			for k := range c.Clusters[i].NodePools[j].Nodes {
				data = append(data,
					[]string{
						c.Clusters[i].NodePools[j].Nodes[k].ID,
						c.Clusters[i].NodePools[j].Nodes[k].Label,
						c.Clusters[i].NodePools[j].Nodes[k].Status,
						c.Clusters[i].NodePools[j].Nodes[k].PrivateIP,
						c.Clusters[i].NodePools[j].Nodes[k].DateCreated,
					},
				)
			}

			if len(c.Clusters[i].NodePools[j].Labels) != 0 {
				data = append(data, []string{" "}, []string{"NODE LABELS"})
				for k := range c.Clusters[i].NodePools[j].Labels {
					label := fmt.Sprintf("%s=%s", k, c.Clusters[i].NodePools[j].Labels[k])
					data = append(data, []string{label})
				}
			}

			data = append(data, []string{" "})
		}
	}

	return data
}

// Paging ...
func (c *ClustersPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(c.Meta).Compose()
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
		[]string{"STATUS", c.Cluster.Status},
		[]string{"VERSION", c.Cluster.Version},
		[]string{"REGION", c.Cluster.Region},
		[]string{"CLUSTER SUBNET", c.Cluster.ClusterSubnet},
		[]string{"SERVICE SUBNET", c.Cluster.ServiceSubnet},
		[]string{"IP", c.Cluster.IP},
		[]string{"ENDPOINT", c.Cluster.Endpoint},
		[]string{"HIGH AVAIL", strconv.FormatBool(c.Cluster.HAControlPlanes)},
		[]string{"FIREWALL GROUP ID", c.Cluster.FirewallGroupID},
		[]string{"DATE CREATED", c.Cluster.DateCreated},
		[]string{" "},
		[]string{"NODE POOLS"},
	)

	for i := range c.Cluster.NodePools {
		data = append(data,
			[]string{"ID", c.Cluster.NodePools[i].ID},
			[]string{"LABEL", c.Cluster.NodePools[i].Label},
			[]string{"TAG", c.Cluster.NodePools[i].Tag},
			[]string{"PLAN", c.Cluster.NodePools[i].Plan},
			[]string{"STATUS", c.Cluster.NodePools[i].Status},
			[]string{"NODE QUANTITY", strconv.Itoa(c.Cluster.NodePools[i].NodeQuantity)},
			[]string{"AUTO SCALER", strconv.FormatBool(c.Cluster.NodePools[i].AutoScaler)},
			[]string{"MIN NODES", strconv.Itoa(c.Cluster.NodePools[i].MinNodes)},
			[]string{"MAX NODES", strconv.Itoa(c.Cluster.NodePools[i].MaxNodes)},
			[]string{"VPC ONLY", strconv.FormatBool(c.Cluster.NodePools[i].VPCOnly)},
			[]string{"DATE CREATED", c.Cluster.NodePools[i].DateCreated},
			[]string{"DATE UPDATED", c.Cluster.NodePools[i].DateUpdated},
			[]string{" "},
			[]string{"NODES"},
		)

		if len(c.Cluster.NodePools[i].Nodes) != 0 {
			// Shouldn't ever be zero
			data = append(data, []string{"ID", "LABEL", "STATUS", "PRIVATE IP", "DATE CREATED"})
		}

		for j := range c.Cluster.NodePools[i].Nodes {
			data = append(data,
				[]string{
					c.Cluster.NodePools[i].Nodes[j].ID,
					c.Cluster.NodePools[i].Nodes[j].Label,
					c.Cluster.NodePools[i].Nodes[j].Status,
					c.Cluster.NodePools[i].Nodes[j].PrivateIP,
					c.Cluster.NodePools[i].Nodes[j].DateCreated,
				},
			)
		}

		if len(c.Cluster.NodePools[i].Labels) != 0 {
			data = append(data, []string{" "}, []string{"NODE LABELS"})
			for k := range c.Cluster.NodePools[i].Labels {
				label := fmt.Sprintf("%s=%s", k, c.Cluster.NodePools[i].Labels[k])
				data = append(data, []string{label})
			}
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
			[]string{"LABEL", n.NodePools[i].Label},
			[]string{"TAG", n.NodePools[i].Tag},
			[]string{"PLAN", n.NodePools[i].Plan},
			[]string{"STATUS", n.NodePools[i].Status},
			[]string{"NODE QUANTITY", strconv.Itoa(n.NodePools[i].NodeQuantity)},
			[]string{"AUTO SCALER", strconv.FormatBool(n.NodePools[i].AutoScaler)},
			[]string{"MIN NODES", strconv.Itoa(n.NodePools[i].MinNodes)},
			[]string{"MAX NODES", strconv.Itoa(n.NodePools[i].MaxNodes)},
			[]string{"VPC ONLY", strconv.FormatBool(n.NodePools[i].VPCOnly)},
			[]string{"DATE CREATED", n.NodePools[i].DateCreated},
			[]string{"DATE UPDATED", n.NodePools[i].DateUpdated},
			[]string{" "},
			[]string{"NODES"},
		)

		if len(n.NodePools[i].Nodes) != 0 {
			// Shouldn't ever be zero
			data = append(data, []string{"ID", "LABEL", "STATUS", "PRIVATE IP", "DATE CREATED"})
		}

		for j := range n.NodePools[i].Nodes {
			data = append(data,
				[]string{
					n.NodePools[i].Nodes[j].ID,
					n.NodePools[i].Nodes[j].Label,
					n.NodePools[i].Nodes[j].Status,
					n.NodePools[i].Nodes[j].PrivateIP,
					n.NodePools[i].Nodes[j].DateCreated,
				},
			)
		}

		if len(n.NodePools[i].Labels) != 0 {
			data = append(data, []string{" "}, []string{"NODE LABELS"})
			for k := range n.NodePools[i].Labels {
				label := fmt.Sprintf("%s=%s", k, n.NodePools[i].Labels[k])
				data = append(data, []string{label})
			}
		}
	}

	return data
}

// Paging ...
func (n *NodePoolsPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(n.Meta).Compose()
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
		[]string{"LABEL", n.NodePool.Label},
		[]string{"TAG", n.NodePool.Tag},
		[]string{"PLAN", n.NodePool.Plan},
		[]string{"STATUS", n.NodePool.Status},
		[]string{"NODE QUANTITY", strconv.Itoa(n.NodePool.NodeQuantity)},
		[]string{"AUTO SCALER", strconv.FormatBool(n.NodePool.AutoScaler)},
		[]string{"MIN NODES", strconv.Itoa(n.NodePool.MinNodes)},
		[]string{"MAX NODES", strconv.Itoa(n.NodePool.MaxNodes)},
		[]string{"VPC ONLY", strconv.FormatBool(n.NodePool.VPCOnly)},
		[]string{"DATE CREATED", n.NodePool.DateCreated},
		[]string{"DATE UPDATED", n.NodePool.DateUpdated},
		[]string{" "},
		[]string{"NODES"},
	)

	if len(n.NodePool.Nodes) != 0 {
		// Shouldn't ever be zero
		data = append(data, []string{"ID", "LABEL", "STATUS", "PRIVATE IP", "DATE CREATED"})
	}

	for i := range n.NodePool.Nodes {
		data = append(data,
			[]string{
				n.NodePool.Nodes[i].ID,
				n.NodePool.Nodes[i].Label,
				n.NodePool.Nodes[i].Status,
				n.NodePool.Nodes[i].PrivateIP,
				n.NodePool.Nodes[i].DateCreated,
			},
		)
	}

	if len(n.NodePool.Labels) != 0 {
		data = append(data, []string{" "}, []string{"NODE LABELS"})
		for k := range n.NodePool.Labels {
			label := fmt.Sprintf("%s=%s", k, n.NodePool.Labels[k])
			data = append(data, []string{label})
		}
	}

	return data
}

// Paging ...
func (n *NodePoolPrinter) Paging() [][]string {
	return nil
}

// ======================================

// NodePoolsSummaryPrinter ...
type NodePoolsSummaryPrinter struct {
	NodePools []govultr.NodePool `json:"node_pools"`
	Meta      *govultr.Meta      `json:"meta"`
}

// JSON ...
func (n *NodePoolsSummaryPrinter) JSON() []byte {
	return printer.MarshalObject(n, "json")
}

// YAML ...
func (n *NodePoolsSummaryPrinter) YAML() []byte {
	return printer.MarshalObject(n, "yaml")
}

// Columns ...
func (n *NodePoolsSummaryPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"PLAN",
		"STATUS",
		"NODE QUANTITY",
		"AUTO SCALER",
		"MIN NODES",
		"MAX NODES",
	}}
}

// Data ...
func (n *NodePoolsSummaryPrinter) Data() [][]string {
	if len(n.NodePools) == 0 {
		// there should always be a node pool
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range n.NodePools {
		data = append(data, []string{
			n.NodePools[i].ID,
			n.NodePools[i].Plan,
			n.NodePools[i].Status,
			strconv.Itoa(n.NodePools[i].NodeQuantity),
			strconv.FormatBool(n.NodePools[i].AutoScaler),
			strconv.Itoa(n.NodePools[i].MinNodes),
			strconv.Itoa(n.NodePools[i].MaxNodes),
		})
	}

	return data
}

// Paging ...
func (n *NodePoolsSummaryPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(n.Meta).Compose()
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

// ======================================

// AppsPrinter ...
type AppsPrinter struct {
	Apps []govultr.ClusterApplication
}

// JSON ...
func (a *AppsPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML ...
func (a *AppsPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
}

// Columns ...
func (a *AppsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"NAME",
		"DESCRIPTION",
		"VERSION",
		"DATE CREATED",
		"DATE UPDATED",
	}}
}

// Data ...
func (a *AppsPrinter) Data() [][]string {
	if len(a.Apps) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range a.Apps {
		data = append(data, []string{
			a.Apps[i].ID,
			utils.TruncateString(a.Apps[i].Name, utils.DefaultRuneLengthShort),
			utils.TruncateString(a.Apps[i].Description, utils.DefaultRuneLengthLong),
			a.Apps[i].Version,
			utils.FormatDateTimeString(a.Apps[i].DateCreated),
			utils.FormatDateTimeString(a.Apps[i].DateUpdated),
		})
	}

	return data
}

// Paging ...
func (a *AppsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// AppPrinter ...
type AppPrinter struct {
	App *govultr.ClusterApplication
}

// JSON ...
func (a *AppPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML ...
func (a *AppPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
}

// Columns ...
func (a *AppPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (a *AppPrinter) Data() [][]string {
	if a.App == nil {
		return [][]string{0: {"No such app ID found"}}
	}

	return [][]string{
		{"ID", a.App.ID},
		{"NAME", a.App.Name},
		{"DESCRIPTION", a.App.Description},
		{"VERSION", a.App.Version},
		{"DATE CREATED", a.App.DateCreated},
		{"DATE UPDATED", a.App.DateUpdated},
	}
}

// Paging ...
func (a *AppPrinter) Paging() [][]string {
	return nil
}

// ======================================

// AppDeploymentsPrinter ...
type AppDeploymentsPrinter struct {
	Deployments []govultr.ClusterApplicationDeployment
}

// JSON ...
func (a *AppDeploymentsPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML ...
func (a *AppDeploymentsPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
}

// Columns ...
func (a *AppDeploymentsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"APP ID",
		"APP NAME",
		"VERSION",
		"STATUS",
		"DATE CREATED",
		"DATE UPDATED",
	}}
}

// Data ...
func (a *AppDeploymentsPrinter) Data() [][]string {
	if len(a.Deployments) <= 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range a.Deployments {
		data = append(data, []string{
			a.Deployments[i].ID,
			a.Deployments[i].AppID,
			a.Deployments[i].AppName,
			a.Deployments[i].Version,
			a.Deployments[i].Status,
			a.Deployments[i].DateCreated,
			a.Deployments[i].DateUpdated,
		})
	}

	return data
}

// Paging ...
func (a *AppDeploymentsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// AppDeploymentPrinter ...
type AppDeploymentPrinter struct {
	Deployment *govultr.ClusterApplicationDeployment
}

// JSON ...
func (a *AppDeploymentPrinter) JSON() []byte {
	return printer.MarshalObject(a, "json")
}

// YAML ...
func (a *AppDeploymentPrinter) YAML() []byte {
	return printer.MarshalObject(a, "yaml")
}

// Columns ...
func (a *AppDeploymentPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (a *AppDeploymentPrinter) Data() [][]string {
	if a.Deployment == nil {
		return [][]string{0: {"The deployment does not exist"}}
	}

	return [][]string{
		{"ID", a.Deployment.ID},
		{"APP ID", a.Deployment.AppID},
		{"NAME", a.Deployment.AppName},
		{"VERSION", a.Deployment.Version},
		{"STATUS", a.Deployment.Status},
		{"DATE CREATED", a.Deployment.DateCreated},
		{"DATE UPDATED", a.Deployment.DateUpdated},
	}
}

// Paging ...
func (a *AppDeploymentPrinter) Paging() [][]string {
	return nil
}

// ======================================

// LabelsPrinter ...
type LabelsPrinter struct {
	Labels []govultr.NodePoolLabel
}

// JSON ...
func (l *LabelsPrinter) JSON() []byte {
	return printer.MarshalObject(l, "json")
}

// YAML ...
func (l *LabelsPrinter) YAML() []byte {
	return printer.MarshalObject(l, "yaml")
}

// Columns ...
func (l *LabelsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"KEY",
		"VALUE",
	}}
}

// Data ...
func (l *LabelsPrinter) Data() [][]string {
	if len(l.Labels) == 0 {
		return [][]string{0: {"---", "---", "---"}}
	}

	var data [][]string
	for i := range l.Labels {
		data = append(data, []string{
			l.Labels[i].ID,
			l.Labels[i].Key,
			l.Labels[i].Value,
		})
	}

	return data
}

// Paging ...
func (l *LabelsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// LabelPrinter ...
type LabelPrinter struct {
	Label *govultr.NodePoolLabel
}

// JSON ...
func (l *LabelPrinter) JSON() []byte {
	return printer.MarshalObject(l, "json")
}

// YAML ...
func (l *LabelPrinter) YAML() []byte {
	return printer.MarshalObject(l, "yaml")
}

// Columns ...
func (l *LabelPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (l *LabelPrinter) Data() [][]string {
	if l.Label == nil {
		return [][]string{0: {"No such label found"}}
	}

	return [][]string{
		{"ID", l.Label.ID},
		{"KEY", l.Label.Key},
		{"VALUE", l.Label.Value},
	}
}

// Paging ...
func (l *LabelPrinter) Paging() [][]string {
	return nil
}

// ======================================

// TaintsPrinter ...
type TaintsPrinter struct {
	Taints []govultr.NodePoolTaint
}

// JSON ...
func (t *TaintsPrinter) JSON() []byte {
	return printer.MarshalObject(t, "json")
}

// YAML ...
func (t *TaintsPrinter) YAML() []byte {
	return printer.MarshalObject(t, "yaml")
}

// Columns ...
func (t *TaintsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"KEY",
		"VALUE",
		"EFFECT",
	}}
}

// Data ...
func (t *TaintsPrinter) Data() [][]string {
	if len(t.Taints) == 0 {
		return [][]string{0: {"---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range t.Taints {
		data = append(data, []string{
			t.Taints[i].ID,
			t.Taints[i].Key,
			t.Taints[i].Value,
			t.Taints[i].Effect,
		})
	}

	return data
}

// Paging ...
func (t *TaintsPrinter) Paging() [][]string {
	return nil
}

// ======================================

// TaintPrinter ...
type TaintPrinter struct {
	Taint *govultr.NodePoolTaint
}

// JSON ...
func (t *TaintPrinter) JSON() []byte {
	return printer.MarshalObject(t, "json")
}

// YAML ...
func (t *TaintPrinter) YAML() []byte {
	return printer.MarshalObject(t, "yaml")
}

// Columns ...
func (t *TaintPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (t *TaintPrinter) Data() [][]string {
	if t.Taint == nil {
		return [][]string{0: {"No such taint found"}}
	}

	return [][]string{
		{"ID", t.Taint.ID},
		{"KEY", t.Taint.Key},
		{"VALUE", t.Taint.Value},
		{"EFFECT", t.Taint.Effect},
	}
}

// Paging ...
func (t *TaintPrinter) Paging() [][]string {
	return nil
}
