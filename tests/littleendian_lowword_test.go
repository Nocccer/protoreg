package tests_test

import (
	"testing"

	"github.com/Nocccer/protoreg/tests"
	"github.com/stretchr/testify/suite"
)

func TestLittleEndianLowWord(t *testing.T) {
	suite.Run(t, new(LittleEndianLowWordTestSuite))
	suite.Run(t, new(LittleEndianLowWordAllCustomTestSuite))
	suite.Run(t, new(LittleEndianLowWordAllCustomExternTestSuite))
}

type LittleEndianLowWordTestSuite struct {
	suite.Suite
	LittleEndianLowWord tests.LittleEndianLowWord
}

func (s *LittleEndianLowWordTestSuite) SetupTest() {
	s.LittleEndianLowWord = tests.LittleEndianLowWord{
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
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
	}
}

func (s *LittleEndianLowWordTestSuite) TestMarshalUnmarshal() {
	reg, err := s.LittleEndianLowWord.Marshal()
	s.Require().NoError(err)

	out := &tests.LittleEndianLowWord{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Empty(out.Ignored)

	out.Ignored = s.LittleEndianLowWord.Ignored // Ignored field is not set by Unmarshal, set it manually for comparison

	s.Equal(s.LittleEndianLowWord, *out)
}

func BenchmarkLittleEndianLowWordMarshal(b *testing.B) {
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
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
	}

	for b.Loop() {
		_, err := test.Marshal()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLittleEndianLowWordUnmarshal(b *testing.B) {
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
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
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

// -----------------------------------------------------------------------

type LittleEndianLowWordAllCustomTestSuite struct {
	suite.Suite
	LittleEndianLowWordAllCustom tests.LittleEndianLowWordAllCustom
}

func (s *LittleEndianLowWordAllCustomTestSuite) SetupTest() {
	s.LittleEndianLowWordAllCustom = tests.LittleEndianLowWordAllCustom{
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
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
	}
}

func (s *LittleEndianLowWordAllCustomTestSuite) TestMarshalUnmarshal() {
	reg, err := s.LittleEndianLowWordAllCustom.Marshal()
	s.Require().NoError(err)

	out := &tests.LittleEndianLowWordAllCustom{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Empty(out.Ignored)

	// Ignored field is not set by Unmarshal, set it manually for comparison
	out.Ignored = s.LittleEndianLowWordAllCustom.Ignored

	s.Equal(s.LittleEndianLowWordAllCustom, *out)
}

func BenchmarkLittleEndianLowWordAllCustomMarshal(b *testing.B) {
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
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
	}

	for b.Loop() {
		_, err := test.Marshal()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLittleEndianLowWordAllCustomUnmarshal(b *testing.B) {
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
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
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

// -----------------------------------------------------------------------

type LittleEndianLowWordAllCustomExternTestSuite struct {
	suite.Suite
	LittleEndianLowWordAllCustomExtern tests.LittleEndianLowWordAllCustomExtern
}

func (s *LittleEndianLowWordAllCustomExternTestSuite) SetupTest() {
	s.LittleEndianLowWordAllCustomExtern = tests.LittleEndianLowWordAllCustomExtern{
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
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
	}
}

func (s *LittleEndianLowWordAllCustomExternTestSuite) TestMarshalUnmarshal() {
	reg, err := s.LittleEndianLowWordAllCustomExtern.Marshal()
	s.Require().NoError(err)

	out := &tests.LittleEndianLowWordAllCustomExtern{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Empty(out.Ignored)

	// Ignored field is not set by Unmarshal, set it manually for comparison
	out.Ignored = s.LittleEndianLowWordAllCustomExtern.Ignored

	s.Equal(s.LittleEndianLowWordAllCustomExtern, *out)
}

func BenchmarkLittleEndianLowWordAllCustomExternMarshal(b *testing.B) {
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
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
	}

	for b.Loop() {
		_, err := test.Marshal()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLittleEndianLowWordAllCustomExternUnmarshal(b *testing.B) {
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
		Bool:          true,
		Bit1:          false,
		Bit3:          false,
		Bit14:         true,
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
