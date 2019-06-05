package printer

import "github.com/vultr/govultr"

func IsoPrivate(iso []govultr.ISO) {
	col := columns{"ID", "File Name", "SIZE", "STATUS", "MD5SUM", "SHA512SUM", "DATE CREATED"}
	display(col)
	for _, i := range iso {
		display(columns{i.ISOID, i.FileName, i.Size, i.Status, i.MD5Sum, i.SHA512Sum, i.DateCreated})
	}
	flush()
}

func IsoPublic(iso []govultr.PublicISO) {
	col := columns{"ID", "Name", "DESCRIPTION"}
	display(col)
	for _, i := range iso {
		display(columns{i.ISOID, i.Name, i.Description})
	}
	flush()
}
