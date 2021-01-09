package hitron

import (
	"net"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return u
}

func mustMAC(s string) net.HardwareAddr {
	m, err := net.ParseMAC(s)
	if err != nil {
		panic(err)
	}

	return m
}

func TestURL(t *testing.T) {
	d := &CableModem{}
	u := d.url("")
	assert.Nil(t, u)
	u = d.url("foo")
	assert.Nil(t, u)

	d = &CableModem{
		base: mustParse("http://example.com/base/"),
	}

	expected := mustParse("http://example.com/base/")
	u = d.url("")
	assert.EqualValues(t, expected, u)

	u = d.url("/")
	assert.EqualValues(t, expected, u)

	expected = mustParse("http://example.com/base/foo")
	u = d.url("/foo")
	assert.EqualValues(t, expected, u)

	expected = mustParse("http://example.com/base/foo")
	u = d.url("foo")
	assert.EqualValues(t, expected, u)

	expected = mustParse("http://foo.example.com/blah")
	u = d.url("http://foo.example.com/blah")
	assert.EqualValues(t, expected, u)
}
