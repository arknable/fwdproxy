package proxy

import (
	"log"
	"net/http"
)

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %10s %s", r.Method, r.RemoteAddr, r.URL.String())

	if r.Method != http.MethodConnect {
		s.serveHTTP(w, r)
		return
	}
	s.serveTLS(w, r)
}
