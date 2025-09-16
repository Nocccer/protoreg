package generator

import (
	"fmt"
	"go/types"
	"strings"
)

func (g *ProtoRegGen) newStringGen(name string, typ types.Type, tags Tags) (NewGenResult, error) {
	field := Field{
		Name:         name,
		Tags:         tags,
		IsCustomType: typ.String() != typ.Underlying().String(),
	}

	if field.IsCustomType {
		field.CustomType = g.extractCustomType(typ.String())
	}

	if field.Tags.Offset == nil {
		return NewGenResult{}, fmt.Errorf(`missing "offset" tag for %s`, name)
	}

	if field.Tags.Size == nil {
		return NewGenResult{}, fmt.Errorf(`missing "size" tag for %s`, name)
	}

	if field.Tags.Char == nil {
		return NewGenResult{}, fmt.Errorf(`missing "char" tag for %s`, name)
	}

	if field.Tags.Encoding == nil {
		field.Tags.Encoding = &g.encoding
	}

	if field.Tags.WordOrder == nil {
		field.Tags.WordOrder = &g.wordOrder
	}

	return NewGenResult{
		Gen: FieldString{
			Field: field,
		},
		Len: *field.Tags.Offset + *field.Tags.Size,
	}, nil
}

type FieldString struct {
	Field
}

func (f FieldString) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Char {
	case Char8:
		sb.WriteString(fmt.Sprintf("\tlength := len(m.%s)\n", f.Name))
		sb.WriteString("\tif length % 2 != 0 { length += 1 }\n")
		sb.WriteString("\tbytes := make([]byte, length)\n")
		sb.WriteString(fmt.Sprintf("\tcopy(bytes, m.%s)\n", f.Name))
		sb.WriteString("\tfor i := 0; i < length; i+=2 {\n")
		sb.WriteString(fmt.Sprintf("\t\tif i >= %d {break}\n", *f.Tags.Size*2))
		if *f.Tags.Encoding == BigEndian {
			sb.WriteString(
				fmt.Sprintf(
					"\t\tbuf[%d+i/2] = uint16(bytes[i]) | uint16(bytes[i+1])<<8\n",
					*f.Tags.Offset,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf(
					"\t\tbuf[%d+i/2] = uint16(bytes[i])<<8 | uint16(bytes[i+1])\n",
					*f.Tags.Offset,
				),
			)
		}
		sb.WriteString("\t}\n")
	case Char16:
		sb.WriteString(fmt.Sprintf("\tfor i := 0; i < len(m.%s); i++ {\n", f.Name))
		sb.WriteString(fmt.Sprintf("\t\tif i >= %d {break}\n", *f.Tags.Size))
		shift := ""
		if *f.Tags.Encoding == LittleEndian {
			shift = "<<8"
		}
		sb.WriteString(
			fmt.Sprintf(
				"\t\tbuf[%d+i] = uint16(m.%s[i])%s\n",
				*f.Tags.Offset,
				f.Name,
				shift,
			),
		)
		sb.WriteString("\t}\n")
	}

	return sb.String()
}

func (f FieldString) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Char {
	case Char8:
		sb.WriteString(fmt.Sprintf("\tbytes = make([]byte, %d)\n", *f.Tags.Size))
		sb.WriteString(
			fmt.Sprintf(
				"\tfor i, v := range buf[%d:%d] {\n",
				*f.Tags.Offset,
				*f.Tags.Offset+*f.Tags.Size,
			),
		)
		if *f.Tags.Encoding == BigEndian {
			sb.WriteString("\t\tlow := byte(v)\n")
		} else {
			sb.WriteString("\t\tlow := byte(v>>8)\n")
		}
		sb.WriteString("\t\tif low == 0 {bytes = bytes[:i*2];break} // stop on empty char\n")
		sb.WriteString("\t\tbytes[i*2] = low\n")
		if *f.Tags.Encoding == BigEndian {
			sb.WriteString("\t\thigh := byte(v >> 8)\n")
		} else {
			sb.WriteString("\t\thigh := byte(v)\n")
		}
		sb.WriteString("\t\tif high == 0 {bytes = bytes[:i*2+1];break} // stop on empty char\n")
		sb.WriteString("\t\tbytes[i*2+1] = high\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tm.%s = string(bytes)\n", f.Name))
	case Char16:
		sb.WriteString(fmt.Sprintf("\tbytes = make([]byte, %d)\n", *f.Tags.Size))
		sb.WriteString(
			fmt.Sprintf(
				"\tfor i, v := range buf[%d:%d] {\n",
				*f.Tags.Offset,
				*f.Tags.Offset+*f.Tags.Size,
			),
		)
		sb.WriteString("\t\tif v == 0 {bytes = bytes[:i];break} // stop on empty char\n")
		shift := ""
		if *f.Tags.Encoding == LittleEndian {
			shift = ">>8"
		}
		sb.WriteString(fmt.Sprintf("\t\tbytes[i] = byte(v%s)\n", shift))
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tm.%s = string(bytes)\n", f.Name))
	}

	return sb.String()
}
