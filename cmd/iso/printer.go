package iso

import (
	"strconv"

	"github.com/vultr/govultr/v3"
	"github.com/vultr/vultr-cli/v3/cmd/printer"
)

// ISOsPrinter...
type ISOsPrinter struct {
	ISOs []govultr.ISO `json:"isos"`
	Meta *govultr.Meta `json:"meta"`
}

// JSON ...
func (i *ISOsPrinter) JSON() []byte {
	return printer.MarshalObject(i, "json")
}

// YAML ...
func (i *ISOsPrinter) YAML() []byte {
	return printer.MarshalObject(i, "yaml")
}

// Columns ...
func (i *ISOsPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"FILE NAME",
		"SIZE",
		"STATUS",
		"MD5SUM",
		"SHA512SUM",
		"DATE CREATED",
	}}
}

// Data ...
func (i *ISOsPrinter) Data() [][]string {
	if len(i.ISOs) == 0 {
		return [][]string{0: {"---", "---", "---", "---", "---", "---", "---"}}

	}

	var data [][]string
	for n := range i.ISOs {
		data = append(data, []string{
			i.ISOs[n].ID,
			i.ISOs[n].FileName,
			strconv.Itoa(i.ISOs[n].Size),
			i.ISOs[n].Status,
			i.ISOs[n].MD5Sum,
			i.ISOs[n].SHA512Sum,
			i.ISOs[n].DateCreated,
		})
	}

	return data
}

// Paging ...
func (i *ISOsPrinter) Paging() [][]string {
	return printer.NewPaging(i.Meta.Total, &i.Meta.Links.Next, &i.Meta.Links.Prev).Compose()
}

// ======================================

// ISOPrinter...
type ISOPrinter struct {
	ISO govultr.ISO `json:"iso"`
}

// JSON ...
func (i *ISOPrinter) JSON() []byte {
	return printer.MarshalObject(i, "json")
}

// YAML ...
func (i *ISOPrinter) YAML() []byte {
	return printer.MarshalObject(i, "yaml")
}

// Columns ...
func (i *ISOPrinter) Columns() [][]string {
	return [][]string{0: {
		"ID",
		"FILE NAME",
		"SIZE",
		"STATUS",
		"MD5SUM",
		"SHA512SUM",
		"DATE CREATED",
	}}
}

// Data ...
func (i *ISOPrinter) Data() [][]string {
	return [][]string{0: {
		i.ISO.ID,
		i.ISO.FileName,
		strconv.Itoa(i.ISO.Size),
		i.ISO.Status,
		i.ISO.MD5Sum,
		i.ISO.SHA512Sum,
		i.ISO.DateCreated,
	}}
}

// Paging ...
func (i *ISOPrinter) Paging() [][]string {
	return nil
}

// ======================================

// PublicISOsPrinter...
type PublicISOsPrinter struct {
	ISOs []govultr.PublicISO `json:"public_isos"`
	Meta *govultr.Meta       `json:"meta"`
}

// JSON ...
func (i *PublicISOsPrinter) JSON() []byte {
	return printer.MarshalObject(i, "json")
}

// YAML ...
func (i *PublicISOsPrinter) YAML() []byte {
	return printer.MarshalObject(i, "yaml")
}

// Columns ...
func (i *PublicISOsPrinter) Columns() [][]string {
	return [][]string{0: {"ID", "NAME", "DESCRIPTION"}}
}

// Data ...
func (i *PublicISOsPrinter) Data() [][]string {
	if len(i.ISOs) == 0 {
		return [][]string{0: {"---", "---", "---"}}

	}

	var data [][]string
	for n := range i.ISOs {
		data = append(data, []string{
			i.ISOs[n].ID,
			i.ISOs[n].Name,
			i.ISOs[n].Description,
		})
	}

	return data
}

// Paging ...
func (i *PublicISOsPrinter) Paging() [][]string {
	return printer.NewPaging(i.Meta.Total, &i.Meta.Links.Next, &i.Meta.Links.Prev).Compose()
}
