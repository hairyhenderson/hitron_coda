package hitron

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

//go:generate gomplate -c .=apilist.yaml -f methods.go.tmpl -o methods.go

// CableModem represents the Hitron CODA Cable Modem/Router
type CableModem struct {
	base        *url.URL
	hc          *http.Client
	credentials credentials
}

// debugTransport - logs the request and response if debug is enabled
type debugTransport struct {
	rt http.RoundTripper
}

func (t *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !slog.Default().Enabled(req.Context(), slog.LevelDebug) {
		return t.rt.RoundTrip(req)
	}

	drq, _ := httputil.DumpRequest(req, true)
	slog.DebugContext(req.Context(), "dumping request", slog.String("request", string(drq)))

	resp, err := t.rt.RoundTrip(req)
	if err == nil {
		drs, _ := httputil.DumpResponse(resp, true)
		slog.DebugContext(req.Context(), "dumping response", slog.String("response", string(drs)))
	}

	return resp, err
}

// New instantiates a default CableModem struct
func New(host, username, password string) (*CableModem, error) {
	u, err := url.Parse(fmt.Sprintf("http://%s/1/Device/", host))
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	tr := http.DefaultTransport
	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		tr = &debugTransport{tr}
	}

	client := &http.Client{
		Jar: jar,
		// Ignore redirects
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: tr,
	}

	creds := credentials{username, password}

	return &CableModem{
		credentials: creds,
		base:        u,
		hc:          client,
	}, nil
}

func (c *CableModem) url(s string) *url.URL {
	if len(s) == 0 || c.base == nil {
		return c.base
	}

	if s[0] == '/' {
		s = s[1:]
	}

	p, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return c.base.ResolveReference(p)
}

func (c *CableModem) getJSON(ctx context.Context, path string, o interface{}) error {
	return c.sendRequest(ctx, http.MethodGet, path, http.NoBody, o)
}

func (c *CableModem) sendRequest(ctx context.Context, method, path string, body, o interface{}) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u := c.url(path).String()

	contentType := ""

	var reqBody io.Reader
	switch b := body.(type) {
	case io.Reader:
		reqBody = b
	case url.Values:
		contentType = "application/x-www-form-urlencoded"
		reqBody = strings.NewReader(b.Encode())
	default:
		return fmt.Errorf("unsupported body type %T", body)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reqBody)
	if err != nil {
		return err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status %d: %s (Header: %v)", resp.StatusCode, string(b), resp.Header)
	}

	err = json.Unmarshal(b, o)
	if err != nil {
		return fmt.Errorf("JSON decoding failed: %w", err)
	}

	return nil
}

func atoi64(s string) int64 {
	i, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)

	return i
}

func atoui64(s string) uint64 {
	i, _ := strconv.ParseUint(strings.TrimSpace(s), 10, 64)

	return i
}

func atof64(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)

	return f
}

//nolint:gomnd
const (
	_byte = 1 << (10 * iota)
	kib
	mib
	gib
	tib
	pib
	eib
)

func formattedBytesToUint64(s string) uint64 {
	i, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	if err == nil {
		return i
	}

	s = strings.TrimSuffix(s, " Bytes")
	if len(s) <= 1 {
		return atoui64(s)
	}

	switch s[len(s)-1] {
	case 'B':
		i = uint64(atof64(s[:len(s)-1]))
	case 'K':
		i = uint64(atof64(s[:len(s)-1]) * kib)
	case 'M':
		i = uint64(atof64(s[:len(s)-1]) * mib)
	case 'G':
		i = uint64(atof64(s[:len(s)-1]) * gib)
	case 'T':
		i = uint64(atof64(s[:len(s)-1]) * tib)
	case 'P':
		i = uint64(atof64(s[:len(s)-1]) * pib)
	case 'E':
		i = uint64(atof64(s[:len(s)-1]) * eib)
	default:
		i = uint64(atof64(s))
	}

	return i
}

func byteSize(bytes uint64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case bytes >= eib:
		unit = "E"
		value /= eib
	case bytes >= pib:
		unit = "P"
		value /= pib
	case bytes >= tib:
		unit = "T"
		value /= tib
	case bytes >= gib:
		unit = "G"
		value /= gib
	case bytes >= mib:
		unit = "M"
		value /= mib
	case bytes >= kib:
		unit = "K"
		value /= kib
	case bytes >= _byte:
		unit = "B"
	case bytes == 0:
		return "0B"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")

	return result + unit
}
