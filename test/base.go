package test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/arknable/upwork-test-proxy/client"
	"github.com/arknable/upwork-test-proxy/config"
	"github.com/arknable/upwork-test-proxy/server"
)

const (
	// URL to be requested to external proxy
	targetURL = "http://google.com"

	// Server address to be used as proxy
	proxyAddress = "127.0.0.1"
)

// Represents an HTTP/HTTPS test
type serverTest struct {
	url            string
	method         string
	proxyAddress   string
	server         *http.Server
	RequestBody    io.Reader
	IsTLS          bool
	UsingBasicAuth bool
	Username       string
	Password       string
	Response       *http.Response
	ResponseBody   string
}

// Creates new HTTP test
func new() *serverTest {
	users := config.Users()
	return &serverTest{
		url:            targetURL,
		method:         http.MethodGet,
		proxyAddress:   proxyAddress,
		UsingBasicAuth: true,
		Username:       users[0],
		Password:       config.Password(users[0]),
	}
}

// Creates new HTTPS test
func newTLS() *serverTest {
	t := new()
	t.IsTLS = true
	return t
}

// Starts server
func (h *serverTest) startServer(errChan chan error) {
	var err error
	if h.IsTLS {
		err = h.server.ListenAndServeTLS(config.CertPath, config.KeyPath)
	} else {
		err = h.server.ListenAndServe()
	}

	if err != nil {
		errChan <- err
	}
}

// Do performs test
func (h *serverTest) Do() error {
	var srv *http.Server
	if h.IsTLS {
		_, srv = server.NewTLS()
	} else {
		srv = server.New()
	}
	h.server = srv
	defer h.server.Close()

	errChan := make(chan error)
	go h.startServer(errChan)
	time.Sleep(500 * time.Millisecond)

	testURL, err := url.Parse(h.url)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(h.method, testURL.String(), h.RequestBody)
	if err != nil {
		return err
	}
	if h.UsingBasicAuth {
		request.SetBasicAuth(h.Username, h.Password)
	}

	proxyAddr := fmt.Sprintf("http://%s:%s", proxyAddress, config.HTTPPort)
	if h.IsTLS {
		proxyAddr = fmt.Sprintf("https://%s:%s", proxyAddress, config.TLSPort)
	}
	cl, err := client.New(proxyAddr)
	if err != nil {
		return err
	}

	resp, err := cl.Do(request)
	if err != nil {
		return err
	}
	h.Response = resp

	body, err := ioutil.ReadAll(h.Response.Body)
	if err != nil {
		return err
	}
	h.ResponseBody = string(body)

	return nil
}

// ResponseContains checks if response body contains given cutset
func (h *serverTest) ResponseContains(cutset string) bool {
	if h.Response != nil {
		return strings.Contains(h.ResponseBody, cutset)
	}
	return false
}

// Close ends test and close attached response
func (h *serverTest) Close() {
	if h.Response != nil {
		h.Response.Body.Close()
	}
}
