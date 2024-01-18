package containerregistry

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
)

// ContainerRegistryPrinter ...
type ContainerRegistryPrinter struct {
	Registry *govultr.ContainerRegistry `json:"registry"`
}

// JSON ...
func (c *ContainerRegistryPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ContainerRegistryPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ContainerRegistryPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (c *ContainerRegistryPrinter) Data() [][]string {
	return [][]string{
		0:  {"ID", c.Registry.ID},
		1:  {"NAME", c.Registry.Name},
		2:  {"PUBLIC", strconv.FormatBool(c.Registry.Public)},
		3:  {"URN", c.Registry.URN},
		4:  {"REGION", c.Registry.Metadata.Region.Name},
		5:  {" "},
		6:  {"ROOT USER"},
		7:  {"ID", strconv.Itoa(c.Registry.RootUser.ID)},
		8:  {"USER NAME", c.Registry.RootUser.UserName},
		9:  {"PASSWORD", c.Registry.RootUser.Password},
		10: {"CREATED", c.Registry.RootUser.DateCreated},
		11: {"MODIFIED", c.Registry.RootUser.DateModified},
		12: {" "},
		13: {"STORAGE"},
		14: {"USED", fmt.Sprintf("%vGB", c.Registry.Storage.Used.GigaBytes)},
		15: {"ALLOWED", fmt.Sprintf("%vGB", c.Registry.Storage.Allowed.GigaBytes)},
		16: {" "},
		17: {"BILLING"},
		18: {"PRICE",
			strconv.FormatFloat(
				float64(c.Registry.Metadata.Subscription.Billing.MonthlyPrice),
				'f',
				utils.FloatPrecision,
				utils.FloatBitDepth,
			),
		},
		19: {"CHARGES",
			strconv.FormatFloat(
				float64(c.Registry.Metadata.Subscription.Billing.PendingCharges),
				'f',
				utils.FloatPrecision,
				utils.FloatBitDepth,
			),
		},
		20: {" "},
		21: {"CREATED", c.Registry.DateCreated},
	}
}

// Paging ...
func (c *ContainerRegistryPrinter) Paging() [][]string {
	return nil
}

// ======================================

// ContainerRegistriesPrinter ...
type ContainerRegistriesPrinter struct {
	Registries []govultr.ContainerRegistry `json:"registries"`
	Meta       *govultr.Meta               `json:"meta"`
}

// JSON ...
func (c *ContainerRegistriesPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ContainerRegistriesPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ContainerRegistriesPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"NAME",
		"URN",
		"USED/ALLOWED",
		"REGION ID",
		"REGION NAME",
		"PUBLIC",
	}}
}

// Data ...
func (c *ContainerRegistriesPrinter) Data() [][]string {
	if len(c.Registries) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range c.Registries {
		usage := fmt.Sprintf("%vGB / %vGB", c.Registries[i].Storage.Used.GigaBytes, c.Registries[i].Storage.Allowed.GigaBytes)
		data = append(data, []string{
			c.Registries[i].ID,
			c.Registries[i].Name,
			c.Registries[i].URN,
			usage,
			strconv.Itoa(c.Registries[i].Metadata.Region.ID),
			c.Registries[i].Metadata.Region.Name,
			strconv.FormatBool(c.Registries[i].Public),
		})
	}

	return data
}

// Paging ...
func (c *ContainerRegistriesPrinter) Paging() [][]string {
	return printer.NewPaging(c.Meta.Total, &c.Meta.Links.Next, &c.Meta.Links.Prev).Compose()
}

// ======================================

// ContainerRegistryPlansPrinter ...
type ContainerRegistryPlansPrinter struct {
	Plans govultr.ContainerRegistryPlanTypes `json:"plans"`
}

// JSON ...
func (c *ContainerRegistryPlansPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ContainerRegistryPlansPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ContainerRegistryPlansPrinter) Columns() [][]string {
	return [][]string{0: {
		"NAME",
		"MAX STORAGE",
		"MONTHLY PRICE",
	}}
}

// Data ...
func (c *ContainerRegistryPlansPrinter) Data() [][]string {
	var data [][]string
	topVals := reflect.ValueOf(c.Plans)
	for i := 0; i < topVals.NumField(); i++ {
		botVals := reflect.ValueOf(topVals.Field(i).Interface())

		data = append(data, []string{
			botVals.FieldByName("VanityName").String(),
			fmt.Sprintf("%vGB", botVals.FieldByName("MaxStorageMB").Int()/1024),
			strconv.FormatInt(botVals.FieldByName("MonthlyPrice").Int(), 10),
		})
	}
	return data
}

// Paging ...
func (c *ContainerRegistryPlansPrinter) Paging() [][]string {
	return nil
}

// ======================================

// ContainerRegistryRegionsPrinter ...
type ContainerRegistryRegionsPrinter struct {
	Regions []govultr.ContainerRegistryRegion `json:"regions"`
	Meta    *govultr.Meta                     `json:"meta"`
}

// JSON ...
func (c *ContainerRegistryRegionsPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ContainerRegistryRegionsPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ContainerRegistryRegionsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"NAME",
		"URN",
		"COUNTRY",
		"REGION",
	}}
}

// Data ...
func (c *ContainerRegistryRegionsPrinter) Data() [][]string {
	if len(c.Regions) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range c.Regions {
		data = append(data, []string{
			strconv.Itoa(c.Regions[i].ID),
			c.Regions[i].Name,
			c.Regions[i].URN,
			c.Regions[i].DataCenter.Country,
			c.Regions[i].DataCenter.Region,
		})
	}

	return data
}

// Paging ...
func (c *ContainerRegistryRegionsPrinter) Paging() [][]string {
	return printer.NewPaging(c.Meta.Total, &c.Meta.Links.Next, &c.Meta.Links.Prev).Compose()
}

// ======================================

// ContainerRegistryRepositoryPrinter ...
type ContainerRegistryRepositoryPrinter struct {
	Repository *govultr.ContainerRegistryRepo `json:"repository"`
}

// JSON ...
func (c *ContainerRegistryRepositoryPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ContainerRegistryRepositoryPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ContainerRegistryRepositoryPrinter) Columns() [][]string {
	return [][]string{0: {
		"NAME",
		"IMAGE",
		"DESCRIPTION",
		"DATE CREATED",
		"DATE MODIFIED",
		"PULLS",
		"ARTIFACTS",
	}}
}

// Data ...
func (c *ContainerRegistryRepositoryPrinter) Data() [][]string {
	return [][]string{0: {
		c.Repository.Name,
		c.Repository.Image,
		c.Repository.Description,
		c.Repository.DateCreated,
		c.Repository.DateModified,
		strconv.Itoa(c.Repository.PullCount),
		strconv.Itoa(c.Repository.ArtifactCount),
	}}
}

// Paging ...
func (c *ContainerRegistryRepositoryPrinter) Paging() [][]string {
	return nil
}

// ======================================

// ContainerRegistryRepositoriesPrinter ...
type ContainerRegistryRepositoriesPrinter struct {
	Repositories []govultr.ContainerRegistryRepo `json:"repositories"`
	Meta         *govultr.Meta                   `json:"meta"`
}

// JSON ...
func (c *ContainerRegistryRepositoriesPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ContainerRegistryRepositoriesPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ContainerRegistryRepositoriesPrinter) Columns() [][]string {
	return [][]string{0: {
		"NAME",
		"IMAGE",
		"DESCRIPTION",
		"DATE CREATED",
		"DATE MODIFIED",
		"PULLS",
		"ARTIFACTS",
	}}
}

// Data ...
func (c *ContainerRegistryRepositoriesPrinter) Data() [][]string {
	if len(c.Repositories) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range c.Repositories {
		data = append(data, []string{
			c.Repositories[i].Name,
			c.Repositories[i].Image,
			c.Repositories[i].Description,
			c.Repositories[i].DateCreated,
			c.Repositories[i].DateModified,
			strconv.Itoa(c.Repositories[i].PullCount),
			strconv.Itoa(c.Repositories[i].ArtifactCount),
		})
	}

	return data
}

// Paging ...
func (c *ContainerRegistryRepositoriesPrinter) Paging() [][]string {
	return printer.NewPaging(c.Meta.Total, &c.Meta.Links.Next, &c.Meta.Links.Prev).Compose()
}

// ======================================

// ContainerRegistryCredentialDockerPrinter ...
type ContainerRegistryCredentialDockerPrinter struct {
	Credential *govultr.ContainerRegistryDockerCredentials `json:"docker_credentials"`
}

// JSON ...
func (c *ContainerRegistryCredentialDockerPrinter) JSON() []byte {
	return printer.MarshalObject(c, "json")
}

// YAML ...
func (c *ContainerRegistryCredentialDockerPrinter) YAML() []byte {
	return printer.MarshalObject(c, "yaml")
}

// Columns ...
func (c *ContainerRegistryCredentialDockerPrinter) Columns() [][]string {
	return nil
}

// Data ...
func (c *ContainerRegistryCredentialDockerPrinter) Data() [][]string {
	return [][]string{0: {c.Credential.String()}}
}

// Paging ...
func (c *ContainerRegistryCredentialDockerPrinter) Paging() [][]string {
	return nil
}
