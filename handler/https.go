package handler

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/arknable/fwdproxy/server"
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
	hostConn, err := net.Dial("tcp", server.ProxyAddress)
	if err != nil {
		log.Println(err)
		return
	}
	reqString := []string{
		fmt.Sprintf("CONNECT %s %s", r.URL.Host, r.Proto),
		"Proxy-Authorization: Basic dGVzdDp0ZXN0cGFzc3dvcmQ=",
		"Proxy-Connection: Keep-Alive",
		"Connection: Keep-Alive",
		"\r\n",
	}

	_, err = fmt.Fprintf(hostConn, strings.Join(reqString, "\r\n"))
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
