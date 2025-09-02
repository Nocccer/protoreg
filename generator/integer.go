package generator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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

type FieldUint16 struct {
	Field
}

func (f FieldUint16) Marshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch f.Encoding {
	case BigEndian:
		if f.IsCustomType {
			sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s)\n", f.Offset, f.Name))
		} else {
			sb.WriteString(fmt.Sprintf("\tbuf[%d] = m.%s\n", f.Offset, f.Name))
		}
	case LittleEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tbuf[%d] = uint16(m.%s>>8) | uint16(m.%s<<8)\n",
					f.Offset,
					f.Name,
					f.Name,
				),
			)
		} else {
			sb.WriteString(fmt.Sprintf("\tbuf[%d] = m.%s>>8 | m.%s<<8\n", f.Offset, f.Name, f.Name))
		}
	}

	return sb.String()
}

func (f FieldUint16) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch f.Encoding {
	case BigEndian:
		if f.IsCustomType {
			sb.WriteString(fmt.Sprintf("\tm.%s = %s(buf[%d])\n", f.Name, f.CustomType, f.Offset))
		} else {
			sb.WriteString(fmt.Sprintf("\tm.%s = buf[%d]\n", f.Name, f.Offset))
		}
	case LittleEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = %s(buf[%d]>>8 | buf[%d]<<8)\n",
					f.Name,
					f.CustomType,
					f.Offset,
					f.Offset,
				),
			)
		} else {
			sb.WriteString(fmt.Sprintf("\tm.%s = buf[%d]>>8 | buf[%d]<<8\n", f.Name, f.Offset, f.Offset))
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

	switch f.Encoding {
	case BigEndian:
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s)\n", f.Offset, f.Name))
	case LittleEndian:
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = uint16(m.%s)>>8 | uint16(m.%s)<<8\n",
				f.Offset,
				f.Name,
				f.Name,
			),
		)
	}

	return sb.String()
}

func (f FieldInt16) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch f.Encoding {
	case BigEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf("\tm.%s = %s(int16(buf[%d]))\n", f.Name, f.CustomType, f.Offset),
			)
		} else {
			sb.WriteString(fmt.Sprintf("\tm.%s = int16(buf[%d])\n", f.Name, f.Offset))
		}
	case LittleEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = %s(int16(buf[%d]>>8 | buf[%d]<<8))\n",
					f.Name,
					f.CustomType,
					f.Offset,
					f.Offset,
				),
			)
		} else {
			sb.WriteString(fmt.Sprintf("\tm.%s = int16(buf[%d]>>8 | buf[%d]<<8)\n", f.Name, f.Offset, f.Offset))
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

	switch f.Encoding {
	case BigEndian:
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s)\n", f.Offset, f.Name))
		sb.WriteString(fmt.Sprintf("\tbuf[%d] = uint16(m.%s>>16)\n", f.Offset+1, f.Name))
	case LittleEndian:
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = uint16(m.%s>>8) | uint16(m.%s<<8)\n",
				f.Offset,
				f.Name,
				f.Name,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				"\tbuf[%d] = uint16(m.%s>>24) | uint16(m.%s<<24)\n",
				f.Offset+1,
				f.Name,
				f.Name,
			),
		)
	}

	return sb.String()
}

func (f FieldUint32) Unmarshaler() string {
	var sb strings.Builder

	sb.WriteString(f.Comment())

	switch f.Encoding {
	case BigEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = %s(buf[%d]) | %s(buf[%d]) << 16\n",
					f.Name,
					f.CustomType,
					f.Offset,
					f.CustomType,
					f.Offset+1,
				),
			)
		} else {
			sb.WriteString(fmt.Sprintf("\tm.%s = uint32(buf[%d]) | uint32(buf[%d]) << 16\n", f.Name, f.Offset, f.Offset+1))
		}
	case LittleEndian:
		if f.IsCustomType {
			sb.WriteString(
				fmt.Sprintf(
					"\tm.%s = %s(buf[%d]>>8 | buf[%d]<<8) | %s(buf[%d]>>8 | buf[%d]<<8)<<16\n",
					f.Name,
					f.CustomType,
					f.Offset,
					f.Offset,
					f.CustomType,
					f.Offset+1,
					f.Offset+1,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf("\tm.%s = uint32(buf[%d]>>8 | buf[%d]<<8) | uint32(buf[%d]>>8 | buf[%d]<<8)<<16\n",
					f.Name,
					f.Offset,
					f.Offset,
					f.Offset+1,
					f.Offset+1,
				),
			)
		}
	}

	return sb.String()
}
