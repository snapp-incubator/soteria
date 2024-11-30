package strconv

import (
	"strconv"
)

// ToString convert json data types into string.
// Because we are working with JWT Tokens these convert
// return empty string on types like object or array. Also,
// they only considering the integers.
func ToString(input any) string {
	switch v := input.(type) {
	case float64:
		return strconv.Itoa(int(v))
	case float32:
		return strconv.Itoa(int(v))
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	default:
		return ""
	}
}
