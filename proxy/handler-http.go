package proxy

import (
	"io"
	"net/http"
)

// Handles HTTP request
func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		response: w,
		request:  r,
	}
	if !ctx.Authenticated() {
		return
	}

	request, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		ctx.ResponseError(err, http.StatusBadRequest)
		return
	}
	copyHeader(r.Header, request.Header)
	request.Header.Del("Proxy-Authorization")

	client := s.NewClient()
	response, err := client.Do(request)
	if err != nil {
		ctx.ResponseError(err, http.StatusServiceUnavailable)
		return
	}
	copyHeader(response.Header, w.Header())
	_, err = io.Copy(w, response.Body)
	if err != nil {
		ctx.ResponseError(err, http.StatusInternalServerError)
		return
	}
}
