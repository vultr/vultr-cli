package printer

import "github.com/vultr/govultr"

func IsoPrivates(iso []govultr.ISO, meta *govultr.Meta) {
	col := columns{"ID", "FILE NAME", "SIZE", "STATUS", "MD5SUM", "SHA512SUM", "DATE CREATED"}
	display(col)
	for _, i := range iso {
		display(columns{i.ID, i.FileName, i.Size, i.Status, i.MD5Sum, i.SHA512Sum, i.DateCreated})
	}

	Meta(meta)
	flush()
}

func IsoPrivate(iso *govultr.ISO) {
	col := columns{"ID", "FILE NAME", "SIZE", "STATUS", "MD5SUM", "SHA512SUM", "DATE CREATED"}
	display(col)
	display(columns{iso.ID, iso.FileName, iso.Size, iso.Status, iso.MD5Sum, iso.SHA512Sum, iso.DateCreated})
	flush()
}

func IsoPublic(iso []govultr.PublicISO, meta *govultr.Meta) {
	col := columns{"ID", "NAME", "DESCRIPTION"}
	display(col)
	for _, i := range iso {
		display(columns{i.ID, i.Name, i.Description})
	}

	Meta(meta)
	flush()
}
