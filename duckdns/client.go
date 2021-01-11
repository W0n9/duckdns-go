package duckdns

import (
	"context"
	"errors"
	"log"
	"net/http"
)

const (
	// Version identifies the current library version.
	// This is a pro-forma convention given that Go dependencies
	// tends to be fetched directly from the repo.
	// It is also used in the user-agent identify the client.
	Version = "1.0.0"

	// defaultBaseURL to the DNSimple production API.
	defaultBaseURL = "https://www.duckdns.org/"

	// userAgent represents the default user agent used
	// when no other user agent is set.
	defaultUserAgent = "duckdns-go/" + Version
)

// Client represents a client to the DuckDNS API.
type Client struct {
	httpClient *http.Client
	BaseURL    string
	UserAgent  string

	Auth    *AuthService
	Domains *DomainsService
	Records *RecordsService

	Debug bool
}

func NewClient(httpClient *http.Client) *Client {
	c := &Client{
		httpClient: httpClient,
		BaseURL:    defaultBaseURL,
		UserAgent:  defaultUserAgent}
	c.Auth = &AuthService{client: c}
	c.Domains = &DomainsService{client: c}
	c.Records = &RecordsService{client: c}
	return c
}

func (c *Client) SetUserAgent(ua string) {
	c.UserAgent = ua
}

func (c *Client) SetDebug(debug bool) {
	c.Debug = debug
}

func (c *Client) get(ctx context.Context, path string) (*http.Response, error) {
	return c.makeRequest(ctx, http.MethodGet, path)
}

func (c *Client) makeRequest(ctx context.Context, method, path string) (*http.Response, error) {

	req, err := c.newRequest(method, path)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		log.Printf("Request (%v): %#v", req.URL, req)
	}

	resp, err := c.request(ctx, req)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		log.Printf("Response: %#v", resp)
	}

	return resp, nil
}

func (c *Client) newRequest(method, path string) (*http.Request, error) {
	url := c.BaseURL + path

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = make(http.Header)
	req.Header.Add("User-Agent", c.UserAgent)

	return req, err
}

func (c *Client) request(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	return resp, err
}
