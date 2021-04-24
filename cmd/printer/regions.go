package printer

import (
	"encoding/json"

	"github.com/go-yaml/yaml"
	"github.com/vultr/govultr/v2"
)

var _ ResourceOutput = &Regions{}

type Regions struct {
	Regions []govultr.Region `json:"regions"`
	Meta    *govultr.Meta
}

func (r *Regions) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

func (r *Regions) Yaml() []byte {
	yam, err := yaml.Marshal(r)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

func (r *Regions) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"ID", "CITY", "COUNTRY", "CONTINENT", "OPTIONS"}}
}
func (r *Regions) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	for k, r := range r.Regions {
		data[k] = []interface{}{r.ID, r.City, r.Country, r.Continent, r.Options}
	}
	return data
}

func (r *Regions) Paging() map[int][]interface{} {
	return map[int][]interface{}{
		0: {"======================================"},
		1: {"TOTAL", "NEXT PAGE", "PREV PAGE"},
		2: {r.Meta.Total, r.Meta.Links.Next, r.Meta.Links.Prev},
	}
}

type RegionsAvailability struct {
	AvailablePlans *govultr.PlanAvailability `json:"available_plans"`
}

func (r *RegionsAvailability) JSON() []byte {
	prettyJSON, err := json.MarshalIndent(r.AvailablePlans, "", "    ")
	if err != nil {
		panic("move this into byte")
	}

	return prettyJSON
}

func (r *RegionsAvailability) Yaml() []byte {
	yam, err := yaml.Marshal(r.AvailablePlans)
	if err != nil {
		panic("move this into byte")
	}
	return yam
}

func (r *RegionsAvailability) Columns() map[int][]interface{} {
	return map[int][]interface{}{0: {"AVAILABLE PLANS"}}
}

func (r *RegionsAvailability) Data() map[int][]interface{} {
	data := map[int][]interface{}{}
	for k, r := range r.AvailablePlans.AvailablePlans {
		data[k] = []interface{}{r}
	}
	return data
}

func (r RegionsAvailability) Paging() map[int][]interface{} {
	return nil
}
