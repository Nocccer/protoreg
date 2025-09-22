package tests_test

import (
	"testing"

	"github.com/Nocccer/protoreg/tests"
	"github.com/Nocccer/protoreg/tests/sub"
	"github.com/stretchr/testify/suite"
)

func TestBigEndianHighWord(t *testing.T) {
	suite.Run(t, new(BigEndianHighWordTestSuite))
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
		CustomUint16:  tests.CustomUint16(123),
		CustomInt16:   sub.CustomInt16(-789),
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
		CustomUint16:  tests.CustomUint16(123),
		CustomInt16:   sub.CustomInt16(-789),
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
		CustomUint16:  tests.CustomUint16(123),
		CustomInt16:   sub.CustomInt16(-789),
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
