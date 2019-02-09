package test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/arknable/upwork-test-proxy/client"
	"github.com/arknable/upwork-test-proxy/config"
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
	url, err := url.Parse(fmt.Sprintf("http://%s", targetURL))
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

	port := config.HTTPPort
	if useTLS {
		port = config.TLSPort
	}
	cl, err := client.NewProxied(fmt.Sprintf("%s://%s:%s", proc, proxyAddress, port), useTLS)
	if err != nil {
		return nil, err
	}
	fmt.Println("----> ", request.URL.String())
	return cl.Do(request)
}
