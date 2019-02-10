package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTP(t *testing.T) {
	test := new()
	defer test.Close()

	if err := test.Do(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, test.Response.StatusCode)
	assert.True(t, test.ResponseContains("<title>Google</title>"))
}
