// Package pinboard is an implementation of the Pinboard V1 API (https://pinboard.in/api)
//
// This package communicates with the Pinboard API in XML for novelty's sake. All authenticated requests
// happen via auth token (except for the initial token retrieval, if a password is supplied). This package
// implements the API as documented, though some small fixes have been made to maintain type cohesion.
// See method comments for exceptions to the API documentation.
//
// Note that the Pinboard API attempts to faithfully re-implement the del.icio.us API and does not behave how
// a modern API may be expected to behave. URLs are not RESTful; every operation is done via GET requests.
// Response/ status is communicated in the response body; only API/HTTP errors (such as throttling or server
// issues) cause an HTTP status code > 299.
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

func (p *Pinboard) Get(u *url.URL) (*http.Response, error) {
	fmt.Println("Calling API with", u.String())
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
