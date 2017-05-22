package tweet

import (
	"fmt"
	"net/http"
)

const tweetLength = 140

// Tweet is a custom struct containing a HTTP client and request.
type Tweet struct {
	client  HTTPClient
	request NewRequest
}

// NewTweet returns a new Tweet with all of its required fields.
func NewTweet(client HTTPClient, req NewRequest) Tweet {
	return Tweet{
		client:  client,
		request: req,
	}
}

// Send takes Twitter app (https://apps.titter.com) credentials and passes
// them to an OAuth generator along with a given status. It then attempts
// to send a POST request to the Twitter API and handle a response.
func (t Tweet) Send(c *Config, status string) (string, error) {
	if len(status) > tweetLength {
		return status, fmt.Errorf("tweet exceeds %d character limit", tweetLength)
	}

	oa := NewOAuthDetails(c, status)

	req, err := t.request(http.MethodPost, fmt.Sprintf(apiURL+apiVersion+statusURI+"?status=%s", encodeStatus(&status)), nil)
	if err != nil {
		return status, fmt.Errorf("error building request: %v", err)
	}
	req.Header.Set(authHeader, fmt.Sprintf("%s", oa))

	res, err := t.client.Do(req)
	if err != nil {
		return status, fmt.Errorf("%v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return status, fmt.Errorf("invalid status code: %s", res.Status)
	}

	return fmt.Sprintf("Tweet successfully sent!"), nil
}