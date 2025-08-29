package generator_test

import (
	"reflect"
	"testing"

	"github.com/Nocccer/protoreg/generator"
)

func TestExtractStringsTags(t *testing.T) {
	tcs := []struct {
		name     string
		tagStr   string
		expField generator.FieldString
		expErr   string
	}{
		{
			name:   "valid tags",
			tagStr: "offset=10,size=8,char=8",
			expField: generator.FieldString{
				Field: generator.Field{Offset: 10, Size: 8},
				Char:  generator.Char8,
			},
			expErr: "",
		},
		{
			name:     "too many tags",
			tagStr:   "offset=10,size=8,char=8,tag=invalid",
			expField: generator.FieldString{},
			expErr:   `invalid tags, expected "offset", "size" and "char"`,
		},
		{
			name:     "missing offset tag",
			tagStr:   "size=8,char=8,tag=invalid",
			expField: generator.FieldString{},
			expErr:   `missing "offset" tag`,
		},
		{
			name:     "missing size tag",
			tagStr:   "offset=10,char=8,tag=invalid",
			expField: generator.FieldString{},
			expErr:   `missing "size" tag`,
		},
		{
			name:     "missing char tag",
			tagStr:   "offset=10,size=8,tag=invalid",
			expField: generator.FieldString{},
			expErr:   `missing "char" tag`,
		},
		{
			name:     "missing offset value",
			tagStr:   "offset,size=8,char=8",
			expField: generator.FieldString{},
			expErr:   `invalid "offset" tag format`,
		},
		{
			name:     "missing size value",
			tagStr:   "offset=10,size,char=8",
			expField: generator.FieldString{},
			expErr:   `invalid "size" tag format`,
		},
		{
			name:     "missing char value",
			tagStr:   "offset=10,size=8,char",
			expField: generator.FieldString{},
			expErr:   `invalid "char" tag format`,
		},
		{
			name:     "offset value not a number",
			tagStr:   "offset=invalid,size=8,char=8",
			expField: generator.FieldString{},
			expErr:   `invalid "offset" value`,
		},
		{
			name:     "size value not a number",
			tagStr:   "offset=10,size=invalid,char=8",
			expField: generator.FieldString{},
			expErr:   `invalid "size" value`,
		},
		{
			name:     "char value not supported",
			tagStr:   "offset=10,size=8,char=7",
			expField: generator.FieldString{},
			expErr:   `invalid "char" value`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			field, err := generator.ExtractStringsTags(tc.tagStr)
			if err != nil && tc.expErr == "" {
				t.Fatalf("Expected no error, got %v", err)
			}

			if err == nil && tc.expErr != "" {
				t.Fatalf("Expected error %q, got nil", tc.expErr)
			}

			if err != nil && err.Error() != tc.expErr {
				t.Errorf("Expected error %q, got %q", tc.expErr, err.Error())
			}

			if !reflect.DeepEqual(field, tc.expField) {
				t.Errorf("Expected field %+v, got %+v", tc.expField, field)
			}
		})
	}
}
