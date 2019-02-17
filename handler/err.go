package handler

import (
	"net/http"

	mylog "github.com/arknable/fwdproxy/log"
)

// Writes error to response
func internalError(res http.ResponseWriter, req *http.Request, err error) {
	http.Error(res, err.Error(), http.StatusInternalServerError)
	mylog.WithRequest(req).Warning(err)
}
