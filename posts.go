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

type PostDates struct {
	Date     time.Time
	NumPosts int
}

type TagSuggestions struct {
	PopularTags     []string
	RecommendedTags []string
}

type RecentPostsFilter struct {
	Tags  []string
	Count int
}

func ParsePostsResponse(resp *http.Response) ([]Post, error) {
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

func (p *Pinboard) LastUpdate() (time.Time, error) {
	u, err := url.Parse(APIBase + "posts/update")

	resp, err := p.Get(u.String())
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

func (p *Pinboard) AddPost(pp Post) error {
	return nil
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

func (p *Pinboard) GetPosts(pf PostsFilter) ([]Post, error) {
	u, _ := url.Parse(APIBase + "posts/get")
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
	resp, err := p.Get(u.String())
	if err != nil {
		return nil, err
	}

	return ParsePostsResponse(resp)
}

func (p *Pinboard) GetPostDates(tags []string) ([]PostDates, error) {
	return nil, nil
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

	return ParsePostsResponse(resp)
}

func (p *Pinboard) GetAllPosts() []Post {
	posts := make([]Post, 3)
	return posts
}

func (p *Pinboard) GetTagSuggestions(postUrl string) TagSuggestions {
	ts := TagSuggestions{}
	return ts
}
