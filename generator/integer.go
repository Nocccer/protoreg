package generator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type FieldUint16 struct {
	Field
}

func (f FieldUint16) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	if f.IsCustomType {
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s)\n", f.Offset, f.Name))
	} else {
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = m.%s\n", f.Offset, f.Name))
	}

	return sb.String()
}

func (f FieldUint16) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	if f.IsCustomType {
		sb.WriteString(fmt.Sprintf("\tm.%s = %s(buf[%d])\n", f.Name, f.CustomType, f.Offset))
	} else {
		sb.WriteString(fmt.Sprintf("\tm.%s = buf[%d]\n", f.Name, f.Offset))
	}

	return sb.String()
}

type FieldInt16 struct {
	Field
}

func (f FieldInt16) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())
	sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s)\n", f.Offset, f.Name))

	return sb.String()
}

func (f FieldInt16) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	if f.IsCustomType {
		sb.WriteString(fmt.Sprintf("\tm.%s = %s(int16(buf[%d]))\n", f.Name, f.CustomType, f.Offset))
	} else {
		sb.WriteString(fmt.Sprintf("\tm.%s = int16(buf[%d])\n", f.Name, f.Offset))
	}

	return sb.String()
}

func ExtractIntegerTags(tagStr string) (Field, error) {
	tags := strings.Split(tagStr, ",")

	if len(tags) > 1 {
		return Field{}, errors.New(`too many tags, only "offset" is needed`)
	}

	kv := strings.Split(tags[0], "=")
	if len(kv) != 2 {
		return Field{}, errors.New(`invalid "offset" tag format`)
	}

	offset, err := strconv.Atoi(kv[1])
	if err != nil {
		return Field{}, errors.New(`invalid "offset" value`)
	}

	return Field{
		Offset: offset,
	}, nil
}
