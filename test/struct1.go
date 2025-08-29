package test

import "github.com/Nocccer/protoreg/test/sub"

//go:generate go run ../cmd/main.go -type=Struct1 -v

type CustomUint16 uint16

type Struct1 struct {
	Ignored uint16
	Field1  string          `protoreg:"offset=0,size=8,char=16"`
	Field2  int16           `protoreg:"offset=8"`
	Field3  CustomUint16    `protoreg:"offset=9"`
	Field4  uint16          `protoreg:"offset=10"`
	Field5  string          `protoreg:"offset=11,size=8,char=8"`
	Field6  sub.CustomInt16 `protoreg:"offset=15"`
}
