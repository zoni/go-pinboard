package pinboard

import (
	"fmt"
	"net/url"
)

// UserSecret returns the user's secret RSS key for viewing private feeds.
func (p *Pinboard) UserSecret() (string, error) {
	u, err := url.Parse(apiBase + "user/secret")
	if err != nil {
		return "", fmt.Errorf("Failed to parse UserSecret url: %v", err)
	}

	resp, err := p.get(u)
	if err != nil {
		return "", fmt.Errorf("Error from UserSecret request: %v", err)
	}

	tmp, err := parseResponse(resp, &result{})
	if err != nil {
		return "", fmt.Errorf("Failed to parse UserSecret response %v", err)
	}
	res := tmp.(*result)
	return res.Result, err
}

// UserApitoken returns the user's API token.
func (p *Pinboard) UserApiToken() (string, error) {
	u, err := url.Parse(apiBase + "user/api_token")
	if err != nil {
		return "", fmt.Errorf("Failed to parse UserApiToken url: %v", err)
	}

	resp, err := p.get(u)
	if err != nil {
		return "", fmt.Errorf("Error from UserApiToken request: %v", err)
	}

	tmp, err := parseResponse(resp, &result{})
	if err != nil {
		return "", fmt.Errorf("Failed to parse UserApiToken response %v", err)
	}
	res := tmp.(*result)
	return res.Result, err
}
