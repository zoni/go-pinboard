package pinboard

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"regexp"
)

type notes struct {
	XMLName xml.Name `xml:"notes" json:"-"`
	Notes   []Note   `xml:"note"`
}

// A Note can represent either a note in the NotesList, or a single note. Depending
// on which type of note is called for different fields will be populated (ex: Created
// and Updated are only returned in the NotesList. Text is only returned by NotesGet).
// Text may be contain newlines.
type Note struct {
	XMLName xml.Name  `xml:"note" json:"-"`
	ID      string    `xml:"id,attr"`
	Title   string    `xml:"title"`
	Hash    string    `xml:"hash"`
	Created notesDate `xml:"created_at"`
	Updated notesDate `xml:"updated_at"`
	Length  int       `xml:"length"`
	Text    string    `xml:"text"`
}

// NotesList returns a list of the user's notes.
func (p *Pinboard) NotesList() ([]Note, error) {
	u, err := url.Parse(apiBase + "notes/list")
	if err != nil {
		return []Note{}, fmt.Errorf("Failed to parse Notes list API URL: %v", err)
	}

	resp, err := p.get(u)
	if err != nil {
		return []Note{}, err
	}

	tmp, err := parseResponse(resp, &notes{})
	if err != nil {
		return []Note{}, fmt.Errorf("Failed to parse Notes response: %v", err)
	}
	no := tmp.(*notes)
	return no.Notes, err
}

// NotesGet returns a single Note.
func (p *Pinboard) NotesGet(noteID string) (Note, error) {
	if m, _ := regexp.Match("[a-z0-9]{20}", []byte(noteID)); !m {
		return Note{}, fmt.Errorf("Note ID must be a 20 character sha1 hash")
	}

	u, err := url.Parse(apiBase + "notes/" + noteID)
	if err != nil {
		return Note{}, fmt.Errorf("Failed to parse note URL: %v", err)
	}

	resp, err := p.get(u)
	if err != nil {
		return Note{}, fmt.Errorf("Error getting note: %v", err)
	}

	tmp, err := parseResponse(resp, &Note{})
	if err != nil {
		return Note{}, fmt.Errorf("Failed to parse Note response: %v", err)
	}
	note := tmp.(*Note)
	return *note, err
}
