// Package channeltest is a stand-in for the dashboard's agent endpoint, for
// tests of whatever dials it.
package channeltest

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/bytepunx/system-flow/flai/internal/channel"
)

// Dashboard accepts flai's connection, does the handshake with Key, and
// hands each proven connection to Proven.
type Dashboard struct {
	URL    string
	Proven chan *Conn
}

// Conn is one proven connection from a flai.
type Conn struct {
	Project string
	Kind    string // what hello said: empty for a project, channel.KindCandidate for one offered for import
	ws      *websocket.Conn
}

type msg struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *channel.Error  `json:"error,omitempty"`
}

// New starts a dashboard that holds key.
func New(t *testing.T, key string) *Dashboard {
	t.Helper()
	d := &Dashboard{Proven: make(chan *Conn, 16)}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != channel.Path {
			http.NotFound(w, r)
			return
		}
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		ctx := r.Context()
		_, data, err := c.Read(ctx)
		if err != nil {
			return
		}
		var hello msg
		var hp struct {
			Nonce   string          `json:"nonce"`
			Project channel.Project `json:"project"`
			Kind    string          `json:"kind"`
		}
		_ = json.Unmarshal(data, &hello)
		_ = json.Unmarshal(hello.Params, &hp)
		b := make([]byte, 16)
		_, _ = rand.Read(b)
		mine := hex.EncodeToString(b)
		res, _ := json.Marshal(map[string]string{"nonce": mine, "proof": channel.Proof([]byte(key), "dashboard", hp.Nonce, mine), "dashboard": "test"})
		out, _ := json.Marshal(msg{JSONRPC: "2.0", ID: hello.ID, Result: res})
		_ = c.Write(ctx, websocket.MessageText, out)
		_, data, err = c.Read(ctx)
		if err != nil {
			return
		}
		var prove msg
		var pp struct{ Proof string }
		_ = json.Unmarshal(data, &prove)
		_ = json.Unmarshal(prove.Params, &pp)
		if !hmac.Equal([]byte(pp.Proof), []byte(channel.Proof([]byte(key), "flai", mine, hp.Nonce))) {
			_ = c.Close(websocket.StatusPolicyViolation, "bad proof")
			return
		}
		d.Proven <- &Conn{Project: hp.Project.Key, Kind: hp.Kind, ws: c}
		<-ctx.Done()
	}))
	t.Cleanup(srv.Close)
	d.URL = srv.URL
	return d
}

// Wait returns the next proven connection.
func (d *Dashboard) Wait(t *testing.T) *Conn {
	t.Helper()
	select {
	case c := <-d.Proven:
		return c
	case <-time.After(5 * time.Second):
		t.Fatal("no flai connected")
		return nil
	}
}

// Ask sends one request and returns the raw result or the error.
func (c *Conn) Ask(t *testing.T, method, params string) (json.RawMessage, *channel.Error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := json.Marshal(msg{JSONRPC: "2.0", ID: json.RawMessage(`1`), Method: method, Params: json.RawMessage(params)})
	if err := c.ws.Write(ctx, websocket.MessageText, req); err != nil {
		t.Fatal(err)
	}
	_, data, err := c.ws.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var m msg
	_ = json.Unmarshal(data, &m)
	return m.Result, m.Error
}

// Next reads the next message flai sends by itself: a notification.
func (c *Conn) Next(t *testing.T) (method string, params json.RawMessage) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, data, err := c.ws.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var m msg
	_ = json.Unmarshal(data, &m)
	return m.Method, m.Params
}

// Close drops the connection from the dashboard's side.
func (c *Conn) Close() { _ = c.ws.Close(websocket.StatusGoingAway, "gone") }
