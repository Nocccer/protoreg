package tests

//go:generate protoreg -type=StringChar8BigEndian,StringChar8LittleEndian,StringUTF16BigEndian,StringUTF16LittleEndian,StringUTF32BigHighWord,StringUTF32LittleLowWord -v

type StringChar8BigEndian struct {
	_     struct{} `protoreg:"encoding=big"`
	Value string   `protoreg:"offset=0,size=2,char=8"`
}

type StringChar8LittleEndian struct {
	_     struct{} `protoreg:"encoding=little"`
	Value string   `protoreg:"offset=0,size=2,char=8"`
}

type StringUTF16BigEndian struct {
	_     struct{} `protoreg:"encoding=big"`
	Value string   `protoreg:"offset=0,size=5,char=16,charencoding=utf8"`
}

type StringUTF16LittleEndian struct {
	_     struct{} `protoreg:"encoding=little"`
	Value string   `protoreg:"offset=0,size=5,char=16,charencoding=utf8"`
}

type StringUTF32BigHighWord struct {
	_     struct{} `protoreg:"encoding=big,wordorder=high"`
	Value string   `protoreg:"offset=0,size=4,char=32,charencoding=utf8"`
}

type StringUTF32LittleLowWord struct {
	_     struct{} `protoreg:"encoding=little,wordorder=low"`
	Value string   `protoreg:"offset=0,size=4,char=32,charencoding=utf8"`
}
