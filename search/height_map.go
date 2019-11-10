package search

import (
	"github.com/chippydip/go-sc2ai/api"
)

// HeightMap ...
type HeightMap struct {
	data api.ImageDataBytes
}

// NewHeightMap ...
func NewHeightMap(start *api.StartRaw) HeightMap {
	return HeightMap{start.GetTerrainHeight().Bytes()}
}

// Width ...
func (hm HeightMap) Width() int32 {
	return hm.data.Width() + 1
}

// Height ...
func (hm HeightMap) Height() int32 {
	return hm.data.Height() + 1
}

// InBounds ...
func (hm HeightMap) InBounds(x, y int32) bool {
	return 0 <= x && x < hm.Width() && 0 <= y && y < hm.Height()
}

// Get ...
func (hm HeightMap) Get(x, y int32) float32 {
	return hm.decode(x, y)
}

// Interpolate returns a bilinear interpolated height value for any fractional map coordinates.
func (hm HeightMap) Interpolate(x, y float32) float32 {
	x0, y0 := int32(x), int32(y)
	f00 := hm.decode(x0, y0)
	f01 := hm.decode(x0, y0+1)
	f10 := hm.decode(x0+1, y0)
	f11 := hm.decode(x0+1, y0+1)

	x, y = x-float32(x0), y-float32(y0)
	return f00*(1-x)*(1-y) + f01*(1-x)*y + f10*x*(1-y) + f11*x*y
}

func (hm HeightMap) decode(x, y int32) float32 {
	if hm.data.InBounds(x, y) {
		return (float32(hm.data.Get(x, y)) - 127) / 8
	}
	return float32(-127) / 8
}
