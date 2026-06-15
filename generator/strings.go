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

	if field.Tags.CharEncoding == nil {
		field.Tags.CharEncoding = new(CharEncodingASCII)
	}

	if *field.Tags.Char == Char8 && *field.Tags.CharEncoding == CharEncodingUTF8 {
		return NewGenResult{}, fmt.Errorf(
			`you can't set "charencoding" to %q if "char" is %q, because utf8 needs atleast two bytes for each character`,
			CharEncodingUTF8,
			Char8,
		)
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
		fmt.Fprintf(&sb, "\tlength := len(m.%s)\n", f.Name)
		sb.WriteString("\tfor i := 0; i < length; i+=2 {\n")
		fmt.Fprintf(&sb, "\t\tif i >= %d {break}\n", *f.Tags.Size*2)
		fmt.Fprintf(&sb, "\t\tb1 := m.%s[i]\n", f.Name)
		sb.WriteString("\t\tvar b2 byte\n")
		fmt.Fprintf(&sb, "\t\tif i+1 < length {b2 = m.%s[i+1]}\n", f.Name)
		if *f.Tags.Encoding == BigEndian {
			fmt.Fprintf(&sb, "\t\tbuf[%d+i/2] = uint16(b1) | uint16(b2)<<8\n",
				*f.Tags.Offset)
		} else {
			fmt.Fprintf(&sb, "\t\tbuf[%d+i/2] = uint16(b1)<<8 | uint16(b2)\n",
				*f.Tags.Offset)
		}
		sb.WriteString("\t}\n")
	case Char16:
		if *f.Tags.CharEncoding == CharEncodingASCII {
			fmt.Fprintf(&sb, "\tfor i := 0; i < len(m.%s); i++ {\n", f.Name)
			fmt.Fprintf(&sb, "\t\tif i >= %d {break}\n", *f.Tags.Size)
			shift := ""
			if *f.Tags.Encoding == LittleEndian {
				shift = "<<8"
			}

			fmt.Fprintf(&sb, "\t\tbuf[%d+i] = uint16(m.%s[i])%s\n",
				*f.Tags.Offset,
				f.Name,
				shift)
		} else {
			sb.WriteString("\ti = 0\n")
			fmt.Fprintf(&sb, "\tfor _, r := range m.%s {\n", f.Name)
			fmt.Fprintf(&sb, "\t\tif i >= %d {break}\n", *f.Tags.Size)

			if *f.Tags.Encoding == LittleEndian {
				fmt.Fprintf(&sb, "\t\tbuf[%d+i] = uint16(r>>8) | uint16(r<<8)\n",
					*f.Tags.Offset)
			} else {
				fmt.Fprintf(&sb, "\t\tbuf[%d+i] = uint16(r)\n",
					*f.Tags.Offset)
			}

			sb.WriteString("\t\ti++\n")
		}

		sb.WriteString("\t}\n")
	}

	return sb.String()
}

func (f FieldString) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Char {
	case Char8:
		fmt.Fprintf(&sb, "\tbytes = make([]byte, %d)\n", *f.Tags.Size*2)
		fmt.Fprintf(&sb, "\tfor i, v := range buf[%d:%d] {\n",
			*f.Tags.Offset,
			*f.Tags.Offset+*f.Tags.Size)
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
	case Char16:
		if *f.Tags.CharEncoding == CharEncodingASCII {
			fmt.Fprintf(&sb, "\tbytes = make([]byte, %d)\n", *f.Tags.Size)
			fmt.Fprintf(&sb, "\tfor i, v := range buf[%d:%d] {\n",
				*f.Tags.Offset,
				*f.Tags.Offset+*f.Tags.Size)
			sb.WriteString("\t\tif v == 0 {bytes = bytes[:i];break} // stop on empty char\n")
			shift := ""
			if *f.Tags.Encoding == LittleEndian {
				shift = ">>8"
			}
			fmt.Fprintf(&sb, "\t\tbytes[i] = byte(v%s)\n", shift)
		} else {
			fmt.Fprintf(&sb, "\trunes = make([]rune, %d)\n", *f.Tags.Size)
			fmt.Fprintf(&sb, "\tfor i, v := range buf[%d:%d] {\n",
				*f.Tags.Offset,
				*f.Tags.Offset+*f.Tags.Size)
			sb.WriteString("\t\tif v == 0 {runes = runes[:i];break} // stop on empty char\n")
			if *f.Tags.Encoding == LittleEndian {
				sb.WriteString("\t\trunes[i] = rune(v>>8) | rune(v<<8)\n")
			} else {
				sb.WriteString("\t\trunes[i] = rune(v)\n")
			}
		}
	}

	varStr := "bytes"
	if *f.Tags.CharEncoding != CharEncodingASCII {
		varStr = "runes"
	}
	sb.WriteString("\t}\n")
	if f.IsCustomType {
		fmt.Fprintf(&sb, "\tm.%s = %s(%s)\n", f.Name, f.CustomType, varStr)
	} else {
		fmt.Fprintf(&sb, "\tm.%s = string(%s)\n", f.Name, varStr)
	}

	return sb.String()
}
