package exampleutil

import "strconv"

// MaskInt64 returns a stable redacted representation of a numeric private ID.
func MaskInt64(value int64) string {
	if value == 0 {
		return "0"
	}

	raw := strconv.FormatInt(value, 10)
	if len(raw) <= 6 {
		return "***"
	}
	return raw[:3] + "***" + raw[len(raw)-3:]
}
