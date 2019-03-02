package botutil

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/buff"
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

// IsTownHall returns true if the unit is a Nexus/CC/OC/PF/Hatch/Lair/Hive.
func (u Unit) IsTownHall() bool {
	switch u.UnitType {
	case unit.Protoss_Nexus,
		unit.Terran_CommandCenter,
		unit.Terran_OrbitalCommand,
		unit.Terran_PlanetaryFortress,
		unit.Zerg_Hatchery,
		unit.Zerg_Lair,
		unit.Zerg_Hive:
		return true
	}
	return false
}

// IsGasBuilding returns true if the unit is an Assimilator/Refinery/Extractory.
func (u Unit) IsGasBuilding() bool {
	switch u.UnitType {
	case unit.Protoss_Assimilator,
		unit.Protoss_AssimilatorRich,
		unit.Terran_Refinery,
		unit.Terran_RefineryRich,
		unit.Zerg_Extractor,
		unit.Zerg_ExtractorRich:
		return true
	}
	return false
}

// IsWorker returns true if the unit is a Probe/SCV/Drone (but not MULE).
func (u Unit) IsWorker() bool {
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

// IsCarryingResources returns true if the unit is carrying minerals or gas.
func (u Unit) IsCarryingResources() bool {
	for _, b := range u.BuffIds {
		switch b {
		case buff.CarryMineralFieldMinerals,
			buff.CarryHighYieldMineralFieldMinerals,
			buff.CarryHarvestableVespeneGeyserGas,
			buff.CarryHarvestableVespeneGeyserGasProtoss,
			buff.CarryHarvestableVespeneGeyserGasZerg:
			return true
		}
	}
	return false
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

// GroundWeaponDamage returns damage per shot the unit can do to ground targets.
func (u Unit) GroundWeaponDamage() float32 {
	return u.weaponDamage(api.Weapon_Ground)
}

// AirWeaponDamage returns damage per shot the unit can do to air targets.
func (u Unit) AirWeaponDamage() float32 {
	return u.weaponDamage(api.Weapon_Air)
}

// WeaponDamage returns damage per shot the unit can do to the given target.
func (u Unit) WeaponDamage(target Unit) float32 {
	if target.IsFlying {
		return u.weaponDamage(api.Weapon_Air)
	}
	return u.weaponDamage(api.Weapon_Ground)
}

func (u Unit) weaponDamage(weaponType api.Weapon_TargetType) float32 {
	maxDamage := float32(0)
	for _, weapon := range u.Weapons {
		if weapon.Type == weaponType || weapon.Type == api.Weapon_Any {
			if weapon.Damage > maxDamage {
				maxDamage = weapon.Damage
			}
		}
	}
	return maxDamage
}

// WeaponRange returns the maximum range to attack the target from. If
// the result is negative the target cannot be attacked.
func (u Unit) WeaponRange(target Unit) float32 {
	if target.IsNil() {
		return -1
	}

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

// IsInWeaponsRange returns true if the unit is within weapons range of the target.
func (u Unit) IsInWeaponsRange(target Unit, gap float32) bool {
	maxRange := u.WeaponRange(target)
	if maxRange < 0 {
		return false
	}

	dist := float32(u.Pos2D().Distance(target.Pos2D()))
	return dist-gap <= maxRange+u.Radius+target.Radius
}

// AttackTarget issues an attack order if the unit isn't already attacking the target.
func (u Unit) AttackTarget(target Unit) {
	if !u.IsIdle() {
		if ability.Remap(u.Orders[0].AbilityId) == ability.Attack &&
			u.Orders[0].GetTargetUnitTag() == target.Tag {
			return
		}
	}

	u.OrderTarget(ability.Attack, target)
}

// AttackMove issues an attack order if the unit isn't already attacking within tollerance of pos.
func (u Unit) AttackMove(pos api.Point2D, tollerance float32) {
	if !u.IsIdle() {
		i := 0
		// If the first order is a targeted attack, examine the second order
		if tag := u.Orders[0].GetTargetUnitTag(); tag != 0 &&
			len(u.Orders) > 1 &&
			ability.Remap(u.Orders[0].AbilityId) == ability.Attack {
			i++
		}
		// If the non-specific order is an attack close enough to pos just use that
		if p := u.Orders[i].GetTargetWorldSpacePos(); p != nil &&
			ability.Remap(u.Orders[i].AbilityId) == ability.Attack &&
			p.ToPoint2D().Distance2(pos) <= tollerance*tollerance {
			return
		}

	}

	u.OrderPos(ability.Attack, pos)
}

// MoveTo issues a move order if the unit isn't already moving within tollerance of pos.
func (u Unit) MoveTo(pos api.Point2D, tollerance float32) {
	if !u.IsIdle() {
		if p := u.Orders[0].GetTargetWorldSpacePos(); p != nil &&
			ability.Remap(u.Orders[0].AbilityId) == ability.Move &&
			p.ToPoint2D().Distance2(pos) <= tollerance*tollerance {
			return
		}
	}

	u.OrderPos(ability.Move, pos)
}
