package handler

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/arknable/fwdproxy/config"
	"github.com/arknable/fwdproxy/server"
	"github.com/arknable/fwdproxy/user"
)

// HandleRequest handles requests
func HandleRequest(res http.ResponseWriter, req *http.Request) {
	var (
		fwdMethod string
		username  string
		password  string
		data      []byte
		err       error
	)

	// Dump request for log information
	data, err = httputil.DumpRequest(req, true)
	if err != nil {
		log.Println("httputil.DumpRequest: ", err)
	} else {
		log.Printf("Request from client: \n%s\n", string(data))
	}

	authOK := false
	if req.Method == http.MethodConnect { // HTTPS
		fwdMethod = http.MethodGet
		header := req.Header.Get("Proxy-Authorization")
		if len(header) > 0 {
			parts := strings.Split(header, " ")
			data, err = base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				log.Println("base64.StdEncoding.Decode: ", err)
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			parts = strings.Split(string(data), ":")
			authOK = true
			username = parts[0]
			password = parts[1]
		}
		req.URL.Scheme = "https" // Somehow scheme always cleared
	} else {
		fwdMethod = req.Method
		username, password, authOK = req.BasicAuth()
	}

	if !authOK {
		log.Println("Unauthorized access, missing authentication")
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
		log.Printf("Invalid authentication: %s/%s\n", username, password)
		http.Error(res, "You have no access to do a request", http.StatusForbidden)
		return
	}

	request, err := http.NewRequest(fwdMethod, req.URL.String(), req.Body)
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

	// Dump request for log information
	data, err = httputil.DumpRequest(request, true)
	if err != nil {
		log.Println("httputil.DumpRequest(ext): ", err)
	} else {
		log.Printf("Request to ext proxy: \n%s\n", string(data))
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

	// Dump response for log information
	data, err = httputil.DumpResponse(resp, false)
	if err != nil {
		log.Println("httputil.DumpResponse: ", err)
	} else {
		log.Printf("Response from ext proxy: \n%s\n", string(data))
	}

	_, err = io.Copy(res, resp.Body)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
