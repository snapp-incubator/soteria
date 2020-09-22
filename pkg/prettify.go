package pkg

import "encoding/json"

func PrettifyStruct(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
