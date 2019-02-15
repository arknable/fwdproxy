package handler

import "net/http"

// HandleRequest handles both HTTP and HTTPS request
func HandleRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodConnect {
		handleHTTP(res, req)
		return
	}
	handleTLS(res, req)
}
