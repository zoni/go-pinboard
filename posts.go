package pinboard

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var validSchemes = []string{
	"http",
	"https",
	"javascript",
	"mailto",
	"ftp",
	"file",
	"feed",
}

type Posts struct {
	XMLName xml.Name  `xml:"posts"`
	User    string    `xml:"user,attr"`
	Date    time.Time `xml:"dt,attr"`
	Posts   []Post    `xml:"post"`
}

type Post struct {
	XMLName     xml.Name  `xml:"post"`
	Url         string    `xml:"href,attr"`
	Description string    `xml:"description,attr"`
	Hash        string    `xml:"hash,attr"`
	Tags        PostTags  `xml:"tag,attr"`
	Extended    string    `xml:"extended,attr"`
	Date        time.Time `xml:"time,attr"`
	Shared      string    `xml:"shared,attr"`
	Meta        string    `xml:"meta,attr"`
}

type PostsLastUpdate struct {
	XMLName    xml.Name  `xml:"update"`
	UpdateTime time.Time `xml:"time,attr"`
}

type PostsFilter struct {
	Tags []string
	Date time.Time
	Url  string
	Meta bool
}

type RecentPostsFilter struct {
	Tags  []string
	Count int
}

func parsePostsResponse(resp *http.Response) ([]Post, error) {
	posts := &Posts{}
	resp_body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(resp_body, &posts)
	if err != nil {
		return nil, err
	}
	return posts.Posts, nil
}

func parsePostDatesResponse(resp *http.Response) ([]PostDate, error) {
	postdates := &PostDates{}
	resp_body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(resp_body, &postdates)
	if err != nil {
		return nil, err
	}
	return postdates.PostDates, nil
}

func (p *Pinboard) LastUpdate() (time.Time, error) {
	u, err := url.Parse(apiBase + "posts/update")

	resp, err := p.Get(u)
	if err != nil {
		return time.Time{}, err
	}

	update := &PostsLastUpdate{}
	resp_body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return time.Time{}, err
	}
	err = xml.Unmarshal(resp_body, &update)
	if err != nil {
		return time.Time{}, err
	}

	return update.UpdateTime, err
}

func (p *Pinboard) AddPost(pp Post, keep bool, toread bool) error {
	u, err := url.Parse(apiBase + "posts/add")
	q := u.Query()

	if len(pp.Url) < 1 {
		return fmt.Errorf("AddPost requires a URL")
	}
	pu, err := url.Parse(pp.Url)
	if err != nil {
		return fmt.Errorf("Error parsing AddPost URL ", err)
	}
	validScheme := false
	for _, v := range validSchemes {
		if strings.ToLower(pu.Scheme) == v {
			validScheme = true
		}
	}
	if !validScheme {
		return fmt.Errorf("Invalid scheme for Pinboard URL. Scheme must be one of %v", validSchemes)
	}

	q.Set("url", pp.Url)

	if len(pp.Description) < 1 || len(pp.Description) > 255 {
		return fmt.Errorf("Pinboard URL descriptions must be between 1 and 255 characters long")
	}

	q.Set("description", pp.Description)

	if len(pp.Extended) > 0 {
		if len(pp.Extended) > 65536 {
			return fmt.Errorf("Pinboard extended descriptions must be less than 65536 characters long")
		}
		q.Set("extended", pp.Extended)
	}

	if len(pp.Tags) > 0 {
		if len(pp.Tags) > 100 {
			return fmt.Errorf("Pinboard posts may only have up to 100 tags")
		}
		for _, v := range pp.Tags {
			q.Add("tag", v)
		}
	}

	if !pp.Date.IsZero() {
		q.Set("dt", pp.Date.UTC().Format(time.RFC3339))
	}

	if keep {
		q.Set("replace", "no")
	}

	if toread {
		q.Set("toread", "yes")
	}

	if len(pp.Shared) > 0 {
		lshared := strings.ToLower(pp.Shared)
		if lshared == "yes" || lshared == "no" {
			q.Set("shared", lshared)
		} else {
			return fmt.Errorf("Shared must be either \"yes\" or \"no\"")
		}
	}

	u.RawQuery = q.Encode()

	_, err = p.Get(u)
	if err != nil {
		return fmt.Errorf("Error adding post: %v", err)
	}

	return nil
}

func (p *Pinboard) DeletePost(dUrl string) error {
	u, err := url.Parse(apiBase + "posts/delete")
	if err != nil {
		return fmt.Errorf("Unable to parse DeletePost url %v", err)
	}

	q := u.Query()
	q.Set("url", dUrl)
	u.RawQuery = q.Encode()

	_, err = p.Get(u)
	if err != nil {
		return fmt.Errorf("Error from DeletePost request %v", err)
	}

	return nil
}

func (p *Pinboard) Posts(pf PostsFilter) ([]Post, error) {
	u, _ := url.Parse(apiBase + "posts/get")
	q := u.Query()

	// Filters
	if len(pf.Tags) > 0 {
		if len(pf.Tags) > 3 {
			return nil, fmt.Errorf("PostsFilter cannot accept more than 3 tags")
		}
		for _, t := range pf.Tags {
			q.Add("tag", t)
		}
	}

	if !pf.Date.IsZero() {
		q.Set("dt", pf.Date.Format("2006-01-02"))
	}

	if len(pf.Url) > 0 {
		q.Set("url", pf.Url)
	}

	if pf.Meta {
		q.Set("meta", "yes")
	}
	u.RawQuery = q.Encode()

	// Get posts
	resp, err := p.Get(u)
	if err != nil {
		return nil, err
	}

	return parsePostsResponse(resp)
}

type PostDates struct {
	XMLName   xml.Name   `xml:"dates"`
	User      string     `xml:"user,attr"`
	Tag       string     `xml:"tag,attr"`
	PostDates []PostDate `xml:"date"`
}

type PostDate struct {
	XMLName xml.Name `xml:"date"`
	Date    UTCDate  `xml:"date,attr"`
	Count   int      `xml:"count,attr"`
}

// PostDates returns an array of posts-per-day optionally filtered by a
// given tag. Contrary to Pinboard's API documentation only a single tag is
// accepted for filtering.
func (p *Pinboard) PostDates(tag string) ([]PostDate, error) {
	u, err := url.Parse(apiBase + "posts/dates")
	q := u.Query()

	if len(tag) > 0 {
		q.Set("tag", tag)
	}
	u.RawQuery = q.Encode()

	resp, err := p.Get(u)
	if err != nil {
		return nil, err
	}

	return parsePostDatesResponse(resp)
}

func (p *Pinboard) RecentPosts(rpf RecentPostsFilter) ([]Post, error) {
	u, err := url.Parse(apiBase + "posts/recent")

	// Filters
	q := u.Query()
	if rpf.Count != 0 {
		if rpf.Count < 0 || rpf.Count > 100 {
			return nil, fmt.Errorf("RecentPostsFilter count must be between 0 and 100")
		}
		q.Set("count", fmt.Sprintf("%d", rpf.Count))
	}

	if len(rpf.Tags) > 0 {
		if len(rpf.Tags) > 3 {
			return nil, fmt.Errorf("RecentPostsFilter cannot accept more than 3 tags")
		}
		for _, t := range rpf.Tags {
			q.Add("tag", t)
		}
	}
	u.RawQuery = q.Encode()

	// Get posts
	resp, err := p.Get(u)
	if err != nil {
		return nil, err
	}

	return parsePostsResponse(resp)
}

type AllPostsFilter struct {
	Tags    []string
	Start   int
	Results int
	From    time.Time
	To      time.Time
	Meta    bool
}

func (p *Pinboard) AllPosts(apf AllPostsFilter) ([]Post, error) {
	u, _ := url.Parse(apiBase + "posts/all")
	q := u.Query()

	// Filters
	if len(apf.Tags) > 0 {
		if len(apf.Tags) > 3 {
			return nil, fmt.Errorf("AllPosts can not accept more than 3 tags")
		}
		for _, t := range apf.Tags {
			q.Add("tag", t)
		}
	}

	if apf.Start > 0 {
		q.Set("start", fmt.Sprintf("%d", apf.Start))
	}

	if apf.Results > 0 {
		q.Set("results", fmt.Sprintf("%d", apf.Results))
	}

	if !apf.From.IsZero() {
		q.Set("fromdt", apf.From.UTC().Format(time.RFC3339))
	}

	if !apf.To.IsZero() {
		q.Set("fromdt", apf.To.UTC().Format(time.RFC3339))
	}

	if apf.Meta {
		q.Set("meta", "yes")
	}

	u.RawQuery = q.Encode()
	resp, err := p.Get(u)
	if err != nil {
		return nil, fmt.Errorf("AllPosts failed to retrieve: %v", err)
	}

	return parsePostsResponse(resp)
}
