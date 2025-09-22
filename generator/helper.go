package generator

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

func ptrTo[T any](v T) *T {
	return &v
}

func splitFunc[S ~[]E, E any](s S, f func(E) bool) (a S, b S) {
	for _, v := range s {
		if f(v) {
			a = append(a, v)
		} else {
			b = append(b, v)
		}
	}

	return a, b
}

func extractProtoRegTag(tagStr string) (string, bool) {
	tags := strings.Trim(tagStr, "`")
	return reflect.StructTag(tags).Lookup("protoreg")
}

func (g *ProtoRegGen) extractCustomType(typePath string) string {
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
		return selAndType[1]
	}

	imp := strings.TrimSuffix(typePath, fmt.Sprintf(".%s", selAndType[1]))

	g.log.Debug(
		"type is in a different package",
		slog.String("import", imp),
	)

	g.imports = append(g.imports, imp)

	return typ
}
