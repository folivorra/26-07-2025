package validation

import (
	"net/http"
	"time"
)

type HTTPValidator struct {
	client *http.Client
}

var _ FileValidator = (*HTTPValidator)(nil)

func NewHTTPValidator(timeout time.Duration) *HTTPValidator {
	return &HTTPValidator{
		client: &http.Client{
			Timeout: timeout,
		},
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
