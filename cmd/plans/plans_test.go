package plans

import (
	"context"
	"github.com/vultr/govultr/v2"
	"reflect"
	"testing"
)

type mockVultrPlans struct {
	client *govultr.Client
}

func (m mockVultrPlans) List(ctx context.Context, planType string, options *govultr.ListOptions) ([]govultr.Plan, *govultr.Meta, error) {
	plans := []govultr.Plan{
		{
			ID:          "vhf-8c-32gb",
			VCPUCount:   8,
			RAM:         32768,
			Disk:        512,
			DiskCount:   1,
			Bandwidth:   6144,
			MonthlyCost: 192,
			Type:        "vhf",
			Locations:   []string{"ewr"},
		},
	}

	meta := &govultr.Meta{
		Total: 1,
		Links: nil,
	}

	return plans, meta, nil
}

func (m mockVultrPlans) ListBareMetal(ctx context.Context, options *govultr.ListOptions) ([]govultr.BareMetalPlan, *govultr.Meta, error) {
	metalPlans := []govultr.BareMetalPlan{
		{
			ID:          "vbm-4c-32gb",
			CPUCount:    4,
			CPUModel:    "E3-1270v6",
			CPUThreads:  8,
			RAM:         8,
			Disk:        20,
			DiskCount:   1,
			Bandwidth:   10,
			MonthlyCost: 120,
			Type:        "NVMe",
			Locations:   []string{"ewr"},
		},
	}
	meta := &govultr.Meta{
		Total: 1,
		Links: nil,
	}
	return metalPlans, meta, nil
}

func TestPlanOptions_List(t *testing.T) {
	client := &govultr.Client{Plan: mockVultrPlans{nil}}
	planOption := NewPlanOptions(client)

	expectedPlan := []govultr.Plan{
		{
			ID:          "vhf-8c-32gb",
			VCPUCount:   8,
			RAM:         32768,
			Disk:        512,
			DiskCount:   1,
			Bandwidth:   6144,
			MonthlyCost: 192,
			Type:        "vhf",
			Locations:   []string{"ewr"},
		},
	}
	expectedMeta := &govultr.Meta{
		Total: 1,
		Links: nil,
	}

	plan, meta, _ := planOption.List()

	if !reflect.DeepEqual(expectedPlan, plan) {
		t.Errorf("PlanOptions.list returned %v expected %v", plan, expectedPlan)
	}

	if !reflect.DeepEqual(expectedMeta, meta) {
		t.Errorf("PlanOptions.list returned %v expected %v", plan, expectedPlan)
	}
}

func TestPlanOptions_Metal(t *testing.T) {
	client := &govultr.Client{Plan: mockVultrPlans{nil}}
	planMetal := NewPlanOptions(client)

	expectedMetal := []govultr.BareMetalPlan{
		{
			ID:          "vbm-4c-32gb",
			CPUCount:    4,
			CPUModel:    "E3-1270v6",
			CPUThreads:  8,
			RAM:         8,
			Disk:        20,
			DiskCount:   1,
			Bandwidth:   10,
			MonthlyCost: 120,
			Type:        "NVMe",
			Locations:   []string{"ewr"},
		},
	}
	expectedMeta := &govultr.Meta{
		Total: 1,
		Links: nil,
	}

	plan, meta, _ := planMetal.MetalList()

	if !reflect.DeepEqual(expectedMetal, plan) {
		t.Errorf("PlanOptions.list returned %v expected %v", plan, expectedMetal)
	}

	if !reflect.DeepEqual(expectedMeta, meta) {
		t.Errorf("PlanOptions.metal returned %v expected %v", meta, expectedMeta)
	}
}

func TestNewCmdPlan(t *testing.T) {
	client := &govultr.Client{Plan: mockVultrPlans{nil}}

	planOptions := NewPlanOptions(client)

	ref := reflect.TypeOf(planOptions)
	if _, ok := ref.MethodByName("List"); !ok {
		t.Errorf("Missing list function")
	}

	if _, ok := ref.MethodByName("MetalList"); !ok {
		t.Errorf("Missing metal function")
	}

	if _, ok := ref.MethodByName("validate"); ok {
		t.Errorf("validate isn't exported shouldn't be accessible")
	}

	pInterface := reflect.TypeOf(new(PlanOptionsInterface)).Elem()
	if !ref.Implements(pInterface){
		t.Errorf("PlanOptions does not implement PlanOptionsInterface")
	}
}


func TestNewPlanOptions(t *testing.T) {
	client := &govultr.Client{Plan: mockVultrPlans{nil}}
	cmd := NewCmdPlan(client)


	if cmd.Short != "get information about Vultr plans" {
		t.Errorf("invalid short")
	}

	if cmd.Use != "plans" {
		t.Errorf("invalid plans")
	}

	alias := []string{"p", "plan"}
	if !reflect.DeepEqual(cmd.Aliases, alias) {
		t.Errorf("expected alias %v got %v", alias, cmd.Aliases)
	}

}