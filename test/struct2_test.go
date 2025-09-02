package test

import (
	"testing"

	"github.com/Nocccer/protoreg/test/sub"
	"github.com/stretchr/testify/suite"
)

func TestStruct2(t *testing.T) {
	suite.Run(t, new(Struct2TestSuite))
}

type Struct2TestSuite struct {
	suite.Suite
	struct2 Struct2
}

func (s *Struct2TestSuite) SetupTest() {
	s.struct2 = Struct2{
		Ignored: 0,
		Field1:  "TestDat",
		Field2:  -42,
		Field3:  CustomUint16(123),
		Field4:  456,
		Field5:  "ASCII",
		Field6:  sub.CustomInt16(-789),
		Field7:  12345,
	}
}

func (s *Struct2TestSuite) TestMarshalUnmarshal() {
	reg, err := s.struct2.Marshal()
	s.Require().NoError(err)

	out := &Struct2{}
	err = out.Unmarshal(reg)
	s.Require().NoError(err)

	s.Equal(s.struct2, *out)
}

func BenchmarkStruct2Marshal(b *testing.B) {
	test := &Struct2{
		Ignored: 0,
		Field1:  "Test Data",
		Field2:  -42,
		Field3:  CustomUint16(123),
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

func BenchmarkStruct2Unmarshal(b *testing.B) {
	test := &Struct2{
		Ignored: 0,
		Field1:  "Test Data",
		Field2:  -42,
		Field3:  CustomUint16(123),
		Field4:  456,
		Field5:  "ASCII",
		Field6:  sub.CustomInt16(-789),
		Field7:  12345,
	}

	reg, err := test.Marshal()
	if err != nil {
		b.Fatal(err)
	}

	test = &Struct2{}

	for b.Loop() {
		err := test.Unmarshal(reg)
		if err != nil {
			b.Fatal(err)
		}
	}
}
