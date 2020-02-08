package search

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/unit"
)

var sizeCache = map[api.UnitTypeID]api.Size2DI{}

// UnitPlacementSize estimates building footprints based on unit radius.
func UnitPlacementSize(u botutil.Unit) api.Size2DI {
	if s, ok := sizeCache[u.UnitType]; ok {
		return s
	}

	// Round coordinate to the nearest half (not needed except for things like the KD8Charge)
	x, y := float32(int32(u.Pos.X*2+0.5))/2, float32(int32(u.Pos.Y*2+0.5))/2
	xEven, yEven := int(u.Pos.X*2+0.5)%2 == 0, int(u.Pos.Y*2+0.5)%2 == 0

	// Compute bounds based on the (bad) radius provided by the game
	xMin, yMin := int32(x-u.Radius+0.5), int32(y-u.Radius+0.5)
	xMax, yMax := int32(x+u.Radius+0.5), int32(y+u.Radius+0.5)

	// Get the real radius in all four directions as calculated above
	rxMin, ryMin := x-float32(xMin), y-float32(yMin)
	rxMax, ryMax := float32(xMax)-x, float32(yMax)-y

	// If the radii are not symetric, take the smaller value
	rx, ry := rxMin, ryMin
	if rxMax < rx {
		rx = rxMax
	}
	if ryMax < ry {
		ry = ryMax
	}

	// Re-compute bounds with the hopefully better radii
	xMin, yMin = int32(u.Pos.X-rx+0.5), int32(u.Pos.Y-ry+0.5)
	xMax, yMax = int32(u.Pos.X+rx+0.5), int32(u.Pos.Y+ry+0.5)

	// Adjust for non-square structures (TODO: should this just special-case Minerals?)
	if xEven != yEven {
		if yEven {
			xMin++
			xMax--
		} else {
			yMin++
			yMax--
		}
	}

	// Cache and return the computed size
	size := api.Size2DI{X: xMax - xMin, Y: yMax - yMin}
	sizeCache[u.UnitType] = size
	log.Printf("%v %v %v -> %v", unit.String(u.UnitType), u.Pos2D(), u.Radius, size)
	return size
}

// PlacementGrid ...
type PlacementGrid struct {
	raw        api.ImageDataBits
	grid       api.ImageDataBits
	structures map[api.UnitTag]structureInfo
}

type structureInfo struct {
	point api.Point2D
	size  api.Size2DI
}

// NewPlacementGrid ...
func NewPlacementGrid(bot *botutil.Bot) *PlacementGrid {
	raw := bot.GameInfo().GetStartRaw().GetPlacementGrid().Bits()
	pg := &PlacementGrid{
		raw:        raw,
		grid:       raw.Copy(),
		structures: map[api.UnitTag]structureInfo{},
	}

	update := func() {
		// Remove any units that are gone or have changed type or position
		for k, v := range pg.structures {
			if u := bot.UnitByTag(k); u.IsNil() || !u.IsStructure() || u.Pos2D() != v.point || UnitPlacementSize(u) != v.size {
				pg.markGrid(v.point, v.size, true)
				delete(pg.structures, k)
			}
		}

		// (Re-)add new units or ones that have changed type or position
		bot.AllUnits().Each(func(u botutil.Unit) {
			if _, ok := pg.structures[u.Tag]; !ok && u.IsStructure() {
				v := structureInfo{u.Pos2D(), UnitPlacementSize(u)}
				pg.markGrid(v.point, v.size, false)
				pg.structures[u.Tag] = v
			}
		})
	}

	bot.OnAfterStep(update)
	update()

	var req []*api.RequestQueryBuildingPlacement
	var lut []api.UnitTypeID

	for k, v := range pg.structures {
		xMin, yMin := int32(v.point.X-float32(v.size.X)/2), int32(v.point.Y-float32(v.size.Y)/2)
		xMax, yMax := xMin+v.size.X, yMin+v.size.Y
		for y := yMin; y < yMax; y++ {
			for x := xMin; x < xMax; x++ {
				if pg.grid.Get(x, y) != raw.Get(x, y) {
					req = append(req, &api.RequestQueryBuildingPlacement{
						AbilityId: ability.Build_SensorTower,
						TargetPos: &api.Point2D{X: float32(x) + 0.5, Y: float32(y) + 0.5},
					})
					lut = append(lut, bot.UnitByTag(k).UnitType)
				}
			}
		}
	}

	resp := bot.Query(api.RequestQuery{
		Placements: req,
	})

	heightMap := NewHeightMap(bot.GameInfo().StartRaw)
	var ok, inval = 0, 0
	for i, r := range resp.GetPlacements() {
		var color *api.Color
		if r.GetResult() == api.ActionResult_Success {
			color = green
			inval++
		} else {
			color = red
			ok++
		}
		v := req[i].TargetPos
		z := heightMap.Interpolate(v.X, v.Y)
		queryBoxes = append(queryBoxes, &api.DebugBox{
			Color: color,
			Min:   &api.Point{X: v.X - 0.375, Y: v.Y - 0.375, Z: z},
			Max:   &api.Point{X: v.X + 0.375, Y: v.Y + 0.375, Z: z + 1},
		})
	}
	log.Printf("ok: %v, inval: %v", ok, inval)

	return pg
}

var queryBoxes []*api.DebugBox

func (pg *PlacementGrid) markGrid(pos api.Point2D, size api.Size2DI, value bool) {
	xMin, yMin := int32(pos.X-float32(size.X)/2), int32(pos.Y-float32(size.Y)/2)
	xMax, yMax := xMin+size.X, yMin+size.Y

	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			pg.grid.Set(x, y, value)
		}
	}
}

func (pg *PlacementGrid) checkGrid(pos api.Point2D, size api.Size2DI, value bool) bool {
	xMin, yMin := int32(pos.X-float32(size.X)/2), int32(pos.Y-float32(size.Y)/2)
	xMax, yMax := xMin+size.X, yMin+size.Y

	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			if pg.grid.Get(x, y) != value {
				return false
			}
		}
	}
	return true
}

// CanPlace checks if a structure of a certain type can currently be places at the given location.
func (pg *PlacementGrid) CanPlace(u botutil.Unit, pos api.Point2D) bool {
	return pg.checkGrid(pos, UnitPlacementSize(u), true)
}

// DebugPrint ...
func (pg *PlacementGrid) DebugPrint(bot *botutil.Bot) {
	heightMap := NewHeightMap(bot.GameInfo().StartRaw)

	var boxes []*api.DebugBox
	for _, v := range pg.structures {
		z := heightMap.Interpolate(v.point.X, v.point.Y)
		boxes = append(boxes, &api.DebugBox{
			Min: &api.Point{X: v.point.X - float32(v.size.X)/2, Y: v.point.Y - float32(v.size.Y)/2, Z: z},
			Max: &api.Point{X: v.point.X + float32(v.size.X)/2, Y: v.point.Y + float32(v.size.Y)/2, Z: z + 1},
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
			Command: &api.DebugCommand_Draw{
				Draw: &api.DebugDraw{
					Boxes: queryBoxes,
				},
			},
		},
	})
}
