package pinboard

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var apiBase = "https://api.pinboard.in/v1/"

type Pinboard struct {
	User  string
	Token string
}

func (p *Pinboard) authQuery() string {
	return fmt.Sprintf("%s:%s", p.User, p.Token)
}

func (p *Pinboard) Get(uri string) (*http.Response, error) {
	u, err := url.Parse(uri)
	fmt.Println("Calling API with ", u.String())
	q := u.Query()
	q.Set("auth_token", p.authQuery())
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		resp_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Error reading error body: %v", err)
		}
		resp.Body.Close()
		return nil, fmt.Errorf("Error from Pinboard API: %v", string(resp_body))
	}
	return resp, err
}
