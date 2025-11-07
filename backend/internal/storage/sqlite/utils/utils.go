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

func ToISO(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format(time.RFC3339Nano)
	return &s
}

func FromISO(s *string) *time.Time {
	if s == nil {
		return nil
	}
	t, err := time.Parse(time.RFC3339Nano, *s)
	if err != nil {
		return nil
	}
	return &t
}
