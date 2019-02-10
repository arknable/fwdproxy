package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTLS(t *testing.T) {
	test := newTLS()
	defer test.Close()

	if err := test.Do(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, test.Response.StatusCode)
	assert.True(t, test.ResponseContains("<title>Google</title>"))
}
