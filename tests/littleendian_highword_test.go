package tests_test

import (
	"testing"

	"github.com/nocccer/protoreg/tests"
	"github.com/nocccer/protoreg/tests/extern"
	"github.com/stretchr/testify/suite"
)

func TestLittleEndianHighWord(t *testing.T) {
	suite.Run(t, new(LittleEndianHighWordTestSuite))
	suite.Run(t, new(LittleEndianHighWordAllCustomTestSuite))
	suite.Run(t, new(LittleEndianHighWordAllCustomExternTestSuite))
}

type LittleEndianHighWordTestSuite struct {
	suite.Suite
	LittleEndianHighWord tests.LittleEndianHighWord
}

func (s *LittleEndianHighWordTestSuite) SetupTest() {
	s.LittleEndianHighWord = tests.LittleEndianHighWord{
		Ignored:       0x4321,
		Uint8High:     0x11,
		Uint8Low:      0x22,
		ByteHigh:      0x33,
		ByteLow:       0x44,
		Int8Low:       -50,
		Int8High:      -100,
		Uint16:        0x1234,
		Int16:         -42,
		Uint32:        0x12345678,
		Int32:         -111222,
		Uint64:        0x1234123412341234,
		Int64:         -1111222233334444,
		Float32:       1234.5678,
		Float64:       -1234.5678,
		StringASCII8:  "TestData1",
		StringASCII16: "TestData2",
		StringUTF816:  "TestDäta3", // add utf8 char
		StringASCII32: "Te",
		StringUTF832:  "T😀",
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
		Uint16Array:   [5]uint16{0x1111, 0x2222, 0x3333, 0x4444, 0x5555},
		Uint32Array:   [5]uint32{0x11111111, 0x22222222, 0x33333333, 0x44444444, 0x55555555},
		Float32Array:  [5]float32{1.1, 2.2, 3.3, 4.4, 5.5},
		BoolArray:     [5]bool{true, false, true, false, true},
	}
}

func (s *LittleEndianHighWordTestSuite) TestMarshalUnmarshal() {
	reg, err := s.LittleEndianHighWord.Marshal()
	s.Require().NoError(err)

	out := &tests.LittleEndianHighWord{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Empty(out.Ignored)

	out.Ignored = s.LittleEndianHighWord.Ignored // Ignored field is not set by Unmarshal, set it manually for comparison

	s.Equal(s.LittleEndianHighWord, *out)
}

func BenchmarkLittleEndianHighWordMarshal(b *testing.B) {
	test := &tests.LittleEndianHighWord{
		Ignored:       0x4321,
		Uint8High:     0x11,
		Uint8Low:      0x22,
		ByteHigh:      0x33,
		ByteLow:       0x44,
		Int8Low:       -50,
		Int8High:      -100,
		Uint16:        0x1234,
		Int16:         -42,
		Uint32:        0x12345678,
		Int32:         -111222,
		Uint64:        0x1234123412341234,
		Int64:         -1111222233334444,
		Float32:       1234.5678,
		Float64:       -1234.5678,
		StringASCII8:  "TestData1",
		StringASCII16: "TestData2",
		StringUTF816:  "TestDäta3", // add utf8 char
		StringASCII32: "Te",
		StringUTF832:  "T😀",
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
		Uint16Array:   [5]uint16{0x1111, 0x2222, 0x3333, 0x4444, 0x5555},
		Uint32Array:   [5]uint32{0x11111111, 0x22222222, 0x33333333, 0x44444444, 0x55555555},
		Float32Array:  [5]float32{1.1, 2.2, 3.3, 4.4, 5.5},
		BoolArray:     [5]bool{true, false, true, false, true},
	}

	for b.Loop() {
		_, err := test.Marshal()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLittleEndianHighWordUnmarshal(b *testing.B) {
	test := &tests.LittleEndianHighWord{
		Ignored:       0x4321,
		Uint8High:     0x11,
		Uint8Low:      0x22,
		ByteHigh:      0x33,
		ByteLow:       0x44,
		Int8Low:       -50,
		Int8High:      -100,
		Uint16:        0x1234,
		Int16:         -42,
		Uint32:        0x12345678,
		Int32:         -111222,
		Uint64:        0x1234123412341234,
		Int64:         -1111222233334444,
		Float32:       1234.5678,
		Float64:       -1234.5678,
		StringASCII8:  "TestData1",
		StringASCII16: "TestData2",
		StringUTF816:  "TestDäta3", // add utf8 char
		StringASCII32: "Te",
		StringUTF832:  "T😀",
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
		Uint16Array:   [5]uint16{0x1111, 0x2222, 0x3333, 0x4444, 0x5555},
		Uint32Array:   [5]uint32{0x11111111, 0x22222222, 0x33333333, 0x44444444, 0x55555555},
		Float32Array:  [5]float32{1.1, 2.2, 3.3, 4.4, 5.5},
		BoolArray:     [5]bool{true, false, true, false, true},
	}

	reg, err := test.Marshal()
	if err != nil {
		b.Fatal(err)
	}

	test = &tests.LittleEndianHighWord{}

	for b.Loop() {
		err := test.Unmarshal(reg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// -------------------------------------------------------------------

type LittleEndianHighWordAllCustomTestSuite struct {
	suite.Suite
	LittleEndianHighWordAllCustom tests.LittleEndianHighWordAllCustom
}

func (s *LittleEndianHighWordAllCustomTestSuite) SetupTest() {
	s.LittleEndianHighWordAllCustom = tests.LittleEndianHighWordAllCustom{
		Ignored:       0x4321,
		Uint8High:     0x11,
		Uint8Low:      0x22,
		ByteHigh:      0x33,
		ByteLow:       0x44,
		Int8Low:       -50,
		Int8High:      -100,
		Uint16:        0x1234,
		Int16:         -42,
		Uint32:        0x12345678,
		Int32:         -111222,
		Uint64:        0x1234123412341234,
		Int64:         -1111222233334444,
		Float32:       1234.5678,
		Float64:       -1234.5678,
		StringASCII8:  "TestData1",
		StringASCII16: "TestData2",
		StringUTF816:  "TestDäta3", // add utf8 char
		StringASCII32: "Te",
		StringUTF832:  "T😀",
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
		Uint16Array:   [5]tests.CustomUint16{0x1111, 0x2222, 0x3333, 0x4444, 0x5555},
		Uint32Array: [5]tests.CustomUint32{
			0x11111111,
			0x22222222,
			0x33333333,
			0x44444444,
			0x55555555,
		},
		Float32Array: [5]tests.CustomFloat32{1.1, 2.2, 3.3, 4.4, 5.5},
		BoolArray:    [5]tests.CustomBool{true, false, true, false, true},
		BitField16:   0x1234,
	}
}

func (s *LittleEndianHighWordAllCustomTestSuite) TestMarshalUnmarshal() {
	reg, err := s.LittleEndianHighWordAllCustom.Marshal()
	s.Require().NoError(err)

	out := &tests.LittleEndianHighWordAllCustom{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Empty(out.Ignored)

	// Ignored field is not set by Unmarshal, set it manually for comparison
	out.Ignored = s.LittleEndianHighWordAllCustom.Ignored

	s.Equal(s.LittleEndianHighWordAllCustom, *out)
}

func BenchmarkLittleEndianHighWordAllCustomMarshal(b *testing.B) {
	test := &tests.LittleEndianHighWordAllCustom{
		Ignored:       0x4321,
		Uint8High:     0x11,
		Uint8Low:      0x22,
		ByteHigh:      0x33,
		ByteLow:       0x44,
		Int8Low:       -50,
		Int8High:      -100,
		Uint16:        0x1234,
		Int16:         -42,
		Uint32:        0x12345678,
		Int32:         -111222,
		Uint64:        0x1234123412341234,
		Int64:         -1111222233334444,
		Float32:       1234.5678,
		Float64:       -1234.5678,
		StringASCII8:  "TestData1",
		StringASCII16: "TestData2",
		StringUTF816:  "TestDäta3", // add utf8 char
		StringASCII32: "Te",
		StringUTF832:  "T😀",
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
		Uint16Array:   [5]tests.CustomUint16{0x1111, 0x2222, 0x3333, 0x4444, 0x5555},
		Uint32Array: [5]tests.CustomUint32{
			0x11111111,
			0x22222222,
			0x33333333,
			0x44444444,
			0x55555555,
		},
		Float32Array: [5]tests.CustomFloat32{1.1, 2.2, 3.3, 4.4, 5.5},
		BoolArray:    [5]tests.CustomBool{true, false, true, false, true},
		BitField16:   0x1234,
	}

	for b.Loop() {
		_, err := test.Marshal()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLittleEndianHighWordAllCustomUnmarshal(b *testing.B) {
	test := &tests.LittleEndianHighWordAllCustom{
		Ignored:       0x4321,
		Uint8High:     0x11,
		Uint8Low:      0x22,
		ByteHigh:      0x33,
		ByteLow:       0x44,
		Int8Low:       -50,
		Int8High:      -100,
		Uint16:        0x1234,
		Int16:         -42,
		Uint32:        0x12345678,
		Int32:         -111222,
		Uint64:        0x1234123412341234,
		Int64:         -1111222233334444,
		Float32:       1234.5678,
		Float64:       -1234.5678,
		StringASCII8:  "TestData1",
		StringASCII16: "TestData2",
		StringUTF816:  "TestDäta3", // add utf8 char
		StringASCII32: "Te",
		StringUTF832:  "T😀",
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
		Uint16Array:   [5]tests.CustomUint16{0x1111, 0x2222, 0x3333, 0x4444, 0x5555},
		Uint32Array: [5]tests.CustomUint32{
			0x11111111,
			0x22222222,
			0x33333333,
			0x44444444,
			0x55555555,
		},
		Float32Array: [5]tests.CustomFloat32{1.1, 2.2, 3.3, 4.4, 5.5},
		BoolArray:    [5]tests.CustomBool{true, false, true, false, true},
		BitField16:   0x1234,
	}

	reg, err := test.Marshal()
	if err != nil {
		b.Fatal(err)
	}

	test = &tests.LittleEndianHighWordAllCustom{}

	for b.Loop() {
		err := test.Unmarshal(reg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// -------------------------------------------------------------------

type LittleEndianHighWordAllCustomExternTestSuite struct {
	suite.Suite
	LittleEndianHighWordAllCustomExtern tests.LittleEndianHighWordAllCustomExtern
}

func (s *LittleEndianHighWordAllCustomExternTestSuite) SetupTest() {
	s.LittleEndianHighWordAllCustomExtern = tests.LittleEndianHighWordAllCustomExtern{
		Ignored:       0x4321,
		Uint8High:     0x11,
		Uint8Low:      0x22,
		ByteHigh:      0x33,
		ByteLow:       0x44,
		Int8Low:       -50,
		Int8High:      -100,
		Uint16:        0x1234,
		Int16:         -42,
		Uint32:        0x12345678,
		Int32:         -111222,
		Uint64:        0x1234123412341234,
		Int64:         -1111222233334444,
		Float32:       1234.5678,
		Float64:       -1234.5678,
		StringASCII8:  "TestData1",
		StringASCII16: "TestData2",
		StringUTF816:  "TestDäta3", // add utf8 char
		StringASCII32: "Te",
		StringUTF832:  "T😀",
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
		Uint16Array:   [5]extern.CustomUint16{0x1111, 0x2222, 0x3333, 0x4444, 0x5555},
		Uint32Array: [5]extern.CustomUint32{
			0x11111111,
			0x22222222,
			0x33333333,
			0x44444444,
			0x55555555,
		},
		Float32Array: [5]extern.CustomFloat32{1.1, 2.2, 3.3, 4.4, 5.5},
		BoolArray:    [5]extern.CustomBool{true, false, true, false, true},
		BitField16:   0x1234,
	}
}

func (s *LittleEndianHighWordAllCustomExternTestSuite) TestMarshalUnmarshal() {
	reg, err := s.LittleEndianHighWordAllCustomExtern.Marshal()
	s.Require().NoError(err)

	out := &tests.LittleEndianHighWordAllCustomExtern{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Empty(out.Ignored)

	// Ignored field is not set by Unmarshal, set it manually for comparison
	out.Ignored = s.LittleEndianHighWordAllCustomExtern.Ignored

	s.Equal(s.LittleEndianHighWordAllCustomExtern, *out)
}

func BenchmarkLittleEndianHighWordAllCustomExternMarshal(b *testing.B) {
	test := &tests.LittleEndianHighWordAllCustomExtern{
		Ignored:       0x4321,
		Uint8High:     0x11,
		Uint8Low:      0x22,
		ByteHigh:      0x33,
		ByteLow:       0x44,
		Int8Low:       -50,
		Int8High:      -100,
		Uint16:        0x1234,
		Int16:         -42,
		Uint32:        0x12345678,
		Int32:         -111222,
		Uint64:        0x1234123412341234,
		Int64:         -1111222233334444,
		Float32:       1234.5678,
		Float64:       -1234.5678,
		StringASCII8:  "TestData1",
		StringASCII16: "TestData2",
		StringUTF816:  "TestDäta3", // add utf8 char
		StringASCII32: "Te",
		StringUTF832:  "T😀",
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
		Uint16Array:   [5]extern.CustomUint16{0x1111, 0x2222, 0x3333, 0x4444, 0x5555},
		Uint32Array: [5]extern.CustomUint32{
			0x11111111,
			0x22222222,
			0x33333333,
			0x44444444,
			0x55555555,
		},
		Float32Array: [5]extern.CustomFloat32{1.1, 2.2, 3.3, 4.4, 5.5},
		BoolArray:    [5]extern.CustomBool{true, false, true, false, true},
		BitField16:   0x1234,
	}

	for b.Loop() {
		_, err := test.Marshal()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLittleEndianHighWordAllCustomExternUnmarshal(b *testing.B) {
	test := &tests.LittleEndianHighWordAllCustomExtern{
		Ignored:       0x4321,
		Uint8High:     0x11,
		Uint8Low:      0x22,
		ByteHigh:      0x33,
		ByteLow:       0x44,
		Int8Low:       -50,
		Int8High:      -100,
		Uint16:        0x1234,
		Int16:         -42,
		Uint32:        0x12345678,
		Int32:         -111222,
		Uint64:        0x1234123412341234,
		Int64:         -1111222233334444,
		Float32:       1234.5678,
		Float64:       -1234.5678,
		StringASCII8:  "TestData1",
		StringASCII16: "TestData2",
		StringUTF816:  "TestDäta3", // add utf8 char
		StringASCII32: "Te",
		StringUTF832:  "T😀",
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
		Uint16Array:   [5]extern.CustomUint16{0x1111, 0x2222, 0x3333, 0x4444, 0x5555},
		Uint32Array: [5]extern.CustomUint32{
			0x11111111,
			0x22222222,
			0x33333333,
			0x44444444,
			0x55555555,
		},
		Float32Array: [5]extern.CustomFloat32{1.1, 2.2, 3.3, 4.4, 5.5},
		BoolArray:    [5]extern.CustomBool{true, false, true, false, true},
		BitField16:   0x1234,
	}

	reg, err := test.Marshal()
	if err != nil {
		b.Fatal(err)
	}

	test = &tests.LittleEndianHighWordAllCustomExtern{}

	for b.Loop() {
		err := test.Unmarshal(reg)
		if err != nil {
			b.Fatal(err)
		}
	}
}
