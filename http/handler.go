package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

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
	proxyURL, err := url.Parse(config.ProxyAddress)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	length, err := res.Write(content)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if length != len(content) {
		log.Println("Proxy response incompletely written")
		http.Error(res, "Failed to forward response", http.StatusInternalServerError)
	}
}
