package quizlet

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
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
func (c *Client) Request(method, endpoint string, body io.Reader) (res *http.Response, err error) {
	req, err := http.NewRequest(method, Base+endpoint, body)
	if err != nil {
		return
	}

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	res, err = c.HTTP.Do(req)
	if err != nil {
		return
	}

	for _, cookie := range res.Cookies() {
		if cookie.Name == "qtkn" {
			c.Headers["CS-Token"] = cookie.Value
			break
		}
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return res, errors.New(res.Status)
	}

	return
}

// RequestJSON makes requests with the value serialized using JSON encoding
func (c *Client) RequestJSON(method, endpoint string, value interface{}) (res *http.Response, err error) {
	body := new(bytes.Buffer)
	if err = json.NewEncoder(body).Encode(value); err != nil {
		return
	}

	return c.Request(method, endpoint, body)
}

// RequestForm makes requests with the values serialized using form URL encoding
func (c *Client) RequestForm(method, endpoint string, values url.Values) (*http.Response, error) {
	return c.Request(method, endpoint, strings.NewReader(values.Encode()))
}

// Login makes a login request to Quizlet with the provided username and password
func (c *Client) Login(username, password string) (err error) {
	_, err = c.Request(http.MethodGet, EndpointLogin, nil)
	if err != nil {
		return
	}

	_, err = c.RequestJSON(http.MethodPost, EndpointLogin, map[string]string{
		"username": username,
		"password": password,
	})
	return
}

// SessionID requests Quizlet to generate a session ID
func (c *Client) SessionID(id string, mode StudyMode) (sessionID int, err error) {
	res, err := c.RequestForm(http.MethodPost, EndpointSessions, url.Values{
		"cstoken":       {c.Headers["CS-Token"]},
		"mode":          {strconv.Itoa(int(mode))},
		"selectedOnly":  {"0"},
		"studyableId":   {id},
		"studyableType": {"1"},
	})
	if err != nil {
		return
	}
	defer res.Body.Close()

	var body struct {
		ID int `json:"id"`
	}
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return
	}

	return body.ID, nil
}
