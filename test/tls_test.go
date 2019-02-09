package test

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/arknable/upwork-test-proxy/config"
	"github.com/arknable/upwork-test-proxy/server"
	"github.com/stretchr/testify/assert"
)

func TestTLS(t *testing.T) {
	_, srv := server.NewTLS()
	go func() {
		err := srv.ListenAndServeTLS(config.CertPath, config.KeyPath)
		if err != nil {
			log.Fatal(err)
		}
	}()
	defer srv.Close()
	time.Sleep(500 * time.Millisecond)

	resp, err := doRequest(http.MethodGet, true, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	sbody := string(body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, strings.Contains(sbody, "<title>Google</title>"))
}
