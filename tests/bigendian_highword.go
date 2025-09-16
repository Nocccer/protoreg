package tests

import "github.com/Nocccer/protoreg/tests/sub"

//go:generate go run ../cmd/main.go -type=BigEndianHighWord -v

type BigEndianHighWord struct {
	Ignored uint16
	// Uint8        uint8           `protoreg:"offset=0,byte=high"`
	// Int8         int8            `protoreg:"offset=0,byte=low"`
	Uint16 uint16 `protoreg:"offset=1"`
	Int16  int16  `protoreg:"offset=2"`
	Uint32 uint32 `protoreg:"offset=3"`
	// Int32        int32           `protoreg:"offset=5"`
	// Uint64       uint64          `protoreg:"offset=7"`
	// Int64        int64           `protoreg:"offset=11"`
	// Float64      float64         `protoreg:"offset=15"`
	// Float32      float32         `protoreg:"offset=19"`
	String8      string          `protoreg:"offset=21,size=8,char=8"`
	String16     string          `protoreg:"offset=29,size=8,char=16"`
	CustomUint16 CustomUint16    `protoreg:"offset=37"`
	CustomInt16  sub.CustomInt16 `protoreg:"offset=38"`
}
