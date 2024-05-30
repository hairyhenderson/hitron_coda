package hitron

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCMReboot(t *testing.T) {
	body := `{"errCode":"000","errMsg":""}`

	srv := staticResponseServer(t, body)

	d := testCableModem(srv)
	ctx := ContextWithDebugLogger(context.Background(), t)

	o, err := d.CMReboot(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, &Error{Code: "000"}, o)
}

func staticResponseServer(t *testing.T, body string) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	return srv
}
