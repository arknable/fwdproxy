package http

import (
	"io/ioutil"
	"net/http"
)

// HandleRequest handles HTTP request
func HandleRequest(res http.ResponseWriter, req *http.Request) {
	request, err := http.NewRequest(req.Method, req.URL.String(), nil)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	length, err := res.Write(content)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if length != len(content) {
		http.Error(res, "Failed to forward response", http.StatusInternalServerError)
	}
}
