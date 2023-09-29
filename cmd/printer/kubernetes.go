package printer

import (
	"github.com/vultr/govultr/v3"
)

func Clusters(cluster []govultr.Cluster, meta *govultr.Meta) {
	defer flush()

	if len(cluster) == 0 {
		displayString("No active kubernetes clusters")
		return
	}

	for i := range cluster {
		display(columns{"ID", cluster[i].ID})
		display(columns{"LABEL", cluster[i].Label})
		display(columns{"DATE CREATED", cluster[i].DateCreated})
		display(columns{"CLUSTER SUBNET", cluster[i].ClusterSubnet})
		display(columns{"SERVICE SUBNET", cluster[i].ServiceSubnet})
		display(columns{"IP", cluster[i].IP})
		display(columns{"ENDPOINT", cluster[i].Endpoint})
		display(columns{"VERSION", cluster[i].Version})
		display(columns{"REGION", cluster[i].Region})
		display(columns{"STATUS", cluster[i].Status})

		display(columns{" "})
		display(columns{"NODE POOLS"})
		for n := range cluster[i].NodePools {
			display(columns{"ID", cluster[i].NodePools[n].ID})
			display(columns{"DATE CREATED", cluster[i].NodePools[n].DateCreated})
			display(columns{"DATE UPDATED", cluster[i].NodePools[n].DateUpdated})
			display(columns{"LABEL", cluster[i].NodePools[n].Label})
			display(columns{"TAG", cluster[i].NodePools[n].Tag})
			display(columns{"PLAN", cluster[i].NodePools[n].Plan})
			display(columns{"STATUS", cluster[i].NodePools[n].Status})
			display(columns{"NODE QUANTITY", cluster[i].NodePools[n].NodeQuantity})
			display(columns{"AUTO SCALER", cluster[i].NodePools[n].AutoScaler})
			display(columns{"MIN NODES", cluster[i].NodePools[n].MinNodes})
			display(columns{"MAX NODES", cluster[i].NodePools[n].MaxNodes})

			display(columns{" "})
			display(columns{"NODES"})

			for j := range cluster[i].NodePools[n].Nodes {
				display(columns{"ID", "DATE CREATED", "LABEL", "STATUS"})
				display(columns{
					cluster[i].NodePools[n].Nodes[j].ID,
					cluster[i].NodePools[n].Nodes[j].DateCreated,
					cluster[i].NodePools[n].Nodes[j].Label,
					cluster[i].NodePools[n].Nodes[j].Status,
				})
			}
			display(columns{" "})
		}

		display(columns{"---------------------------"})
	}

	Meta(meta)
}

func Cluster(k *govultr.Cluster) {
	defer flush()

	display(columns{"ID", k.ID})
	display(columns{"LABEL", k.Label})
	display(columns{"DATE CREATED", k.DateCreated})
	display(columns{"CLUSTER SUBNET", k.ClusterSubnet})
	display(columns{"SERVICE SUBNET", k.ServiceSubnet})
	display(columns{"IP", k.IP})
	display(columns{"ENDPOINT", k.Endpoint})
	display(columns{"VERSION", k.Version})
	display(columns{"REGION", k.Region})
	display(columns{"STATUS", k.Status})

	display(columns{" "})
	display(columns{"NODE POOLS"})
	for i := range k.NodePools {
		display(columns{"ID", k.NodePools[i].ID})
		display(columns{"DATE CREATED", k.NodePools[i].DateCreated})
		display(columns{"DATE UPDATED", k.NodePools[i].DateUpdated})
		display(columns{"LABEL", k.NodePools[i].Label})
		display(columns{"TAG", k.NodePools[i].Tag})
		display(columns{"PLAN", k.NodePools[i].Plan})
		display(columns{"STATUS", k.NodePools[i].Status})
		display(columns{"NODE QUANTITY", k.NodePools[i].NodeQuantity})
		display(columns{"AUTO SCALER", k.NodePools[i].AutoScaler})
		display(columns{"MIN NODES", k.NodePools[i].MinNodes})
		display(columns{"MAX NODES", k.NodePools[i].MaxNodes})

		display(columns{" "})
		display(columns{"NODES"})

		for n := range k.NodePools[i].Nodes {
			display(columns{"ID", "DATE CREATED", "LABEL", "STATUS"})
			display(columns{
				k.NodePools[i].Nodes[n].ID,
				k.NodePools[i].Nodes[n].DateCreated,
				k.NodePools[i].Nodes[n].Label,
				k.NodePools[i].Nodes[n].Status,
			})
		}
		display(columns{" "})
	}
}

func NodePools(nodepool []govultr.NodePool, meta *govultr.Meta) {
	defer flush()

	if len(nodepool) == 0 {
		// this shouldn't be possible since at least one nodepool is required
		displayString("No active nodepools on cluster")
		return
	}

	for i := range nodepool {

		display(columns{"ID", nodepool[i].ID})
		display(columns{"DATE CREATED", nodepool[i].DateCreated})
		display(columns{"DATE UPDATED", nodepool[i].DateUpdated})
		display(columns{"LABEL", nodepool[i].Label})
		display(columns{"TAG", nodepool[i].Tag})
		display(columns{"PLAN", nodepool[i].Plan})
		display(columns{"STATUS", nodepool[i].Status})
		display(columns{"NODE QUANTITY", nodepool[i].NodeQuantity})
		display(columns{"AUTO SCALER", nodepool[i].AutoScaler})
		display(columns{"MIN NODES", nodepool[i].MinNodes})
		display(columns{"MAX NODES", nodepool[i].MaxNodes})

		display(columns{" "})
		display(columns{"NODES"})

		display(columns{"ID", "DATE CREATED", "LABEL", "STATUS"})
		for n := range nodepool[i].Nodes {
			display(columns{
				nodepool[i].Nodes[n].ID,
				nodepool[i].Nodes[n].DateCreated,
				nodepool[i].Nodes[n].Label,
				nodepool[i].Nodes[n].Status,
			})
		}
		display(columns{"---------------------------"})
	}

	Meta(meta)
}

func NodePool(np *govultr.NodePool) {
	defer flush()

	display(columns{"ID", np.ID})
	display(columns{"DATE CREATED", np.DateCreated})
	display(columns{"DATE UPDATED", np.DateUpdated})
	display(columns{"LABEL", np.Label})
	display(columns{"TAG", np.Tag})
	display(columns{"PLAN", np.Plan})
	display(columns{"STATUS", np.Status})
	display(columns{"NODE QUANTITY", np.NodeQuantity})
	display(columns{"AUTO SCALER", np.AutoScaler})
	display(columns{"MIN NODES", np.MinNodes})
	display(columns{"MAX NODES", np.MaxNodes})

	display(columns{" "})
	display(columns{"NODES"})

	display(columns{"ID", "DATE CREATED", "LABEL", "STATUS"})
	for i := range np.Nodes {
		display(columns{
			np.Nodes[i].ID,
			np.Nodes[i].DateCreated,
			np.Nodes[i].Label,
			np.Nodes[i].Status,
		})
	}
}

func ClustersSummary(clusters []govultr.Cluster, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "LABEL", "STATUS", "REGION", "VERSION", "NODEPOOL#", "NODE#"})

	if len(clusters) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range clusters {
		nodePoolCount := len(clusters[i].NodePools)
		var nodeCount int = 0

		for n := range clusters[i].NodePools {
			nodeCount += len(clusters[i].NodePools[n].Nodes)
		}

		display(columns{
			clusters[i].ID,
			clusters[i].Label,
			clusters[i].Status,
			clusters[i].Region,
			clusters[i].Version,
			nodePoolCount,
			nodeCount,
		})
	}

	Meta(meta)
}

func K8Versions(versions *govultr.Versions) {
	defer flush()

	display(columns{"VERSIONS"})

	if len(versions.Versions) == 0 {
		display(columns{"---"})
		return
	}

	for i := range versions.Versions {
		display(columns{versions.Versions[i]})
	}

}

func K8Upgrades(upgrades []string) {
	defer flush()

	display(columns{"UPGRADES"})

	if len(upgrades) == 0 {
		display(columns{"---"})
		return
	}

	for i := range upgrades {
		display(columns{upgrades[i]})
	}
}
