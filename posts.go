package pinboard

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

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
	Tags        string    `xml:"tag,attr"`
	Extended    string    `xml:"extended,attr"`
	Date        time.Time `xml:"time,attr"`
	Shared      string    `xml:"shared,attr"`
}

type PostFilter struct {
	Tags []string
	Dt   time.Time
	Url  string
	Meta bool
}

type RecentPostsFilter struct {
	Tags  []string
	Count int
}

func ParseResponse(resp *http.Response) ([]Post, error) {
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
	return posts.Posts, err
}

func (p *Pinboard) GetPosts(pf PostFilter) ([]Post, error) {
	u, _ := url.Parse(APIBase + "posts/get")
	resp, err := p.Get(u.String())
	if err != nil {
		return nil, err
	}
	return ParseResponse(resp)
}

func (p *Pinboard) GetRecentPosts(rpf RecentPostsFilter) ([]Post, error) {
	u, err := url.Parse(APIBase + "posts/recent")

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
	resp, err := p.Get(u.String())
	if err != nil {
		return nil, err
	}

	return ParseResponse(resp)
}

func (p *Pinboard) DeletePost(dUrl string) error {
	u, err := url.Parse(APIBase + "posts/delete")
	if err != nil {
		return fmt.Errorf("Unable to parse delete url %v", err)
	}

	q := u.Query()
	q.Set("url", dUrl)
	u.RawQuery = q.Encode()

	_, err = p.Get(u.String())
	if err != nil {
		return fmt.Errorf("Error from delete request %v", err)
	}

	return nil
}

func (p *Pinboard) GetAllPosts() []Post {
	posts := make([]Post, 3)
	return posts
}
