package pinboard

import (
	"fmt"
	"net/url"
)

func (p *Pinboard) UserSecret() (string, error) {
	u, err := url.Parse(apiBase + "user/secret")
	if err != nil {
		return "", fmt.Errorf("Failed to parse UserSecret url: %v", err)
	}

	resp, err := p.get(u)
	if err != nil {
		return "", fmt.Errorf("Error from UserSecret request: %v", err)
	}

	tmp, err := parseResponse(resp, &Result{})
	if err != nil {
		return "", fmt.Errorf("Failed to parse UserSecret response %v", err)
	}
	result := tmp.(*Result)
	return result.Result, err
}

func (p *Pinboard) UserApiToken() (string, error) {
	u, err := url.Parse(apiBase + "user/api_token")
	if err != nil {
		return "", fmt.Errorf("Failed to parse UserApiToken url: %v", err)
	}

	resp, err := p.get(u)
	if err != nil {
		return "", fmt.Errorf("Error from UserApiToken request: %v", err)
	}

	tmp, err := parseResponse(resp, &Result{})
	if err != nil {
		return "", fmt.Errorf("Failed to parse UserApiToken response %v", err)
	}
	result := tmp.(*Result)
	return result.Result, err
}
