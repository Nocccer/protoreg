package tests

import "strconv"

type (
	CustomUint8   uint8
	CustomInt8    int8
	CustomByte    byte
	CustomUint16  uint16
	CustomInt16   int16
	CustomUint32  uint32
	CustomInt32   int32
	CustomUint64  uint64
	CustomInt64   int64
	CustomFloat32 float32
	CustomFloat64 float64
	CustomString  string
	CustomBool    bool
)

func (c CustomUint16) String() string {
	return strconv.Itoa(int(c))
}

type BitField16[T UInt16Stringer] uint16

type UInt16Stringer interface {
	~uint16
	String() string
}
