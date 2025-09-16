package tests_test

import (
	"testing"

	"github.com/Nocccer/protoreg/tests"
	"github.com/Nocccer/protoreg/tests/sub"
	"github.com/stretchr/testify/suite"
)

func TestLittleEndianHighWord(t *testing.T) {
	suite.Run(t, new(LittleEndianHighWordTestSuite))
}

type LittleEndianHighWordTestSuite struct {
	suite.Suite
	LittleEndianHighWord tests.LittleEndianHighWord
}

func (s *LittleEndianHighWordTestSuite) SetupTest() {
	s.LittleEndianHighWord = tests.LittleEndianHighWord{
		Ignored: 0,
		Field1:  "TestDat",
		Field2:  -42,
		Field3:  tests.CustomUint16(123),
		Field4:  456,
		Field5:  "ASCII",
		Field6:  sub.CustomInt16(-789),
		Field7:  12345,
	}
}

func (s *LittleEndianHighWordTestSuite) TestMarshalUnmarshal() {
	reg, err := s.LittleEndianHighWord.Marshal()
	s.Require().NoError(err)

	out := &tests.LittleEndianHighWord{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Equal(s.LittleEndianHighWord, *out)
}

func BenchmarkLittleEndianHighWordMarshal(b *testing.B) {
	test := &tests.LittleEndianHighWord{
		Ignored: 0,
		Field1:  "Test Data",
		Field2:  -42,
		Field3:  tests.CustomUint16(123),
		Field4:  456,
		Field5:  "ASCII",
		Field6:  sub.CustomInt16(-789),
		Field7:  12345,
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
		Ignored: 0,
		Field1:  "Test Data",
		Field2:  -42,
		Field3:  tests.CustomUint16(123),
		Field4:  456,
		Field5:  "ASCII",
		Field6:  sub.CustomInt16(-789),
		Field7:  12345,
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
