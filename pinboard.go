// Package pinboard is an implementation of the Pinboard V1 API (https://pinboard.in/api)
//
// This package implements  the API as documented, though some fixes have been made to
// maintain type cohesion. See method comments for exceptions to the API documentation.
package pinboard

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var apiBase = "https://api.pinboard.in/v1/"

// A Pinboard represents a client for the Pinboard V1 API. Authentication can use passwords
// or tokens. Token auth is recommended for good password hygiene.
type Pinboard struct {
	User     string
	Password string
	Token    string
}

func (p *Pinboard) authQuery(u *url.URL) error {
	if len(p.User) < 1 {
		return fmt.Errorf("Pinboard requires a Username and either a Password or Token for authentication")
	}

	if len(p.Token) < 1 {
		if len(p.Password) < 1 {
			return fmt.Errorf("Pinboard requires either a Password or Token for authentication")
		}
		u.User = url.UserPassword(p.User, p.Password)
		return nil
	}

	q := u.Query()
	q.Set("auth_token", fmt.Sprintf("%s:%s", p.User, p.Token))
	u.RawQuery = q.Encode()
	return nil
}

func (p *Pinboard) get(u *url.URL) (*http.Response, error) {
	err := p.authQuery(u)
	if err != nil {
		return nil, fmt.Errorf("Pinboard failed to generate an auth query param", err)
	}

	fmt.Println("Calling API with", u.String())
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

func parseResponse(resp *http.Response, to interface{}) (interface{}, error) {
	resp_body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(resp_body, to)
	if err != nil {
		return nil, err
	}
	return to, nil
}

type Result struct {
	XMLName xml.Name `xml:"result"`
	Result  string   `xml:",innerxml"`
}
