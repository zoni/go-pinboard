package pinboard

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
)

type Tags struct {
	XMLName xml.Name `xml:"tags"`
	Tags    []Tag    `xml:"tag"`
}

type Tag struct {
	XMLName xml.Name `xml:"tag"`
	Count   int      `xml:"count,attr"`
	Tag     string   `xml:"tag,attr"`
}

func (p *Pinboard) Tags() ([]Tag, error) {
	u, err := url.Parse(apiBase + "tags/get")
	if err != nil {
		return []Tag{}, fmt.Errorf("Failed to parse Tags API URL: %v", err)
	}

	resp, err := p.Get(u)
	if err != nil {
		return []Tag{}, err
	}

	tmp, err := parseResponse(resp, &Tags{})
	if err != nil {
		return []Tag{}, fmt.Errorf("Failed to parse Tags response %v", err)
	}
	t := tmp.(*Tags)

	return t.Tags, err
}

// DeleteTag deletes the given tag from a user's Pinboard account. There is no
// central store for tags, they are simply removed from every post in a user's
// account. This API endpoint has no meaningful response, so an error is returned
// only if the HTTP request fails.
func (p *Pinboard) DeleteTag(tag string) error {
	u, err := url.Parse(apiBase + "tags/delete")
	if err != nil {
		return fmt.Errorf("Failed to parse DeleteTag API URL: %v", err)
	}
	q := u.Query()

	if len(tag) < 1 || len(tag) > 255 {
		return fmt.Errorf("Tags must be between 1 and 255 characters in length")
	}
	q.Set("tag", tag)

	u.RawQuery = q.Encode()

	_, err = p.Get(u)
	if err != nil {
		return fmt.Errorf("Error from DeleteTag request %v", err)
	}

	return nil
}

func (p *Pinboard) RenameTag(old, new string) error {
	u, err := url.Parse(apiBase + "tags/rename")
	if err != nil {
		return fmt.Errorf("Failed to parse RenameTag API URL: %v", err)
	}
	q := u.Query()

	if len(old) < 1 || len(new) < 1 {
		return fmt.Errorf("Both old and new tag must not be empty string for RenameTag")
	}

	q.Set("old", old)
	q.Set("new", new)

	u.RawQuery = q.Encode()

	_, err = p.Get(u)
	if err != nil {
		return fmt.Errorf("Error from RenameTag request %v", err)
	}

	return nil
}

type TagSuggestions struct {
	XMLName     xml.Name `xml:"suggested"`
	Popular     []string `xml:"popular"`
	Recommended []string `xml:"recommended"`
}

func (p *Pinboard) TagSuggestions(postUrl string) (TagSuggestions, error) {
	u, _ := url.Parse(apiBase + "posts/suggest")
	q := u.Query()

	pu, _ := url.Parse(postUrl)
	validScheme := false
	for _, v := range validSchemes {
		if strings.ToLower(pu.Scheme) == v {
			validScheme = true
		}
	}
	if !validScheme {
		return TagSuggestions{}, fmt.Errorf("Invalid scheme for Pinboard URL. Scheme must be one of %v", validSchemes)
	}

	q.Set("url", postUrl)
	u.RawQuery = q.Encode()

	resp, err := p.Get(u)
	if err != nil {
		return TagSuggestions{}, err
	}

	tmp, err := parseResponse(resp, &TagSuggestions{})
	if err != nil {
		return TagSuggestions{}, fmt.Errorf("Failed to parse TagSuggestions response %v", err)
	}
	t := tmp.(*TagSuggestions)

	return *t, err
}
