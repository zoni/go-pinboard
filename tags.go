package pinboard

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
)

type tags struct {
	XMLName xml.Name `xml:"tags"`
	Tags    []Tag    `xml:"tag"`
}

// A Tag in Pinboard is a simple string applied to posts for organizational purposes.
// Users can create Private Tags (tags only visible to the posting user) by prepending
// the tag with a period (.). Tags returned from the Tags endpoints contain a count of
// how often they're used.
type Tag struct {
	XMLName xml.Name `xml:"tag"`
	Count   int      `xml:"count,attr"`
	Tag     string   `xml:"tag,attr"`
}

// TagsGet returns a list of []Tag corresponding to the tags in the user's account.
func (p *Pinboard) TagsGet() ([]Tag, error) {
	u, err := url.Parse(apiBase + "tags/get")
	if err != nil {
		return []Tag{}, fmt.Errorf("Failed to parse Tags API URL: %v", err)
	}

	resp, err := p.get(u)
	if err != nil {
		return []Tag{}, err
	}

	tmp, err := parseResponse(resp, &tags{})
	if err != nil {
		return []Tag{}, fmt.Errorf("Failed to parse Tags response %v", err)
	}
	t := tmp.(*tags)

	return t.Tags, err
}

// TagsDelete deletes the given tag from a user's Pinboard account. There is no
// central store for tags, they are simply removed from every post in a user's
// account. This API endpoint has no meaningful response, so an error is returned
// only if the HTTP request fails.
func (p *Pinboard) TagsDelete(tag string) error {
	u, err := url.Parse(apiBase + "tags/delete")
	if err != nil {
		return fmt.Errorf("Failed to parse TagsDelete API URL: %v", err)
	}
	q := u.Query()

	if len(tag) < 1 || len(tag) > 255 {
		return fmt.Errorf("Tags must be between 1 and 255 characters in length")
	}
	q.Set("tag", tag)

	u.RawQuery = q.Encode()

	_, err = p.get(u)
	if err != nil {
		return fmt.Errorf("Error from TagsDelete request %v", err)
	}

	return nil
}

// TagsRename renames a tag by changing that tag on every post in the user's account.
func (p *Pinboard) TagsRename(old, new string) error {
	u, err := url.Parse(apiBase + "tags/rename")
	if err != nil {
		return fmt.Errorf("Failed to parse TagsRename API URL: %v", err)
	}
	q := u.Query()

	if len(old) < 1 || len(new) < 1 {
		return fmt.Errorf("Both old and new tag must not be empty string for TagsRename")
	}

	q.Set("old", old)
	q.Set("new", new)

	u.RawQuery = q.Encode()

	_, err = p.get(u)
	if err != nil {
		return fmt.Errorf("Error from TagsRename request %v", err)
	}

	return nil
}

// Pinboard returns two types of tag suggestions, popular tags are tags from the community.
// Recommended tags are based on the user's existing tags
type TagSuggestions struct {
	XMLName     xml.Name `xml:"suggested"`
	Popular     []string `xml:"popular"`
	Recommended []string `xml:"recommended"`
}

// TagsSuggestions returns tag suggestions for the given URL. Note: Currently only recommended
// tags are actually returned from the Pinboard API
func (p *Pinboard) TagsSuggestions(postUrl string) (TagSuggestions, error) {
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

	resp, err := p.get(u)
	if err != nil {
		return TagSuggestions{}, err
	}

	tmp, err := parseResponse(resp, &TagSuggestions{})
	if err != nil {
		return TagSuggestions{}, fmt.Errorf("Failed to parse TagsSuggestions response %v", err)
	}
	t := tmp.(*TagSuggestions)

	return *t, err
}
