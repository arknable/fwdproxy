package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

// Context represents a request processing
type Context struct {
	// Client request
	request *http.Request

	// Request response
	response http.ResponseWriter

	// Client connection, used on CONNECT.
	clientConn net.Conn

	// Proxy connection, used on CONNECT.
	proxyConn net.Conn
}

// ResponseError writes error response
func (c *Context) ResponseError(err error, status int) {
	if c.request.Method != http.MethodConnect {
		http.Error(c.response, err.Error(), status)
		return
	}

	if c.clientConn != nil {
		_, err = fmt.Fprintf(
			c.clientConn,
			"%s %v %s\n%s\r\n\r\n",
			c.request.Proto,
			status,
			http.StatusText(status),
			err,
		)
		if err != nil {
			log.Println(err)
		}

		return
	}

	http.Error(c.response, err.Error(), status)
	log.Println(err)
}
