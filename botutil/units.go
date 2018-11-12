package botutil

import (
	"math"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

// Units ...
type Units struct {
	units []*api.Unit
	data  []*api.UnitTypeData
}

func NewUnits(info client.AgentInfo) *Units {
	u := &Units{}
	update := func() {
		u.units = info.Observation().GetObservation().GetRawData().GetUnits()
		u.data = info.Data().GetUnits()
	}
	update()
	info.OnAfterStep(update)
	return u
}

func (units *Units) wrap(u *api.Unit) Unit {
	if u == nil {
		return Unit{}
	}
	return Unit{u, units.data[u.UnitType]}
}

// Len ...
func (units *Units) Len() int {
	return len(units.units)
}

// Slice ...
func (units *Units) Slice() []*api.Unit {
	return units.units
}

// Append ...
func (units *Units) Append(u Unit) {
	units.units = append(units.units, u.Unit)
}

// Tags ...
func (units *Units) Tags() []api.UnitTag {
	tags := make([]api.UnitTag, len(units.units))
	for i, u := range units.units {
		tags[i] = u.Tag
	}
	return tags
}

// Each ...
func (units *Units) Each(f func(Unit)) {
	for _, u := range units.units {
		f(units.wrap(u))
	}
}

// Choose ...
func (units *Units) Choose(filter func(Unit) bool) Units {
	var result []*api.Unit
	for _, u := range units.units {
		if filter(units.wrap(u)) {
			result = append(result, u)
		}
	}
	return Units{result, units.data}
}

// Drop ...
func (units *Units) Drop(filter func(Unit) bool) Units {
	var result []*api.Unit
	for _, u := range units.units {
		if !filter(units.wrap(u)) {
			result = append(result, u)
		}
	}
	return Units{result, units.data}
}

// First returns the first unit matching the given filter from the latest observation.
func (units *Units) First(filter func(Unit) bool) Unit {
	for _, u := range units.units {
		u := units.wrap(u)
		if filter(u) {
			return u
		}
	}
	return Unit{}
}

// Closest returns the closest unit from the latest observation.
func (units *Units) Closest(pos api.Point2D) Unit {
	minDist := float32(math.Inf(1))
	var closest *api.Unit
	for _, u := range units.units {
		dist := pos.Distance2(u.Pos.ToPoint2D())
		if dist < minDist {
			closest = u
			minDist = dist
		}
	}
	return units.wrap(closest)
}

// ClosestWithFilter returns the closest unit matching the given filter from the latest observation.
func (units *Units) ClosestWithFilter(pos api.Point2D, filter func(Unit) bool) Unit {
	minDist := float32(math.Inf(1))
	var closest *api.Unit
	for _, u := range units.units {
		if !filter(units.wrap(u)) {
			continue
		}
		dist := pos.Distance2(u.Pos.ToPoint2D())
		if dist < minDist {
			closest = u
			minDist = dist
		}
	}
	return units.wrap(closest)
}

// CountSelfType ...
func (units *Units) CountSelfType(unitType api.UnitTypeID) int {
	n := 0
	for _, u := range units.units {
		if u.Alliance == api.Alliance_Self && u.UnitType == unitType {
			n++
		}
	}
	return n
}

// CountSelfTypeInProduction ...
func (units *Units) CountSelfTypeInProduction(unitType api.UnitTypeID) int {
	n, abil := 0, units.data[unitType].AbilityId
	for _, u := range units.units {
		if u.Alliance == api.Alliance_Self {
			for _, order := range u.Orders {
				if order.AbilityId == abil {
					n++
				}
			}
		}
	}
	return n
}

// CountSelfTypeAll ...
func (units *Units) CountSelfTypeAll(unitType api.UnitTypeID) int {
	return units.CountSelfType(unitType) + units.CountSelfTypeInProduction(unitType)
}
