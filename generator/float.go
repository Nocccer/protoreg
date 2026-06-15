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
	if *field.Tags.Encoding == LittleEndian {
		g.imports = append(g.imports, "math/bits")
	}

	switch typ.Underlying().String() {
	case "float32":
		field.Tags.Size = new(2)
		return NewGenResult{
			Gen: FieldFloat32{Field: field},
			Len: *field.Tags.Offset + *field.Tags.Size,
		}, nil
	case "float64":
		field.Tags.Size = new(4)
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
	fmt.Fprintf(&sb, "\ttmp32 = math.Float32bits(m.%s)\n", f.Name)

	switch *f.Tags.Encoding {
	case BigEndian:
		fmt.Fprintf(&sb, "\tbuf[%d] = uint16(tmp32)\n", *f.Tags.Offset)
		fmt.Fprintf(&sb, "\tbuf[%d] = uint16(tmp32>>16)\n", *f.Tags.Offset+1)
	case LittleEndian:
		fmt.Fprintf(&sb, "\tbuf[%d] = bits.ReverseBytes16(uint16(tmp32))\n",
			*f.Tags.Offset)
		fmt.Fprintf(&sb, "\tbuf[%d] = bits.ReverseBytes16(uint16(tmp32>>16))\n",
			*f.Tags.Offset+1)
	}

	return sb.String()
}

func (f FieldFloat32) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		fmt.Fprintf(&sb, "\ttmp32 = uint32(buf[%d]) | uint32(buf[%d]) << 16\n",
			*f.Tags.Offset,
			*f.Tags.Offset+1)
	case LittleEndian:
		fmt.Fprintf(
			&sb,
			"\ttmp32 = uint32(bits.ReverseBytes16(buf[%d])) | uint32(bits.ReverseBytes16(buf[%d]))<<16\n",
			*f.Tags.Offset,
			*f.Tags.Offset+1,
		)
	}

	if f.IsCustomType {
		fmt.Fprintf(&sb, "\tm.%s = %s(math.Float32frombits(tmp32))\n", f.Name, f.CustomType)
	} else {
		fmt.Fprintf(&sb, "\tm.%s = math.Float32frombits(tmp32)\n", f.Name)
	}

	return sb.String()
}

type FieldFloat64 struct {
	Field
}

func (f FieldFloat64) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())
	fmt.Fprintf(&sb, "\ttmp64 = math.Float64bits(m.%s)\n", f.Name)

	switch *f.Tags.Encoding {
	case BigEndian:
		fmt.Fprintf(&sb, "\tbuf[%d] = uint16(tmp64)\n", *f.Tags.Offset)
		fmt.Fprintf(&sb, "\tbuf[%d] = uint16(tmp64>>16)\n", *f.Tags.Offset+1)
		fmt.Fprintf(&sb, "\tbuf[%d] = uint16(tmp64>>32)\n", *f.Tags.Offset+2)
		fmt.Fprintf(&sb, "\tbuf[%d] = uint16(tmp64>>48)\n", *f.Tags.Offset+3)
	case LittleEndian:
		fmt.Fprintf(&sb, "\tbuf[%d] = bits.ReverseBytes16(uint16(tmp64))\n",
			*f.Tags.Offset)
		fmt.Fprintf(&sb, "\tbuf[%d] = bits.ReverseBytes16(uint16(tmp64>>16))\n",
			*f.Tags.Offset+1)
		fmt.Fprintf(&sb, "\tbuf[%d] = bits.ReverseBytes16(uint16(tmp64>>32))\n",
			*f.Tags.Offset+2)
		fmt.Fprintf(&sb, "\tbuf[%d] = bits.ReverseBytes16(uint16(tmp64>>48))\n",
			*f.Tags.Offset+3)
	}

	return sb.String()
}

func (f FieldFloat64) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		fmt.Fprintf(
			&sb,
			"\ttmp64 = uint64(buf[%d]) | uint64(buf[%d]) << 16 | uint64(buf[%d]) << 32 | uint64(buf[%d]) << 48\n",
			*f.Tags.Offset,
			*f.Tags.Offset+1,
			*f.Tags.Offset+2,
			*f.Tags.Offset+3,
		)
	case LittleEndian:
		fmt.Fprintf(
			&sb,
			"\ttmp64 = uint64(bits.ReverseBytes16(buf[%d])) | uint64(bits.ReverseBytes16(buf[%d]))<<16 | "+
				"uint64(bits.ReverseBytes16(buf[%d]))<<32 | uint64(bits.ReverseBytes16(buf[%d]))<<48\n",
			*f.Tags.Offset,
			*f.Tags.Offset+1,
			*f.Tags.Offset+2,
			*f.Tags.Offset+3,
		)
	}

	if f.IsCustomType {
		fmt.Fprintf(&sb, "\tm.%s = %s(math.Float64frombits(tmp64))\n", f.Name, f.CustomType)
	} else {
		fmt.Fprintf(&sb, "\tm.%s = math.Float64frombits(tmp64)\n", f.Name)
	}

	return sb.String()
}
