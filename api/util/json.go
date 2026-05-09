package util

import "encoding/json"

func JSONIgnoreErr(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
