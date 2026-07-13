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
			field.Tags.Byte = new(Low)
		}
		field.Tags.Size = new(1)
	case "uint16", "int16":
		if *field.Tags.Encoding == LittleEndian {
			g.imports = append(g.imports, "math/bits")
		}
		field.Tags.Size = new(1)
	case "uint32", "int32":
		if *field.Tags.Encoding == LittleEndian {
			g.imports = append(g.imports, "math/bits")
		}
		field.Tags.Size = new(2)
	case "uint64", "int64":
		if *field.Tags.Encoding == LittleEndian {
			g.imports = append(g.imports, "math/bits")
		}
		field.Tags.Size = new(4)
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

	fmt.Fprintf(&sb, "\tbuf[%d] = buf[%d] | uint16(m.%s)%s\n",
		*f.Tags.Offset,
		*f.Tags.Offset,
		f.Name,
		shift)

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
		fmt.Fprintf(&sb, "\tm.%s = %s(buf[%d]%s)\n", f.Name, f.CustomType, *f.Tags.Offset, shift)
	} else {
		fmt.Fprintf(&sb, "\tm.%s = uint8(buf[%d]%s)\n", f.Name, *f.Tags.Offset, shift)
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

	fmt.Fprintf(&sb, "\tbuf[%d] = buf[%d] | uint16(uint8(m.%s))%s\n",
		*f.Tags.Offset,
		*f.Tags.Offset,
		f.Name,
		shift)

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
		fmt.Fprintf(&sb, "\tm.%s = %s(buf[%d]%s)\n", f.Name, f.CustomType, *f.Tags.Offset, shift)
	} else {
		fmt.Fprintf(&sb, "\tm.%s = int8(buf[%d]%s)\n", f.Name, *f.Tags.Offset, shift)
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
			fmt.Fprintf(&sb, "\tbuf[%d] = uint16(m.%s)\n", *f.Tags.Offset, f.Name)
		} else {
			fmt.Fprintf(&sb, "\tbuf[%d] = m.%s\n", *f.Tags.Offset, f.Name)
		}
	case LittleEndian:
		if f.IsCustomType {
			fmt.Fprintf(&sb, "\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s))\n",
				*f.Tags.Offset,
				f.Name)
		} else {
			fmt.Fprintf(&sb, "\tbuf[%d] = bits.ReverseBytes16(m.%s)\n", *f.Tags.Offset, f.Name)
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
			fmt.Fprintf(&sb, "\tm.%s = %s(buf[%d])\n", f.Name, f.CustomType, *f.Tags.Offset)
		} else {
			fmt.Fprintf(&sb, "\tm.%s = buf[%d]\n", f.Name, *f.Tags.Offset)
		}
	case LittleEndian:
		if f.IsCustomType {
			fmt.Fprintf(&sb, "\tm.%s = %s(bits.ReverseBytes16(buf[%d]))\n",
				f.Name,
				f.CustomType,
				*f.Tags.Offset)
		} else {
			fmt.Fprintf(&sb, "\tm.%s = bits.ReverseBytes16(buf[%d])\n", f.Name, *f.Tags.Offset)
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
		fmt.Fprintf(&sb, "\tbuf[%d] = uint16(m.%s)\n", *f.Tags.Offset, f.Name)
	case LittleEndian:
		fmt.Fprintf(&sb, "\tbuf[%d] = bits.ReverseBytes16(uint16(m.%s))\n",
			*f.Tags.Offset,
			f.Name)
	}

	return sb.String()
}

func (f FieldInt16) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		if f.IsCustomType {
			fmt.Fprintf(&sb, "\tm.%s = %s(int16(buf[%d]))\n", f.Name, f.CustomType, *f.Tags.Offset)
		} else {
			fmt.Fprintf(&sb, "\tm.%s = int16(buf[%d])\n", f.Name, *f.Tags.Offset)
		}
	case LittleEndian:
		if f.IsCustomType {
			fmt.Fprintf(&sb, "\tm.%s = %s(bits.ReverseBytes16(buf[%d]))\n",
				f.Name,
				f.CustomType,
				*f.Tags.Offset)
		} else {
			fmt.Fprintf(&sb, "\tm.%s = int16(bits.ReverseBytes16(buf[%d]))\n",
				f.Name,
				*f.Tags.Offset)
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

	offsets := f.calcWordOffsets32()

	// Low-Word (Bits 0-15)
	fmt.Fprintf(
		&sb,
		"\tbuf[%d] = %s\n",
		offsets[0],
		f.encodeWord16(fmt.Sprintf("m.%s", f.Name), ""),
	)
	// High-Word (Bits 16-31)
	fmt.Fprintf(
		&sb,
		"\tbuf[%d] = %s\n",
		offsets[1],
		f.encodeWord16(fmt.Sprintf("m.%s", f.Name), ">>16"),
	)

	return sb.String()
}

func (f FieldUint32) Unmarshaler() string {
	var sb strings.Builder
	sb.WriteString(f.Comment())

	offsets := f.calcWordOffsets32()

	lowWordCode := f.decodeWord16(offsets[0])
	highWordCode := f.decodeWord16(offsets[1])

	if f.IsCustomType {
		fmt.Fprintf(&sb, "\tm.%s = %s(%s) | %s(%s) << 16\n",
			f.Name,
			f.CustomType, lowWordCode,
			f.CustomType, highWordCode)
	} else {
		fmt.Fprintf(&sb, "\tm.%s = uint32(%s) | uint32(%s) << 16\n",
			f.Name,
			lowWordCode,
			highWordCode)
	}

	return sb.String()
}

type FieldInt32 struct {
	Field
}

func (f FieldInt32) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	offsets := f.calcWordOffsets32()
	fmt.Fprintf(
		&sb,
		"\tbuf[%d] = %s\n",
		offsets[0],
		f.encodeWord16(fmt.Sprintf("m.%s", f.Name), ""),
	)
	fmt.Fprintf(
		&sb,
		"\tbuf[%d] = %s\n",
		offsets[1],
		f.encodeWord16(fmt.Sprintf("m.%s", f.Name), ">>16"),
	)

	return sb.String()
}

func (f FieldInt32) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	offsets := f.calcWordOffsets32()

	lowWordCode := f.decodeWord16(offsets[0])
	highWordCode := f.decodeWord16(offsets[1])

	if f.IsCustomType {
		fmt.Fprintf(&sb, "\tm.%s = %s(int32(%s) | int32(%s)<<16)\n",
			f.Name,
			f.CustomType,
			lowWordCode,
			highWordCode)
	} else {
		fmt.Fprintf(&sb, "\tm.%s = int32(%s) | int32(%s)<<16\n",
			f.Name,
			lowWordCode,
			highWordCode)
	}

	return sb.String()
}

type FieldUint64 struct {
	Field
}

func (f FieldUint64) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	offsets := f.calcWordOffsets64()

	fmt.Fprintf(
		&sb,
		"\tbuf[%d] = %s\n",
		offsets[0],
		f.encodeWord16(fmt.Sprintf("m.%s", f.Name), ""),
	)
	fmt.Fprintf(
		&sb,
		"\tbuf[%d] = %s\n",
		offsets[1],
		f.encodeWord16(fmt.Sprintf("m.%s", f.Name), ">>16"),
	)
	fmt.Fprintf(
		&sb,
		"\tbuf[%d] = %s\n",
		offsets[2],
		f.encodeWord16(fmt.Sprintf("m.%s", f.Name), ">>32"),
	)
	fmt.Fprintf(
		&sb,
		"\tbuf[%d] = %s\n",
		offsets[3],
		f.encodeWord16(fmt.Sprintf("m.%s", f.Name), ">>48"),
	)

	return sb.String()
}

func (f FieldUint64) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	offsets := f.calcWordOffsets64()

	w0 := f.decodeWord16(offsets[0])
	w1 := f.decodeWord16(offsets[1])
	w2 := f.decodeWord16(offsets[2])
	w3 := f.decodeWord16(offsets[3])

	if f.IsCustomType {
		fmt.Fprintf(&sb,
			"\tm.%s = %s(%s) | %s(%s) << 16 | %s(%s) << 32 | %s(%s) << 48\n",
			f.Name,
			f.CustomType, w0,
			f.CustomType, w1,
			f.CustomType, w2,
			f.CustomType, w3)
	} else {
		fmt.Fprintf(&sb,
			"\tm.%s = uint64(%s) | uint64(%s) << 16 | uint64(%s) << 32 | uint64(%s) << 48\n",
			f.Name, w0, w1, w2, w3)
	}

	return sb.String()
}

type FieldInt64 struct {
	Field
}

func (f FieldInt64) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	offsets := f.calcWordOffsets64()

	fmt.Fprintf(
		&sb,
		"\tbuf[%d] = %s\n",
		offsets[0],
		f.encodeWord16(fmt.Sprintf("m.%s", f.Name), ""),
	)
	fmt.Fprintf(
		&sb,
		"\tbuf[%d] = %s\n",
		offsets[1],
		f.encodeWord16(fmt.Sprintf("m.%s", f.Name), ">>16"),
	)
	fmt.Fprintf(
		&sb,
		"\tbuf[%d] = %s\n",
		offsets[2],
		f.encodeWord16(fmt.Sprintf("m.%s", f.Name), ">>32"),
	)
	fmt.Fprintf(
		&sb,
		"\tbuf[%d] = %s\n",
		offsets[3],
		f.encodeWord16(fmt.Sprintf("m.%s", f.Name), ">>48"),
	)

	return sb.String()
}

func (f FieldInt64) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	offsets := f.calcWordOffsets64()

	w0 := f.decodeWord16(offsets[0])
	w1 := f.decodeWord16(offsets[1])
	w2 := f.decodeWord16(offsets[2])
	w3 := f.decodeWord16(offsets[3])

	if f.IsCustomType {
		fmt.Fprintf(&sb,
			"\tm.%s = %s(int64(%s) | int64(%s) << 16 | int64(%s) << 32 | int64(%s) << 48)\n",
			f.Name,
			f.CustomType,
			w0,
			w1,
			w2,
			w3)
	} else {
		fmt.Fprintf(&sb,
			"\tm.%s = int64(%s) | int64(%s) << 16 | int64(%s) << 32 | int64(%s) << 48\n",
			f.Name, w0, w1, w2, w3)
	}

	return sb.String()
}
