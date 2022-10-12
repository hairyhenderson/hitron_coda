package hitron

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c credentials) String() string {
	b, _ := json.Marshal(&c)

	return string(b)
}

// Login to a new session
func (c *CableModem) Login(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u := c.url("/Users/Login").String()

	data := url.Values{
		"model": []string{c.credentials.String()},
	}
	rbody := strings.NewReader(data.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, rbody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Logout from the current session
func (c *CableModem) Logout(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u := c.url("/Users/Logout").String()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return err
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
