package pinboard

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Pinboard struct {
	User  string
	Token string
}

func (p *Pinboard) AuthQuery() string {
	return fmt.Sprintf("%s:%s", p.User, p.Token)
}

func (p *Pinboard) Get(uri string) (*http.Response, error) {
	u, err := url.Parse(uri)
	q := u.Query()
	q.Set("auth_token", p.AuthQuery())
	u.RawQuery = q.Encode()
	log.Println("Calling API with URL", u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	return resp, err
}
