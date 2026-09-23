package integration

import (
	"testing"

	unionapi "github.com/ogen-go/ogen/internal/integration/test_disjoint_field_types"
	"github.com/stretchr/testify/require"
)

func TestDisjointFieldTypesRoundTrip(t *testing.T) {
	for _, input := range []string{
		`{"record_type":"draft"}`,
		`{"record_type":"template"}`,
		`{"record_type":"draft","review":null}`,
		`{"record_type":"template","review":null}`,
		`{"record_type":"draft","review":{"kind":"review"}}`,
		`{"review":{"kind":"review"},"record_type":"draft"}`,
	} {
		t.Run(input, func(t *testing.T) {
			var value unionapi.Reference
			require.NoError(t, value.UnmarshalJSON([]byte(input)))
			require.NoError(t, value.Validate())
			encoded, err := value.MarshalJSON()
			require.NoError(t, err)
			var decoded unionapi.Reference
			require.NoError(t, decoded.UnmarshalJSON(encoded))
			require.Equal(t, value, decoded)
		})
	}
}

func TestDisjointFieldTypesRejectInvalidVariants(t *testing.T) {
	for _, input := range []string{
		`{"record_type":"template","review":{"kind":"review"}}`,
		`{"review":{"kind":"review"},"record_type":"template"}`,
		`{"record_type":"draft","review":{"kind":"wrong"}}`,
		`{"record_type":"draft","review":{}}`,
		`{"record_type":"draft","review":[]}`,
		`{"record_type":"draft","review":"review"}`,
		`{"record_type":"unknown","review":null}`,
		`{"review":{"kind":"review"}}`,
		`{"record_type":"draft","review":{"kind":"review","extra":true}}`,
	} {
		t.Run(input, func(t *testing.T) {
			var value unionapi.Reference
			err := value.UnmarshalJSON([]byte(input))
			if err == nil {
				// Ogen deliberately separates enum validation from JSON decoding.
				// Generated HTTP clients run both boundaries.
				err = value.Validate()
			}
			require.Error(t, err)
		})
	}
}
