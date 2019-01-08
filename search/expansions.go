package search

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
	"github.com/chippydip/go-sc2ai/client"
)

// CalculateExpansionLocations groups resources into clusters and determines the best town hall location for each cluster.
// The Center() point of each cluster is the optimal town hall location. If debug is true then the results will also
// be visualized in-game (until new debug info is drawn).
func CalculateExpansionLocations(bot *botutil.Bot, debug bool) []UnitCluster {
	// Start by finding resource clusters
	clusters := Cluster(bot.Neutral.Resources(), 15)

	// Add resource-restrictions to the placement grid
	placement := bot.GameInfo().StartRaw.PlacementGrid.Copy().Bytes()
	bot.Neutral.Minerals().Each(func(u botutil.Unit) {
		markUnbuildable(placement, int(u.Pos.X-0.5), int(u.Pos.Y), 2, 1)
	})
	bot.Neutral.Vespene().Each(func(u botutil.Unit) {
		markUnbuildable(placement, int(u.Pos.X-1), int(u.Pos.Y-1), 3, 3)
	})

	// Mark locations which *can't* have an expansion centers
	for y := 0; y < placement.Height(); y++ {
		for x := 0; x < placement.Width(); x++ {
			if placement.Get(x, y) < 128 {
				expandUnbuildable(placement, x, y)
			}
		}
	}

	// Find the nearest remaining square to each cluster's CoM
	for i, cluster := range clusters {
		pt := cluster.Center()
		px, py := int(pt.X), int(pt.Y)
		r2Min, xBest, yBest := 256, -1, -1
		for r := 0; r*r <= r2Min; r++ { // search radius
			xMin, xMax, yMin, yMax := px-r, px+r, py-r, py+r
			for y := yMin; y <= yMax; y++ {
				for x := xMin; x <= xMax; x++ {
					// This is slightly inefficient, but much easier than repeating the same loop 4x for the edges
					if (x == xMin || x == xMax || y == yMin || y == yMax) && placement.Get(x, y) == 255 {
						dx, dy := x-px, y-py
						r2 := dx*dx + dy*dy
						if r2 < r2Min {
							r2Min = r2
							xBest = x
							yBest = y
						}
					}
				}
			}
		}

		// Update the Center to be the detected location rather than the actual CoM (just don't add new units)
		clusters[i].sum = api.Vec2D{X: float32(xBest) + 0.5, Y: float32(yBest) + 0.5}.Mul(float32(cluster.units.Len()))
	}

	if debug {
		debugPrint(clusters, placement, bot)
	}

	return clusters
}

// markUnbuildable marks a w x h area around px, py (minus corners) as unbuildable (red)
func markUnbuildable(placement api.ImageDataBytes, px, py, w, h int) {
	xMin, xMax := px-3, px+w+2
	yMin, yMax := py-3, py+h+2

	for y := yMin; y <= yMax; y++ {
		for x := xMin; x <= xMax; x++ {
			if (y == yMin || y == yMax) && (x == xMin || x == xMax) {
				continue // skip corners
			}
			if placement.Get(x, y) == 255 {
				placement.Set(x, y, 1)
			}
		}
	}
}

// expandUnbuildable marks any tile within 2 units of px, py as unbuildable (blue)
func expandUnbuildable(placement api.ImageDataBytes, px, py int) {
	xMin, xMax := px-2, px+2
	yMin, yMax := py-2, py+2

	for y := yMin; y <= yMax; y++ {
		for x := xMin; x <= xMax; x++ {
			if placement.Get(x, y) == 255 {
				placement.Set(x, y, 128)
			}
		}
	}
}

// debugPrint shows debug info about the expansion search procedure in-game
func debugPrint(clusters []UnitCluster, placement api.ImageDataBytes, bot client.AgentInfo) {
	info := bot.GameInfo()
	heightMap := info.StartRaw.TerrainHeight.Bytes()
	pathable := info.StartRaw.PathingGrid.Bytes()

	var boxes []*api.DebugBox

	// Debug placement grid
	for y := 0; y < placement.Height(); y++ {
		for x := 0; x < placement.Width(); x++ {
			color := mapColor(placement.Get(x, y), pathable.Get(x, y))
			if color != nil {
				//z := float32(int(0.75*(float32(heightMap.Get(x, y))-127)+0.5)) + 0.01
				z := (float32(heightMap.Get(x, y))/254)*200 - 100
				boxes = append(boxes, &api.DebugBox{
					Color: color,
					Min:   &api.Point{X: float32(x) + 0.25, Y: float32(y) + 0.25, Z: z},
					Max:   &api.Point{X: float32(x) + 0.75, Y: float32(y) + 0.75, Z: z},
				})
			}
		}
	}

	// Expansion locations
	for _, cluster := range clusters {
		pt := cluster.Center()
		z := (float32(heightMap.Get(int(pt.X), int(pt.Y)))/254)*200 - 100
		boxes = append(boxes, &api.DebugBox{
			Color: green,
			Min:   &api.Point{X: pt.X - 2.5, Y: pt.Y - 2.5, Z: z},
			Max:   &api.Point{X: pt.X + 2.5, Y: pt.Y + 2.5, Z: z},
		}, &api.DebugBox{
			Color: green,
			Min:   &api.Point{X: pt.X - 0.05, Y: pt.Y - 0.05, Z: z},
			Max:   &api.Point{X: pt.X + 0.05, Y: pt.Y + 0.05, Z: z},
		})

	}

	bot.SendDebugCommands([]*api.DebugCommand{
		&api.DebugCommand{
			Command: &api.DebugCommand_Draw{
				Draw: &api.DebugDraw{
					Boxes: boxes,
				},
			},
		},
		&api.DebugCommand{
			Command: &api.DebugCommand_GameState{
				GameState: api.DebugGameState_show_map,
			},
		},
	})

}

// Re-use these colors so we don't have to keep allocating them
var (
	white = &api.Color{R: 255, G: 255, B: 255}
	red   = &api.Color{R: 255, G: 1, B: 1}
	blue  = &api.Color{R: 1, G: 1, B: 255}
	green = &api.Color{R: 1, G: 255, B: 1}
)

// mapColor converts a building placement value into a display color
func mapColor(value byte, pathable byte) *api.Color {
	switch value {
	case 255:
		return white // center buildable
	case 128:
		return blue // too close for center
	case 1:
		return red // too close to resources
	}
	if pathable == 0 {
		return green // not buildable, but pathable
	}
	return nil
}
