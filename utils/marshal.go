package utils

import (
	"encoding/json"
)

func MustMarshal(v interface{}) []byte {
	bytes, _ := json.Marshal(v)
	return bytes
}
