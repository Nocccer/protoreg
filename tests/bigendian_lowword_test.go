package tests_test

import (
	"testing"

	"github.com/Nocccer/protoreg/tests"
	"github.com/stretchr/testify/suite"
)

func TestBigEndianLowWord(t *testing.T) {
	suite.Run(t, new(BigEndianLowWordTestSuite))
	suite.Run(t, new(BigEndianLowWordAllCustomTestSuite))
	suite.Run(t, new(BigEndianLowWordAllCustomExternTestSuite))
}

type BigEndianLowWordTestSuite struct {
	suite.Suite
	BigEndianLowWord tests.BigEndianHighWord
}

func (s *BigEndianLowWordTestSuite) SetupTest() {
	s.BigEndianLowWord = tests.BigEndianHighWord{
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

func (s *BigEndianLowWordTestSuite) TestMarshalUnmarshal() {
	reg, err := s.BigEndianLowWord.Marshal()
	s.Require().NoError(err)

	out := &tests.BigEndianHighWord{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Empty(out.Ignored)

	out.Ignored = s.BigEndianLowWord.Ignored // Ignored field is not set by Unmarshal, set it manually for comparison

	s.Equal(s.BigEndianLowWord, *out)
}

func BenchmarkBigEndianLowWordMarshal(b *testing.B) {
	test := &tests.BigEndianLowWord{
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

func BenchmarkBigEndianLowWordUnmarshal(b *testing.B) {
	test := &tests.BigEndianLowWord{
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

	test = &tests.BigEndianLowWord{}

	for b.Loop() {
		err := test.Unmarshal(reg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ---------------------------------------------------------------

type BigEndianLowWordAllCustomTestSuite struct {
	suite.Suite
	BigEndianLowWordAllCustom tests.BigEndianHighWordAllCustom
}

func (s *BigEndianLowWordAllCustomTestSuite) SetupTest() {
	s.BigEndianLowWordAllCustom = tests.BigEndianHighWordAllCustom{
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

func (s *BigEndianLowWordAllCustomTestSuite) TestMarshalUnmarshal() {
	reg, err := s.BigEndianLowWordAllCustom.Marshal()
	s.Require().NoError(err)

	out := &tests.BigEndianHighWordAllCustom{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Empty(out.Ignored)

	// Ignored field is not set by Unmarshal, set it manually for comparison
	out.Ignored = s.BigEndianLowWordAllCustom.Ignored

	s.Equal(s.BigEndianLowWordAllCustom, *out)
}

func BenchmarkBigEndianLowWordAllCustomMarshal(b *testing.B) {
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

func BenchmarkBigEndianLowWordAllCustomUnmarshal(b *testing.B) {
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

// ---------------------------------------------------------------

type BigEndianLowWordAllCustomExternTestSuite struct {
	suite.Suite
	BigEndianLowWordAllCustomExtern tests.BigEndianHighWordAllCustomExtern
}

func (s *BigEndianLowWordAllCustomExternTestSuite) SetupTest() {
	s.BigEndianLowWordAllCustomExtern = tests.BigEndianHighWordAllCustomExtern{
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

func (s *BigEndianLowWordAllCustomExternTestSuite) TestMarshalUnmarshal() {
	reg, err := s.BigEndianLowWordAllCustomExtern.Marshal()
	s.Require().NoError(err)

	out := &tests.BigEndianHighWordAllCustomExtern{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Empty(out.Ignored)

	// Ignored field is not set by Unmarshal, set it manually for comparison
	out.Ignored = s.BigEndianLowWordAllCustomExtern.Ignored

	s.Equal(s.BigEndianLowWordAllCustomExtern, *out)
}

func BenchmarkBigEndianLowWordAllCustomExternMarshal(b *testing.B) {
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

func BenchmarkBigEndianLowWordAllCustomExternUnmarshal(b *testing.B) {
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
