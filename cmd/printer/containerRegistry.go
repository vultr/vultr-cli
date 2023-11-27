package printer

import (
	"fmt"
	"reflect"

	"github.com/vultr/govultr/v3"
)

func ContainerRegistry(cr *govultr.ContainerRegistry) {
	defer flush()

	display(columns{"ID", cr.ID})
	display(columns{"NAME", cr.Name})
	display(columns{"PUBLIC", cr.Public})
	display(columns{"URN", cr.URN})
	display(columns{"REGION", cr.Metadata.Region.Name})

	display(columns{" "})

	display(columns{"ROOT USER"})
	display(columns{"ID", cr.RootUser.ID})
	display(columns{"USER NAME", cr.RootUser.UserName})
	display(columns{"PASSWORD", cr.RootUser.Password})
	display(columns{"CREATED", cr.RootUser.DateCreated})
	display(columns{"MODIFIED", cr.RootUser.DateModified})

	display(columns{" "})

	display(columns{"STORAGE"})
	display(columns{"USED", fmt.Sprintf("%vGB", cr.Storage.Used.GigaBytes)})
	display(columns{"ALLOWED", fmt.Sprintf("%vGB", cr.Storage.Allowed.GigaBytes)})

	display(columns{" "})

	display(columns{"BILLING"})
	display(columns{"PRICE", cr.Metadata.Subscription.Billing.MonthlyPrice})
	display(columns{"CHARGES", cr.Metadata.Subscription.Billing.PendingCharges})

	display(columns{" "})

	display(columns{"CREATED", cr.DateCreated})
}

func ContainerRegistryList(crs []govultr.ContainerRegistry, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "NAME", "URN", "USED/ALLOWED", "REGION ID", "REGION NAME"})
	if len(crs) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range crs {
		usage := fmt.Sprintf("%vGB / %vGB", crs[i].Storage.Used.GigaBytes, crs[i].Storage.Allowed.GigaBytes)

		display(columns{
			crs[i].ID,
			crs[i].Name,
			crs[i].URN,
			usage,
			crs[i].Metadata.Region.ID,
			crs[i].Metadata.Region.Name,
		})
	}

	Meta(meta)
}

func ContainerRegistryPlans(plans govultr.ContainerRegistryPlanTypes) {
	defer flush()

	display(columns{"NAME", "MAX STORAGE", "MONTHLY PRICE"})

	topVals := reflect.ValueOf(plans)
	for i := 0; i < topVals.NumField(); i++ {

		botVals := reflect.ValueOf(topVals.Field(i).Interface())

		display(columns{
			botVals.FieldByName("VanityName").String(),
			fmt.Sprintf("%vGB", botVals.FieldByName("MaxStorageMB").Int()/1024),
			botVals.FieldByName("MonthlyPrice"),
		})
	}
}

func ContainerRegistryRegions(plans []govultr.ContainerRegistryRegion) {
	defer flush()

	display(columns{"ID", "NAME", "URN", "COUNTRY", "REGION"})
	if len(plans) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---"})
		return
	}

	for i := range plans {
		display(columns{
			plans[i].ID,
			plans[i].Name,
			plans[i].URN,
			plans[i].DataCenter.Country,
			plans[i].DataCenter.Region,
		})
	}
}

func ContainerRegistryRepository(repo *govultr.ContainerRegistryRepo) {
	defer flush()

	display(columns{"NAME", "IMAGE", "DESCRIPTION", "DATE CREATED", "DATE MODIFIED", "PULLS", "ARTIFACTS"})

	display(columns{
		repo.Name,
		repo.Image,
		repo.Description,
		repo.DateCreated,
		repo.DateModified,
		repo.PullCount,
		repo.ArtifactCount,
	})
}

func ContainerRegistryRepositoryList(repos []govultr.ContainerRegistryRepo, meta *govultr.Meta) {
	defer flush()

	display(columns{"NAME", "IMAGE", "DESCRIPTION", "DATE CREATED", "DATE MODIFIED", "PULLS", "ARTIFACTS"})
	if len(repos) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range repos {
		display(columns{
			repos[i].Name,
			repos[i].Image,
			repos[i].Description,
			repos[i].DateCreated,
			repos[i].DateModified,
			repos[i].PullCount,
			repos[i].ArtifactCount,
		})
	}

	Meta(meta)
}

func ContainerRegistryDockerCredentials(creds *govultr.ContainerRegistryDockerCredentials) {
	defer flush()

	display(columns{creds.String()})
}
