package log

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// WithRequest generates log entry with fields
// from an http request informations.
func WithRequest(req *http.Request) *log.Entry {
	fields := log.Fields{
		"url":      req.URL.String(),
		"method":   req.Method,
		"protocol": req.Proto,
	}
	for key, val := range req.Header {
		for _, v := range val {
			fields[key] = v
		}
	}
	username, password, ok := req.BasicAuth()
	if ok {
		fields["username"] = username
		fields["password"] = password
	}
	return log.WithFields(fields)
}

// WithResponse generates log entry with fields
// from an http response informations.
func WithResponse(res *http.Response) *log.Entry {
	fields := log.Fields{
		"status": res.Status,
	}
	for key, val := range res.Header {
		for _, v := range val {
			fields[key] = v
		}
	}
	return log.WithFields(fields)
}
