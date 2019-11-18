package pinboard

import "strings"
import "time"
import "encoding/json"

// postTags is a type for parsing tags returned by the Pinboard API. The API returns
// tags as a space delimeted list, however we want to return them to the user as a
// slice of []string
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

func (t postTags) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(t))
}

// utcDate is a type for parsing _some_ of the dates returned by the Pinboard API.
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

func (u utcDate) MarshalJSON() ([]byte, error) {
	s := u.UTC().Format("2006-01-02")
	return []byte(s), nil
}

// notesDate is a type for parsing the datetime stamps in the notes list
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
