package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const httpURL = "http://google.com"

func TestHTTPNoAuth(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, httpURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleRequest)
	handler(res, req)
	assert.Equal(t, http.StatusUnauthorized, res.Code)
}

func TestHTTPInvalidAuth(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, httpURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("tos", "tospassword")
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleRequest)
	handler(res, req)
	assert.Equal(t, http.StatusForbidden, res.Code)
}
