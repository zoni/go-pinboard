package pinboard

import (
	"encoding/xml"
	"io/ioutil"
	"log"
)

type Posts struct {
	XMLName xml.Name `xml:"posts"`
	User    string   `xml:"user,attr"`
	Date    string   `xml:"dt,attr"`
	Posts   []Post   `xml:"post"`
}

type Post struct {
	XMLName     xml.Name `xml:"post"`
	Url         string   `xml:"href,attr"`
	Description string   `xml:"description,attr"`
	Hash        string   `xml:"hash,attr"`
	Tags        string   `xml:"tag,attr"`
	Extended    string   `xml:"extended,attr"`
	Date        string   `xml:"time,attr"`
	Shared      string   `xml:"shared,attr"`
}

func (p *Pinboard) GetRecentPosts() ([]Post, error) {
	posts := &Posts{}
	url := "https://api.pinboard.in/v1/posts/recent"
	resp, err := p.Get(url)
	if err != nil {
		return nil, err
	}
	resp_body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	log.Println("Got resp_body length ", len(resp_body))
	log.Println(string(resp_body))
	err = xml.Unmarshal(resp_body, &posts)
	if err != nil {
		return nil, err
	}
	log.Println("Got user ", posts.User)
	return posts.Posts, err
}

func (p *Pinboard) GetAllPosts() []Post {
	posts := make([]Post, 3)
	return posts
}
