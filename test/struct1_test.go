package test

import (
	"testing"

	"github.com/Nocccer/protoreg/test/sub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshal(t *testing.T) {
	in := &Struct1{
		Ignored: 0,
		Field1:  "TestDat",
		Field2:  -42,
		Field3:  CustomUint16(123),
		Field4:  456,
		Field5:  "ASCII",
		Field6:  sub.CustomInt16(-789),
	}

	reg, err := in.Marshal()
	require.NoError(t, err)

	out := &Struct1{}
	err = out.Unmarshal(reg)
	require.NoError(t, err)

	assert.Equal(t, in, out)
}

func BenchmarkMarshal(b *testing.B) {
	test := &Struct1{
		Ignored: 0,
		Field1:  "Test Data",
		Field2:  -42,
		Field3:  CustomUint16(123),
		Field4:  456,
		Field5:  "ASCII",
		Field6:  sub.CustomInt16(-789),
	}

	for b.Loop() {
		_, err := test.Marshal()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	reg := []uint16{
		84,
		101,
		115,
		116,
		32,
		68,
		97,
		116,
		65494,
		123,
		456,
		16723,
		17225,
		18688,
		0,
		64747,
		0,
		0,
		0,
	}

	test := &Struct1{}

	for b.Loop() {
		err := test.Unmarshal(reg)
		if err != nil {
			b.Fatal(err)
		}
	}
}
