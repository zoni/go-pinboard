package pinboard

import "strings"

type Tags []string

func (t Tags) MarshalText() ([]byte, error) {
	s := strings.Join(t, " ")
	return []byte(s), nil
}

func (t *Tags) UnmarshalText(text []byte) error {
	tags := strings.Fields(string(text))
	for _, v := range tags {
		*t = append(*t, v)
	}

	return nil
}
