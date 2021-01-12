package duckdns

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
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

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; want != got {
		t.Errorf("Request METHOD expected to be `%v`, got `%v`", want, got)
	}
}

func testQuery(t *testing.T, r *http.Request, want url.Values) {
	if got := r.URL.Query(); !reflect.DeepEqual(want, got) {
		t.Errorf("Request METHOD expected to be `%v`, got `%v`", want, got)
	}
}

func testHeader(t *testing.T, r *http.Request, name, want string) {
	if got := r.Header.Get(name); want != got {
		t.Errorf("Request() %v expected to be `%#v`, got `%#v`", name, want, got)
	}
}

func testHeaders(t *testing.T, r *http.Request) {
	testHeader(t, r, "User-Agent", defaultUserAgent)
}

func readHTTPFixture(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile("../fixtures.http" + filename)
	if err != nil {
		t.Fatalf("Unable to read HTTP fixture: %v", err)
	}
	s := string(data[:])
	return s
}

func httpResponseFixture(t *testing.T, filename string) *http.Response {
	resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(readHTTPFixture(t, filename))), nil)
	if err != nil {
		t.Fatalf("Unable to create http.Response from fixture: %v", err)
	}
	// resp.Body.Close()
	return resp
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

func TestUpdateIP(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/update?", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/ip/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		testQuery(t, r, url.Values{})

		w.WriteHeader(httpResponse.StatusCode)
		_, _ = io.Copy(w, httpResponse.Body)
	})

	resp, err := client.UpdateIP()
	if err != nil {
		t.Fatalf("UpdateIP() returned error: %v", err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if want, got := "OK", string(bodyBytes); want != got {
		t.Errorf("UpdateIP() expected to return %v, got %v", want, got)
	}
}
