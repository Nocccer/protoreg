package generator_test

import (
	"reflect"
	"testing"

	"github.com/Nocccer/protoreg/generator"
)

func TestExtractIntegerTags(t *testing.T) {
	tcs := []struct {
		name     string
		tagStr   string
		expField generator.Field
		expErr   string
	}{
		{
			name:     "valid tags",
			tagStr:   "offset=10",
			expField: generator.Field{Offset: 10},
			expErr:   "",
		},
		{
			name:     "to many tags",
			tagStr:   "offset=10,size=8",
			expField: generator.Field{},
			expErr:   `too many tags, only "offset" is needed`,
		},
		{
			name:     "missing offset value",
			tagStr:   "offset",
			expField: generator.Field{},
			expErr:   `invalid "offset" tag format`,
		},
		{
			name:     "offset value not a number",
			tagStr:   "offset=invalid",
			expField: generator.Field{},
			expErr:   `invalid "offset" value`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			field, err := generator.ExtractIntegerTags(tc.tagStr)
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
