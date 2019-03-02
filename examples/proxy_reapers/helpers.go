package main

import (
	"math"

	"github.com/chippydip/go-sc2ai/api"
)

func closestToPos(points []api.Point2D, pos api.Point2D) api.Point2D {
	minDist := float32(math.Inf(1))
	var closest api.Point2D
	for _, p := range points {
		dist := pos.Distance2(p)
		if dist < minDist {
			closest = p
			minDist = dist
		}
	}
	return closest
}
