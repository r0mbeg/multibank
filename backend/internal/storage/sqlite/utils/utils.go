// internal/storage/utils/utils.go

package utils

import "time"

const TsLayout = "2006-01-02 15:04:05"
const TsLayoutNs = "2006-01-02 15:04:05.99999"

// helpers for date parsing
func ParseTS(s string) (time.Time, error) {
	t, err := time.Parse(TsLayout, s)
	if err == nil {
		return t, nil
	}
	return time.Parse(TsLayoutNs, s)
}
