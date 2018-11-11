package api

import (
	"encoding/binary"
	"log"
)

// Copy returns an ImageData with a separate Data slice copied from the original.
func (img *ImageData) Copy() *ImageData {
	data := make([]byte, len(img.Data))
	copy(data, img.Data)

	return &ImageData{
		BitsPerPixel: img.BitsPerPixel,
		Size_:        &Size2DI{X: img.Size_.X, Y: img.Size_.Y},
		Data:         data,
	}
}

// assertBPP checks for the expected pixel size and panics if it doesn't match.
func (img *ImageData) assertBPP(count int32) {
	if img.BitsPerPixel != count {
		log.Panicf("bad BitsPerPixel, expected %v got %v", count, img.BitsPerPixel)
	}
}

// Bits returns a bit-indexed version of the ImageData.
// It panics if ImageData.BitsPerPixel != 1.
func (img *ImageData) Bits() ImageDataBits {
	img.assertBPP(1)

	return ImageDataBits{imageData{*img.Size_, img.Data}}
}

// ImageDataBits is a bit-indexed version of ImageData.
type ImageDataBits struct {
	imageData
}

// Get returns the bit value at (x, y).
// True if the bit is set, false if not or if (x, y) is out of bounds.
func (img *ImageDataBits) Get(x, y int) bool {
	if img.inBounds(x, y) {
		i := img.offset(x, y)
		i, bit := i/8, byte(1<<(uint(i)%8))
		return img.data[i]&bit != 0
	}
	return false
}

// Set updates the bit at (x, y) to the given value.
// If (x, y) is out of bounds it does nothing.
func (img *ImageDataBits) Set(x, y int, value bool) {
	if img.inBounds(x, y) {
		i := img.offset(x, y)
		i, bit := i/8, byte(1<<(uint(i)%8))
		if value {
			img.data[i] |= bit // set
		} else {
			img.data[i] &^= bit // clear
		}
	}
}

// Bytes returns a byte-indexed version of the ImageData.
// It panics if ImageData.BitsPerPixel != 8.
func (img *ImageData) Bytes() ImageDataBytes {
	img.assertBPP(8)

	return ImageDataBytes{imageData{*img.Size_, img.Data}}
}

// ImageDataBytes is a byte-indexed version of ImageData.
type ImageDataBytes struct {
	imageData
}

// Get returns the byte value at (x, y).
// If (x, y) is out of bounds it returns 0.
func (img *ImageDataBytes) Get(x, y int) byte {
	if img.inBounds(x, y) {
		return img.data[img.offset(x, y)]
	}
	return 0
}

// Set updates the byte at (x, y) to the given value.
// If (x, y) is out of bounds it does nothing.
func (img *ImageDataBytes) Set(x, y int, value byte) {
	if img.inBounds(x, y) {
		img.data[img.offset(x, y)] = value
	}
}

// Ints returns an int32-indexed version of the ImageData.
// It panics if ImageData.BitsPerPixel != 32.
func (img *ImageData) Ints() ImageDataInt32 {
	img.assertBPP(32)

	return ImageDataInt32{imageData{*img.Size_, img.Data}}
}

// ImageDataInt32 is an int32-indexed version of ImageData.
// LittleEndian byte ordering is assumed.
type ImageDataInt32 struct {
	imageData
}

// Get returns the int32 value at (x, y).
// If (x, y) is out of bounds it returns 0.
func (img *ImageDataInt32) Get(x, y int) int32 {
	if img.inBounds(x, y) {
		i := img.offset(x, y)
		// Assuming this is LE byte order...
		return int32(binary.LittleEndian.Uint32(img.data[4*i : 4*(i+1)]))
	}
	return 0
}

// Set updates the int32 at (x, y) to the given value.
// If (x, y) is out of bounds it does nothing.
func (img *ImageDataInt32) Set(x, y int, value int32) {
	if img.inBounds(x, y) {
		i := img.offset(x, y)
		// Assuming this is LE byte order...
		binary.LittleEndian.PutUint32(img.data[4*i:4*(i+1)], uint32(value))
	}
}

// imageData contains common data and methods used by the typed versions above.
type imageData struct {
	size Size2DI
	data []byte
}

// Width is the horizontal size of the data.
func (img *imageData) Width() int {
	return int(img.size.X)
}

// Height is the vertical size of the data.
func (img *imageData) Height() int {
	return int(img.size.Y)
}

// inBounds checks that the coordinates fall within the valid range for the image.
func (img *imageData) inBounds(x, y int) bool {
	return 0 <= x && x < int(img.size.X) && 0 <= y && y < int(img.size.Y)
}

// offset converts XY coordinates into a linear offset into the ImageData.
func (img *imageData) offset(x, y int) int {
	// Image data is stored with an upper left origin
	return x + (int(img.size.Y)-1-y)*int(img.size.X)
}
