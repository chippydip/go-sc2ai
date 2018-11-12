package botutil

import "github.com/chippydip/go-sc2ai/api"

type Unit struct {
	*api.Unit
	*api.UnitTypeData
}

func (u *Unit) IsNil() bool {
	return u.Unit == nil
}

// UnitHasAttribute ...
// func (a *Agent) UnitHasAttribute(u *api.Unit, attr api.Attribute) bool {
// 	for _, a := range a.info.Data().Units[u.UnitType].Attributes {
// 		if a == attr {
// 			return true
// 		}
// 	}
// 	return false
// }

func (u *Unit) IsSelf() bool {
	return u.Alliance == api.Alliance_Self
}

func (u *Unit) IsEnemy() bool {
	return u.Alliance == api.Alliance_Enemy
}

func (u *Unit) IsType(unitType api.UnitTypeID) bool {
	return u.UnitType == unitType
}

func (u *Unit) IsAnyType(types ...api.UnitTypeID) bool {
	for _, unitType := range types {
		if u.UnitType == unitType {
			return true
		}
	}
	return false
}
