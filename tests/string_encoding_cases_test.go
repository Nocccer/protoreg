package tests_test

import (
	"testing"

	"github.com/nocccer/protoreg/tests"
	"github.com/stretchr/testify/require"
)

func TestStringChar8EndianEncoding(t *testing.T) {
	const input = "ABCD"

	t.Run("big endian", func(t *testing.T) {
		in := tests.StringChar8BigEndian{Value: input}
		buf, err := in.Marshal()
		require.NoError(t, err)
		require.Equal(t, []uint16{0x4142, 0x4344}, buf)

		var out tests.StringChar8BigEndian
		require.NoError(t, out.Unmarshal(buf))
		require.Equal(t, input, out.Value)
	})

	t.Run("little endian", func(t *testing.T) {
		in := tests.StringChar8LittleEndian{Value: input}
		buf, err := in.Marshal()
		require.NoError(t, err)
		require.Equal(t, []uint16{0x4241, 0x4443}, buf)

		var out tests.StringChar8LittleEndian
		require.NoError(t, out.Unmarshal(buf))
		require.Equal(t, input, out.Value)
	})
}

func TestStringUTF16OneRegisterPerRune(t *testing.T) {
	const input = "AäB"

	t.Run("big endian", func(t *testing.T) {
		in := tests.StringUTF16BigEndian{Value: input}
		buf, err := in.Marshal()
		require.NoError(t, err)
		require.Equal(t, []uint16{0x0041, 0x00E4, 0x0042, 0x0000, 0x0000}, buf)

		var out tests.StringUTF16BigEndian
		require.NoError(t, out.Unmarshal(buf))
		require.Equal(t, input, out.Value)
	})

	t.Run("little endian", func(t *testing.T) {
		in := tests.StringUTF16LittleEndian{Value: input}
		buf, err := in.Marshal()
		require.NoError(t, err)
		require.Equal(t, []uint16{0x4100, 0xE400, 0x4200, 0x0000, 0x0000}, buf)

		var out tests.StringUTF16LittleEndian
		require.NoError(t, out.Unmarshal(buf))
		require.Equal(t, input, out.Value)
	})
}

func TestStringUTF32EncodingAndWordOrder(t *testing.T) {
	const input = "😀A"

	t.Run("big endian high word", func(t *testing.T) {
		in := tests.StringUTF32BigHighWord{Value: input}
		buf, err := in.Marshal()
		require.NoError(t, err)
		require.Equal(t, []uint16{0x0001, 0xF600, 0x0000, 0x0041}, buf)

		var out tests.StringUTF32BigHighWord
		require.NoError(t, out.Unmarshal(buf))
		require.Equal(t, input, out.Value)
	})

	t.Run("little endian low word", func(t *testing.T) {
		in := tests.StringUTF32LittleLowWord{Value: input}
		buf, err := in.Marshal()
		require.NoError(t, err)
		require.Equal(t, []uint16{0x00F6, 0x0100, 0x4100, 0x0000}, buf)

		var out tests.StringUTF32LittleLowWord
		require.NoError(t, out.Unmarshal(buf))
		require.Equal(t, input, out.Value)
	})
}
