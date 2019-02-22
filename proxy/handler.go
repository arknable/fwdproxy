package proxy

import (
	"net/http"
)

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodConnect {
		s.serveHTTP(w, r)
		return
	}
	s.serveTLS(w, r)
}
