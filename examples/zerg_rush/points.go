package main

import (
	"github.com/chippydip/go-sc2ai/api"
)

// User should check that he receives not nil
func ClosestUnit(a api.Point2D, units []*api.Unit) *api.Unit {
	var closest *api.Unit
	for _, unit := range units {
		if closest == nil ||
			a.Sub(closest.Pos.ToPoint2D()).LenSqr() > a.Sub(unit.Pos.ToPoint2D()).LenSqr() {
			closest = unit
		}
	}
	return closest
}
