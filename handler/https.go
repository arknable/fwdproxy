package handler

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

// Handles HTTPS request
func serveTLS(w http.ResponseWriter, r *http.Request) {
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
	hostConn, err := net.Dial("tcp", r.URL.Host)
	if err != nil {
		log.Println(err)
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
	clientConn.Close()
	hostConn.Close()
}

func transfer(waiter *sync.WaitGroup, src io.ReadCloser, dest io.WriteCloser) {
	defer waiter.Done()
	_, err := io.Copy(dest, src)
	if err != nil {
		log.Println(err)
		return
	}
}
