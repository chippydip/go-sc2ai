package search

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
)

// ComputeOpenness ...
func ComputeOpenness(bot *botutil.Bot) api.ImageDataBytes {
	depth, _ := computeDepth(bot)

	// Convert depth to openness
	for y := int32(0); y < depth.Height(); y++ {
		for x := int32(0); x < depth.Width(); x++ {
			depth.Set(x, y, 255-depth.Get(x, y))
		}
	}

	return depth
}

func computeDepth(bot *botutil.Bot) (api.ImageDataBytes, byte) {
	placement := bot.GameInfo().StartRaw.PlacementGrid.Bits()     // false == blocked, true == buildable
	depth := bot.GameInfo().StartRaw.PathingGrid.Bits().ToBytes() // false == pathable, true == blocked

	processed := map[api.PointI]bool{}
	todo := map[api.PointI]bool{}
	curr := map[api.PointI]bool{}

	for y := int32(0); y < depth.Height(); y++ {
		for x := int32(0); x < depth.Width(); x++ {
			v := depth.Get(x, y)
			if v == 255 {
				if placement.Get(x, y) {
					depth.Set(x, y, 0)
				} else {
					todo[api.PointI{X: int32(x), Y: int32(y)}] = true
				}
			}
		}
	}

	min := byte(255)
	for len(todo) > 0 {
		curr, todo = todo, curr // swap
		for k := range todo {
			delete(todo, k)
		}

		for k := range curr {
			neighbors := k.Offset4By(1)

			if depth.Get(k.X, k.Y) != 255 {
				var max byte
				for _, n := range neighbors {
					if d := depth.Get(n.X, n.Y); d > max {
						max = d
					}
				}
				max--
				depth.Set(k.X, k.Y, max)
				if max < min {
					min = max
				}
			}
			processed[k] = true

			for _, n := range neighbors {
				if depth.InBounds(n.X, n.Y) && !processed[n] {
					todo[n] = true
				}
			}
		}
	}

	if len(processed) != int(depth.Width()*depth.Height()) {
		log.Panicf("Only process %v of %v cells", len(processed), depth.Width()*depth.Height())
	}

	return depth, min
}
