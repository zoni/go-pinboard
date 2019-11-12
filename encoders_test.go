package pinboard

import (
	"encoding/xml"
	"reflect"
	"testing"
	"time"
)

func TestPostTagsUnmarshal(t *testing.T) {
	type tagsTest struct {
		XMLName xml.Name
		Tags    postTags `xml:"tags,attr"`
	}
	got := &tagsTest{}
	want := postTags{"foo", "bar", "baz"}
	body := `<post tags="foo bar baz"/>`
	err := xml.Unmarshal([]byte(body), got)
	if err != nil {
		t.Errorf("Failed to %v unmarshal body", err)
	}
	if !reflect.DeepEqual(want, got.Tags) {
		t.Errorf("Wanted %v, got %v", want, got.Tags)
	}
}

func TestUtcDateUnmarshal(t *testing.T) {
	type utcDateTest struct {
		XMLName xml.Name
		Created utcDate `xml:"created_at,attr"`
	}
	got := &utcDateTest{}
	want := utcDate{time.Date(1985, time.June, 27, 0, 0, 0, 0, time.UTC)}
	body := `<post created_at="1985-06-27"/>`
	err := xml.Unmarshal([]byte(body), got)
	if err != nil {
		t.Errorf("Failed to %v unmarshal body", err)
	}
	if !reflect.DeepEqual(want, got.Created) {
		t.Errorf("Wanted %v, got %v", want, got.Created)
	}
}

func TestNotesDateUnmarshal(t *testing.T) {
	type notesDateTest struct {
		XMLName xml.Name
		Created notesDate `xml:"created_at,attr"`
	}
	got := &notesDateTest{}
	want := notesDate{time.Date(1985, time.June, 27, 15, 13, 33, 0, time.UTC)}
	body := `<post created_at="1985-06-27 15:13:33"/>`
	err := xml.Unmarshal([]byte(body), got)
	if err != nil {
		t.Errorf("Failed to %v unmarshal body", err)
	}
	if !reflect.DeepEqual(want, got.Created) {
		t.Errorf("Wanted %v, got %v", want, got.Created)
	}
}
