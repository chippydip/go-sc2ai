// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: common.proto

package api

import (
	encoding_binary "encoding/binary"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Race int32

const (
	Race_NoRace  Race = 0
	Race_Terran  Race = 1
	Race_Zerg    Race = 2
	Race_Protoss Race = 3
	Race_Random  Race = 4
)

var Race_name = map[int32]string{
	0: "NoRace",
	1: "Terran",
	2: "Zerg",
	3: "Protoss",
	4: "Random",
}

var Race_value = map[string]int32{
	"NoRace":  0,
	"Terran":  1,
	"Zerg":    2,
	"Protoss": 3,
	"Random":  4,
}

func (x Race) String() string {
	return proto.EnumName(Race_name, int32(x))
}

func (Race) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0}
}

type AvailableAbility struct {
	AbilityId     AbilityID `protobuf:"varint,1,opt,name=ability_id,json=abilityId,proto3,casttype=AbilityID" json:"ability_id,omitempty"`
	RequiresPoint bool      `protobuf:"varint,2,opt,name=requires_point,json=requiresPoint,proto3" json:"requires_point,omitempty"`
}

func (m *AvailableAbility) Reset()         { *m = AvailableAbility{} }
func (m *AvailableAbility) String() string { return proto.CompactTextString(m) }
func (*AvailableAbility) ProtoMessage()    {}
func (*AvailableAbility) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0}
}
func (m *AvailableAbility) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AvailableAbility) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AvailableAbility.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AvailableAbility) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AvailableAbility.Merge(m, src)
}
func (m *AvailableAbility) XXX_Size() int {
	return m.Size()
}
func (m *AvailableAbility) XXX_DiscardUnknown() {
	xxx_messageInfo_AvailableAbility.DiscardUnknown(m)
}

var xxx_messageInfo_AvailableAbility proto.InternalMessageInfo

func (m *AvailableAbility) GetAbilityId() AbilityID {
	if m != nil {
		return m.AbilityId
	}
	return 0
}

func (m *AvailableAbility) GetRequiresPoint() bool {
	if m != nil {
		return m.RequiresPoint
	}
	return false
}

type ImageData struct {
	BitsPerPixel int32    `protobuf:"varint,1,opt,name=bits_per_pixel,json=bitsPerPixel,proto3" json:"bits_per_pixel,omitempty"`
	Size_        *Size2DI `protobuf:"bytes,2,opt,name=size,proto3" json:"size,omitempty"`
	Data         []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *ImageData) Reset()         { *m = ImageData{} }
func (m *ImageData) String() string { return proto.CompactTextString(m) }
func (*ImageData) ProtoMessage()    {}
func (*ImageData) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{1}
}
func (m *ImageData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ImageData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ImageData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ImageData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageData.Merge(m, src)
}
func (m *ImageData) XXX_Size() int {
	return m.Size()
}
func (m *ImageData) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageData.DiscardUnknown(m)
}

var xxx_messageInfo_ImageData proto.InternalMessageInfo

func (m *ImageData) GetBitsPerPixel() int32 {
	if m != nil {
		return m.BitsPerPixel
	}
	return 0
}

func (m *ImageData) GetSize_() *Size2DI {
	if m != nil {
		return m.Size_
	}
	return nil
}

func (m *ImageData) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// Point on the screen/minimap (e.g., 0..64).
// Note: bottom left of the screen is 0, 0.
type PointI struct {
	X int32 `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	Y int32 `protobuf:"varint,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (m *PointI) Reset()         { *m = PointI{} }
func (m *PointI) String() string { return proto.CompactTextString(m) }
func (*PointI) ProtoMessage()    {}
func (*PointI) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{2}
}
func (m *PointI) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PointI) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PointI.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PointI) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PointI.Merge(m, src)
}
func (m *PointI) XXX_Size() int {
	return m.Size()
}
func (m *PointI) XXX_DiscardUnknown() {
	xxx_messageInfo_PointI.DiscardUnknown(m)
}

var xxx_messageInfo_PointI proto.InternalMessageInfo

func (m *PointI) GetX() int32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *PointI) GetY() int32 {
	if m != nil {
		return m.Y
	}
	return 0
}

// Screen space rectangular area.
type RectangleI struct {
	P0 *PointI `protobuf:"bytes,1,opt,name=p0,proto3" json:"p0,omitempty"`
	P1 *PointI `protobuf:"bytes,2,opt,name=p1,proto3" json:"p1,omitempty"`
}

func (m *RectangleI) Reset()         { *m = RectangleI{} }
func (m *RectangleI) String() string { return proto.CompactTextString(m) }
func (*RectangleI) ProtoMessage()    {}
func (*RectangleI) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{3}
}
func (m *RectangleI) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RectangleI) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RectangleI.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RectangleI) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RectangleI.Merge(m, src)
}
func (m *RectangleI) XXX_Size() int {
	return m.Size()
}
func (m *RectangleI) XXX_DiscardUnknown() {
	xxx_messageInfo_RectangleI.DiscardUnknown(m)
}

var xxx_messageInfo_RectangleI proto.InternalMessageInfo

func (m *RectangleI) GetP0() *PointI {
	if m != nil {
		return m.P0
	}
	return nil
}

func (m *RectangleI) GetP1() *PointI {
	if m != nil {
		return m.P1
	}
	return nil
}

// Point on the game board, 0..255.
// Note: bottom left of the screen is 0, 0.
type Point2D struct {
	X float32 `protobuf:"fixed32,1,opt,name=x,proto3" json:"x,omitempty"`
	Y float32 `protobuf:"fixed32,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (m *Point2D) Reset()         { *m = Point2D{} }
func (m *Point2D) String() string { return proto.CompactTextString(m) }
func (*Point2D) ProtoMessage()    {}
func (*Point2D) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{4}
}
func (m *Point2D) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Point2D) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Point2D.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Point2D) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Point2D.Merge(m, src)
}
func (m *Point2D) XXX_Size() int {
	return m.Size()
}
func (m *Point2D) XXX_DiscardUnknown() {
	xxx_messageInfo_Point2D.DiscardUnknown(m)
}

var xxx_messageInfo_Point2D proto.InternalMessageInfo

func (m *Point2D) GetX() float32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Point2D) GetY() float32 {
	if m != nil {
		return m.Y
	}
	return 0
}

// Point on the game board, 0..255.
// Note: bottom left of the screen is 0, 0.
type Point struct {
	X float32 `protobuf:"fixed32,1,opt,name=x,proto3" json:"x,omitempty"`
	Y float32 `protobuf:"fixed32,2,opt,name=y,proto3" json:"y,omitempty"`
	Z float32 `protobuf:"fixed32,3,opt,name=z,proto3" json:"z,omitempty"`
}

func (m *Point) Reset()         { *m = Point{} }
func (m *Point) String() string { return proto.CompactTextString(m) }
func (*Point) ProtoMessage()    {}
func (*Point) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{5}
}
func (m *Point) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Point) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Point.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Point) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Point.Merge(m, src)
}
func (m *Point) XXX_Size() int {
	return m.Size()
}
func (m *Point) XXX_DiscardUnknown() {
	xxx_messageInfo_Point.DiscardUnknown(m)
}

var xxx_messageInfo_Point proto.InternalMessageInfo

func (m *Point) GetX() float32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Point) GetY() float32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Point) GetZ() float32 {
	if m != nil {
		return m.Z
	}
	return 0
}

// Screen dimensions.
type Size2DI struct {
	X int32 `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	Y int32 `protobuf:"varint,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (m *Size2DI) Reset()         { *m = Size2DI{} }
func (m *Size2DI) String() string { return proto.CompactTextString(m) }
func (*Size2DI) ProtoMessage()    {}
func (*Size2DI) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{6}
}
func (m *Size2DI) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Size2DI) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Size2DI.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Size2DI) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Size2DI.Merge(m, src)
}
func (m *Size2DI) XXX_Size() int {
	return m.Size()
}
func (m *Size2DI) XXX_DiscardUnknown() {
	xxx_messageInfo_Size2DI.DiscardUnknown(m)
}

var xxx_messageInfo_Size2DI proto.InternalMessageInfo

func (m *Size2DI) GetX() int32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Size2DI) GetY() int32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func init() {
	proto.RegisterEnum("SC2APIProtocol.Race", Race_name, Race_value)
	proto.RegisterType((*AvailableAbility)(nil), "SC2APIProtocol.AvailableAbility")
	proto.RegisterType((*ImageData)(nil), "SC2APIProtocol.ImageData")
	proto.RegisterType((*PointI)(nil), "SC2APIProtocol.PointI")
	proto.RegisterType((*RectangleI)(nil), "SC2APIProtocol.RectangleI")
	proto.RegisterType((*Point2D)(nil), "SC2APIProtocol.Point2D")
	proto.RegisterType((*Point)(nil), "SC2APIProtocol.Point")
	proto.RegisterType((*Size2DI)(nil), "SC2APIProtocol.Size2DI")
}

func init() { proto.RegisterFile("common.proto", fileDescriptor_555bd8c177793206) }

var fileDescriptor_555bd8c177793206 = []byte{
	// 412 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xcf, 0x6f, 0x94, 0x40,
	0x14, 0xc7, 0x77, 0x28, 0xbb, 0xed, 0xbe, 0xdd, 0x6e, 0xc8, 0x1c, 0x74, 0xe3, 0x01, 0x09, 0x69,
	0xcd, 0x46, 0x0d, 0x76, 0xf1, 0xe8, 0x89, 0xca, 0x85, 0x8b, 0x21, 0x53, 0x4f, 0x8d, 0x09, 0x19,
	0x60, 0x42, 0x26, 0x01, 0x06, 0x07, 0x34, 0xbb, 0xfc, 0x15, 0xfe, 0x59, 0x1e, 0x7b, 0xf4, 0x64,
	0xcc, 0xee, 0x7f, 0xe1, 0xc9, 0xcc, 0x80, 0x31, 0x35, 0x69, 0x6f, 0xef, 0x3d, 0x3e, 0xf9, 0xfe,
	0x20, 0x03, 0xcb, 0x4c, 0x54, 0x95, 0xa8, 0xbd, 0x46, 0x8a, 0x4e, 0xe0, 0xd5, 0xcd, 0x7b, 0x3f,
	0x88, 0xa3, 0x58, 0x2d, 0x99, 0x28, 0x9f, 0x41, 0x21, 0x0a, 0x31, 0x7c, 0x73, 0x0b, 0xb0, 0x82,
	0xaf, 0x94, 0x97, 0x34, 0x2d, 0x59, 0x90, 0xf2, 0x92, 0x77, 0x7b, 0xfc, 0x1a, 0x80, 0x0e, 0x63,
	0xc2, 0xf3, 0x35, 0x72, 0xd0, 0x66, 0x7a, 0x7d, 0xfe, 0xfb, 0xe7, 0xf3, 0xf9, 0x08, 0x44, 0x21,
	0x99, 0x8f, 0x40, 0x94, 0xe3, 0x4b, 0x58, 0x49, 0xf6, 0xf9, 0x0b, 0x97, 0xac, 0x4d, 0x1a, 0xc1,
	0xeb, 0x6e, 0x6d, 0x38, 0x68, 0x73, 0x46, 0xce, 0xff, 0x5e, 0x63, 0x75, 0x74, 0x25, 0xcc, 0xa3,
	0x8a, 0x16, 0x2c, 0xa4, 0x1d, 0xc5, 0x17, 0xb0, 0x4a, 0x79, 0xd7, 0x26, 0x0d, 0x93, 0x49, 0xc3,
	0x77, 0xac, 0x1c, 0x5c, 0xc8, 0x52, 0x5d, 0x63, 0x26, 0x63, 0x75, 0xc3, 0xaf, 0xc0, 0x6c, 0x79,
	0xcf, 0xb4, 0xde, 0xc2, 0x7f, 0xea, 0xdd, 0xaf, 0xe1, 0xdd, 0xf0, 0x9e, 0xf9, 0x61, 0x44, 0x34,
	0x84, 0x31, 0x98, 0x39, 0xed, 0xe8, 0xfa, 0xc4, 0x41, 0x9b, 0x25, 0xd1, 0xb3, 0x7b, 0x01, 0x33,
	0x6d, 0x1e, 0xe1, 0x25, 0xa0, 0xdd, 0xe8, 0x81, 0x76, 0x6a, 0xdb, 0x6b, 0xd5, 0x29, 0x41, 0x7b,
	0xf7, 0x13, 0x00, 0x61, 0x59, 0x47, 0xeb, 0xa2, 0x64, 0x11, 0x7e, 0x01, 0x46, 0x73, 0xa5, 0xd1,
	0x85, 0xff, 0xe4, 0x7f, 0xcb, 0x41, 0x8d, 0x18, 0xcd, 0x95, 0xe6, 0xb6, 0x63, 0xb4, 0x87, 0xb9,
	0xad, 0x7b, 0x09, 0xa7, 0x7a, 0xf3, 0xc3, 0x7f, 0x21, 0x8c, 0x7b, 0x21, 0x0c, 0x15, 0x62, 0x0b,
	0x53, 0x8d, 0x3d, 0x06, 0xa9, 0xad, 0xd7, 0x05, 0x0d, 0x82, 0x7a, 0xa5, 0x3c, 0xfe, 0x82, 0xc7,
	0xea, 0xbd, 0x0c, 0xc0, 0x24, 0x34, 0x63, 0x18, 0x60, 0xf6, 0x41, 0xa8, 0xc9, 0x9a, 0xa8, 0xf9,
	0x23, 0x93, 0x92, 0xd6, 0x16, 0xc2, 0x67, 0x60, 0xde, 0x32, 0x59, 0x58, 0x06, 0x5e, 0xc0, 0xa9,
	0x6e, 0xd0, 0xb6, 0xd6, 0x89, 0x42, 0x08, 0xad, 0x73, 0x51, 0x59, 0xe6, 0xb5, 0xf3, 0xfd, 0x60,
	0xa3, 0xbb, 0x83, 0x8d, 0x7e, 0x1d, 0x6c, 0xf4, 0xed, 0x68, 0x4f, 0xee, 0x8e, 0xf6, 0xe4, 0xc7,
	0xd1, 0x9e, 0xdc, 0xce, 0xbc, 0x37, 0xef, 0x68, 0xc3, 0xd3, 0x99, 0x7e, 0x4d, 0x6f, 0xff, 0x04,
	0x00, 0x00, 0xff, 0xff, 0xcf, 0xa5, 0x6d, 0x28, 0x79, 0x02, 0x00, 0x00,
}

func (m *AvailableAbility) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AvailableAbility) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AvailableAbility) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RequiresPoint {
		i--
		if m.RequiresPoint {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.AbilityId != 0 {
		i = encodeVarintCommon(dAtA, i, uint64(m.AbilityId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ImageData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ImageData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ImageData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintCommon(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Size_ != nil {
		{
			size, err := m.Size_.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintCommon(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.BitsPerPixel != 0 {
		i = encodeVarintCommon(dAtA, i, uint64(m.BitsPerPixel))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *PointI) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PointI) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PointI) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Y != 0 {
		i = encodeVarintCommon(dAtA, i, uint64(m.Y))
		i--
		dAtA[i] = 0x10
	}
	if m.X != 0 {
		i = encodeVarintCommon(dAtA, i, uint64(m.X))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *RectangleI) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RectangleI) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RectangleI) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.P1 != nil {
		{
			size, err := m.P1.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintCommon(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.P0 != nil {
		{
			size, err := m.P0.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintCommon(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Point2D) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Point2D) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Point2D) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Y != 0 {
		i -= 4
		encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(math.Float32bits(float32(m.Y))))
		i--
		dAtA[i] = 0x15
	}
	if m.X != 0 {
		i -= 4
		encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(math.Float32bits(float32(m.X))))
		i--
		dAtA[i] = 0xd
	}
	return len(dAtA) - i, nil
}

func (m *Point) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Point) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Point) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Z != 0 {
		i -= 4
		encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(math.Float32bits(float32(m.Z))))
		i--
		dAtA[i] = 0x1d
	}
	if m.Y != 0 {
		i -= 4
		encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(math.Float32bits(float32(m.Y))))
		i--
		dAtA[i] = 0x15
	}
	if m.X != 0 {
		i -= 4
		encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(math.Float32bits(float32(m.X))))
		i--
		dAtA[i] = 0xd
	}
	return len(dAtA) - i, nil
}

func (m *Size2DI) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Size2DI) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Size2DI) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Y != 0 {
		i = encodeVarintCommon(dAtA, i, uint64(m.Y))
		i--
		dAtA[i] = 0x10
	}
	if m.X != 0 {
		i = encodeVarintCommon(dAtA, i, uint64(m.X))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintCommon(dAtA []byte, offset int, v uint64) int {
	offset -= sovCommon(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *AvailableAbility) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AbilityId != 0 {
		n += 1 + sovCommon(uint64(m.AbilityId))
	}
	if m.RequiresPoint {
		n += 2
	}
	return n
}

func (m *ImageData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BitsPerPixel != 0 {
		n += 1 + sovCommon(uint64(m.BitsPerPixel))
	}
	if m.Size_ != nil {
		l = m.Size_.Size()
		n += 1 + l + sovCommon(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovCommon(uint64(l))
	}
	return n
}

func (m *PointI) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.X != 0 {
		n += 1 + sovCommon(uint64(m.X))
	}
	if m.Y != 0 {
		n += 1 + sovCommon(uint64(m.Y))
	}
	return n
}

func (m *RectangleI) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.P0 != nil {
		l = m.P0.Size()
		n += 1 + l + sovCommon(uint64(l))
	}
	if m.P1 != nil {
		l = m.P1.Size()
		n += 1 + l + sovCommon(uint64(l))
	}
	return n
}

func (m *Point2D) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.X != 0 {
		n += 5
	}
	if m.Y != 0 {
		n += 5
	}
	return n
}

func (m *Point) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.X != 0 {
		n += 5
	}
	if m.Y != 0 {
		n += 5
	}
	if m.Z != 0 {
		n += 5
	}
	return n
}

func (m *Size2DI) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.X != 0 {
		n += 1 + sovCommon(uint64(m.X))
	}
	if m.Y != 0 {
		n += 1 + sovCommon(uint64(m.Y))
	}
	return n
}

func sovCommon(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCommon(x uint64) (n int) {
	return sovCommon(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AvailableAbility) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommon
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AvailableAbility: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AvailableAbility: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AbilityId", wireType)
			}
			m.AbilityId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AbilityId |= AbilityID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequiresPoint", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.RequiresPoint = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipCommon(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommon
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ImageData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommon
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ImageData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ImageData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BitsPerPixel", wireType)
			}
			m.BitsPerPixel = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BitsPerPixel |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Size_", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Size_ == nil {
				m.Size_ = &Size2DI{}
			}
			if err := m.Size_.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommon(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommon
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PointI) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommon
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PointI: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PointI: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field X", wireType)
			}
			m.X = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.X |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Y", wireType)
			}
			m.Y = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Y |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCommon(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommon
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RectangleI) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommon
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RectangleI: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RectangleI: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field P0", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.P0 == nil {
				m.P0 = &PointI{}
			}
			if err := m.P0.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field P1", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.P1 == nil {
				m.P1 = &PointI{}
			}
			if err := m.P1.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommon(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommon
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Point2D) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommon
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Point2D: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Point2D: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field X", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
			iNdEx += 4
			m.X = float32(math.Float32frombits(v))
		case 2:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field Y", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
			iNdEx += 4
			m.Y = float32(math.Float32frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipCommon(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommon
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Point) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommon
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Point: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Point: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field X", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
			iNdEx += 4
			m.X = float32(math.Float32frombits(v))
		case 2:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field Y", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
			iNdEx += 4
			m.Y = float32(math.Float32frombits(v))
		case 3:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field Z", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
			iNdEx += 4
			m.Z = float32(math.Float32frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipCommon(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommon
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Size2DI) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommon
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Size2DI: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Size2DI: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field X", wireType)
			}
			m.X = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.X |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Y", wireType)
			}
			m.Y = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Y |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCommon(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommon
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipCommon(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCommon
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthCommon
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCommon
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCommon
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCommon        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCommon          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCommon = fmt.Errorf("proto: unexpected end of group")
)
