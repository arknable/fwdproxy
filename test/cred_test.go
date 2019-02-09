package test

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/arknable/upwork-test-proxy/server"
	"github.com/stretchr/testify/assert"
)

func TestNoCred(t *testing.T) {
	srv := server.New()
	go srv.ListenAndServe()
	defer srv.Close()
	time.Sleep(500 * time.Millisecond)

	fn := func(r *http.Request) {
		r.SetBasicAuth("", "")
	}
	resp, err := doRequest(http.MethodGet, false, nil, fn)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestInvalidUser(t *testing.T) {
	srv := server.New()
	go srv.ListenAndServe()
	defer srv.Close()
	time.Sleep(500 * time.Millisecond)

	fn := func(r *http.Request) {
		_, password, _ := r.BasicAuth()
		r.SetBasicAuth("foo", password)
	}
	resp, err := doRequest(http.MethodGet, false, nil, fn)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestInvalidPassword(t *testing.T) {
	srv := server.New()
	go srv.ListenAndServe()
	defer srv.Close()
	time.Sleep(500 * time.Millisecond)

	fn := func(r *http.Request) {
		username, _, _ := r.BasicAuth()
		r.SetBasicAuth(username, "n0tv4lidpwd")
	}
	resp, err := doRequest(http.MethodGet, false, nil, fn)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}
