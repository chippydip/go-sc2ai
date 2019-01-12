package botutil

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/enums/ability"
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

// IsSelf checks if this is our unit.
func (u Unit) IsSelf() bool {
	return u.Alliance == api.Alliance_Self
}

// IsEnemy checks if this is an enemy unit.
func (u Unit) IsEnemy() bool {
	return u.Alliance == api.Alliance_Enemy
}

// IsType checks if the unit is the specified type.
func (u Unit) IsType(unitType api.UnitTypeID) bool {
	return u.UnitType == unitType
}

// IsAnyType checks if the unit is one of the given types.
func (u Unit) IsAnyType(types ...api.UnitTypeID) bool {
	for _, unitType := range types {
		if u.UnitType == unitType {
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

// TODO: Generate command alias maps so this can be more general?
var isHarvestGather = map[api.AbilityID]bool{
	ability.Harvest_Gather:       true,
	ability.Harvest_Gather_Drone: true,
	ability.Harvest_Gather_Mule:  true,
	ability.Harvest_Gather_Probe: true,
	ability.Harvest_Gather_SCV:   true,
}

// IsGathering returns true if the unit is currently gathering.
func (u Unit) IsGathering() bool {
	return !u.IsIdle() && isHarvestGather[u.Orders[0].AbilityId]
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
