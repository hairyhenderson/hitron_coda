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

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))

	defer srv.Close()

	d := testCableModem(srv)
	ctx := ContextWithDebugLogger(context.Background(), t)

	o, err := d.CMReboot(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, &Error{Code: "000"}, o)
}
