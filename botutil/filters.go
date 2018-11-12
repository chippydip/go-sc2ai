package botutil

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/enums/unit"
)

// IsType ...
func IsType(t api.UnitTypeID) func(Unit) bool {
	return func(u Unit) bool {
		return u.UnitType == t
	}
}

// IsSelfType ...
func IsSelfType(t api.UnitTypeID) func(Unit) bool {
	return func(u Unit) bool {
		return u.Alliance == api.Alliance_Self && u.UnitType == t
	}
}

// IsMineral ...
func IsMineral(u Unit) bool {
	return u.HasMinerals
}

// IsGeyser ...
func IsGeyser(u Unit) bool {
	return u.HasVespene
}

func isUnitTypeInMap(m map[api.UnitTypeID]struct{}) func(Unit) bool {
	return func(u Unit) bool {
		_, ok := m[u.UnitType]
		return ok
	}
}

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

func HasAttribute(attribute api.Attribute) func(u Unit) bool {
	return func(u Unit) bool {
		for _, attr := range u.Attributes {
			if attr == attribute {
				return true
			}
		}
		return false
	}
}

var IsStructure = HasAttribute(api.Attribute_Structure)
