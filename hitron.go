package hitron

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

//go:generate gomplate -c .=apilist.yaml -f methods.go.tmpl -o methods.go

// CableModem represents the Hitron CODA Cable Modem/Router
type CableModem struct {
	credentials credentials
	base        *url.URL

	hc *http.Client
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

	client := &http.Client{
		Jar: jar,
		// Ignore redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
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

type debugLogger interface {
	Logf(format string, args ...interface{})
}

type debugLoggerKey struct{}

// ContextWithDebugLogger - add a logger for debugging the client
func ContextWithDebugLogger(ctx context.Context, l interface {
	Logf(format string, args ...interface{})
}) context.Context {
	return context.WithValue(ctx, debugLoggerKey{}, l)
}

type debugLoggerFunc func(format string, args ...interface{})

func (f debugLoggerFunc) Logf(format string, args ...interface{}) {
	f(format, args...)
}

func debugLoggerFromContext(ctx context.Context) debugLogger {
	if l := ctx.Value(debugLoggerKey{}); l != nil {
		dl, ok := l.(debugLogger)
		if ok {
			return dl
		}
	}

	return debugLoggerFunc(func(f string, args ...interface{}) {})
}

func (c *CableModem) getJSON(ctx context.Context, path string, o interface{}) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u := c.url(path).String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	// DEBUG
	debugLoggerFromContext(ctx).Logf("JSON response: %s", string(b))

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
