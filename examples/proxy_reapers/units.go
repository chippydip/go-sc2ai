package main

import (
	"github.com/chippydip/go-sc2ai/api"
)

type Units []*Unit

func (u *Units) Add(unit *Unit) {
	if unit != nil {
		*u = append(*u, unit)
	}
}

func (u *Units) AddFromApi(unit *api.Unit) {
	if unit != nil {
		*u = append(*u, NewUnit(unit))
	}
}

func (u Units) Len() int {
	return len(u)
}

// While using each method you should check that it returns not nil
func (u Units) First() *Unit {
	if len(u) == 0 {
		return nil
	}
	return u[0]
}

func (u Units) FirstFiltered(filter func(*Unit) bool) *Unit {
	for _, unit := range u {
		if filter(unit) {
			return unit
		}
	}
	return nil
}

func (u Units) ByTag(tag api.UnitTag) *Unit {
	for _, unit := range u {
		if unit.Tag == tag {
			return unit
		}
	}
	return nil
}

func (u Units) ClosestTo(p api.Point2D) *Unit {
	var closest *Unit
	for _, unit := range u {
		if closest == nil || p.Distance2(closest.Pos.ToPoint2D()) > p.Distance2(unit.Pos.ToPoint2D()) {
			closest = unit
		}
	}
	return closest
}

func (u Units) CloserThan(dist float64, pos api.Point2D) Units {
	dist2 := float32(dist * dist)
	units := Units{}
	for _, unit := range u {
		if unit.Pos.ToPoint2D().Distance2(pos) <= dist2 {
			units.Add(unit)
		}
	}
	return units
}

func (u Units) Center() api.Point2D {
	points := []api.Point2D{}
	for _, unit := range u {
		points = append(points, unit.Pos.ToPoint2D())
	}
	return CenterP2D(points)
}
