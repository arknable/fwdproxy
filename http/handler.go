package http

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/arknable/upwork-test-proxy/client"
	"github.com/arknable/upwork-test-proxy/config"
)

// HandleRequest handles HTTP request
func HandleRequest(res http.ResponseWriter, req *http.Request) {
	request, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Proxy authentication
	credential := fmt.Sprintf("%s:%s", config.ProxyUsername, config.ProxyPassword)
	authorization := base64.StdEncoding.EncodeToString([]byte(credential))
	request.Header.Add("Proxy-Authorization", fmt.Sprintf("Basic %s", authorization))

	client, err := client.NewProxied(config.ProxyAddress)
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for hkey, hval := range resp.Header {
		for _, v := range hval {
			res.Header().Add(hkey, v)
		}
	}
	res.WriteHeader(resp.StatusCode)
	_, err = io.Copy(res, resp.Body)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
