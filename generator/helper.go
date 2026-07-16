package generator

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

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

func extractTagByKey(tagStr string, key string) (string, bool) {
	tags := strings.Trim(tagStr, "`")
	return reflect.StructTag(tags).Lookup(key)
}

func (g *ProtoRegGen) extractCustomType(typePath string) string {
	g.log.Debug(
		"extract custom type",
	)

	// catch generic
	var gtyp string
	var genericType string
	if strings.Contains(typePath, "[") {
		typePath, genericType, _ = strings.Cut(typePath, "[")
		gtyp = g.extractCustomType(genericType[:len(genericType)-1]) // remove closing bracket
	}

	parts := strings.Split(typePath, "/")
	typ := parts[len(parts)-1]
	pkgAndType := strings.Split(typ, ".")
	if pkgAndType[0] == g.pkg {
		g.log.Debug(
			"type is in the same package",
			slog.String("pkg", g.pkg),
			slog.String("pkg-type", typePath),
		)
		typ = pkgAndType[1]
	} else {
		imp := strings.TrimSuffix(typePath, fmt.Sprintf(".%s", pkgAndType[1]))

		g.log.Debug(
			"type is in a different package",
			slog.String("import", imp),
		)

		g.imports = append(g.imports, imp)
	}

	if gtyp != "" {
		return fmt.Sprintf("%s[%s]", typ, gtyp)
	}

	return typ
}
