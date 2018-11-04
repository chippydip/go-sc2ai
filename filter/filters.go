package filter

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/enums/unit"
)

// Filter ...
type Filter func(*api.Unit) bool

// IsType ...
func IsType(t api.UnitTypeID) Filter {
	return func(u *api.Unit) bool {
		return u.UnitType == t
	}
}

// IsSelfType ...
func IsSelfType(t api.UnitTypeID) Filter {
	return func(u *api.Unit) bool {
		return u.Alliance == api.Alliance_Self && u.UnitType == t
	}
}

func isUnitTypeInMap(m map[api.UnitTypeID]struct{}) Filter {
	return func(u *api.Unit) bool {
		_, ok := m[u.UnitType]
		return ok
	}
}

// IsMineral ...
var IsMineral = isUnitTypeInMap(map[api.UnitTypeID]struct{}{
	unit.Neutral_MineralField:                 {},
	unit.Neutral_MineralField750:              {},
	unit.Neutral_BattleStationMineralField:    {},
	unit.Neutral_BattleStationMineralField750: {},
	unit.Neutral_LabMineralField:              {},
	unit.Neutral_LabMineralField750:           {},
	unit.Neutral_PurifierRichMineralField:     {},
	unit.Neutral_PurifierRichMineralField750:  {},
	unit.Neutral_PurifierMineralField:         {},
	unit.Neutral_PurifierMineralField750:      {},
	unit.Neutral_RichMineralField:             {},
	unit.Neutral_RichMineralField750:          {},
})

// IsGeyser ...
var IsGeyser = isUnitTypeInMap(map[api.UnitTypeID]struct{}{
	unit.Neutral_ProtossVespeneGeyser:  {},
	unit.Neutral_PurifierVespeneGeyser: {},
	unit.Neutral_RichVespeneGeyser:     {},
	unit.Neutral_ShakurasVespeneGeyser: {},
	unit.Neutral_SpacePlatformGeyser:   {},
	unit.Neutral_VespeneGeyser:         {},
})

// IsTownHall ...
var IsTownHall = isUnitTypeInMap(map[api.UnitTypeID]struct{}{
	unit.Protoss_Nexus:               {},
	unit.Terran_CommandCenter:        {},
	unit.Terran_CommandCenterFlying:  {},
	unit.Terran_OrbitalCommand:       {},
	unit.Terran_OrbitalCommandFlying: {},
	unit.Zerg_Hatchery:               {},
	unit.Zerg_Lair:                   {},
	unit.Zerg_Hive:                   {},
})

// IsGasBuilding ...
var IsGasBuilding = isUnitTypeInMap(map[api.UnitTypeID]struct{}{
	unit.Protoss_Assimilator: {},
	unit.Terran_Refinery:     {},
	unit.Zerg_Extractor:      {},
})

// IsWorker ...
var IsWorker = isUnitTypeInMap(map[api.UnitTypeID]struct{}{
	unit.Protoss_Probe:      {},
	unit.Terran_SCV:         {},
	unit.Terran_MULE:        {},
	unit.Zerg_Drone:         {},
	unit.Zerg_DroneBurrowed: {},
})
