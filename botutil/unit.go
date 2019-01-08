package botutil

import "github.com/chippydip/go-sc2ai/api"

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
