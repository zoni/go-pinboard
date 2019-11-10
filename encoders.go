package pinboard

import "strings"
import "time"

type PostTags []string

func (t PostTags) MarshalText() ([]byte, error) {
	s := strings.Join(t, " ")
	return []byte(s), nil
}

func (t *PostTags) UnmarshalText(text []byte) error {
	tags := strings.Fields(string(text))
	for _, v := range tags {
		*t = append(*t, v)
	}

	return nil
}

type UTCDate struct {
	time.Time
}

func (u UTCDate) MarshalText() ([]byte, error) {
	s := u.UTC().Format("2006-01-02")
	return []byte(s), nil
}

func (u *UTCDate) UnmarshalText(text []byte) error {
	d, err := time.Parse("2006-01-02", string(text))
	*u = UTCDate{d}
	return err
}

type NotesDate struct {
	time.Time
}

func (n NotesDate) MarshalText() ([]byte, error) {
	s := n.UTC().Format("2006-01-02 15:04:05")
	return []byte(s), nil
}

func (n *NotesDate) UnmarshalText(text []byte) error {
	d, err := time.Parse("2006-01-02 15:04:05", string(text))
	*n = NotesDate{d}
	return err
}
