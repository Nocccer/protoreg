package generator

import (
	"fmt"
)

type Field struct {
	Name         string
	Tags         Tags
	IsCustomType bool
	CustomType   string
}

func (f Field) Comment() string {
	return fmt.Sprintf("\t// %s\n", f.Name)
}

// calcWordOffsets32 returns offsets for low and high word based on WordOrder
// Returns [lowWordOffset, highWordOffset]
func (f Field) calcWordOffsets32() [2]int {
	low := *f.Tags.Offset
	high := *f.Tags.Offset + 1
	if *f.Tags.WordOrder == LowWordFirst {
		return [2]int{low, high}
	}
	// HighWordFirst (default)
	return [2]int{high, low}
}

// calcWordOffsets64 returns offsets for all 4 words based on WordOrder
// Returns [word0Offset, word1Offset, word2Offset, word3Offset] in correct order
func (f Field) calcWordOffsets64() [4]int {
	offsets := [4]int{
		*f.Tags.Offset,
		*f.Tags.Offset + 1,
		*f.Tags.Offset + 2,
		*f.Tags.Offset + 3,
	}

	if *f.Tags.WordOrder == LowWordFirst {
		return offsets
	}
	// HighWordFirst: reverse order
	return [4]int{offsets[3], offsets[2], offsets[1], offsets[0]}
}

// encodeWord16 returns encoding logic for a 16-bit word
func (f Field) encodeWord16(value string, shift string) string {
	switch *f.Tags.Encoding {
	case LittleEndian:
		return fmt.Sprintf("bits.ReverseBytes16(uint16(%s%s))", value, shift)
	default:
		return fmt.Sprintf("uint16(%s%s)", value, shift)
	}
}

// decodeWord16 returns decoding logic for a 16-bit word
func (f Field) decodeWord16(bufIndex int) string {
	switch *f.Tags.Encoding {
	case LittleEndian:
		return fmt.Sprintf("bits.ReverseBytes16(buf[%d])", bufIndex)
	default:
		return fmt.Sprintf("buf[%d]", bufIndex)
	}
}
