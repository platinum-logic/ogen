package gen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ogen-go/ogen/gen/ir"
	"github.com/ogen-go/ogen/jsonschema"
)

func reviewUnionSchema(fieldName string, ordinaryRequired, reviewedRequired, reviewedNullable bool) *jsonschema.Schema {
	closed := false
	variant := func(kind, value *jsonschema.Schema, required bool) *jsonschema.Schema {
		return &jsonschema.Schema{
			Type: jsonschema.Object, AdditionalProperties: &closed,
			Properties: []jsonschema.Property{
				{Name: "record_type", Required: true, Schema: kind},
				{Name: fieldName, Required: required, Schema: value},
			},
		}
	}
	return &jsonschema.Schema{OneOf: []*jsonschema.Schema{
		variant(&jsonschema.Schema{Type: jsonschema.String, Enum: []any{"draft", "template"}}, &jsonschema.Schema{Type: jsonschema.Null}, ordinaryRequired),
		variant(&jsonschema.Schema{Type: jsonschema.String, Const: "draft", ConstSet: true}, &jsonschema.Schema{
			Type: jsonschema.Object, Nullable: reviewedNullable, AdditionalProperties: &closed,
			Properties: []jsonschema.Property{{Name: "kind", Required: true, Schema: &jsonschema.Schema{Type: jsonschema.String}}},
		}, reviewedRequired),
	}}
}

func TestDisjointFieldTypesWithOverlappingRecordKinds(t *testing.T) {
	for _, field := range []string{"review", "value"} {
		for _, ordinaryRequired := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/required=%v", field, ordinaryRequired), func(t *testing.T) {
				generator := newSchemaGen(func(jsonschema.Ref) (*ir.Type, bool) { return nil, false })
				typ, err := generator.generate("Reference", reviewUnionSchema(field, ordinaryRequired, true, false), false)
				require.NoError(t, err)
				require.Equal(t, ir.KindSum, typ.Kind)
				require.Len(t, typ.SumSpec.UniqueFields, 1)
				require.Contains(t, typ.SumSpec.UniqueFields, field)
				require.NotContains(t, typ.SumSpec.UniqueFields, "record_type")
				if ordinaryRequired {
					require.Empty(t, typ.SumSpec.DefaultMapping)
				} else {
					require.Equal(t, typ.SumOf[0].Name, typ.SumSpec.DefaultMapping)
				}
			})
		}
	}
}

func TestOverlappingNullOrMissingFieldsStayUnsupported(t *testing.T) {
	for _, test := range []struct {
		name               string
		required, nullable bool
	}{
		{"both optional", false, false},
		{"nullable object overlaps null", true, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			generator := newSchemaGen(func(jsonschema.Ref) (*ir.Type, bool) { return nil, false })
			_, err := generator.generate("Reference", reviewUnionSchema("review", false, test.required, test.nullable), false)
			require.Error(t, err)
		})
	}
}
