package generator

import (
	"fmt"
	"strconv"
	"strings"
)

type Byte string

const (
	High Byte = "high"
	Low  Byte = "low"
)

func (b Byte) Validate() error {
	switch b {
	case High, Low:
		return nil
	}
	return fmt.Errorf(`invalid "byte" value %q`, b)
}

type Char string

const (
	Char8  Char = "8"
	Char16 Char = "16"
)

func (c Char) Validate() error {
	switch c {
	case Char8, Char16:
		return nil
	}
	return fmt.Errorf(`invalid "char" value %q`, c)
}

type Encoding string

const (
	BigEndian    Encoding = "big"
	LittleEndian Encoding = "little"
)

func (e Encoding) Validate() error {
	switch e {
	case BigEndian, LittleEndian:
		return nil
	}
	return fmt.Errorf(`invalid "encoding" value %q`, e)
}

type WordOrder string

const (
	HighWordFirst WordOrder = "high"
	LowWordFirst  WordOrder = "low"
)

func (w WordOrder) Validate() error {
	switch w {
	case HighWordFirst, LowWordFirst:
		return nil
	}
	return fmt.Errorf(`invalid "wordorder" value %q`, w)
}

type Tags struct {
	Encoding  *Encoding
	WordOrder *WordOrder
	Offset    *int
	Size      *int
	Char      *Char
	Byte      *Byte
	Bit       *int
}

func extractTags(tagStr string) (Tags, error) {
	var t Tags
	tags := strings.Split(tagStr, ",")

	for _, tag := range tags {
		kv := strings.Split(tag, "=")

		if len(kv) != 2 {
			return Tags{}, fmt.Errorf("invalid tag format: %q", tag)
		}

		switch kv[0] {
		case "encoding":
			t.Encoding = ptrTo(Encoding(kv[1]))
			if err := t.Encoding.Validate(); err != nil {
				return Tags{}, err
			}
		case "wordorder":
			t.WordOrder = ptrTo(WordOrder(kv[1]))
			if err := t.WordOrder.Validate(); err != nil {
				return Tags{}, err
			}
		case "offset":
			offset, err := strconv.Atoi(kv[1])
			if err != nil || offset < 0 {
				return Tags{}, fmt.Errorf("invalid offset value: %q", kv[1])
			}
			t.Offset = &offset
		case "size":
			size, err := strconv.Atoi(kv[1])
			if err != nil || size <= 0 {
				return Tags{}, fmt.Errorf("invalid size value: %q", kv[1])
			}
			t.Size = &size
		case "char":
			t.Char = ptrTo(Char(kv[1]))
			if err := t.Char.Validate(); err != nil {
				return Tags{}, err
			}
		case "byte":
			t.Byte = ptrTo(Byte(kv[1]))
			if err := t.Byte.Validate(); err != nil {
				return Tags{}, err
			}
		case "bit":
			bit, err := strconv.Atoi(kv[1])
			if err != nil || bit < 0 {
				return Tags{}, fmt.Errorf("invalid bit value: %q", kv[1])
			}
			t.Bit = &bit
		}
	}

	return t, nil
}
