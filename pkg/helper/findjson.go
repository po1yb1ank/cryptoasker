package helper

import (
	"encoding/json"
)

// Outer ...
type Outer struct {
	RAW map[string]map[string]Value `json:"RAW"`
}

type Value struct {
	CHANGE24HOUR    string `json:"CHANGE_24_HOUR,omitempty"`
	CHANGEPCT24HOUR string `json:"CHANGEPCT_24_HOUR,omitempty"`
	OPEN24HOUR      string `json:"OPEN_24_HOUR,omitempty"`
	VOLUME24HOUR    string `json:"VOLUME_24_HOUR,omitempty"`
	VOLUME24HOURTO  string `json:"VOLUME_24_HOURTO,omitempty"`
	LOW24HOUR       string `json:"LOW_24_HOUR,omitempty"`
	HIGH24HOUR      string `json:"HIGH_24_HOUR,omitempty"`
	PRICE           string `json:"PRICE,omitempty"`
	SUPPLY          string `json:"SUPPLY,omitempty"`
	MKTCAP          string `json:"MKTCAP,omitempty"`
}

const rawtag = "RAW"

func GetPairsFromRawJSON(fsyms, tsyms []string, raw json.RawMessage) (json.RawMessage, error) {
	m := make(map[string]map[string]map[string]Value)
	err := json.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}

	newm := make(map[string]map[string]Value, len(m))
	for _, fsym := range fsyms {
		for _, tsym := range tsyms {
			if newm[fsym] == nil {
				newm[fsym] = make(map[string]Value)
			}
			newm[fsym][tsym] = m[rawtag][fsym][tsym]
		}
	}
	o := Outer{
		RAW: newm,
	}
	b, err := json.Marshal(o)
	return b, err
}
