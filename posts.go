package pinboard

import (
	"encoding/xml"
	"io/ioutil"
	//"log"
	"net/http"
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
	tags []string
	dt   time.Time
	url  string
	meta bool
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
	url := "https://api.pinboard.in/v1/posts/get"
	resp, err := p.Get(url)
	if err != nil {
		return nil, err
	}
	return ParseResponse(resp)
}

func (p *Pinboard) GetRecentPosts() ([]Post, error) {
	url := "https://api.pinboard.in/v1/posts/recent"
	resp, err := p.Get(url)
	if err != nil {
		return nil, err
	}
	return ParseResponse(resp)
}

func (p *Pinboard) GetAllPosts() []Post {
	posts := make([]Post, 3)
	return posts
}
