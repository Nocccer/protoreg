package tests_test

import (
	"testing"

	"github.com/Nocccer/protoreg/tests"
	"github.com/stretchr/testify/suite"
)

func TestBigEndianHighWord(t *testing.T) {
	suite.Run(t, new(BigEndianHighWordTestSuite))
	suite.Run(t, new(BigEndianHighWordAllCustomTestSuite))
	suite.Run(t, new(BigEndianHighWordAllCustomExternTestSuite))
}

type BigEndianHighWordTestSuite struct {
	suite.Suite
	BigEndianHighWord tests.BigEndianHighWord
}

func (s *BigEndianHighWordTestSuite) SetupTest() {
	s.BigEndianHighWord = tests.BigEndianHighWord{
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

func (s *BigEndianHighWordTestSuite) TestMarshalUnmarshal() {
	reg, err := s.BigEndianHighWord.Marshal()
	s.Require().NoError(err)

	out := &tests.BigEndianHighWord{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Empty(out.Ignored)

	out.Ignored = s.BigEndianHighWord.Ignored // Ignored field is not set by Unmarshal, set it manually for comparison

	s.Equal(s.BigEndianHighWord, *out)
}

func BenchmarkBigEndianHighWordMarshal(b *testing.B) {
	test := &tests.BigEndianHighWord{
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

func BenchmarkBigEndianHighWordUnmarshal(b *testing.B) {
	test := &tests.BigEndianHighWord{
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

	test = &tests.BigEndianHighWord{}

	for b.Loop() {
		err := test.Unmarshal(reg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type BigEndianHighWordAllCustomTestSuite struct {
	suite.Suite
	BigEndianHighWordAllCustom tests.BigEndianHighWordAllCustom
}

func (s *BigEndianHighWordAllCustomTestSuite) SetupTest() {
	s.BigEndianHighWordAllCustom = tests.BigEndianHighWordAllCustom{
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

// ----------------------------------------------------------------------

func (s *BigEndianHighWordAllCustomTestSuite) TestMarshalUnmarshal() {
	reg, err := s.BigEndianHighWordAllCustom.Marshal()
	s.Require().NoError(err)

	out := &tests.BigEndianHighWordAllCustom{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Empty(out.Ignored)

	// Ignored field is not set by Unmarshal, set it manually for comparison
	out.Ignored = s.BigEndianHighWordAllCustom.Ignored

	s.Equal(s.BigEndianHighWordAllCustom, *out)
}

func BenchmarkBigEndianHighWordAllCustomMarshal(b *testing.B) {
	test := &tests.BigEndianHighWordAllCustom{
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

func BenchmarkBigEndianHighWordAllCustomUnmarshal(b *testing.B) {
	test := &tests.BigEndianHighWordAllCustom{
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

	test = &tests.BigEndianHighWordAllCustom{}

	for b.Loop() {
		err := test.Unmarshal(reg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ----------------------------------------------------------------------

type BigEndianHighWordAllCustomExternTestSuite struct {
	suite.Suite
	BigEndianHighWordAllCustomExtern tests.BigEndianHighWordAllCustomExtern
}

func (s *BigEndianHighWordAllCustomExternTestSuite) SetupTest() {
	s.BigEndianHighWordAllCustomExtern = tests.BigEndianHighWordAllCustomExtern{
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

func (s *BigEndianHighWordAllCustomExternTestSuite) TestMarshalUnmarshal() {
	reg, err := s.BigEndianHighWordAllCustomExtern.Marshal()
	s.Require().NoError(err)

	out := &tests.BigEndianHighWordAllCustomExtern{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Empty(out.Ignored)

	// Ignored field is not set by Unmarshal, set it manually for comparison
	out.Ignored = s.BigEndianHighWordAllCustomExtern.Ignored

	s.Equal(s.BigEndianHighWordAllCustomExtern, *out)
}

func BenchmarkBigEndianHighWordAllCustomExternMarshal(b *testing.B) {
	test := &tests.BigEndianHighWordAllCustomExtern{
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

func BenchmarkBigEndianHighWordAllCustomExternUnmarshal(b *testing.B) {
	test := &tests.BigEndianHighWordAllCustomExtern{
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

	test = &tests.BigEndianHighWordAllCustomExtern{}

	for b.Loop() {
		err := test.Unmarshal(reg)
		if err != nil {
			b.Fatal(err)
		}
	}
}
