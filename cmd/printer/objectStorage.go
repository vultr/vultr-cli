package printer

import (
	"github.com/vultr/govultr/v3"
)

func ObjectStorages(obj []govultr.ObjectStorage, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "REGION", "OBJSTORECLUSTER ID", "STATUS", "LABEL", "DATE CREATED", "S3 HOSTNAME", "S3 ACCESS KEY", "S3 SECRET KEY"})

	if len(obj) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "S3 HOSTNAME", "S3 ACCESS KEY", "S3 SECRET KEY"})
		Meta(meta)
		return
	}

	for i := range obj {
		vals := columns{
			obj[i].ID,
			obj[i].Region,
			obj[i].ObjectStoreClusterID,
			obj[i].Status,
			obj[i].Label,
			obj[i].DateCreated,
			obj[i].S3Keys.S3Hostname,
			obj[i].S3Keys.S3AccessKey,
			obj[i].S3Keys.S3SecretKey,
		}

		display(vals)
	}

	Meta(meta)
}

func SingleObjectStorage(obj *govultr.ObjectStorage) {
	defer flush()

	display(columns{"ID", "REGION", "OBJSTORECLUSTER ID", "STATUS", "LABEL", "DATE CREATED", "S3 HOSTNAME", "S3 ACCESS KEY", "S3 SECRET KEY"})
	display(columns{obj.ID, obj.Region, obj.ObjectStoreClusterID, obj.Status, obj.Label, obj.DateCreated, obj.S3Keys.S3Hostname, obj.S3Keys.S3AccessKey, obj.S3Keys.S3SecretKey})
}

func ObjectStorageClusterList(cluster []govultr.ObjectStorageCluster, meta *govultr.Meta) {
	defer flush()

	display(columns{"OBJSTORECLUSTER", "REGION ID", "HOSTNAME", "DEPLOY"})

	if len(cluster) == 0 {
		display(columns{"---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range cluster {
		display(columns{
			cluster[i].ID,
			cluster[i].Region,
			cluster[i].Hostname,
			cluster[i].Deploy,
		})
	}

	Meta(meta)
}

func ObjStorageS3KeyRegenerate(key *govultr.S3Keys) {
	defer flush()

	display(columns{"S3 HOSTNAME", "S3 ACCESS KEY", "S3 SECRET KEY"})
	display(columns{key.S3Hostname, key.S3AccessKey, key.S3SecretKey})
}
