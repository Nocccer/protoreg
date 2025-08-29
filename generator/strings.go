package generator

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

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
	return errors.New(`invalid "char" value`)
}

type FieldString struct {
	Field
	Char Char
}

func (f FieldString) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch f.Char {
	case Char8:
		sb.WriteString(fmt.Sprintf("\tlength := len(m.%s)\n", f.Name))
		sb.WriteString("\tif length % 2 != 0 { length += 1 }\n")
		sb.WriteString("\tbytes := make([]byte, length)\n")
		sb.WriteString(fmt.Sprintf("\tcopy(bytes, m.%s)\n", f.Name))
		sb.WriteString("\tfor i := 0; i < length; i+=2 {\n")
		sb.WriteString(fmt.Sprintf("\t\tif i >= %d {break}\n", f.Size*2))
		sb.WriteString(
			fmt.Sprintf(
				"\t\tbuf[%d+i/2] = uint16(bytes[i+1]) | uint16(bytes[i])<<8\n",
				f.Offset,
			),
		)
		sb.WriteString("\t}\n")
	case Char16:
		sb.WriteString(fmt.Sprintf("\tfor i := 0; i < len(m.%s); i++ {\n", f.Name))
		sb.WriteString(fmt.Sprintf("\t\tif i >= %d {break}\n", f.Size))
		sb.WriteString(
			fmt.Sprintf(
				"\t\tbuf[%d+i] = uint16(m.%s[i])\n",
				f.Offset,
				f.Name,
			),
		)
		sb.WriteString("\t}\n")
	}

	return sb.String()
}

func (f FieldString) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch f.Char {
	case Char8:
		sb.WriteString(fmt.Sprintf("\tbytes = make([]byte, %d)\n", f.Size))
		sb.WriteString(fmt.Sprintf("\tfor i, v := range buf[%d:%d] {\n", f.Offset, f.Offset+f.Size))
		sb.WriteString("\t\tlow := byte(v >> 8)\n")
		sb.WriteString("\t\tif low == 0 {bytes = bytes[:i*2];break} // stop on empty char\n")
		sb.WriteString("\t\tbytes[i*2] = low\n")
		sb.WriteString("\t\thigh := byte(v)\n")
		sb.WriteString("\t\tif high == 0 {bytes = bytes[:i*2+1];break} // stop on empty char\n")
		sb.WriteString("\t\tbytes[i*2+1] = high\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tm.%s = string(bytes)\n", f.Name))
	case Char16:
		sb.WriteString(fmt.Sprintf("\tbytes = make([]byte, %d)\n", f.Size))
		sb.WriteString(fmt.Sprintf("\tfor i, v := range buf[%d:%d] {\n", f.Offset, f.Offset+f.Size))
		sb.WriteString("\t\tif v == 0 {bytes = bytes[:i];break} // stop on empty char\n")
		sb.WriteString("\t\tbytes[i] = byte(v)\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tm.%s = string(bytes)\n", f.Name))
	}

	return sb.String()
}

func ExtractStringsTags(tagStr string) (FieldString, error) {
	var field FieldString
	var err error

	tags := strings.Split(tagStr, ",")

	if len(tags) != 3 {
		return FieldString{}, errors.New(`invalid tags, expected "offset", "size" and "char"`)
	}

	offsetIndex := slices.IndexFunc(tags, func(s string) bool {
		return strings.HasPrefix(s, "offset")
	})
	if offsetIndex == -1 {
		return FieldString{}, errors.New(`missing "offset" tag`)
	}

	kv := strings.Split(tags[offsetIndex], "=")
	if len(kv) != 2 {
		return FieldString{}, errors.New(`invalid "offset" tag format`)
	}

	field.Offset, err = strconv.Atoi(kv[1])
	if err != nil {
		return FieldString{}, errors.New(`invalid "offset" value`)
	}

	sizeIndex := slices.IndexFunc(tags, func(s string) bool {
		return strings.HasPrefix(s, "size")
	})
	if sizeIndex == -1 {
		return FieldString{}, errors.New(`missing "size" tag`)
	}

	kv = strings.Split(tags[sizeIndex], "=")
	if len(kv) != 2 {
		return FieldString{}, errors.New(`invalid "size" tag format`)
	}

	field.Size, err = strconv.Atoi(kv[1])
	if err != nil {
		return FieldString{}, errors.New(`invalid "size" value`)
	}

	charIndex := slices.IndexFunc(tags, func(s string) bool {
		return strings.HasPrefix(s, "char")
	})
	if charIndex == -1 {
		return FieldString{}, errors.New(`missing "char" tag`)
	}

	kv = strings.Split(tags[charIndex], "=")
	if len(kv) != 2 {
		return FieldString{}, errors.New(`invalid "char" tag format`)
	}

	field.Char = Char(kv[1])
	if err := field.Char.Validate(); err != nil {
		return FieldString{}, err
	}

	return field, nil
}
