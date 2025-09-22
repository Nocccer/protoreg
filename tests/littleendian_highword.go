package tests

import "github.com/Nocccer/protoreg/tests/sub"

//go:generate go run ../cmd/main.go -type=LittleEndianHighWord -v

type LittleEndianHighWord struct {
	_             struct{} `protoreg:"encoding=little"`
	Ignored       uint16
	Uint8Low      uint8           `protoreg:"offset=0,byte=low"`
	Uint8High     uint8           `protoreg:"offset=0,byte=high"`
	ByteLow       byte            `protoreg:"offset=1,byte=low"`
	ByteHigh      byte            `protoreg:"offset=1,byte=high"`
	Int8Low       int8            `protoreg:"offset=2,byte=low"`
	Int8High      int8            `protoreg:"offset=2,byte=high"`
	Uint16        uint16          `protoreg:"offset=3"`
	Int16         int16           `protoreg:"offset=4"`
	Uint32        uint32          `protoreg:"offset=5"`
	Int32         int32           `protoreg:"offset=7"`
	Uint64        uint64          `protoreg:"offset=9"`
	Int64         int64           `protoreg:"offset=13"`
	Float32       float32         `protoreg:"offset=17"`
	Float64       float64         `protoreg:"offset=21"`
	StringASCII8  string          `protoreg:"offset=25,size=9,char=8"`
	StringASCII16 string          `protoreg:"offset=34,size=9,char=16"`
	StringUTF816  string          `protoreg:"offset=43,size=9,char=16,charencoding=utf8"`
	CustomUint16  CustomUint16    `protoreg:"offset=52"`
	CustomInt16   sub.CustomInt16 `protoreg:"offset=53"`
}
