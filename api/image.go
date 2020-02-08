package api

import (
	"encoding/binary"
	"log"
)

// Copy returns an ImageData with a separate Data slice copied from the original.
func (img ImageData) Copy() ImageData {
	data := make([]byte, len(img.Data))
	copy(data, img.Data)

	return ImageData{
		BitsPerPixel: img.BitsPerPixel,
		Size_:        &Size2DI{X: img.Size_.X, Y: img.Size_.Y},
		Data:         data,
	}
}

// assertBPP checks for the expected pixel size and panics if it doesn't match.
func (img ImageData) assertBPP(count int32) {
	if img.BitsPerPixel != count {
		log.Panicf("bad BitsPerPixel, expected %v got %v", count, img.BitsPerPixel)
	}
}

// Bits returns a bit-indexed version of the ImageData.
// It panics if ImageData.BitsPerPixel != 1.
func (img ImageData) Bits() ImageDataBits {
	img.assertBPP(1)

	return ImageDataBits{imageData{*img.Size_, img.Data}}
}

// ImageDataBits is a bit-indexed version of ImageData.
type ImageDataBits struct {
	imageData
}

// NewImageDataBits returns an empty bit-indexed ImageData of the given size.
func NewImageDataBits(w, h int32) ImageDataBits {
	size := Size2DI{int32(w), int32(h)}
	data := make([]byte, (w*h+7)/8)

	return ImageDataBits{imageData{size, data}}
}

// Get returns the bit value at (x, y).
// True if the bit is set, false if not or if (x, y) is out of bounds.
func (img ImageDataBits) Get(x, y int32) bool {
	if img.InBounds(x, y) {
		i := img.offset(x, y)
		i, bit := i/8, byte(1<<(7-(uint(i)%8)))
		return img.data[i]&bit != 0
	}
	return false
}

// Copy returns an ImageDataBits with a separate data slice copied from the original.
func (img ImageDataBits) Copy() ImageDataBits {
	data := make([]byte, len(img.data))
	copy(data, img.data)

	return ImageDataBits{imageData{img.size, data}}
}

// Set updates the bit at (x, y) to the given value.
// If (x, y) is out of bounds it does nothing.
func (img ImageDataBits) Set(x, y int32, value bool) {
	if img.InBounds(x, y) {
		i := img.offset(x, y)
		i, bit := i/8, byte(1<<(7-(uint(i)%8)))
		if value {
			img.data[i] |= bit // set
		} else {
			img.data[i] &^= bit // clear
		}
	}
}

// ToBytes converts a bitmap into a bytemap with false -> 0 and true -> 255.
func (img ImageDataBits) ToBytes() ImageDataBytes {
	bytes := ImageDataBytes{imageData{
		img.size,
		make([]byte, img.Width()*img.Height()),
	}}

	for y := int32(0); y < img.Height(); y++ {
		for x := int32(0); x < img.Width(); x++ {
			if img.Get(x, y) {
				bytes.Set(x, y, 255)
			}
		}
	}

	return bytes
}

// Bytes returns a byte-indexed version of the ImageData.
// It panics if ImageData.BitsPerPixel != 8.
func (img ImageData) Bytes() ImageDataBytes {
	img.assertBPP(8)

	return ImageDataBytes{imageData{*img.Size_, img.Data}}
}

// ImageDataBytes is a byte-indexed version of ImageData.
type ImageDataBytes struct {
	imageData
}

// NewImageDataBytes returns an empty byte-indexed ImageData of the given size.
func NewImageDataBytes(w, h int32) ImageDataBytes {
	size := Size2DI{int32(w), int32(h)}
	data := make([]byte, w*h)

	return ImageDataBytes{imageData{size, data}}
}

// Copy returns an ImageDataBytes with a separate data slice copied from the original.
func (img ImageDataBytes) Copy() ImageDataBytes {
	data := make([]byte, len(img.data))
	copy(data, img.data)

	return ImageDataBytes{imageData{img.size, data}}
}

// Get returns the byte value at (x, y).
// If (x, y) is out of bounds it returns 0.
func (img ImageDataBytes) Get(x, y int32) byte {
	if img.InBounds(x, y) {
		return img.data[img.offset(x, y)]
	}
	return 0
}

// Set updates the byte at (x, y) to the given value.
// If (x, y) is out of bounds it does nothing.
func (img ImageDataBytes) Set(x, y int32, value byte) {
	if img.InBounds(x, y) {
		img.data[img.offset(x, y)] = value
	}
}

// Ints returns an int32-indexed version of the ImageData.
// It panics if ImageData.BitsPerPixel != 32.
func (img ImageData) Ints() ImageDataInt32 {
	img.assertBPP(32)

	return ImageDataInt32{imageData{*img.Size_, img.Data}}
}

// ImageDataInt32 is an int32-indexed version of ImageData.
// LittleEndian byte ordering is assumed.
type ImageDataInt32 struct {
	imageData
}

// NewImageDataInts returns an empty int32-indexed ImageData of the given size.
func NewImageDataInts(w, h int32) ImageDataInt32 {
	size := Size2DI{int32(w), int32(h)}
	data := make([]byte, w*h*4)

	return ImageDataInt32{imageData{size, data}}
}

// Copy returns an ImageDataInt32 with a separate data slice copied from the original.
func (img ImageDataInt32) Copy() ImageDataInt32 {
	data := make([]byte, len(img.data))
	copy(data, img.data)

	return ImageDataInt32{imageData{img.size, data}}
}

// Get returns the int32 value at (x, y).
// If (x, y) is out of bounds it returns 0.
func (img ImageDataInt32) Get(x, y int32) int32 {
	if img.InBounds(x, y) {
		i := img.offset(x, y)
		// Assuming this is LE byte order...
		return int32(binary.LittleEndian.Uint32(img.data[4*i : 4*(i+1)]))
	}
	return 0
}

// Set updates the int32 at (x, y) to the given value.
// If (x, y) is out of bounds it does nothing.
func (img ImageDataInt32) Set(x, y int32, value int32) {
	if img.InBounds(x, y) {
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
func (img imageData) Width() int32 {
	return img.size.X
}

// Height is the vertical size of the data.
func (img imageData) Height() int32 {
	return img.size.Y
}

// InBounds checks that the coordinates fall within the valid range for the image.
func (img imageData) InBounds(x, y int32) bool {
	return 0 <= x && x < img.Width() && 0 <= y && y < img.Height()
}

// offset converts XY coordinates into a linear offset into the ImageData.
func (img imageData) offset(x, y int32) int32 {
	// Image data is stored with an upper left origin
	return x + y*img.Width()
}
