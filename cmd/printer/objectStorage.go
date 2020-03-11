package printer

import (
	"github.com/vultr/govultr"
)

func ObjectStorage(obj []govultr.ObjectStorage, options *govultr.ObjectListOptions) {
	col := columns{"ID", "REGION ID", "OBJSTORECLUSTER ID", "STATUS", "LABEL", "DATE CREATED"}
	if options.IncludeS3 {
		col = columns{"ID", "REGION ID", "OBJSTORECLUSTER ID", "STATUS", "LABEL", "DATE CREATED", "S3 HOSTNAME", "S3 ACCESS KEY", "S3 SECRET KEY"}
	}
	display(col)

	for _, o := range obj {
		vals := columns{o.ID, o.RegionID, o.ObjectStoreClusterID, o.Status, o.Label, o.DateCreated}
		if options.IncludeS3 {
			vals = columns{o.ID, o.RegionID, o.ObjectStoreClusterID, o.Status, o.Label, o.DateCreated, o.S3Keys.S3Hostname, o.S3Keys.S3AccessKey, o.S3Keys.S3SecretKey}
		}
		display(vals)
	}
	flush()
}

func SingleObjectStorage(obj *govultr.ObjectStorage, options *govultr.ObjectListOptions) {
	col := columns{"ID", "REGION ID", "OBJSTORECLUSTER ID", "STATUS", "LABEL", "DATE CREATED"}
	if options.IncludeS3 {
		col = columns{"ID", "REGION ID", "OBJSTORECLUSTER ID", "STATUS", "LABEL", "DATE CREATED", "S3 HOSTNAME", "S3 ACCESS KEY", "S3 SECRET KEY"}
	}
	display(col)

	vals := columns{obj.ID, obj.RegionID, obj.ObjectStoreClusterID, obj.Status, obj.Label, obj.DateCreated}
	if options.IncludeS3 {
		vals = columns{obj.ID, obj.RegionID, obj.ObjectStoreClusterID, obj.Status, obj.Label, obj.DateCreated, obj.S3Keys.S3Hostname, obj.S3Keys.S3AccessKey, obj.S3Keys.S3SecretKey}
	}
	display(vals)
	flush()
}

func ObjectStorageClusterList(cluster []govultr.ObjectStorageCluster) {
	display(columns{"OBJSTORECLUSTER ID", "REGION ID", "LOCATION", "HOSTNAME", "DEPLOY"})
	for _, c := range cluster {
		display(columns{c.ObjectStoreClusterID, c.RegionID, c.Location, c.Hostname, c.Deploy})
	}
	flush()
}

func ObjStorageS3KeyRegenerate(key *govultr.S3Keys) {
	display(columns{"S3 HOSTNAME", "S3 ACCESS KEY", "S3 SECRET KEY"})
	display(columns{key.S3Hostname, key.S3AccessKey, key.S3SecretKey})
	flush()
}
