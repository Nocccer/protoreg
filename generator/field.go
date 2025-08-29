package generator

import (
	"fmt"
	"go/ast"
	"go/types"
	"log/slog"
	"strings"
)

type Encoding string

const (
	BigEndian    Encoding = "big"
	LittleEndian Encoding = "little"
)

type WordOrder string

const (
	HighWordFirst WordOrder = "high"
	LowWordFirst  WordOrder = "low"
)

type Field struct {
	Name         string
	Offset       int
	Size         int
	Encoding     Encoding
	WordOrder    WordOrder
	IsCustomType bool
	CustomType   string
}

func (f Field) Comment() string {
	return fmt.Sprintf("\t// %s\n", f.Name)
}

type ExtractResult struct {
	Gen    Generator
	Len    int
	Import string
}

func (g *ProtoRegGen) extractField(
	name string,
	typ ast.Expr,
	tagStr string,
	typesInfo *types.Info,
) (ExtractResult, error) {
	t, ok := typesInfo.Types[typ]
	if !ok {
		return ExtractResult{}, fmt.Errorf("unknown type: %v", typ)
	}

	g.log.Debug(
		"extract field",
		slog.String("field", name),
		slog.Any("type", t.Type.String()),
		slog.Any("tags", tagStr),
	)

	switch t.Type.Underlying().String() {
	case "int16":
		field, err := ExtractIntegerTags(tagStr)
		if err != nil {
			return ExtractResult{}, fmt.Errorf(
				"failed to extract integer tags for %s: %v",
				name,
				err,
			)
		}
		field.Name = name
		field.Size = 1
		field.IsCustomType = t.Type.String() != "int16"
		var imports string
		if field.IsCustomType {
			field.CustomType, imports = g.extractCustomType(t.Type.String())
		}
		return ExtractResult{
			Gen:    FieldInt16{Field: field},
			Len:    field.Offset + 1,
			Import: imports,
		}, nil
	case "uint16":
		field, err := ExtractIntegerTags(tagStr)
		if err != nil {
			return ExtractResult{}, fmt.Errorf(
				"failed to extract integer tags for %s: %v",
				name,
				err,
			)
		}
		field.Name = name
		field.Size = 1
		field.IsCustomType = t.Type.String() != "uint16"
		var imports string
		if field.IsCustomType {
			field.CustomType, imports = g.extractCustomType(t.Type.String())
		}
		return ExtractResult{
			Gen:    FieldUint16{Field: field},
			Len:    field.Offset + 1,
			Import: imports,
		}, nil
	case "string":
		field, err := ExtractStringsTags(tagStr)
		if err != nil {
			return ExtractResult{}, fmt.Errorf(
				"failed to extract string tags for %s: %v",
				name,
				err,
			)
		}
		field.Name = name
		field.IsCustomType = t.Type.String() != "string"
		var imports string
		if field.IsCustomType {
			field.CustomType, imports = g.extractCustomType(t.Type.String())
		}
		return ExtractResult{
			Gen:    field,
			Len:    field.Offset + field.Size,
			Import: imports,
		}, nil
	default:
		return ExtractResult{}, fmt.Errorf(
			"Field %q unsupported underlying type: %q",
			name,
			t.Type.Underlying().String(),
		)
	}
}

func (g *ProtoRegGen) extractCustomType(typePath string) (string, string) {
	g.log.Debug(
		"extract custom type",
	)

	parts := strings.Split(typePath, "/")
	typ := parts[len(parts)-1]
	selAndType := strings.Split(typ, ".")
	if selAndType[0] == g.pkg {
		g.log.Debug(
			"type is in the same package",
			slog.String("pkg", g.pkg),
			slog.String("type", typePath),
		)
		return selAndType[1], ""
	}

	imp := strings.TrimSuffix(typePath, fmt.Sprintf(".%s", selAndType[1]))

	g.log.Debug(
		"type is in a different package",
		slog.String("import", imp),
	)
	return typ, imp
}
