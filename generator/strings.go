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

	if *field.Tags.Char == Char32 && *field.Tags.Size%2 != 0 {
		return NewGenResult{}, fmt.Errorf(
			`"size" for %s with "char=32" must be divisible by 2.`,
			name,
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
			fmt.Fprintf(&sb, "\t\tbuf[%d+i/2] = uint16(b1) << 8 | uint16(b2)\n",
				*f.Tags.Offset)
		} else {
			fmt.Fprintf(&sb, "\t\tbuf[%d+i/2] = uint16(b1) | uint16(b2) << 8\n",
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
			sb.WriteString("\t}\n")
		} else {
			sb.WriteString("\ti = 0\n")
			fmt.Fprintf(&sb, "\tfor _, r := range m.%s {\n", f.Name)
			fmt.Fprintf(&sb, "\t\tif i >= %d {break}\n", *f.Tags.Size)
			if *f.Tags.Encoding == LittleEndian {
				fmt.Fprintf(&sb, "\t\tbuf[%d+i] = uint16(r >> 8) | uint16(r << 8)\n",
					*f.Tags.Offset)
			} else {
				fmt.Fprintf(&sb, "\t\tbuf[%d+i] = uint16(r)\n",
					*f.Tags.Offset)
			}
			sb.WriteString("\t\ti++\n")
			sb.WriteString("\t}\n")
		}
	case Char32:
		charsPerField := *f.Tags.Size / 2
		if *f.Tags.CharEncoding == CharEncodingASCII {
			fmt.Fprintf(&sb, "\tfor i := 0; i < len(m.%s); i++ {\n", f.Name)
			fmt.Fprintf(&sb, "\t\tif i >= %d {break}\n", charsPerField)
			fmt.Fprintf(&sb, "\t\ttmp32 = uint32(m.%s[i])\n", f.Name)

			if *f.Tags.WordOrder == LowWordFirst {
				if *f.Tags.Encoding == LittleEndian {
					fmt.Fprintf(
						&sb,
						"\t\tbuf[%d+i*2] = uint16(tmp32)>>8 | uint16(tmp32)<<8\n",
						*f.Tags.Offset,
					)
					fmt.Fprintf(
						&sb,
						"\t\tbuf[%d+i*2+1] = uint16(tmp32>>16)>>8 | uint16(tmp32>>16)<<8\n",
						*f.Tags.Offset,
					)
				} else {
					fmt.Fprintf(&sb, "\t\tbuf[%d+i*2] = uint16(tmp32)\n", *f.Tags.Offset)
					fmt.Fprintf(&sb, "\t\tbuf[%d+i*2+1] = uint16(tmp32 >> 16)\n", *f.Tags.Offset)
				}
			} else {
				if *f.Tags.Encoding == LittleEndian {
					fmt.Fprintf(
						&sb,
						"\t\tbuf[%d+i*2] = uint16(tmp32>>16)>>8 | uint16(tmp32>>16)<<8\n",
						*f.Tags.Offset,
					)
					fmt.Fprintf(
						&sb,
						"\t\tbuf[%d+i*2+1] = uint16(tmp32)>>8 | uint16(tmp32)<<8\n",
						*f.Tags.Offset,
					)
				} else {
					fmt.Fprintf(&sb, "\t\tbuf[%d+i*2] = uint16(tmp32 >> 16)\n", *f.Tags.Offset)
					fmt.Fprintf(&sb, "\t\tbuf[%d+i*2+1] = uint16(tmp32)\n", *f.Tags.Offset)
				}
			}

			sb.WriteString("\t}\n")
		} else {
			sb.WriteString("\ti = 0\n")
			fmt.Fprintf(&sb, "\tfor _, r := range m.%s {\n", f.Name)
			fmt.Fprintf(&sb, "\t\tif i >= %d {break}\n", charsPerField)
			sb.WriteString("\t\ttmp32 = uint32(r)\n")

			if *f.Tags.WordOrder == LowWordFirst {
				if *f.Tags.Encoding == LittleEndian {
					fmt.Fprintf(
						&sb,
						"\t\tbuf[%d+i*2] = uint16(tmp32)>>8 | uint16(tmp32)<<8\n",
						*f.Tags.Offset,
					)
					fmt.Fprintf(
						&sb,
						"\t\tbuf[%d+i*2+1] = uint16(tmp32>>16)>>8 | uint16(tmp32>>16)<<8\n",
						*f.Tags.Offset,
					)
				} else {
					fmt.Fprintf(&sb, "\t\tbuf[%d+i*2] = uint16(tmp32)\n", *f.Tags.Offset)
					fmt.Fprintf(&sb, "\t\tbuf[%d+i*2+1] = uint16(tmp32 >> 16)\n", *f.Tags.Offset)
				}
			} else {
				if *f.Tags.Encoding == LittleEndian {
					fmt.Fprintf(
						&sb,
						"\t\tbuf[%d+i*2] = uint16(tmp32>>16)>>8 | uint16(tmp32>>16)<<8\n",
						*f.Tags.Offset,
					)
					fmt.Fprintf(
						&sb,
						"\t\tbuf[%d+i*2+1] = uint16(tmp32)>>8 | uint16(tmp32)<<8\n",
						*f.Tags.Offset,
					)
				} else {
					fmt.Fprintf(&sb, "\t\tbuf[%d+i*2] = uint16(tmp32 >> 16)\n", *f.Tags.Offset)
					fmt.Fprintf(&sb, "\t\tbuf[%d+i*2+1] = uint16(tmp32)\n", *f.Tags.Offset)
				}
			}

			sb.WriteString("\t\ti++\n")
			sb.WriteString("\t}\n")
		}
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
			sb.WriteString("\t\tfirst := byte(v >> 8)\n")
		} else {
			sb.WriteString("\t\tfirst := byte(v)\n")
		}
		sb.WriteString("\t\tif first == 0 {bytes = bytes[:i*2];break} // stop on empty char\n")
		sb.WriteString("\t\tbytes[i*2] = first\n")
		if *f.Tags.Encoding == BigEndian {
			sb.WriteString("\t\tsecond := byte(v)\n")
		} else {
			sb.WriteString("\t\tsecond := byte(v >> 8)\n")
		}
		sb.WriteString("\t\tif second == 0 {bytes = bytes[:i*2+1];break} // stop on empty char\n")
		sb.WriteString("\t\tbytes[i*2+1] = second\n")
		sb.WriteString("\t}\n")
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
			sb.WriteString("\t}\n")
		} else {
			fmt.Fprintf(&sb, "\trunes = make([]rune, %d)\n", *f.Tags.Size)
			fmt.Fprintf(&sb, "\tfor i, v := range buf[%d:%d] {\n",
				*f.Tags.Offset,
				*f.Tags.Offset+*f.Tags.Size)
			sb.WriteString("\t\tif v == 0 {runes = runes[:i];break} // stop on empty char\n")
			if *f.Tags.Encoding == LittleEndian {
				sb.WriteString("\t\trunes[i] = rune(v >> 8) | rune(v << 8)\n")
			} else {
				sb.WriteString("\t\trunes[i] = rune(v)\n")
			}
			sb.WriteString("\t}\n")
		}
	case Char32:
		charsPerField := *f.Tags.Size / 2
		if *f.Tags.CharEncoding == CharEncodingASCII {
			fmt.Fprintf(&sb, "\tbytes = make([]byte, 0, %d)\n", charsPerField)
		} else {
			fmt.Fprintf(&sb, "\trunes = make([]rune, 0, %d)\n", charsPerField)
		}

		fmt.Fprintf(&sb, "\tfor i := 0; i < %d; i++ {\n", charsPerField)

		if *f.Tags.WordOrder == LowWordFirst {
			if *f.Tags.Encoding == LittleEndian {
				fmt.Fprintf(
					&sb,
					"\t\tfirst := uint16(buf[%d+i*2] >> 8) | uint16(buf[%d+i*2] << 8)\n",
					*f.Tags.Offset,
					*f.Tags.Offset,
				)
				fmt.Fprintf(
					&sb,
					"\t\tsecond := uint16(buf[%d+i*2+1] >> 8) | uint16(buf[%d+i*2+1] << 8)\n",
					*f.Tags.Offset,
					*f.Tags.Offset,
				)
			} else {
				fmt.Fprintf(&sb, "\t\tfirst := buf[%d+i*2]\n", *f.Tags.Offset)
				fmt.Fprintf(&sb, "\t\tsecond := buf[%d+i*2+1]\n", *f.Tags.Offset)
			}
		} else {
			if *f.Tags.Encoding == LittleEndian {
				fmt.Fprintf(
					&sb,
					"\t\tfirst := uint16(buf[%d+i*2] >> 8) | uint16(buf[%d+i*2] << 8)\n",
					*f.Tags.Offset,
					*f.Tags.Offset,
				)
				fmt.Fprintf(
					&sb,
					"\t\tsecond := uint16(buf[%d+i*2+1] >> 8) | uint16(buf[%d+i*2+1] << 8)\n",
					*f.Tags.Offset,
					*f.Tags.Offset,
				)
			} else {
				fmt.Fprintf(&sb, "\t\tfirst := buf[%d+i*2]\n", *f.Tags.Offset)
				fmt.Fprintf(&sb, "\t\tsecond := buf[%d+i*2+1]\n", *f.Tags.Offset)
			}
		}

		if *f.Tags.WordOrder == LowWordFirst {
			sb.WriteString("\t\ttmp32 = uint32(first) | uint32(second) << 16\n")
		} else {
			sb.WriteString("\t\ttmp32 = uint32(second) | uint32(first) << 16\n")
		}
		sb.WriteString("\t\tif tmp32 == 0 {break} // stop on empty char\n")

		if *f.Tags.CharEncoding == CharEncodingASCII {
			sb.WriteString("\t\tbytes = append(bytes, byte(tmp32))\n")
		} else {
			sb.WriteString("\t\trunes = append(runes, rune(tmp32))\n")
		}
		sb.WriteString("\t}\n")
	}

	varStr := "bytes"
	if *f.Tags.CharEncoding != CharEncodingASCII {
		varStr = "runes"
	}
	if f.IsCustomType {
		fmt.Fprintf(&sb, "\tm.%s = %s(%s)\n", f.Name, f.CustomType, varStr)
	} else {
		fmt.Fprintf(&sb, "\tm.%s = string(%s)\n", f.Name, varStr)
	}

	return sb.String()
}
