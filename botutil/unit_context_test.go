package botutil_test

import (
	"testing"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
	"github.com/chippydip/go-sc2ai/enums/terran"
)

type groupingInfo struct {
	mockAgentInfo
	units []*api.Unit
}

func (i *groupingInfo) Data() *api.ResponseData {
	data := make([]*api.UnitTypeData, terran.Barracks+1)
	data[terran.Barracks] = &api.UnitTypeData{Attributes: []api.Attribute{api.Attribute_Structure}}
	return &api.ResponseData{Units: data}
}

func (i *groupingInfo) Observation() *api.ResponseObservation {
	return &api.ResponseObservation{
		Observation: &api.Observation{
			RawData: &api.ObservationRaw{Units: i.units},
		},
	}
}

func (i *groupingInfo) Query(query api.RequestQuery) *api.ResponseQuery {
	abilities := make([]*api.ResponseQueryAvailableAbilities, len(query.Abilities))
	for j, q := range query.Abilities {
		abilities[j] = &api.ResponseQueryAvailableAbilities{UnitTag: q.UnitTag}
	}
	return &api.ResponseQuery{Abilities: abilities}
}

func TestUnitContextMixedFlyingSameType(t *testing.T) {
	for _, alliance := range []api.Alliance{api.Alliance_Self, api.Alliance_Enemy} {
		info := &groupingInfo{units: []*api.Unit{
			{Tag: 1, UnitType: terran.Barracks, Alliance: alliance},
			{Tag: 2, UnitType: terran.Barracks, Alliance: alliance, IsFlying: true},
			{Tag: 3, UnitType: terran.Barracks, Alliance: alliance},
		}}
		ctx := botutil.NewUnitContext(info, nil)

		units := ctx.Self[terran.Barracks]
		if alliance == api.Alliance_Enemy {
			units = ctx.Enemy[terran.Barracks]
		}

		if units.Len() != 3 {
			t.Fatalf("%v: got %d Barracks, want 3", alliance, units.Len())
		}
		seen := map[api.UnitTag]bool{}
		units.Each(func(u botutil.Unit) {
			if u.IsNil() {
				t.Errorf("%v: nil unit in Barracks list", alliance)
				return
			}
			seen[u.Tag] = true
		})
		for tag := api.UnitTag(1); tag <= 3; tag++ {
			if !seen[tag] {
				t.Errorf("%v: Barracks %d missing", alliance, tag)
			}
		}
	}
}
