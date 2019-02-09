package test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/arknable/upwork-test-proxy/client"
	"github.com/arknable/upwork-test-proxy/config"
	"github.com/arknable/upwork-test-proxy/server"
)

const (
	// URL to be requested to external proxy
	targetURL = "google.com"

	// Server address to be used as proxy
	proxyAddress = "127.0.0.1"
)

func doRequest(method string, useTLS bool, body io.Reader, fn func(*http.Request)) (*http.Response, error) {
	proc := "http"
	if useTLS {
		proc = "https"
	}
	url, err := url.Parse(fmt.Sprintf("%s://%s", proc, targetURL))
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}
	users := config.Users()
	request.SetBasicAuth(users[0], config.Password(users[0]))

	if fn != nil {
		fn(request)
	}

	cl, err := client.NewProxied(fmt.Sprintf("%s://%s:%s", proc, proxyAddress, server.HTTPPort))
	if err != nil {
		return nil, err
	}

	return cl.Do(request)
}
