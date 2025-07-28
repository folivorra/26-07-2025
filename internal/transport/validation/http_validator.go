package validation

import "net/http"

type HTTPValidator struct {
	client *http.Client
}

func NewHTTPValidator(client *http.Client) *HTTPValidator {
	return &HTTPValidator{
		client: client,
	}
}

func (v *HTTPValidator) IsReachable(url string) bool {
	resp, err := v.client.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
