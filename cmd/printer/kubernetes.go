package printer

import (
	"github.com/vultr/govultr/v3"
)

func Clusters(cluster []govultr.Cluster, meta *govultr.Meta) {
	for _, k := range cluster {
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
		for _, np := range k.NodePools {
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

			for _, n := range np.Nodes {
				display(columns{"ID", "DATE CREATED", "LABEL", "STATUS"})
				display(columns{n.ID, n.DateCreated, n.Label, n.Status})
			}
			display(columns{" "})
		}

		display(columns{"---------------------------"})
	}

	Meta(meta)
	flush()
}

func Cluster(k *govultr.Cluster) {
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
	for _, np := range k.NodePools {
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

		for _, n := range np.Nodes {
			display(columns{"ID", "DATE CREATED", "LABEL", "STATUS"})
			display(columns{n.ID, n.DateCreated, n.Label, n.Status})
		}
		display(columns{" "})
	}

	flush()
}

func NodePools(nodepool []govultr.NodePool, meta *govultr.Meta) {
	for _, np := range nodepool {

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

		for _, n := range np.Nodes {
			display(columns{"ID", "DATE CREATED", "LABEL", "STATUS"})
			display(columns{n.ID, n.DateCreated, n.Label, n.Status})
		}
		display(columns{"---------------------------"})
	}

	Meta(meta)
	flush()
}

func NodePool(np *govultr.NodePool) {
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

	for _, n := range np.Nodes {
		display(columns{"ID", "DATE CREATED", "LABEL", "STATUS"})
		display(columns{n.ID, n.DateCreated, n.Label, n.Status})
	}

	flush()
}

func ClustersSummary(clusters []govultr.Cluster, meta *govultr.Meta) {
	display(columns{"ID", "LABEL", "STATUS", "REGION", "VERSION", "NODEPOOL#", "NODE#"})

	for _, k := range clusters {
		nodePoolCount := len(k.NodePools)
		var nodeCount int = 0

		for _, np := range k.NodePools {
			nodeCount += len(np.Nodes)
		}

		display(columns{k.ID, k.Label, k.Status, k.Region, k.Version, nodePoolCount, nodeCount})
	}

	Meta(meta)
	flush()
}

func K8Versions(versions *govultr.Versions) {
	display(columns{"VERSIONS"})
	for _, v := range versions.Versions {
		display(columns{v})
	}

	flush()
}

func K8Upgrades(upgrades []string) {
	display(columns{"UPGRADES"})
	for _, v := range upgrades {
		display(columns{v})
	}

	flush()
}
