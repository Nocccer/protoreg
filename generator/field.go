package generator

import (
	"fmt"
	"go/ast"
	"go/types"
	"log/slog"
	"strings"
)

type Field struct {
	Name         string
	Tags         Tags
	IsCustomType bool
	CustomType   string
}

func (f Field) Comment() string {
	return fmt.Sprintf("\t// %s\n", f.Name)
}

func (g *ProtoRegGen) extractField(
	name string,
	typ ast.Expr,
	tagStr string,
	typesInfo *types.Info,
) (NewGenResult, error) {
	t, ok := typesInfo.Types[typ]
	if !ok {
		return NewGenResult{}, fmt.Errorf("unknown type: %v", typ)
	}

	g.log.Debug(
		"extract field",
		slog.String("field", name),
		slog.Any("type", t.Type.String()),
		slog.Any("tags", tagStr),
	)

	tag, err := extractTags(tagStr)
	if err != nil {
		return NewGenResult{}, fmt.Errorf("failed to extract tags for %s: %v", name, err)
	}

	if strings.Contains(t.Type.Underlying().String(), "int") ||
		strings.Contains(t.Type.Underlying().String(), "byte") {
		return g.newIntegerGen(name, t.Type, tag)
	} else if t.Type.Underlying().String() == "string" {
		return g.newStringGen(name, t.Type, tag)
	}

	return NewGenResult{}, fmt.Errorf(
		"unsupported underlying field type: %s",
		t.Type.Underlying().String(),
	)
}

func (g *ProtoRegGen) extractOpts(tagStr string) error {
	// Extract encoding and word order from the tag string
	parts := strings.Split(tagStr, ",")
	for _, part := range parts {
		switch {
		case strings.HasPrefix(part, "encoding="):
			g.encoding = Encoding(strings.TrimPrefix(part, "encoding="))
		case strings.HasPrefix(part, "wordOrder="):
			g.wordOrder = WordOrder(strings.TrimPrefix(part, "wordOrder="))
		}
	}

	if g.encoding != BigEndian && g.encoding != LittleEndian {
		return fmt.Errorf("invalid encoding: %v", g.encoding)
	}
	if g.wordOrder != HighWordFirst && g.wordOrder != LowWordFirst {
		return fmt.Errorf("invalid word order: %v", g.wordOrder)
	}

	return nil
}
