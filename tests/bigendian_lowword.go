package tests

import "github.com/nocccer/protoreg/tests/extern"

//go:generate go run ../main.go -type=BigEndianLowWord,BigEndianLowWordAllCustom,BigEndianLowWordAllCustomExtern -v

type BigEndianLowWord struct {
	_             struct{} `protoreg:"wordorder=low"`
	Ignored       uint16
	Uint8Low      uint8      `protoreg:"offset=0,byte=low"`
	Uint8High     uint8      `protoreg:"offset=0,byte=high"`
	ByteLow       byte       `protoreg:"offset=1,byte=low"`
	ByteHigh      byte       `protoreg:"offset=1,byte=high"`
	Int8Low       int8       `protoreg:"offset=2,byte=low"`
	Int8High      int8       `protoreg:"offset=2,byte=high"`
	Uint16        uint16     `protoreg:"offset=3"`
	Int16         int16      `protoreg:"offset=4"`
	Uint32        uint32     `protoreg:"offset=5"`
	Int32         int32      `protoreg:"offset=7"`
	Uint64        uint64     `protoreg:"offset=9"`
	Int64         int64      `protoreg:"offset=13"`
	Float32       float32    `protoreg:"offset=17"`
	Float64       float64    `protoreg:"offset=21"`
	StringASCII8  string     `protoreg:"offset=25,size=9,char=8"`
	StringASCII16 string     `protoreg:"offset=34,size=9,char=16"`
	StringUTF816  string     `protoreg:"offset=43,size=9,char=16,charencoding=utf8"`
	Bool          bool       `protoreg:"offset=52"`
	Bit1          bool       `protoreg:"offset=53,bit=1"`
	Bit3          bool       `protoreg:"offset=53,bit=3"`
	Bit14         bool       `protoreg:"offset=53,bit=14"`
	Uint16Array   [5]uint16  `protoreg:"offset=54"`
	Uint32Array   [5]uint32  `protoreg:"offset=59"`
	Float32Array  [5]float32 `protoreg:"offset=79"`
	BoolArray     [5]bool    `protoreg:"offset=89"`
}

type BigEndianLowWordAllCustom struct {
	_             struct{} `protoreg:"wordorder=low"`
	Ignored       CustomUint16
	Uint8Low      CustomUint8      `protoreg:"offset=0,byte=low"`
	Uint8High     CustomUint8      `protoreg:"offset=0,byte=high"`
	ByteLow       CustomByte       `protoreg:"offset=1,byte=low"`
	ByteHigh      CustomByte       `protoreg:"offset=1,byte=high"`
	Int8Low       CustomInt8       `protoreg:"offset=2,byte=low"`
	Int8High      CustomInt8       `protoreg:"offset=2,byte=high"`
	Uint16        CustomUint16     `protoreg:"offset=3"`
	Int16         CustomInt16      `protoreg:"offset=4"`
	Uint32        CustomUint32     `protoreg:"offset=5"`
	Int32         CustomInt32      `protoreg:"offset=7"`
	Uint64        CustomUint64     `protoreg:"offset=9"`
	Int64         CustomInt64      `protoreg:"offset=13"`
	Float32       CustomFloat32    `protoreg:"offset=17"`
	Float64       CustomFloat64    `protoreg:"offset=21"`
	StringASCII8  CustomString     `protoreg:"offset=25,size=9,char=8"`
	StringASCII16 CustomString     `protoreg:"offset=34,size=9,char=16"`
	StringUTF816  CustomString     `protoreg:"offset=43,size=9,char=16,charencoding=utf8"`
	Bool          CustomBool       `protoreg:"offset=52"`
	Bit1          CustomBool       `protoreg:"offset=53,bit=1"`
	Bit3          CustomBool       `protoreg:"offset=53,bit=3"`
	Bit14         CustomBool       `protoreg:"offset=53,bit=14"`
	Uint16Array   [5]CustomUint16  `protoreg:"offset=54"`
	Uint32Array   [5]CustomUint32  `protoreg:"offset=59"`
	Float32Array  [5]CustomFloat32 `protoreg:"offset=79"`
	BoolArray     [5]CustomBool    `protoreg:"offset=89"`
}

type BigEndianLowWordAllCustomExtern struct {
	_             struct{} `protoreg:"wordorder=low"`
	Ignored       extern.CustomUint16
	Uint8Low      extern.CustomUint8      `protoreg:"offset=0,byte=low"`
	Uint8High     extern.CustomUint8      `protoreg:"offset=0,byte=high"`
	ByteLow       extern.CustomByte       `protoreg:"offset=1,byte=low"`
	ByteHigh      extern.CustomByte       `protoreg:"offset=1,byte=high"`
	Int8Low       extern.CustomInt8       `protoreg:"offset=2,byte=low"`
	Int8High      extern.CustomInt8       `protoreg:"offset=2,byte=high"`
	Uint16        extern.CustomUint16     `protoreg:"offset=3"`
	Int16         extern.CustomInt16      `protoreg:"offset=4"`
	Uint32        extern.CustomUint32     `protoreg:"offset=5"`
	Int32         extern.CustomInt32      `protoreg:"offset=7"`
	Uint64        extern.CustomUint64     `protoreg:"offset=9"`
	Int64         extern.CustomInt64      `protoreg:"offset=13"`
	Float32       extern.CustomFloat32    `protoreg:"offset=17"`
	Float64       extern.CustomFloat64    `protoreg:"offset=21"`
	StringASCII8  extern.CustomString     `protoreg:"offset=25,size=9,char=8"`
	StringASCII16 extern.CustomString     `protoreg:"offset=34,size=9,char=16"`
	StringUTF816  extern.CustomString     `protoreg:"offset=43,size=9,char=16,charencoding=utf8"`
	Bool          extern.CustomBool       `protoreg:"offset=52"`
	Bit1          extern.CustomBool       `protoreg:"offset=53,bit=1"`
	Bit3          extern.CustomBool       `protoreg:"offset=53,bit=3"`
	Bit14         extern.CustomBool       `protoreg:"offset=53,bit=14"`
	Uint16Array   [5]extern.CustomUint16  `protoreg:"offset=54"`
	Uint32Array   [5]extern.CustomUint32  `protoreg:"offset=59"`
	Float32Array  [5]extern.CustomFloat32 `protoreg:"offset=79"`
	BoolArray     [5]extern.CustomBool    `protoreg:"offset=89"`
}
