package printer

import (
	"github.com/vultr/govultr"
)

func ObjectStorages(obj []govultr.ObjectStorage, meta *govultr.Meta) {
	display(columns{"ID", "REGION", "OBJSTORECLUSTER ID", "STATUS", "LABEL", "DATE CREATED", "S3 HOSTNAME", "S3 ACCESS KEY", "S3 SECRET KEY"})
	for _, o := range obj {
		vals := columns{o.ID, o.Region, o.ObjectStoreClusterID, o.Status, o.Label, o.DateCreated, o.S3Keys.S3Hostname, o.S3Keys.S3AccessKey, o.S3Keys.S3SecretKey}
		display(vals)
	}

	Meta(meta)
	flush()
}

func SingleObjectStorage(obj *govultr.ObjectStorage) {
	display(columns{"ID", "REGION", "OBJSTORECLUSTER ID", "STATUS", "LABEL", "DATE CREATED", "S3 HOSTNAME", "S3 ACCESS KEY", "S3 SECRET KEY"})
	display(columns{obj.ID, obj.Region, obj.ObjectStoreClusterID, obj.Status, obj.Label, obj.DateCreated, obj.S3Keys.S3Hostname, obj.S3Keys.S3AccessKey, obj.S3Keys.S3SecretKey})

	flush()
}

func ObjectStorageClusterList(cluster []govultr.ObjectStorageCluster, meta *govultr.Meta) {
	display(columns{"OBJSTORECLUSTER", "REGION ID", "HOSTNAME", "DEPLOY"})
	for _, c := range cluster {
		display(columns{c.ID, c.Region, c.Hostname, c.Deploy})
	}

	Meta(meta)
	flush()
}

func ObjStorageS3KeyRegenerate(key *govultr.S3Keys) {
	display(columns{"S3 HOSTNAME", "S3 ACCESS KEY", "S3 SECRET KEY"})
	display(columns{key.S3Hostname, key.S3AccessKey, key.S3SecretKey})
	flush()
}
