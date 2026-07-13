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
	fmt.Fprintf(&sb, "\ttmp32 = math.Float32bits(float32(m.%s))\n", f.Name)

	offsets := f.calcWordOffsets32()
	fmt.Fprintf(&sb, "\tbuf[%d] = %s\n", offsets[0], f.encodeWord16("tmp32", ""))
	fmt.Fprintf(&sb, "\tbuf[%d] = %s\n", offsets[1], f.encodeWord16("tmp32", " >> 16"))

	return sb.String()
}

func (f FieldFloat32) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	offsets := f.calcWordOffsets32()
	w0 := f.decodeWord16(offsets[0])
	w1 := f.decodeWord16(offsets[1])

	fmt.Fprintf(&sb, "\ttmp32 = uint32(%s) | uint32(%s) << 16\n", w0, w1)

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
	fmt.Fprintf(&sb, "\ttmp64 = math.Float64bits(float64(m.%s))\n", f.Name)

	offsets := f.calcWordOffsets64()
	fmt.Fprintf(&sb, "\tbuf[%d] = %s\n", offsets[0], f.encodeWord16("tmp64", ""))
	fmt.Fprintf(&sb, "\tbuf[%d] = %s\n", offsets[1], f.encodeWord16("tmp64", " >> 16"))
	fmt.Fprintf(&sb, "\tbuf[%d] = %s\n", offsets[2], f.encodeWord16("tmp64", " >> 32"))
	fmt.Fprintf(&sb, "\tbuf[%d] = %s\n", offsets[3], f.encodeWord16("tmp64", " >> 48"))

	return sb.String()
}

func (f FieldFloat64) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	offsets := f.calcWordOffsets64()
	w0 := f.decodeWord16(offsets[0])
	w1 := f.decodeWord16(offsets[1])
	w2 := f.decodeWord16(offsets[2])
	w3 := f.decodeWord16(offsets[3])

	fmt.Fprintf(&sb,
		"\ttmp64 = uint64(%s) | uint64(%s) << 16 | uint64(%s) << 32 | uint64(%s) << 48\n",
		w0, w1, w2, w3)

	if f.IsCustomType {
		fmt.Fprintf(&sb, "\tm.%s = %s(math.Float64frombits(tmp64))\n", f.Name, f.CustomType)
	} else {
		fmt.Fprintf(&sb, "\tm.%s = math.Float64frombits(tmp64)\n", f.Name)
	}

	return sb.String()
}
