package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoCred(t *testing.T) {
	test := new()
	test.UsingBasicAuth = false
	defer test.Close()

	if err := test.Do(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusUnauthorized, test.Response.StatusCode)
}

func TestInvalidUser(t *testing.T) {
	test := new()
	test.Username = "foo"
	defer test.Close()

	if err := test.Do(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusForbidden, test.Response.StatusCode)
}

func TestInvalidPassword(t *testing.T) {
	test := new()
	test.Password = "n0tv4lidpwd"
	defer test.Close()

	if err := test.Do(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusForbidden, test.Response.StatusCode)
}
