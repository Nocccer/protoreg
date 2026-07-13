package generator

import (
	"fmt"
	"go/types"
	"strings"
)

func (g *ProtoRegGen) newBoolGen(name string, typ types.Type, tags Tags) (NewGenResult, error) {
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
	field.Tags.Size = new(1)

	if field.Tags.Encoding == nil {
		field.Tags.Encoding = &g.encoding
	}

	if field.Tags.WordOrder == nil {
		field.Tags.WordOrder = &g.wordOrder
	}

	if field.Tags.Bit == nil {
		field.Tags.Bit = new(0)
	}

	g.imports = append(g.imports, "math")
	if *field.Tags.Encoding == LittleEndian {
		g.imports = append(g.imports, "math/bits")
	}

	return NewGenResult{
		Gen: FieldBool{Field: field},
		Len: *field.Tags.Offset + *field.Tags.Size,
	}, nil
}

type FieldBool struct {
	Field
}

func (f FieldBool) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())
	fmt.Fprintf(&sb, "\ttmp16 = 0; if m.%s {tmp16 = 1}\n", f.Name)

	switch *f.Tags.Encoding {
	case BigEndian:
		fmt.Fprintf(&sb, "\tbuf[%d] |= tmp16 << %d\n", *f.Tags.Offset, *f.Tags.Bit)
	case LittleEndian:
		fmt.Fprintf(
			&sb,
			"\tbuf[%d] |= bits.ReverseBytes16(tmp16 << %d)\n",
			*f.Tags.Offset,
			*f.Tags.Bit,
		)
	}

	return sb.String()
}

func (f FieldBool) Unmarshaler() string {
	var sb strings.Builder

	mask := uint16(1 << *f.Tags.Bit)

	sb.WriteString(f.Comment())

	switch *f.Tags.Encoding {
	case BigEndian:
		fmt.Fprintf(&sb, "\ttmp16 = buf[%d] & 0x%04X\n", *f.Tags.Offset, mask)
	case LittleEndian:
		fmt.Fprintf(
			&sb,
			"\ttmp16 = bits.ReverseBytes16(buf[%d]) & 0x%04X\n",
			*f.Tags.Offset,
			mask,
		)
	}

	fmt.Fprintf(&sb, "\tif tmp16 != 0 { m.%s = true } else { m.%s = false }\n", f.Name, f.Name)

	return sb.String()
}
