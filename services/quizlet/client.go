package quizlet

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
)

// Client represents a Quizlet client
type Client struct {
	HTTP    *http.Client
	Headers map[string]string
}

// New creats a new Quizlet client
func New() (c *Client, err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}

	return &Client{
		HTTP: &http.Client{Jar: jar},
		Headers: map[string]string{
			"Accept":           "application/json",
			"Accept-Language":  "en-US,en;q=0.9",
			"Content-Type":     "application/json",
			"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.67 Safari/537.36",
			"Origin":           Base,
			"X-Requested-With": "XMLHttpRequest",
		},
	}, nil
}

// Request makes a request to the provided Quizlet endpoint
func (c *Client) Request(method, endpoint string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, Base+endpoint, body)
	if err != nil {
		return
	}

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	resp, err = c.HTTP.Do(req)
	if err != nil {
		return
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "qtkn" {
			c.Headers["cs-token"] = cookie.Value
			break
		}
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return resp, errors.New(resp.Status)
	}

	return
}

// Login makes a login request to Quizlet with the provided username and password
func (c *Client) Login(username, password string) (err error) {
	_, err = c.Request(http.MethodGet, "/login", nil)
	if err != nil {
		return
	}

	body, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return
	}

	_, err = c.Request(http.MethodPost, "/login", bytes.NewReader(body))
	return
}
