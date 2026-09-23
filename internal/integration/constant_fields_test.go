package integration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	constapi "github.com/ogen-go/ogen/internal/integration/test_constant_fields"
)

func TestDecodeConstantFields(t *testing.T) {
	const valid = `{"state":"ready","code":42,"ratio":1.5,"enabled":true,"empty":null}`
	for name, input := range map[string]string{
		"all constants":        valid,
		"escaped equal string": strings.Replace(valid, `"ready"`, `"r\u0065ady"`, 1),
		"equal decimal":        strings.Replace(valid, "1.5", "1.500", 1),
		"equal exponent":       strings.Replace(valid, "1.5", "15e-1", 1),
		"optional equal":       strings.TrimSuffix(valid, "}") + `,"optional":"optional"}`,
	} {
		t.Run(name, func(t *testing.T) {
			var value constapi.Probe
			require.NoError(t, value.UnmarshalJSON([]byte(input)))
		})
	}
	for name, input := range map[string]string{
		"string mismatch":      strings.Replace(valid, `"ready"`, `"other"`, 1),
		"string case mismatch": strings.Replace(valid, `"ready"`, `"READY"`, 1),
		"integer mismatch":     strings.Replace(valid, "42", "43", 1),
		"number mismatch":      strings.Replace(valid, "1.5", "1.6", 1),
		"bool mismatch":        strings.Replace(valid, "true", "false", 1),
		"null mismatch":        strings.Replace(valid, "null", `"null"`, 1),
		"wrong JSON type":      strings.Replace(valid, `"ready"`, "null", 1),
		"missing required":     strings.Replace(valid, `"state":"ready",`, "", 1),
		"optional mismatch":    strings.TrimSuffix(valid, "}") + `,"optional":"other"}`,
	} {
		t.Run(name, func(t *testing.T) {
			var value constapi.Probe
			require.Error(t, value.UnmarshalJSON([]byte(input)))
		})
	}
}
