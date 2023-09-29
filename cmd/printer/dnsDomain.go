package printer

import "github.com/vultr/govultr/v3"

func SecInfo(info []string) {
	defer flush()

	display(columns{"DNSSEC INFO"})

	if len(info) == 0 {
		display(columns{"---"})
		Meta(meta)
		return
	}

	for i := range info {
		display(columns{
			info[i],
		})
	}
}

func DomainList(domain []govultr.Domain, meta *govultr.Meta) {
	defer flush()

	display(columns{"DOMAIN", "DATE CREATED", "DNS SEC"})

	if len(domain) == 0 {
		display(columns{"---", "---", "---"})
		Meta(meta)
		return
	}

	for i := range domain {
		display(columns{
			domain[i].Domain,
			domain[i].DateCreated,
			domain[i].DNSSec,
		})
	}

	Meta(meta)
}

func Domain(domain *govultr.Domain) {
	defer flush()

	display(columns{"DOMAIN", "DATE CREATED", "DNS SEC"})
	display(columns{domain.Domain, domain.DateCreated, domain.DNSSec})
}

func SoaInfo(soa *govultr.Soa) {
	defer flush()

	display(columns{"NS PRIMARY", "EMAIL"})
	display(columns{soa.NSPrimary, soa.Email})
}
