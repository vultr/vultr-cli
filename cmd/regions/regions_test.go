package regions

import (
	"context"
	"reflect"
	"testing"

	"github.com/vultr/vultr-cli/pkg/cli"

	"github.com/vultr/govultr/v3"
)

type mockVultrRegions struct {
	client *govultr.Client
}

func (m mockVultrRegions) Availability(ctx context.Context, regionID string, planType string) (*govultr.PlanAvailability, error) {
	return &govultr.PlanAvailability{AvailablePlans: []string{"1"}}, nil
}

func (m mockVultrRegions) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.Region, *govultr.Meta, error) {
	return []govultr.Region{{
			ID:        "ewr",
			City:      "NJ",
			Country:   "US",
			Continent: "NA",
			Options:   []string{"test"},
		}}, &govultr.Meta{
			Total: 0,
			Links: nil,
		}, nil
}

func TestOptions_List(t *testing.T) {
	avail := NewRegionOptions(&cli.Base{Client: &govultr.Client{Region: mockVultrRegions{nil}}})

	expected := []govultr.Region{{
		ID:        "ewr",
		City:      "NJ",
		Country:   "US",
		Continent: "NA",
		Options:   []string{"test"},
	}}
	expectedMeta := &govultr.Meta{
		Total: 0,
		Links: nil,
	}

	regions, meta, _ := avail.List()
	if !reflect.DeepEqual(expected, regions) {
		t.Errorf("RegionOptions.list returned %v expected %v", regions, expected)
	}

	if !reflect.DeepEqual(expectedMeta, meta) {
		t.Errorf("RegionOptions.list returned %v expected %v", meta, expectedMeta)
	}
}

func TestOptions_Availability(t *testing.T) {
	avail := NewRegionOptions(&cli.Base{Client: &govultr.Client{Region: mockVultrRegions{nil}}, Args: []string{"test"}})
	expected := &govultr.PlanAvailability{AvailablePlans: []string{"1"}}

	a, _ := avail.Availability()
	if !reflect.DeepEqual(expected, a) {
		t.Errorf("RegionAvailability.list returned %v expected %v", a, expected)
	}
}

func TestNewCmdRegion(t *testing.T) {
	cmd := NewCmdRegion(&cli.Base{Client: &govultr.Client{Region: mockVultrRegions{nil}}})

	if cmd.Short != "get regions" {
		t.Errorf("invalid short")
	}

	if cmd.Use != "regions" {
		t.Errorf("invalid regions")
	}

	alias := []string{"r", "region"}
	if !reflect.DeepEqual(cmd.Aliases, alias) {
		t.Errorf("expected alias %v got %v", alias, cmd.Aliases)
	}
}

func TestNewRegionOptions(t *testing.T) {
	options := NewRegionOptions(&cli.Base{Client: &govultr.Client{Region: mockVultrRegions{nil}}})

	ref := reflect.TypeOf(options)
	if _, ok := ref.MethodByName("List"); !ok {
		t.Errorf("Missing list function")
	}

	if _, ok := ref.MethodByName("validate"); ok {
		t.Errorf("validate isn't exported shouldn't be accessible")
	}

	rInterface := reflect.TypeOf(new(Interface)).Elem()
	if !ref.Implements(rInterface) {
		t.Errorf("Options does not implement Interface")
	}
}
