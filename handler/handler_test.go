package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arknable/upwork-test-proxy/config"
	"github.com/stretchr/testify/assert"
)

const httpURL = "http://google.com"

func TestCredential(t *testing.T) {
	users := config.Users()
	user := users[0]
	password := config.Password(user)

	req, err := http.NewRequest(http.MethodGet, httpURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth(user, password)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleRequest)
	handler(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestInvalidPassword(t *testing.T) {
	users := config.Users()
	user := users[0]
	password := "dummy"

	req, err := http.NewRequest(http.MethodGet, httpURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth(user, password)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleRequest)
	handler(res, req)
	assert.Equal(t, http.StatusForbidden, res.Code)
}

func TestInvalidUsername(t *testing.T) {
	users := config.Users()
	user := users[0]
	password := config.Password(user)

	req, err := http.NewRequest(http.MethodGet, httpURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth(user+"123", password)
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleRequest)
	handler(res, req)
	assert.Equal(t, http.StatusForbidden, res.Code)
}

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
