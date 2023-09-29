package printer

import "github.com/vultr/govultr/v3"

func IsoPrivates(iso []govultr.ISO, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "FILE NAME", "SIZE", "STATUS", "MD5SUM", "SHA512SUM", "DATE CREATED"})

	if len(iso) == 0 {
		display(columns{"---", "---", "---", "---", "---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range iso {
		display(columns{
			iso[i].ID,
			iso[i].FileName,
			iso[i].Size,
			iso[i].Status,
			iso[i].MD5Sum,
			iso[i].SHA512Sum,
			iso[i].DateCreated,
		})
	}

	Meta(meta)
}

func IsoPrivate(iso *govultr.ISO) {
	defer flush()

	display(columns{"ID", "FILE NAME", "SIZE", "STATUS", "MD5SUM", "SHA512SUM", "DATE CREATED"})
	display(columns{iso.ID, iso.FileName, iso.Size, iso.Status, iso.MD5Sum, iso.SHA512Sum, iso.DateCreated})
}

func IsoPublic(iso []govultr.PublicISO, meta *govultr.Meta) {
	defer flush()

	display(columns{"ID", "NAME", "DESCRIPTION"})

	if len(iso) == 0 {
		display(columns{"---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range iso {
		display(columns{
			iso[i].ID,
			iso[i].Name,
			iso[i].Description,
		})
	}

	Meta(meta)
}
