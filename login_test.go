package hitron

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	creds := credentials{
		Username: "cusadmin",
		Password: "supersecretpassword",
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/Users/Login", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		err := r.ParseForm()
		assert.NoError(t, err)

		model := r.Form.Get("model")
		assert.NotEmpty(t, model)

		c := credentials{}
		err = json.Unmarshal([]byte(model), &c)
		assert.NoError(t, err)
		assert.Equal(t, creds, c)

		cookie := &http.Cookie{
			Name:     "PHPSESSID",
			Value:    "1234567890",
			Path:     "/",
			HttpOnly: true,
		}
		w.Header().Set("Set-Cookie", cookie.String())
		w.WriteHeader(http.StatusOK)
	}))

	defer srv.Close()

	hc := srv.Client()
	jar, err := cookiejar.New(nil)
	assert.NoError(t, err)

	hc.Jar = jar

	u := mustParse(srv.URL)
	d := &CableModem{credentials: creds, base: u, hc: hc}

	err = d.Login(context.Background())
	assert.NoError(t, err)

	cookies := jar.Cookies(u)
	assert.Len(t, cookies, 1)
	assert.EqualValues(t, &http.Cookie{
		Name:  "PHPSESSID",
		Value: "1234567890",
	}, cookies[0])
}

//nolint:funlen
func TestLogout(t *testing.T) {
	creds := credentials{
		Username: "cusadmin",
		Password: "supersecretpassword",
	}
	n := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch n {
		case 0:
			assert.Equal(t, "/Users/Login", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			err := r.ParseForm()
			assert.NoError(t, err)

			model := r.Form.Get("model")
			assert.NotEmpty(t, model)

			c := credentials{}
			err = json.Unmarshal([]byte(model), &c)
			assert.NoError(t, err)
			assert.Equal(t, creds, c)

			cookie := &http.Cookie{
				Name:     "PHPSESSID",
				Value:    "1234567890",
				Path:     "/",
				HttpOnly: true,
			}
			w.Header().Set("Set-Cookie", cookie.String())
			w.WriteHeader(http.StatusOK)
		case 1:
			assert.Equal(t, "/Users/Logout", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			cookie := &http.Cookie{
				Name:     "PHPSESSID",
				Value:    "9999999999",
				Path:     "/",
				HttpOnly: true,
			}
			w.Header().Set("Set-Cookie", cookie.String())
			w.WriteHeader(http.StatusOK)
		}

		n++
	}))

	defer srv.Close()

	hc := srv.Client()
	jar, err := cookiejar.New(nil)
	assert.NoError(t, err)

	hc.Jar = jar

	u := mustParse(srv.URL)
	d := &CableModem{credentials: creds, base: u, hc: hc}

	err = d.Login(context.Background())
	assert.NoError(t, err)

	cookies := jar.Cookies(u)
	assert.Len(t, cookies, 1)
	assert.EqualValues(t, &http.Cookie{
		Name:  "PHPSESSID",
		Value: "1234567890",
	}, cookies[0])

	err = d.Logout(context.Background())
	assert.NoError(t, err)

	cookies = jar.Cookies(u)
	assert.Len(t, cookies, 1)
	assert.EqualValues(t, &http.Cookie{
		Name:  "PHPSESSID",
		Value: "9999999999",
	}, cookies[0])
}
