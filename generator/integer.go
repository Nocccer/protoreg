package generator

import (
	"fmt"
	"go/types"
	"strings"
)

func (g *ProtoRegGen) newIntegerGen(name string, typ types.Type, tags Tags) (NewGenResult, error) {
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

	if field.Tags.Size != nil {
		return NewGenResult{}, fmt.Errorf(
			`"size" tag is not applicable for %s`,
			name,
		)
	}

	if field.Tags.Encoding == nil {
		field.Tags.Encoding = &g.encoding
	}

	if field.Tags.WordOrder == nil {
		field.Tags.WordOrder = &g.wordOrder
	}

	switch typ.Underlying().String() {
	case "uint8", "int8", "byte":
		if field.Tags.Byte == nil {
			return NewGenResult{}, fmt.Errorf(
				`missing "byte" tag for %s`,
				name,
			)
		}
		field.Tags.Size = ptrTo(1)
	case "uint16", "int16":
		if *field.Tags.Encoding == LittleEndian {
			g.imports = append(g.imports, "math/bits")
		}
		field.Tags.Size = ptrTo(1)
	case "uint32", "int32":
		if *field.Tags.Encoding == LittleEndian {
			g.imports = append(g.imports, "math/bits")
		}
		field.Tags.Size = ptrTo(2)
	case "uint64", "int64":
		if *field.Tags.Encoding == LittleEndian {
			g.imports = append(g.imports, "math/bits")
		}
		field.Tags.Size = ptrTo(4)
	default:
		return NewGenResult{}, fmt.Errorf("unsupported integer type: %s", typ.String())
	}

	switch typ.Underlying().String() {
	case "uint8", "byte":
		return NewGenResult{
			Gen: FieldUint8{Field: field},
			Len: *field.Tags.Offset + *field.Tags.Size,
		}, nil
	case "int8":
		return NewGenResult{
			Gen: FieldInt8{Field: field},
			Len: *field.Tags.Offset + *field.Tags.Size,
		}, nil
	case "uint16":
		return NewGenResult{
			Gen: FieldUint16{Field: field},
			Len: *field.Tags.Offset + *field.Tags.Size,
		}, nil
	case "int16":
		return NewGenResult{
			Gen: FieldInt16{Field: field},
			Len: *field.Tags.Offset + *field.Tags.Size,
		}, nil
	case "uint32":
		return NewGenResult{
			Gen: FieldUint32{Field: field},
			Len: *field.Tags.Offset + *field.Tags.Size*2,
		}, nil
	case "int32":
		return NewGenResult{
			Gen: FieldInt32{Field: field},
			Len: *field.Tags.Offset + *field.Tags.Size*2,
		}, nil
	case "uint64":
		return NewGenResult{
			Gen: FieldUint64{Field: field},
			Len: *field.Tags.Offset + *field.Tags.Size*4,
		}, nil
	case "int64":
		return NewGenResult{
			Gen: FieldInt64{Field: field},
			Len: *field.Tags.Offset + *field.Tags.Size*4,
		}, nil
	default:
		return NewGenResult{}, fmt.Errorf("unsupported integer type: %s", typ.String())
	}
}

type FieldUint8 struct {
	Field
}

func (f FieldUint8) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	shift := ""
	if (*f.Tags.Encoding == BigEndian && *f.Tags.Byte == High) ||
		(*f.Tags.Encoding == LittleEndian && *f.Tags.Byte == Low) {
		shift = "<<8"
	}

	sb.WriteString(
		fmt.Sprintf(
			"\tbuf[%d] = buf[%d] | uint16(m.%s)%s\n",
			*f.Tags.Offset,
			*f.Tags.Offset,
			f.Name,
			shift,
		),
	)

	return sb.String()
}

func (f FieldUint8) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	shift := ""
	if (*f.Tags.Encoding == BigEndian && *f.Tags.Byte == High) ||
		(*f.Tags.Encoding == LittleEndian && *f.Tags.Byte == Low) {
		shift = ">>8"
	}

	if f.IsCustomType {
		sb.WriteString(
			fmt.Sprintf("\tm.%s = %s(buf[%d]%s)\n", f.Name, f.CustomType, *f.Tags.Offset, shift),
		)
	} else {
		sb.WriteString(
			fmt.Sprintf("\tm.%s = uint8(buf[%d]%s)\n", f.Name, *f.Tags.Offset, shift),
		)
	}

	return sb.String()
}

type FieldInt8 struct {
	Field
}

func (f FieldInt8) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	shift := ""
	if (*f.Tags.Encoding == BigEndian && *f.Tags.Byte == High) ||
		(*f.Tags.Encoding == LittleEndian && *f.Tags.Byte == Low) {
		shift = "<<8"
	}

	sb.WriteString(
		fmt.Sprintf(
			"\tbuf[%d] = buf[%d] | uint16(uint8(m.%s))%s\n",
			*f.Tags.Offset,
			*f.Tags.Offset,
			f.Name,
			shift,
		),
	)

	return sb.String()
}

func (f FieldInt8) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	shift := ""
	if (*f.Tags.Encoding == BigEndian && *f.Tags.Byte == High) ||
		(*f.Tags.Encoding == LittleEndian && *f.Tags.Byte == Low) {
		shift = ">>8"
	}

	if f.IsCustomType {
		sb.WriteString(
			fmt.Sprintf("\tm.%s = %s(buf[%d]%s)\n", f.Name, f.CustomType, *f.Tags.Offset, shift),
		)
	} else {
		sb.WriteString(
			fmt.Sprintf("\tm.%s = int8(buf[%d]%s)\n", f.Name, *f.Tags.Offset, shift),
		)
	}

	return sb.String()
}

type FieldUint16 struct {
	Field
}

func (f FieldUint16) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		if f.IsCustomType {
			sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s)\n", *f.Tags.Offset, f.Name))
		} else {
			sb.WriteString(fmt.Sprintf("\tbuf[%d] = m.%s\n", *f.Tags.Offset, f.Name))
		}
	case LittleEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s))\n",
					*f.Tags.Offset,
					f.Name,
				),
			)
		} else {
			sb.WriteString(fmt.Sprintf("\tbuf[%d] = bits.ReverseBytes16(m.%s)\n", *f.Tags.Offset, f.Name))
		}
	}

	return sb.String()
}

func (f FieldUint16) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf("\tm.%s = %s(buf[%d])\n", f.Name, f.CustomType, *f.Tags.Offset),
			)
		} else {
			sb.WriteString(fmt.Sprintf("\tm.%s = buf[%d]\n", f.Name, *f.Tags.Offset))
		}
	case LittleEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = %s(bits.ReverseBytes16(buf[%d]))\n",
					f.Name,
					f.CustomType,
					*f.Tags.Offset,
				),
			)
		} else {
			sb.WriteString(fmt.Sprintf("\tm.%s = bits.ReverseBytes16(buf[%d])\n", f.Name, *f.Tags.Offset))
		}
	}

	return sb.String()
}

type FieldInt16 struct {
	Field
}

func (f FieldInt16) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s)\n", *f.Tags.Offset, f.Name))
	case LittleEndian:
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s))\n",
				*f.Tags.Offset,
				f.Name,
			),
		)
	}

	return sb.String()
}

func (f FieldInt16) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf("\tm.%s = %s(int16(buf[%d]))\n", f.Name, f.CustomType, *f.Tags.Offset),
			)
		} else {
			sb.WriteString(fmt.Sprintf("\tm.%s = int16(buf[%d])\n", f.Name, *f.Tags.Offset))
		}
	case LittleEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = %s(bits.ReverseBytes16(buf[%d]))\n",
					f.Name,
					f.CustomType,
					*f.Tags.Offset,
				),
			)
		} else {
			sb.WriteString(fmt.Sprintf("\tm.%s = int16(bits.ReverseBytes16(buf[%d]))\n", f.Name, *f.Tags.Offset))
		}
	}

	return sb.String()
}

type FieldUint32 struct {
	Field
}

func (f FieldUint32) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s)\n", *f.Tags.Offset, f.Name))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s>>16)\n", *f.Tags.Offset+1, f.Name))
	case LittleEndian:
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s))\n",
				*f.Tags.Offset,
				f.Name,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s>>16))\n",
				*f.Tags.Offset+1,
				f.Name,
			),
		)
	}

	return sb.String()
}

func (f FieldUint32) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = %s(buf[%d]) | %s(buf[%d]) << 16\n",
					f.Name,
					f.CustomType,
					*f.Tags.Offset,
					f.CustomType,
					*f.Tags.Offset+1,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf("\tm.%s = uint32(buf[%d]) | uint32(buf[%d]) << 16\n", f.Name, *f.Tags.Offset, *f.Tags.Offset+1),
			)
		}
	case LittleEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = %s(bits.ReverseBytes16(buf[%d])) | %s(bits.ReverseBytes16(buf[%d]))<<16\n",
					f.Name,
					f.CustomType,
					*f.Tags.Offset,
					f.CustomType,
					*f.Tags.Offset+1,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf("\tm.%s = uint32(bits.ReverseBytes16(buf[%d])) | uint32(bits.ReverseBytes16(buf[%d]))<<16\n",
					f.Name,
					*f.Tags.Offset,
					*f.Tags.Offset+1,
				),
			)
		}
	}

	return sb.String()
}

type FieldInt32 struct {
	Field
}

func (f FieldInt32) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s)\n", *f.Tags.Offset, f.Name))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s>>16)\n", *f.Tags.Offset+1, f.Name))
	case LittleEndian:
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s))\n",
				*f.Tags.Offset,
				f.Name,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s>>16))\n",
				*f.Tags.Offset+1,
				f.Name,
			),
		)
	}

	return sb.String()
}

func (f FieldInt32) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = %s(buf[%d]) | %s(buf[%d]) << 16\n",
					f.Name,
					f.CustomType,
					*f.Tags.Offset,
					f.CustomType,
					*f.Tags.Offset+1,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf("\tm.%s = int32(buf[%d]) | int32(buf[%d]) << 16\n", f.Name, *f.Tags.Offset, *f.Tags.Offset+1),
			)
		}
	case LittleEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = %s(bits.ReverseBytes16(buf[%d])) | %s(bits.ReverseBytes16(buf[%d]))<<16\n",
					f.Name,
					f.CustomType,
					*f.Tags.Offset,
					f.CustomType,
					*f.Tags.Offset+1,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf("\tm.%s = int32(bits.ReverseBytes16(buf[%d])) | int32(bits.ReverseBytes16(buf[%d]))<<16\n",
					f.Name,
					*f.Tags.Offset,
					*f.Tags.Offset+1,
				),
			)
		}
	}

	return sb.String()
}

type FieldUint64 struct {
	Field
}

func (f FieldUint64) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s)\n", *f.Tags.Offset, f.Name))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s>>16)\n", *f.Tags.Offset+1, f.Name))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s>>32)\n", *f.Tags.Offset+2, f.Name))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s>>48)\n", *f.Tags.Offset+3, f.Name))
	case LittleEndian:
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s))\n",
				*f.Tags.Offset,
				f.Name,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s>>16))\n",
				*f.Tags.Offset+1,
				f.Name,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s>>32))\n",
				*f.Tags.Offset+2,
				f.Name,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s>>48))\n",
				*f.Tags.Offset+3,
				f.Name,
			),
		)
	}

	return sb.String()
}

func (f FieldUint64) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = %s(buf[%d]) | %s(buf[%d]) << 16 | %s(buf[%d]) << 32 | %s(buf[%d]) << 48\n",
					f.Name,
					f.CustomType,
					*f.Tags.Offset,
					f.CustomType,
					*f.Tags.Offset+1,
					f.CustomType,
					*f.Tags.Offset+2,
					f.CustomType,
					*f.Tags.Offset+3,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = uint64(buf[%d]) | uint64(buf[%d]) << 16 | uint64(buf[%d]) << 32 | uint64(buf[%d]) << 48\n",
					f.Name,
					*f.Tags.Offset,
					*f.Tags.Offset+1,
					*f.Tags.Offset+2,
					*f.Tags.Offset+3,
				),
			)
		}
	case LittleEndian:
		if f.IsCustomType {
			sb.WriteString(
				//nolint:lll
				fmt.Sprintf(
					"\tm.%s = %s(bits.ReverseBytes16(buf[%d])) | %s(bits.ReverseBytes16(buf[%d]))<<16 | %s(bits.ReverseBytes16(buf[%d]))<<32 | %s(bits.ReverseBytes16(buf[%d]))<<48\n",
					f.Name,
					f.CustomType,
					*f.Tags.Offset,
					f.CustomType,
					*f.Tags.Offset+1,
					f.CustomType,
					*f.Tags.Offset+2,
					f.CustomType,
					*f.Tags.Offset+3,
				),
			)
		} else {
			sb.WriteString(
				//nolint:lll
				fmt.Sprintf("\tm.%s = uint64(bits.ReverseBytes16(buf[%d])) | uint64(bits.ReverseBytes16(buf[%d]))<<16 | uint64(bits.ReverseBytes16(buf[%d]))<<32 | uint64(bits.ReverseBytes16(buf[%d]))<<48\n",
					f.Name,
					*f.Tags.Offset,
					*f.Tags.Offset+1,
					*f.Tags.Offset+2,
					*f.Tags.Offset+3,
				),
			)
		}
	}

	return sb.String()
}

type FieldInt64 struct {
	Field
}

func (f FieldInt64) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s)\n", *f.Tags.Offset, f.Name))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s>>16)\n", *f.Tags.Offset+1, f.Name))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s>>32)\n", *f.Tags.Offset+2, f.Name))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s>>48)\n", *f.Tags.Offset+3, f.Name))
	case LittleEndian:
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s))\n",
				*f.Tags.Offset,
				f.Name,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s>>16))\n",
				*f.Tags.Offset+1,
				f.Name,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s>>32))\n",
				*f.Tags.Offset+2,
				f.Name,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s>>48))\n",
				*f.Tags.Offset+3,
				f.Name,
			),
		)
	}

	return sb.String()
}

func (f FieldInt64) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = %s(buf[%d]) | %s(buf[%d]) << 16 | %s(buf[%d]) << 32 | %s(buf[%d]) << 48\n",
					f.Name,
					f.CustomType,
					*f.Tags.Offset,
					f.CustomType,
					*f.Tags.Offset+1,
					f.CustomType,
					*f.Tags.Offset+2,
					f.CustomType,
					*f.Tags.Offset+3,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = int64(buf[%d]) | int64(buf[%d]) << 16 | int64(buf[%d]) << 32 | int64(buf[%d]) << 48\n",
					f.Name,
					*f.Tags.Offset,
					*f.Tags.Offset+1,
					*f.Tags.Offset+2,
					*f.Tags.Offset+3,
				),
			)
		}
	case LittleEndian:
		if f.IsCustomType {
			sb.WriteString(
				//nolint:lll
				fmt.Sprintf(
					"\tm.%s = %s(bits.ReverseBytes16(buf[%d])) | %s(bits.ReverseBytes16(buf[%d]))<<16 | %s(bits.ReverseBytes16(buf[%d]))<<32 | %s(bits.ReverseBytes16(buf[%d]))<<48\n",
					f.Name,
					f.CustomType,
					*f.Tags.Offset,
					f.CustomType,
					*f.Tags.Offset+1,
					f.CustomType,
					*f.Tags.Offset+2,
					f.CustomType,
					*f.Tags.Offset+3,
				),
			)
		} else {
			sb.WriteString(
				//nolint:lll
				fmt.Sprintf("\tm.%s = int64(bits.ReverseBytes16(buf[%d])) | int64(bits.ReverseBytes16(buf[%d]))<<16 | int64(bits.ReverseBytes16(buf[%d]))<<32 | int64(bits.ReverseBytes16(buf[%d]))<<48\n",
					f.Name,
					*f.Tags.Offset,
					*f.Tags.Offset+1,
					*f.Tags.Offset+2,
					*f.Tags.Offset+3,
				),
			)
		}
	}

	return sb.String()
}
