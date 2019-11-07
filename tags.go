package pinboard

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
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

func ParseTagsResponse(resp *http.Response) (Tags, error) {
	t := Tags{}
	resp_body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return t, err
	}
	err = xml.Unmarshal(resp_body, &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (p *Pinboard) GetTags() (Tags, error) {
	u, err := url.Parse(APIBase + "tags/get")
	if err != nil {
		return Tags{}, fmt.Errorf("Failed to parse GetTags API URL: %v", err)
	}

	resp, err := p.Get(u.String())
	if err != nil {
		return Tags{}, err
	}

	return ParseTagsResponse(resp)
}

// DeleteTag deletes the given tag from a user's Pinboard account. There is no
// central store for tags, they are simply removed from every post in a user's
// account. This API endpoint has no meaningful response, so an error is returned
// only if the HTTP request fails.
func (p *Pinboard) DeleteTag(tag string) error {
	u, err := url.Parse(APIBase + "tags/delete")
	if err != nil {
		return fmt.Errorf("Failed to parse DeleteTag API URL: %v", err)
	}
	q := u.Query()

	if len(tag) < 1 || len(tag) > 255 {
		return fmt.Errorf("Tags must be between 1 and 255 characters in length")
	}
	q.Set("tag", tag)

	u.RawQuery = q.Encode()

	_, err = p.Get(u.String())
	if err != nil {
		return fmt.Errorf("Error from DeleteTag request %v", err)
	}

	return nil
}

type TagSuggestions struct {
	XMLName     xml.Name `xml:"suggested"`
	Popular     []string `xml:"popular"`
	Recommended []string `xml:"recommended"`
}

func ParseSuggestedTagsResponse(resp *http.Response) (TagSuggestions, error) {
	ts := TagSuggestions{}
	resp_body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ts, err
	}
	err = xml.Unmarshal(resp_body, &ts)
	if err != nil {
		return ts, err
	}
	return ts, nil
}

func (p *Pinboard) GetTagSuggestions(postUrl string) (TagSuggestions, error) {
	u, _ := url.Parse(APIBase + "posts/suggest")
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

	resp, err := p.Get(u.String())
	if err != nil {
		return TagSuggestions{}, err
	}

	return ParseSuggestedTagsResponse(resp)
}
