package botutil_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
	"github.com/chippydip/go-sc2ai/enums/neutral"
	"github.com/chippydip/go-sc2ai/enums/zerg"
)

var benchUnits []*api.Unit
var unitData = make([]*api.UnitTypeData, 1000)

func init() {
	// Add the type data that we need
	unitData[zerg.Drone] = &api.UnitTypeData{Weapons: []*api.Weapon{nil}}
	unitData[zerg.Larva] = &api.UnitTypeData{}
	unitData[zerg.Overlord] = &api.UnitTypeData{}
	unitData[zerg.Hatchery] = &api.UnitTypeData{Attributes: []api.Attribute{api.Attribute_Structure}}
	unitData[zerg.Zergling] = &api.UnitTypeData{Weapons: []*api.Weapon{nil}}

	unitData[neutral.MineralField] = &api.UnitTypeData{HasMinerals: true, Attributes: []api.Attribute{api.Attribute_Structure}}
	unitData[neutral.MineralField750] = &api.UnitTypeData{HasMinerals: true, Attributes: []api.Attribute{api.Attribute_Structure}}
	unitData[neutral.VespeneGeyser] = &api.UnitTypeData{HasVespene: true, Attributes: []api.Attribute{api.Attribute_Structure}}

	// Build the unit list
	for i := 0; i < 12; i++ {
		benchUnits = append(benchUnits, &api.Unit{UnitType: zerg.Drone, Alliance: api.Alliance_Self})
	}
	for i := 0; i < 3; i++ {
		benchUnits = append(benchUnits, &api.Unit{UnitType: zerg.Larva, Alliance: api.Alliance_Self})
	}
	for i := 0; i < 6; i++ {
		benchUnits = append(benchUnits, &api.Unit{UnitType: zerg.Overlord, Alliance: api.Alliance_Self, IsFlying: true})
	}
	for i := 0; i < 2; i++ {
		benchUnits = append(benchUnits, &api.Unit{UnitType: zerg.Hatchery, Alliance: api.Alliance_Self})
	}
	for i := 0; i < 100; i++ {
		benchUnits = append(benchUnits, &api.Unit{UnitType: zerg.Zergling, Alliance: api.Alliance_Self})
	}

	for i := 0; i < 100; i++ {
		benchUnits = append(benchUnits, &api.Unit{UnitType: zerg.Zergling, Alliance: api.Alliance_Enemy})
	}

	for i := 0; i < 14; i++ {
		for j := 0; j < 4; j++ {
			benchUnits = append(benchUnits, &api.Unit{UnitType: neutral.MineralField, Alliance: api.Alliance_Neutral})
			benchUnits = append(benchUnits, &api.Unit{UnitType: neutral.MineralField750, Alliance: api.Alliance_Neutral})
		}
		for j := 0; j < 2; j++ {
			benchUnits = append(benchUnits, &api.Unit{UnitType: neutral.VespeneGeyser, Alliance: api.Alliance_Neutral})
		}
	}

	// Shuffle the order to make sure sorting is fair
	rand.Shuffle(len(benchUnits), func(i, j int) {
		benchUnits[i], benchUnits[j] = benchUnits[j], benchUnits[i]
	})
}

var data = &api.ResponseData{Units: unitData}
var obs = &api.ResponseObservation{
	Observation: &api.Observation{
		RawData: &api.ObservationRaw{},
	},
}
var step func()

type info struct {
	mockAgentInfo
}

func (i *info) Data() *api.ResponseData               { return data }
func (i *info) Observation() *api.ResponseObservation { return obs }
func (i *info) OnAfterStep(f func())                  { step = f }

func BenchmarkCopyOnly(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Copy so we don't modify the input list
		raw := make([]*api.Unit, len(benchUnits))
		copy(raw, benchUnits)
		obs.Observation.RawData.Units = raw
	}
}

func BenchmarkUnitsByStep(b *testing.B) {
	// Copy so we don't modify the input list
	raw := make([]*api.Unit, len(benchUnits))
	copy(raw, benchUnits)
	obs.Observation.RawData.Units = raw

	botutil.NewUnitContext(&info{}, nil)

	for i := 0; i < b.N; i++ {
		rand.Shuffle(len(raw), func(i, j int) {
			raw[i], raw[j] = raw[j], raw[i]
		})

		step()
	}
}

func BenchmarkUnitsByTypesStep(b *testing.B) {
	// Copy so we don't modify the input list
	raw := make([]*api.Unit, len(benchUnits))
	copy(raw, benchUnits)
	obs.Observation.RawData.Units = raw

	parseUnits()

	for i := 0; i < b.N; i++ {
		rand.Shuffle(len(raw), func(i, j int) {
			raw[i], raw[j] = raw[j], raw[i]
		})

		parseUnits()
	}
}

type apiUnitsByTypes map[api.UnitTypeID][]*api.Unit

type myUnit struct{ api.Unit }
type myUnitsByTypes map[api.UnitTypeID][]*myUnit

type unitsByTypes apiUnitsByTypes

//type unitsByTypes myUnitsByTypes

func parseUnits() {
	selfUnits := unitsByTypes{}
	mineralFields := unitsByTypes{}
	vespeneGeysers := unitsByTypes{}
	neutralUnits := unitsByTypes{}
	enemyUnits := unitsByTypes{}
	for _, unit := range (&info{}).Observation().Observation.RawData.Units {
		var units *unitsByTypes
		switch unit.Alliance {
		case api.Alliance_Self:
			units = &selfUnits
		case api.Alliance_Enemy:
			units = &enemyUnits
		case api.Alliance_Neutral:
			if unit.MineralContents > 0 {
				units = &mineralFields
			} else if unit.VespeneContents > 0 {
				units = &vespeneGeysers
			} else {
				units = &neutralUnits
			}
		default:
			fmt.Fprintln(os.Stderr, "Not supported alliance: ", unit)
			continue
		}
		//units.AddFromApi(unit.UnitType, unit)
		(*units)[unit.UnitType] = append((*units)[unit.UnitType], unit)
		//(*units)[unit.UnitType] = append((*units)[unit.UnitType], &myUnit{*unit})
	}
}
