package pinboard

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"regexp"
)

type notes struct {
	XMLName xml.Name `xml:"notes"`
	Notes   []Note   `xml:"note"`
}

type Note struct {
	XMLName xml.Name  `xml:"note"`
	ID      string    `xml:"id,attr"`
	Title   string    `xml:"title"`
	Hash    string    `xml:"hash"`
	Created notesDate `xml:"created_at"`
	Updated notesDate `xml:"updated_at"`
	Length  int       `xml:"length"`
	Text    string    `xml:"text"`
}

func (p *Pinboard) Notes() ([]Note, error) {
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

func (p *Pinboard) Note(noteID string) (Note, error) {
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
