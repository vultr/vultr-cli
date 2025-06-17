package objectstorage

import (
	"strconv"
	"strings"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
	"github.com/vultr/vultr-cli/v3/cmd/utils"
)

// ObjectStoragesPrinter ...
type ObjectStoragesPrinter struct {
	ObjectStorages []govultr.ObjectStorage `json:"object_storages"`
	Meta           *govultr.Meta           `json:"meta"`
}

// JSON ...
func (o *ObjectStoragesPrinter) JSON() []byte {
	return printer.MarshalObject(o, "json")
}

// YAML ...
func (o *ObjectStoragesPrinter) YAML() []byte {
	return printer.MarshalObject(o, "yaml")
}

// Columns ...
func (o *ObjectStoragesPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"REGION",
		"CLUSTER ID",
		"STATUS",
		"LABEL",
		"DATE CREATED",
		"S3 HOSTNAME",
		"S3 ACCESS KEY",
		"S3 SECRET KEY",
	}}
}

// Data ...
func (o *ObjectStoragesPrinter) Data() [][]string {
	if len(o.ObjectStorages) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range o.ObjectStorages {
		data = append(data, []string{
			o.ObjectStorages[i].ID,
			o.ObjectStorages[i].Region,
			strconv.Itoa(o.ObjectStorages[i].ObjectStoreClusterID),
			o.ObjectStorages[i].Status,
			o.ObjectStorages[i].Label,
			o.ObjectStorages[i].DateCreated,
			o.ObjectStorages[i].S3Hostname,
			o.ObjectStorages[i].S3AccessKey,
			o.ObjectStorages[i].S3SecretKey,
		})
	}

	return data
}

// Paging ...
func (o *ObjectStoragesPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(o.Meta).Compose()
}

// ======================================

// ObjectStoragePrinter ...
type ObjectStoragePrinter struct {
	ObjectStorage *govultr.ObjectStorage `json:"object_storage"`
}

// JSON ...
func (o *ObjectStoragePrinter) JSON() []byte {
	return printer.MarshalObject(o, "json")
}

// YAML ...
func (o *ObjectStoragePrinter) YAML() []byte {
	return printer.MarshalObject(o, "yaml")
}

// Columns ...
func (o *ObjectStoragePrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"REGION",
		"CLUSTER ID",
		"STATUS",
		"LABEL",
		"DATE CREATED",
		"S3 HOSTNAME",
		"S3 ACCESS KEY",
		"S3 SECRET KEY",
	}}
}

// Data ...
func (o *ObjectStoragePrinter) Data() [][]string {
	return [][]string{0: {
		o.ObjectStorage.ID,
		o.ObjectStorage.Region,
		strconv.Itoa(o.ObjectStorage.ObjectStoreClusterID),
		o.ObjectStorage.Status,
		o.ObjectStorage.Label,
		o.ObjectStorage.DateCreated,
		o.ObjectStorage.S3Hostname,
		o.ObjectStorage.S3AccessKey,
		o.ObjectStorage.S3SecretKey,
	}}
}

// Paging ...
func (o *ObjectStoragePrinter) Paging() [][]string {
	return nil
}

// ======================================

// ObjectStorageClustersPrinter ...
type ObjectStorageClustersPrinter struct {
	Clusters []govultr.ObjectStorageCluster `json:"clusters"`
	Meta     *govultr.Meta                  `json:"meta"`
}

// JSON ...
func (o *ObjectStorageClustersPrinter) JSON() []byte {
	return printer.MarshalObject(o, "json")
}

// YAML ...
func (o *ObjectStorageClustersPrinter) YAML() []byte {
	return printer.MarshalObject(o, "yaml")
}

// Columns ...
func (o *ObjectStorageClustersPrinter) Columns() [][]string {
	return [][]string{0: {
		"CLUSTER ID",
		"REGION ID",
		"HOSTNAME",
		"DEPLOY",
	}}
}

// Data ...
func (o *ObjectStorageClustersPrinter) Data() [][]string {
	if len(o.Clusters) == 0 {
		return [][]string{0: {"---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range o.Clusters {
		data = append(data, []string{
			strconv.Itoa(o.Clusters[i].ID),
			o.Clusters[i].Region,
			o.Clusters[i].Hostname,
			o.Clusters[i].Deploy,
		})
	}

	return data
}

// Paging ...
func (o *ObjectStorageClustersPrinter) Paging() [][]string {
	return printer.NewPagingFromMeta(o.Meta).Compose()
}

// ======================================

// ObjectStorageKeysPrinter ...
type ObjectStorageKeysPrinter struct {
	Keys *govultr.S3Keys `json:"s3_credentials"`
}

// JSON ...
func (o *ObjectStorageKeysPrinter) JSON() []byte {
	return printer.MarshalObject(o, "json")
}

// YAML ...
func (o *ObjectStorageKeysPrinter) YAML() []byte {
	return printer.MarshalObject(o, "yaml")
}

// Columns ...
func (o *ObjectStorageKeysPrinter) Columns() [][]string {
	return [][]string{0: {
		"S3 HOSTNAME",
		"S3 ACCESS KEY",
		"S3 SECRET KEY",
	}}
}

// Data ...
func (o *ObjectStorageKeysPrinter) Data() [][]string {
	return [][]string{0: {
		o.Keys.S3Hostname,
		o.Keys.S3AccessKey,
		o.Keys.S3SecretKey,
	}}
}

// Paging ...
func (o *ObjectStorageKeysPrinter) Paging() [][]string {
	return nil
}

// ======================================

// ObjectStorageTiersPrinter ...
type ObjectStorageTiersPrinter struct {
	Tiers []govultr.ObjectStorageTier `json:"tiers"`
}

// JSON ...
func (o *ObjectStorageTiersPrinter) JSON() []byte {
	return printer.MarshalObject(o, "json")
}

// YAML ...
func (o *ObjectStorageTiersPrinter) YAML() []byte {
	return printer.MarshalObject(o, "yaml")
}

// Columns ...
func (o *ObjectStorageTiersPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"SLUG",
		"PRICE",
		"PRICE BANDWIDTH",
		"PRICE DISK",
		"RATELIMIT OPS",
		"RATELIMIT BYTES",
		"AVAILABLE REGIONS",
	}}
}

// Data ...
func (o *ObjectStorageTiersPrinter) Data() [][]string {
	if len(o.Tiers) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range o.Tiers {
		var regions []string
		for j := range o.Tiers[i].Locations {
			regions = append(regions, o.Tiers[i].Locations[j].Region)
		}
		data = append(data, []string{
			strconv.Itoa(o.Tiers[i].ID),
			o.Tiers[i].Slug,
			strconv.FormatFloat(
				float64(o.Tiers[i].Price),
				'f',
				utils.FloatPrecision,
				utils.FloatBitDepth,
			),
			strconv.FormatFloat(
				float64(o.Tiers[i].PriceBandwidthGB),
				'f',
				utils.FloatPrecision,
				utils.FloatBitDepth,
			),
			strconv.FormatFloat(
				float64(o.Tiers[i].PriceDiskGB),
				'f',
				utils.FloatPrecision,
				utils.FloatBitDepth,
			),
			strconv.Itoa(o.Tiers[i].RateLimitOpsSec),
			strconv.Itoa(o.Tiers[i].RateLimitBytesSec),
			strings.Join(regions, ", "),
		})
	}

	return data
}

// Paging ...
func (o *ObjectStorageTiersPrinter) Paging() [][]string {
	return nil
}

// ======================================

// ObjectStorageTiersPrinter ...
type ObjectStorageClusterTiersPrinter struct {
	Tiers []govultr.ObjectStorageTier `json:"tiers"`
}

// JSON ...
func (o *ObjectStorageClusterTiersPrinter) JSON() []byte {
	return printer.MarshalObject(o, "json")
}

// YAML ...
func (o *ObjectStorageClusterTiersPrinter) YAML() []byte {
	return printer.MarshalObject(o, "yaml")
}

// Columns ...
func (o *ObjectStorageClusterTiersPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"SLUG",
		"PRICE",
		"PRICE BANDWIDTH",
		"PRICE DISK",
		"RATELIMIT OPS",
		"RATELIMIT BYTES",
	}}
}

// Data ...
func (o *ObjectStorageClusterTiersPrinter) Data() [][]string {
	if len(o.Tiers) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---"}}
	}

	var data [][]string
	for i := range o.Tiers {
		data = append(data, []string{
			strconv.Itoa(o.Tiers[i].ID),
			o.Tiers[i].Slug,
			strconv.FormatFloat(
				float64(o.Tiers[i].Price),
				'f',
				utils.FloatPrecision,
				utils.FloatBitDepth,
			),
			strconv.FormatFloat(
				float64(o.Tiers[i].PriceBandwidthGB),
				'f',
				utils.FloatPrecision,
				utils.FloatBitDepth,
			),
			strconv.FormatFloat(
				float64(o.Tiers[i].PriceDiskGB),
				'f',
				utils.FloatPrecision,
				utils.FloatBitDepth,
			),
			strconv.Itoa(o.Tiers[i].RateLimitOpsSec),
			strconv.Itoa(o.Tiers[i].RateLimitBytesSec),
		})
	}

	return data
}

// Paging ...
func (o *ObjectStorageClusterTiersPrinter) Paging() [][]string {
	return nil
}
