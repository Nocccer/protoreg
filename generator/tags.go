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

type CharEncoding string

const (
	CharEncodingASCII CharEncoding = "ascii" // default
	CharEncodingUTF8  CharEncoding = "utf8"
)

func (c CharEncoding) Validate() error {
	switch c {
	case CharEncodingASCII, CharEncodingUTF8:
		return nil
	}
	return fmt.Errorf(`invalid "charencoding" value %q`, c)
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
	Encoding     *Encoding
	WordOrder    *WordOrder
	Offset       *int
	Size         *int
	Char         *Char
	CharEncoding *CharEncoding
	Byte         *Byte
	Bit          *int
}

func (t Tags) DeepCopy() Tags {
	var copy Tags

	if t.Encoding != nil {
		copy.Encoding = new(*t.Encoding)
	}
	if t.WordOrder != nil {
		copy.WordOrder = new(*t.WordOrder)
	}
	if t.Offset != nil {
		copy.Offset = new(*t.Offset)
	}
	if t.Size != nil {
		copy.Size = new(*t.Size)
	}
	if t.Char != nil {
		copy.Char = new(*t.Char)
	}
	if t.CharEncoding != nil {
		copy.CharEncoding = new(*t.CharEncoding)
	}
	if t.Byte != nil {
		copy.Byte = new(*t.Byte)
	}
	if t.Bit != nil {
		copy.Bit = new(*t.Bit)
	}

	return copy
}

func extractTags(tagStr string) (Tags, error) {
	var t Tags
	tags := strings.SplitSeq(tagStr, ",")

	for tag := range tags {
		kv := strings.Split(tag, "=")

		if len(kv) != 2 {
			return Tags{}, fmt.Errorf("invalid tag format: %q", tag)
		}

		switch kv[0] {
		case "encoding":
			t.Encoding = new(Encoding(kv[1]))
			if err := t.Encoding.Validate(); err != nil {
				return Tags{}, err
			}
		case "wordorder":
			t.WordOrder = new(WordOrder(kv[1]))
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
			t.Char = new(Char(kv[1]))
			if err := t.Char.Validate(); err != nil {
				return Tags{}, err
			}
		case "charencoding":
			t.CharEncoding = new(CharEncoding(kv[1]))
			if err := t.CharEncoding.Validate(); err != nil {
				return Tags{}, err
			}
		case "byte":
			t.Byte = new(Byte(kv[1]))
			if err := t.Byte.Validate(); err != nil {
				return Tags{}, err
			}
		case "bit":
			bit, err := strconv.Atoi(kv[1])
			if err != nil || bit < 0 || bit > 15 {
				return Tags{}, fmt.Errorf("invalid bit value: %q", kv[1])
			}
			t.Bit = &bit
		}
	}

	return t, nil
}
