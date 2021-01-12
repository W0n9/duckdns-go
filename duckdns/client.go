package duckdns

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	Version = "1.0.0"

	defaultBaseURL = "https://www.duckdns.org"
	domainStub     = "/update?domains="
	tokenStub      = "&token="
	ip4Stub        = "&ip="
	ip6Stub        = "&ipv6="
	txtStub        = "&txt="
	verboseStub    = "&verbose="
	clearStub      = "&clear="

	defaultUserAgent = "duckdns-go/" + Version
)

type Config struct {
	DomainNames []string
	Token       string
	IPv4        string
	IPv6        string
}

func (c *Config) Valid() bool {
	if c.Token != "" && len(c.DomainNames) > 0 {
		return true
	}
	return false
}

type Client struct {
	httpClient *http.Client
	BaseURL    string
	UserAgent  string

	Config *Config

	Verbose bool
}

func NewClient(httpClient *http.Client, config *Config) *Client {
	if !config.Valid() {
		log.Fatal("Configuration is not valid")
	}

	c := &Client{httpClient: httpClient, BaseURL: defaultBaseURL, UserAgent: defaultUserAgent}
	c.Config = config
	return c
}

func (c *Client) SetUserAgent(ua string) {
	c.UserAgent = ua
}

func (c *Client) SetVerbose(verbose bool) {
	c.Verbose = verbose
}

func (c *Client) makeGetRequest(ctx context.Context, path string) (*http.Response, error) {

	req, err := c.newRequest(http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	if c.Verbose {
		log.Printf("Request (%v): %#v", req.URL, req)
	}

	resp, err := c.request(ctx, req)
	if err != nil {
		return nil, err
	}

	if c.Verbose {
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

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body response")
		return nil, err
	}

	if strings.Contains(string(bodyBytes), "KO") {
		log.Printf("Error returned from DuckDNS")
		return nil, err
	}

	return resp, err
}

//Update IPv4 and/or without IP address
func (c *Client) UpdateIP() (*http.Response, error) {
	subdomains := strings.Join(c.Config.DomainNames, ",")
	url := fmt.Sprintf("%s%s%s%s%s%s%s%s", c.BaseURL, domainStub, subdomains, tokenStub, c.Config.Token, ip4Stub, verboseStub, strconv.FormatBool(c.Verbose))
	resp, err := c.makeGetRequest(context.Background(), url)

	return resp, err
}

//Update IPv4 and/or with IP address
func (c *Client) UpdateIPWithValues(ipv4, ipv6 string) (*http.Response, error) {
	subdomains := strings.Join(c.Config.DomainNames, ",")
	url := fmt.Sprintf("%s%s%s%s%s%s", c.BaseURL, domainStub, subdomains, tokenStub, c.Config.Token, ip4Stub)
	if ipv6 == "" {
		url = fmt.Sprintf("%s%s%s%s", url, ipv4, verboseStub, strconv.FormatBool(c.Verbose))
	} else {
		url = fmt.Sprintf("%s%s%s%s%s%s", url, ipv4, ip6Stub, ipv6, verboseStub, strconv.FormatBool(c.Verbose))
	}
	resp, err := c.makeGetRequest(context.Background(), url)

	return resp, err
}

//Clear IP
func (c *Client) ClearIP() (*http.Response, error) {
	subdomains := strings.Join(c.Config.DomainNames, ",")
	url := fmt.Sprintf("%s%s%s%s%s%s%s%s%s", c.BaseURL, domainStub, subdomains, tokenStub, c.Config.Token, verboseStub, strconv.FormatBool(c.Verbose), clearStub, "true")
	resp, err := c.makeGetRequest(context.Background(), url)

	return resp, err
}

//Update TXT record
func (c *Client) UpdateRecord(record string) (*http.Response, error) {
	subdomains := strings.Join(c.Config.DomainNames, ",")
	url := fmt.Sprintf("%s%s%s%s%s%s%s%s%s", c.BaseURL, domainStub, subdomains, tokenStub, c.Config.Token, txtStub, record, verboseStub, strconv.FormatBool(c.Verbose))
	resp, err := c.makeGetRequest(context.Background(), url)

	return resp, err
}

//Clear TXT record
func (c *Client) ClearRecord(record string) (*http.Response, error) {
	subdomains := strings.Join(c.Config.DomainNames, ",")
	url := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s", c.BaseURL, domainStub, subdomains, tokenStub, c.Config.Token, txtStub, record, verboseStub, strconv.FormatBool(c.Verbose), clearStub, "true")
	resp, err := c.makeGetRequest(context.Background(), url)

	return resp, err
}

//Get TXT record
func (c *Client) GetRecord() (string, error) {
	subdomains := c.Config.DomainNames[0]
	txt, err := net.LookupTXT(subdomains)
	if err != nil {
		return "", fmt.Errorf("Unable to get txt record, %v", err)
	}

	if len(txt) == 0 {
		return "", nil
	}

	//duckdns should have only 1 record
	return txt[0], nil
}
