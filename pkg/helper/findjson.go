package helper

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

func getPairsFromRawJSON(fsyms, tsyms []string, raw json.RawMessage) json.RawMessage {
	var rawResult json.RawMessage
	//var tsymsRequested, fsymsRequested []json.RawMessage
	for _, fsym := range fsyms {
		data := gjson.Get(string(raw), "RAW."+fsym)
		for _, tsym := range tsyms {
			data.Get(tsym)
		}
	}

	return rawResult
}
