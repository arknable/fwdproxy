package handler

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/arknable/fwdproxy/config"
	"github.com/arknable/fwdproxy/server"
	"github.com/arknable/fwdproxy/user"
)

// HandleRequest handles requests
func HandleRequest(res http.ResponseWriter, req *http.Request) {
	if req.URL.Scheme == "http" {
		username, password, ok := req.BasicAuth()
		if !ok || (len(strings.Trim(username, " ")) == 0) {
			http.Error(res, "Restricted access only", http.StatusUnauthorized)
			return
		}
		valid, err := user.Repo().Validate(username, password)
		if err != nil {
			http.Error(res, "Failed to validate user", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if !valid {
			http.Error(res, "You have no access to do a request", http.StatusForbidden)
			return
		}
	} else {
		req.URL.Scheme = "https" // Somehow Scheme becomes empty on TLS
	}

	log.Printf("New request to %s [%s]\n", req.URL.String(), req.Method)

	request, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		log.Println("http.NewRequest: ", err)
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

	client, err := server.NewClient(config.ProxyAddress)
	resp, err := client.Do(request)
	if err != nil {
		log.Println("client.Do: ", err)
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
