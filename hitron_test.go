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

func TestFormattedBytesToUint64(t *testing.T) {
	assert.Equal(t, uint64(0), formattedBytesToUint64(""))
	assert.Equal(t, uint64(0), formattedBytesToUint64("0"))
	assert.Equal(t, uint64(0), formattedBytesToUint64("0 Bytes"))
	assert.Equal(t, uint64(0), formattedBytesToUint64("0M Bytes"))
	assert.Equal(t, uint64(0), formattedBytesToUint64("0.0G Bytes"))

	assert.Equal(t, uint64(1), formattedBytesToUint64("1"))
	assert.Equal(t, uint64(1), formattedBytesToUint64("1B"))
	assert.Equal(t, uint64(1), formattedBytesToUint64("1B Bytes"))
	assert.Equal(t, uint64(42), formattedBytesToUint64("42 Bytes"))
	assert.Equal(t, uint64(kib), formattedBytesToUint64("1.0K Bytes"))
	assert.Equal(t, uint64(2*mib), formattedBytesToUint64("2.0M Bytes"))

	// 18.65 * 1TiB == 20505891858022.4, truncated to 20505891858022
	assert.Equal(t, uint64(20505891858022), formattedBytesToUint64("18.65T Bytes"))
}

func TestByteSize(t *testing.T) {
	maxUint64 := ^uint64(0)

	//nolint:govet
	testdata := []struct {
		in       uint64
		expected string
	}{
		{maxUint64, "16E"},
		{10 * eib, "10E"},
		{uint64(10.5 * eib), "10.5E"},
		{10 * pib, "10P"},
		{uint64(10.5 * pib), "10.5P"},
		{10 * tib, "10T"},
		{uint64(10.5 * tib), "10.5T"},
		{10 * gib, "10G"},
		{uint64(10.5 * gib), "10.5G"},
		{100 * mib, "100M"},
		{uint64(100.5 * mib), "100.5M"},
		{100 * kib, "100K"},
		{uint64(100.5 * kib), "100.5K"},
		{1, "1B"},
		{0, "0B"},
	}

	for _, d := range testdata {
		assert.Equal(t, d.expected, byteSize(d.in))
	}
}
