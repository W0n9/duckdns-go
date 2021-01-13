package duckdns

import (
	"bufio"
	"context"
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
	server = httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.Start()

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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		v := url.Values{}
		v.Set("domains", "example")
		v.Add("ip", "")
		v.Add("token", "example-token")
		testQuery(t, r, v)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	resp, err := client.UpdateIP(context.Background())
	if err != nil {
		t.Fatalf("UpdateIP() returned error: %v", err)
	}

	if want, got := "OK", resp.Data; want != got {
		t.Errorf("UpdateIP() expected to return %v, got %v", want, got)
	}
}

func TestUpdateIPVerbose(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/success-verbose.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		v := url.Values{}
		v.Set("domains", "example")
		v.Add("ip", "")
		v.Add("token", "example-token")
		v.Add("verbose", "true")
		testQuery(t, r, v)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	client.Config.Verbose = true
	resp, err := client.UpdateIP(context.Background())
	if err != nil {
		t.Fatalf("TestUpdateIPVerbose() returned error: %v", err)
	}

	split := strings.Split(resp.Data, "\n")
	if want, got := "OK", split[0]; want != got {
		t.Errorf("TestUpdateIPVerbose() expected to return %v, got %v", want, got)
	}
}

func TestUpdateIPWithValues(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		v := url.Values{}
		v.Set("domains", "example")
		v.Add("ip", "10.10.10.253")
		v.Add("ipv6", "0:0:0:0:0:ffff:a0a:afd")
		v.Add("token", "example-token")
		testQuery(t, r, v)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	resp, err := client.UpdateIPWithValues(context.Background(), "10.10.10.253", "0:0:0:0:0:ffff:a0a:afd")
	if err != nil {
		t.Fatalf("UpdateIPWithValues() returned error: %v", err)
	}

	if want, got := "OK", resp.Data; want != got {
		t.Errorf("UpdateIPWithValues() expected to return %v, got %v", want, got)
	}
}

func TestClearIP(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		v := url.Values{}
		v.Set("domains", "example")
		v.Add("token", "example-token")
		v.Add("clear", "true")
		testQuery(t, r, v)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	resp, err := client.ClearIP(context.Background())
	if err != nil {
		t.Fatalf("ClearIP() returned error: %v", err)
	}

	if want, got := "OK", resp.Data; want != got {
		t.Errorf("ClearIP() expected to return %v, got %v", want, got)
	}
}

func TestUpdateRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()
	record := "docusign=1b0a6754-49b1-4db5-8540-d2c12664b289"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		v := url.Values{}
		v.Set("domains", "example")
		v.Add("token", "example-token")
		v.Add("txt", record)
		testQuery(t, r, v)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	resp, err := client.UpdateRecord(context.Background(), record)
	if err != nil {
		t.Fatalf("UpdateRecord() returned error: %v", err)
	}

	if want, got := "OK", resp.Data; want != got {
		t.Errorf("UpdateRecord() expected to return %v, got %v", want, got)
	}
}

func TestClearRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()
	record := "docusign=1b0a6754-49b1-4db5-8540-d2c12664b289"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)
		v := url.Values{}
		v.Set("domains", "example")
		v.Add("token", "example-token")
		v.Add("txt", record)
		v.Add("clear", "true")
		testQuery(t, r, v)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	resp, err := client.ClearRecord(context.Background(), record)
	if err != nil {
		t.Fatalf("UpdateRecord() returned error: %v", err)
	}

	if want, got := "OK", resp.Data; want != got {
		t.Errorf("UpdateRecord() expected to return %v, got %v", want, got)
	}
}
