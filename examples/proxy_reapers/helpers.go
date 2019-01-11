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

func neighbours4(p api.Point2D, offset float32) []api.Point2D {
	return []api.Point2D{
		{X: p.X, Y: p.Y + offset},
		{X: p.X + offset, Y: p.Y},
		{X: p.X, Y: p.Y - offset},
		{X: p.X - offset, Y: p.Y},
	}
}

func neighbours8(p api.Point2D, offset float32) []api.Point2D {
	return append(neighbours4(p, offset), []api.Point2D{
		{X: p.X + offset, Y: p.Y + offset},
		{X: p.X + offset, Y: p.Y - offset},
		{X: p.X - offset, Y: p.Y - offset},
		{X: p.X - offset, Y: p.Y + offset},
	}...)
}
