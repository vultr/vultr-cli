package printer

import "github.com/vultr/govultr"

func SecInfo(info []string) {
	col := columns{"DNSSEC INFO"}
	display(col)
	for _, i := range info {
		display(columns{i})
	}
	flush()
}

func DomainList(domain []govultr.DNSDomain) {
	col := columns{"DOMAIN", "DATE CREATED"}
	display(col)
	for _, d := range domain {
		display(columns{d.Domain, d.DateCreated})
	}
	flush()
}

func SoaInfo(soa *govultr.Soa) {
	col := columns{"NS PRIMARY", "EMAIL"}
	display(col)
	display(columns{soa.NsPrimary, soa.Email})
	flush()
}
