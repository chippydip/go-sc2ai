package botutil

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/unit"
)

// Unit combines the api Unit with it's UnitTypeData and adds some additional convenience methods.
type Unit struct {
	ctx *UnitContext
	*api.Unit
	*api.UnitTypeData
}

// IsNil checks if the underlying Unit pointer is nil.
func (u Unit) IsNil() bool {
	return u.Unit == nil
}

// IsVisible checks if DisplayType is Visible.
func (u Unit) IsVisible() bool {
	return u.DisplayType == api.DisplayType_Visible
}

// IsSnapshot checks if DisplayType is Snapshot.
func (u Unit) IsSnapshot() bool {
	return u.DisplayType == api.DisplayType_Snapshot
}

// IsHidden checks if DisplayType is Hidden.
func (u Unit) IsHidden() bool {
	return u.DisplayType == api.DisplayType_Hidden
}

// HasAttribute checks if this unit has the specified attribute.
func (u Unit) HasAttribute(attr api.Attribute) bool {
	for _, a := range u.Attributes {
		if a == attr {
			return true
		}
	}
	return false
}

// IsStructure checks if the unit is a building (has the Structure attribute).
func (u Unit) IsStructure() bool {
	return u.HasAttribute(api.Attribute_Structure)
}

// Pos2D returns the x/y location of the unit.
func (u Unit) Pos2D() api.Point2D {
	return u.Pos.ToPoint2D()
}

// IsBuilt returns true if the unit is done building.
func (u Unit) IsBuilt() bool {
	return u.BuildProgress == 1
}

// IsIdle returns true if the unit has no orders.
func (u Unit) IsIdle() bool {
	return len(u.Orders) == 0
}

// IsTownHall return true if the unit is a Nexus/CC/OC/PF/Hatch/Lair/Hive.
func (u Unit) IsTownHall() bool {
	switch u.UnitType {
	case unit.Protoss_Nexus,
		unit.Terran_CommandCenter,
		unit.Terran_OrbitalCommand,
		unit.Terran_PlanetaryFortress,
		unit.Zerg_Hatchery,
		unit.Zerg_Lair:
		return true
	}
	return false
}

// IsHarvester returns true if the unit is a SCV/Probe/Drone.
func (u Unit) IsHarvester() bool {
	switch u.UnitType {
	case unit.Protoss_Probe,
		unit.Terran_SCV,
		unit.Zerg_Drone:
		return true
	}
	return false
}

// IsGathering returns true if the unit is currently gathering.
func (u Unit) IsGathering() bool {
	return !u.IsIdle() && ability.Remap(u.Orders[0].AbilityId) == ability.Harvest_Gather
}

// HasBuff ...
func (u Unit) HasBuff(buffID api.BuffID) bool {
	for _, b := range u.BuffIds {
		if b == buffID {
			return true
		}
	}
	return false
}

// HasEnergy ...
func (u Unit) HasEnergy(energy float32) bool {
	return u.Energy >= energy
}

// WeaponRange returns the maximum range to attack the target from. If
// the result is negative the target cannot be attacked.
func (u Unit) WeaponRange(target Unit) float32 {
	weaponType := api.Weapon_Ground
	if target.IsFlying {
		weaponType = api.Weapon_Air
	}

	maxRange := float32(-1)
	for _, weapon := range u.Weapons {
		if weapon.Type == weaponType || weapon.Type == api.Weapon_Any {
			if weapon.Damage > 0 && weapon.Range > maxRange {
				maxRange = weapon.Range
			}
		}
	}
	return maxRange
}

// InRange returns true if the unit is within weapons range of the target.
func (u Unit) InRange(target Unit, gap float32) bool {
	maxRange := u.WeaponRange(target)
	if maxRange < 0 {
		return false
	}

	dist := float32(u.Pos2D().Distance(target.Pos2D()))
	return dist-gap <= maxRange+u.Radius+target.Radius
}
