package proxy

import (
	"fmt"
	"io"
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

// ResponseRaw writes given message to client connection as is
func (c *Context) ResponseRaw(message string) {
	if c.clientConn == nil {
		return
	}
	_, err := fmt.Fprintf(c.clientConn, "%s\r\n\r\n", message)
	if err != nil {
		log.Println(err)
	}
}

// CopyResponse copies response information from given response
func (c *Context) CopyResponse(resp *http.Response) {
	dest := c.response.Header()
	for key, val := range resp.Header {
		for _, v := range val {
			dest.Add(key, v)
		}
	}
	_, err := io.Copy(c.response, resp.Body)
	if err != nil {
		c.ResponseError(err, http.StatusInternalServerError)
		return
	}
}
