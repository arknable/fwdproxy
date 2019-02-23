package proxy

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/arknable/fwdproxy/env"
)

// Handles HTTPS request
func (s *Server) serveTLS(w http.ResponseWriter, r *http.Request) {
	// if err := authenticate(r); err != nil {
	// 	status := http.StatusInternalServerError
	// 	if err == ErrAuthRequired {
	// 		status = http.StatusUnauthorized
	// 	} else if err == ErrForbidden {
	// 		status = http.StatusForbidden
	// 	}
	// 	http.Error(w, err.Error(), status)
	// 	return
	// }

	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Println("Hijacking not supported")
		return
	}
	clientConn, _, err := hj.Hijack()
	if err != nil {
		log.Println(err)
		return
	}
	defer clientConn.Close()

	proxyConfig := env.Configuration().ExtProxy
	hostConn, err := net.Dial("tcp", net.JoinHostPort(proxyConfig.Address, proxyConfig.Port))
	if err != nil {
		log.Println(err)
		return
	}
	defer hostConn.Close()

	reqStrings := []string{
		fmt.Sprintf("CONNECT %s %s", r.URL.Host, r.Proto),
		fmt.Sprintf("Proxy-Authorization: Basic %s", s.proxyAuthEncoded),
		"Proxy-Connection: Keep-Alive",
		"Connection: Keep-Alive",
		"\r\n",
	}
	_, err = fmt.Fprintf(hostConn, strings.Join(reqStrings, "\r\n"))
	if err != nil {
		log.Println(err)
		return
	}
	status, err := bufio.NewReader(hostConn).ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}

	if !strings.Contains(status, "200") {
		log.Println("Failed to initiate CONNECT with external proxy")
		return
	}

	_, err = fmt.Fprintf(clientConn, "%s %v %s\r\n\r\n", r.Proto, http.StatusOK, http.StatusText(http.StatusOK))
	if err != nil {
		log.Println(err)
		return
	}

	waiter := &sync.WaitGroup{}
	waiter.Add(2)
	go transfer(waiter, clientConn, hostConn)
	go transfer(waiter, hostConn, clientConn)
	waiter.Wait()
}

func transfer(waiter *sync.WaitGroup, src io.ReadCloser, dest io.WriteCloser) {
	defer waiter.Done()
	_, err := io.Copy(dest, src)
	if err != nil {
		log.Println(err)
		return
	}
}
