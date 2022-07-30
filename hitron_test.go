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

func TestFormattedBytesToInt64(t *testing.T) {
	assert.Equal(t, int64(0), formattedBytesToInt64(""))
	assert.Equal(t, int64(0), formattedBytesToInt64("0"))
	assert.Equal(t, int64(0), formattedBytesToInt64("0 Bytes"))
	assert.Equal(t, int64(0), formattedBytesToInt64("0M Bytes"))
	assert.Equal(t, int64(0), formattedBytesToInt64("0.0G Bytes"))

	assert.Equal(t, int64(1), formattedBytesToInt64("1"))
	assert.Equal(t, int64(1), formattedBytesToInt64("1B"))
	assert.Equal(t, int64(1), formattedBytesToInt64("1B Bytes"))
	assert.Equal(t, int64(42), formattedBytesToInt64("42 Bytes"))
	assert.Equal(t, int64(1024), formattedBytesToInt64("1.0K Bytes"))
	assert.Equal(t, int64(2*1024*1024), formattedBytesToInt64("2.0M Bytes"))

	// 18.65 * 1TiB == 20505891858022.4, truncated to 20505891858022
	assert.Equal(t, int64(20505891858022), formattedBytesToInt64("18.65T Bytes"))
}
