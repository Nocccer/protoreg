package generator

import (
	"fmt"
	"go/types"
	"strings"
)

func (g *ProtoRegGen) newFloatGen(name string, typ types.Type, tags Tags) (NewGenResult, error) {
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

	g.imports = append(g.imports, "math")

	switch typ.Underlying().String() {
	case "float32":
		field.Tags.Size = ptrTo(2)
		return NewGenResult{
			Gen: FieldFloat32{Field: field},
			Len: *field.Tags.Offset + *field.Tags.Size,
		}, nil
	case "float64":
		field.Tags.Size = ptrTo(4)
		return NewGenResult{
			Gen: FieldFloat64{Field: field},
			Len: *field.Tags.Offset + *field.Tags.Size,
		}, nil
	default:
		return NewGenResult{}, fmt.Errorf("unsupported float type: %s", typ.String())
	}
}

type FieldFloat32 struct {
	Field
}

func (f FieldFloat32) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())
	sb.WriteString(fmt.Sprintf("\ttmp32 = math.Float32bits(m.%s)\n", f.Name))

	switch *f.Tags.Encoding {
	case BigEndian:
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(tmp32)\n", *f.Tags.Offset))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(tmp32>>16)\n", *f.Tags.Offset+1))
	case LittleEndian:
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = uint16(tmp32>>8) | uint16(tmp32<<8)\n",
				*f.Tags.Offset,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = uint16(tmp32>>24) | uint16(tmp32<<24)\n",
				*f.Tags.Offset+1,
			),
		)
	}

	return sb.String()
}

func (f FieldFloat32) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		sb.WriteString(
			fmt.Sprintf(
				"\ttmp32 = uint32(buf[%d]) | uint32(buf[%d]) << 16\n",
				*f.Tags.Offset,
				*f.Tags.Offset+1,
			),
		)
	case LittleEndian:
		sb.WriteString(
			fmt.Sprintf(
				"\ttmp32 = uint32(buf[%d]>>8 | buf[%d]<<8) | uint32(buf[%d]>>8 | buf[%d]<<8)<<16\n",
				*f.Tags.Offset,
				*f.Tags.Offset,
				*f.Tags.Offset+1,
				*f.Tags.Offset+1,
			),
		)
	}

	if f.IsCustomType {
		sb.WriteString(
			fmt.Sprintf("\tm.%s = %s(math.Float32frombits(tmp32))\n", f.Name, f.CustomType),
		)
	} else {
		sb.WriteString(
			fmt.Sprintf("\tm.%s = math.Float32frombits(tmp32)\n", f.Name),
		)
	}

	return sb.String()
}

type FieldFloat64 struct {
	Field
}

func (f FieldFloat64) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())
	sb.WriteString(fmt.Sprintf("\ttmp64 = math.Float64bits(m.%s)\n", f.Name))

	switch *f.Tags.Encoding {
	case BigEndian:
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(tmp64)\n", *f.Tags.Offset))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(tmp64>>16)\n", *f.Tags.Offset+1))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(tmp64>>32)\n", *f.Tags.Offset+2))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(tmp64>>48)\n", *f.Tags.Offset+3))
	case LittleEndian:
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = uint16(tmp64>>8) | uint16(tmp64<<8)\n",
				*f.Tags.Offset,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = uint16(tmp64>>24) | uint16(tmp64<<24)\n",
				*f.Tags.Offset+1,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = uint16(tmp64>>40) | uint16(tmp64<<40)\n",
				*f.Tags.Offset+2,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = uint16(tmp64>>56) | uint16(tmp64<<56)\n",
				*f.Tags.Offset+3,
			),
		)
	}

	return sb.String()
}

func (f FieldFloat64) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		sb.WriteString(
			fmt.Sprintf(
				"\ttmp64 = uint64(buf[%d]) | uint64(buf[%d]) << 16 | uint64(buf[%d]) << 32 | uint64(buf[%d]) << 48\n",
				*f.Tags.Offset,
				*f.Tags.Offset+1,
				*f.Tags.Offset+2,
				*f.Tags.Offset+3,
			),
		)
	case LittleEndian:
		sb.WriteString(
			// nolint:lll
			fmt.Sprintf(
				"\ttmp64 = uint64(buf[%d]>>8 | buf[%d]<<8) | uint64(buf[%d]>>8 | buf[%d]<<8)<<16 | uint64(buf[%d]>>8 | buf[%d]<<8)<<32 | uint64(buf[%d]>>8 | buf[%d]<<8)<<48\n",
				*f.Tags.Offset,
				*f.Tags.Offset,
				*f.Tags.Offset+1,
				*f.Tags.Offset+1,
				*f.Tags.Offset+2,
				*f.Tags.Offset+2,
				*f.Tags.Offset+3,
				*f.Tags.Offset+3,
			),
		)
	}

	if f.IsCustomType {
		sb.WriteString(
			fmt.Sprintf("\tm.%s = %s(math.Float64frombits(tmp64))\n", f.Name, f.CustomType),
		)
	} else {
		sb.WriteString(
			fmt.Sprintf("\tm.%s = math.Float64frombits(tmp64)\n", f.Name),
		)
	}

	return sb.String()
}
