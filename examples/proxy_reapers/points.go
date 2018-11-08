package main

import (
	"github.com/chippydip/go-sc2ai/api"
	"math"
)

func FloorP2D(p api.Point2D) api.Point2D {
	return api.Point2D{X: float32(int(p.X)), Y: float32(int(p.Y))}
}

func CenterP2D(ps []api.Point2D) api.Point2D {
	if len(ps) == 0 {
		return api.Point2D{}
	}
	sum := api.Vec2D{}
	for _, p := range ps {
		sum = sum.Add(api.Vec2D(p))
	}
	return api.Point2D(sum.Mul(1 / float32(len(ps))))
}

func ClosestToP2D(pts []api.Point2D, p api.Point2D) api.Point2D {
	var closest api.Point2D
	for _, pt := range pts {
		if closest.X == 0 && closest.Y == 0 || p.Distance2(closest) > p.Distance2(pt) {
			closest = pt
		}
	}
	return closest
}

func SquareDirection(v api.Vec2D) api.Vec2D {
	return api.Vec2D{X: float32(math.Copysign(1, float64(v.X))), Y: float32(math.Copysign(1, float64(v.Y)))}
}

func Neighbours4(p api.Point2D, offset float32) []api.Point2D {
	return []api.Point2D{{p.X, p.Y + offset}, {p.X + offset, p.Y}, {p.X, p.Y - offset}, {p.X - offset, p.Y}}
}

func Neighbours8(p api.Point2D, offset float32) []api.Point2D {
	return append(Neighbours4(p, offset), []api.Point2D{{p.X + offset, p.Y + offset}, {p.X + offset, p.Y - offset},
		{p.X - offset, p.Y - offset}, {p.X - offset, p.Y + offset}}...)
}

func (bot *proxyReapers) heightAt(p api.Point2D) byte {
	m := bot.info.GameInfo().StartRaw.TerrainHeight
	// m.BitsPerPixel == 8

	addr := int(p.X) + int(p.Y)*int(m.Size_.X)
	if addr > len(m.Data)-1 {
		return 0
	}
	return m.Data[addr]
}

func (bot *proxyReapers) isBuildable(p api.Point2D) bool {
	m := bot.info.GameInfo().StartRaw.PlacementGrid

	addr := int(p.X) + int(p.Y)*int(m.Size_.X)
	if addr > len(m.Data)-1 || p.X < 0 || p.Y < 0 {
		return false
	}
	return m.Data[addr] != 0
}

func (bot *proxyReapers) is3x3buildable(pos api.Point2D) bool {
	if !bot.isBuildable(pos) {
		return false
	}
	for _, p := range Neighbours8(pos, 1) {
		if !bot.isBuildable(p) {
			return false
		}
	}
	return true
}
