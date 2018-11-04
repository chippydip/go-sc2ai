package search

import (
	"math"

	"github.com/chippydip/go-sc2ai/api"
)

// ClosestUnit ...
func ClosestUnit(pos api.Point2D, units ...*api.Unit) *api.Unit {
	minDist := float32(math.Inf(1))
	var closest *api.Unit
	for _, u := range units {
		dist := pos.Distance2(u.Pos.ToPoint2D())
		if dist < minDist {
			closest = u
			minDist = dist
		}
	}
	return closest
}

// ClosestUnitWithFilter ...
func ClosestUnitWithFilter(pos api.Point2D, filter func(*api.Unit) bool, units ...*api.Unit) *api.Unit {
	minDist := float32(math.Inf(1))
	var closest *api.Unit
	for _, u := range units {
		if !filter(u) {
			continue
		}
		dist := pos.Distance2(u.Pos.ToPoint2D())
		if dist < minDist {
			closest = u
			minDist = dist
		}
	}
	return closest
}
