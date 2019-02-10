package http

import (
	"io"
	"log"
	"net/http"

	"github.com/arknable/fwdproxy/client"
	"github.com/arknable/fwdproxy/config"
)

// HandleRequest handles HTTP request
func HandleRequest(res http.ResponseWriter, req *http.Request) {
	request, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	request.Header.Set("Host", req.Host)
	request.Header.Set("X-Forwarded-For", req.RemoteAddr)
	for hkey, hval := range req.Header {
		for _, v := range hval {
			request.Header.Add(hkey, v)
		}
	}

	client, err := client.New(config.ProxyAddress)
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
