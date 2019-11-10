package pinboard

import "strings"
import "time"

type postTags []string

func (t postTags) MarshalText() ([]byte, error) {
	s := strings.Join(t, " ")
	return []byte(s), nil
}

func (t *postTags) UnmarshalText(text []byte) error {
	tags := strings.Fields(string(text))
	for _, v := range tags {
		*t = append(*t, v)
	}

	return nil
}

type utcDate struct {
	time.Time
}

func (u utcDate) MarshalText() ([]byte, error) {
	s := u.UTC().Format("2006-01-02")
	return []byte(s), nil
}

func (u *utcDate) UnmarshalText(text []byte) error {
	d, err := time.Parse("2006-01-02", string(text))
	*u = utcDate{d}
	return err
}

type notesDate struct {
	time.Time
}

func (n notesDate) MarshalText() ([]byte, error) {
	s := n.UTC().Format("2006-01-02 15:04:05")
	return []byte(s), nil
}

func (n *notesDate) UnmarshalText(text []byte) error {
	d, err := time.Parse("2006-01-02 15:04:05", string(text))
	*n = notesDate{d}
	return err
}
