package duckdns

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setupMockServer() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	config := &Config{}
	config.Token = "example-token"
	config.DomainNames = []string{"example"}
	client = NewClient(http.DefaultClient, config)
	client.BaseURL = server.URL
}

func teardownMockServer() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	config := &Config{}
	config.Token = "example-token"
	config.DomainNames = []string{"example"}

	c := NewClient(http.DefaultClient, config)

	if c.BaseURL != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL, defaultBaseURL)
	}
}
func TestClient_SetUserAgent(t *testing.T) {
	config := &Config{}
	config.Token = "example-token"
	config.DomainNames = []string{"example"}
	c := NewClient(http.DefaultClient, config)
	customAgent := "custom-agent/0.1"

	c.SetUserAgent(customAgent)
	if want, got := "custom-agent/0.1", c.UserAgent; want != got {
		t.Errorf("UserAgent not assigned, expected %v, got %v", want, got)
	}

	req, _ := c.newRequest("GET", "/foo")

	if want, got := "custom-agent/0.1", req.Header.Get("User-Agent"); want != got {
		t.Errorf("Incorrect User-Agent Header, expected %v, got %v", want, got)
	}
}

func TestClient_NewRequest(t *testing.T) {
	config := &Config{}
	config.Token = "example-token"
	config.DomainNames = []string{"example"}
	c := NewClient(http.DefaultClient, config)
	c.BaseURL = "https://go.example.com"

	inURL, outURL := "/foo", "https://go.example.com/foo"
	req, _ := c.newRequest("GET", inURL)

	// test that relative URL was expanded with the proper BaseURL
	if req.URL.String() != outURL {
		t.Errorf("Incorrect request URL, expected %v, got %v", outURL, req.URL.String())
	}

	// test that default user-agent is attached to the request
	ua := req.Header.Get("User-Agent")
	if ua != defaultUserAgent {
		t.Errorf("Incorrect request User-Agent, expected %v, got %v", defaultUserAgent, ua)
	}
}
