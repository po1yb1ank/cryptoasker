package helper

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getPairsFromRawJSON(t *testing.T) {
	tests := []struct {
		name  string
		fsyms []string
		tsyms []string
		raw   string
		want  string
	}{
		{
			name:  "test build",
			fsyms: []string{"USD"},
			tsyms: []string{"BTC"},
			raw: `{
  "RAW": {
    "USD": {
      "BTC": {
        "PRICE": "USD"
      },
      "ETH": {
        "PRICE": "USD"
      }
    },
    "RUB": {
      "BTC": {
        "PRICE": "USD"
      },
      "ETH": {
        "PRICE": "USD"
      }
    }
  }
}`,
			want: `{"RAW":{"USD":{"BTC":{"PRICE":"USD"}}}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPairsFromRawJSON(tt.fsyms, tt.tsyms, json.RawMessage(tt.raw))

			assert.Equal(t, json.RawMessage(tt.want), got)
			assert.NoError(t, err)
		})
	}
}
