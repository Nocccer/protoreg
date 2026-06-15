package tests

//go:generate go run ../cmd/main.go -type=CustomStructKey -key=myreg -v

type CustomStructKey struct {
	UInt16 uint16 `myreg:"offset=0"`
}

//go:generate go run ../cmd/main.go -type=OnlyMarshaler,OnlyUnmarshaler,CustomFuncNames -v

type OnlyMarshaler struct {
	_      struct{} `protoreg:"mode=marshal"`
	UInt16 uint16   `protoreg:"offset=0"`
}

type OnlyUnmarshaler struct {
	_      struct{} `protoreg:"mode=unmarshal"`
	UInt16 uint16   `protoreg:"offset=0"`
}

type CustomFuncNames struct {
	_      struct{} `protoreg:"marshalfunc=CustomMarshal,unmarshalfunc=CustomUnmarshal"`
	UInt16 uint16   `protoreg:"offset=0"`
}
