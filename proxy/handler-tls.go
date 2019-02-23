package proxy

import (
	"bufio"
	"errors"
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
	ctx := &Context{
		response: w,
		request:  r,
	}
	if !ctx.Authenticated() {
		return
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		ctx.ResponseError(errors.New("Hijacking not supported"), http.StatusServiceUnavailable)
		return
	}
	clientConn, _, err := hj.Hijack()
	if err != nil {
		ctx.ResponseError(err, http.StatusServiceUnavailable)
		return
	}
	ctx.clientConn = clientConn
	defer clientConn.Close()

	proxyConfig := env.Configuration().ExtProxy
	proxyConn, err := net.Dial("tcp", net.JoinHostPort(proxyConfig.Address, proxyConfig.Port))
	if err != nil {
		ctx.ResponseError(err, http.StatusServiceUnavailable)
		return
	}
	ctx.proxyConn = proxyConn
	defer proxyConn.Close()

	reqStrings := []string{
		fmt.Sprintf("CONNECT %s %s", r.URL.Host, r.Proto),
		fmt.Sprintf("Proxy-Authorization: Basic %s", s.proxyAuthEncoded),
		"Proxy-Connection: Keep-Alive",
		"Connection: Keep-Alive",
		"\r\n",
	}
	_, err = fmt.Fprintf(proxyConn, strings.Join(reqStrings, "\r\n"))
	if err != nil {
		ctx.ResponseError(err, http.StatusInternalServerError)
		return
	}
	status, err := bufio.NewReader(proxyConn).ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}

	if !strings.Contains(status, "200") {
		_, err = fmt.Fprintf(clientConn, status)
		if err != nil {
			ctx.ResponseError(err, http.StatusInternalServerError)
		}
		return
	}

	_, err = fmt.Fprintf(clientConn, "%s %v %s\r\n\r\n", r.Proto, http.StatusOK, http.StatusText(http.StatusOK))
	if err != nil {
		ctx.ResponseError(err, http.StatusInternalServerError)
		return
	}

	waiter := &sync.WaitGroup{}
	waiter.Add(2)
	go transfer(ctx, waiter, clientConn, proxyConn)
	go transfer(ctx, waiter, proxyConn, clientConn)
	waiter.Wait()
}

func transfer(ctx *Context, waiter *sync.WaitGroup, src io.ReadCloser, dest io.WriteCloser) {
	defer waiter.Done()
	_, err := io.Copy(dest, src)
	if err != nil {
		ctx.ResponseError(err, http.StatusInternalServerError)
		return
	}
}
