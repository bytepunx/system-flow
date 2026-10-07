package host

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client speaks to a running host's control API.
type Client struct {
	URL   string
	Token string
	HTTP  *http.Client
}

// FromEnv is the client a child of the host is given, and false for a
// process the host did not start.
func FromEnv() (*Client, bool) {
	url, token := os.Getenv(URLEnv), os.Getenv(TokenEnv)
	if url == "" || token == "" {
		return nil, false
	}
	return &Client{URL: url, Token: token}, true
}

func (c *Client) do(ctx context.Context, method, path string, in, out any) error {
	var body *bytes.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	} else {
		body = bytes.NewReader(nil)
	}
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(c.URL, "/")+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(c.Token))
	req.Header.Set("Content-Type", "application/json")
	hc := c.HTTP
	if hc == nil {
		hc = &http.Client{Timeout: 6 * time.Minute}
	}
	resp, err := hc.Do(req)
	if err != nil {
		return fmt.Errorf("flai host at %s: %w", c.URL, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		var e apiError
		if json.NewDecoder(resp.Body).Decode(&e) == nil && e.Error != "" {
			return errors.New(e.Error)
		}
		return fmt.Errorf("flai host at %s answered %s", c.URL, resp.Status)
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// Status is the host as it is now, children and all.
func (c *Client) Status(ctx context.Context) (Status, error) {
	var st Status
	return st, c.do(ctx, http.MethodGet, "/status", nil, &st)
}

// Act starts, stops, or restarts serve, mcp, or all, and answers the status after.
func (c *Client) Act(ctx context.Context, process, action string) (Status, error) {
	var st Status
	return st, c.do(ctx, http.MethodPost, "/process", ProcessRequest{Process: process, Action: action}, &st)
}

// KeepMCP makes the host keep exactly these projects' MCP servers.
func (c *Client) KeepMCP(ctx context.Context, roots []string) (Status, error) {
	if roots == nil {
		roots = []string{}
	}
	var st Status
	return st, c.do(ctx, http.MethodPut, "/mcp", MCPRequest{Roots: roots}, &st)
}

// Check asks whether a newer flai is published.
func (c *Client) Check(ctx context.Context) (json.RawMessage, error) {
	var out json.RawMessage
	return out, c.do(ctx, http.MethodGet, "/check", nil, &out)
}

// Versions lists the published flai releases, newest first, as flai
// self-upgrade --list --json prints them.
func (c *Client) Versions(ctx context.Context) (json.RawMessage, error) {
	var out json.RawMessage
	return out, c.do(ctx, http.MethodGet, "/versions", nil, &out)
}

// Upgrade installs the newest flai, or the published release version names
// when it is not ""; the host restarts on it when it did.
func (c *Client) Upgrade(ctx context.Context, version string) (json.RawMessage, error) {
	var in any
	if version != "" {
		in = UpgradeRequest{Version: version}
	}
	var out json.RawMessage
	return out, c.do(ctx, http.MethodPost, "/upgrade", in, &out)
}
