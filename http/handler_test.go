package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const url = "http://google.com"

func TestHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleRequest)
	handler(res, req)
	content := res.Body.String()

	assert.Equal(t, http.StatusOK, res.Code)
	assert.True(t, len(content) > 0)
	assert.True(t, strings.Contains(content, "<title>Google</title>"))
}
