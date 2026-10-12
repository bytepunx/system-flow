// Package channeltest is a stand-in for the dashboard's agent endpoint, for
// tests of whatever dials it.
package channeltest

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/bytepunx/system-flow/flai/internal/channel"
)

// Dashboard accepts flai's connection, does the handshake with Key, and
// hands each proven connection to Proven. It answers hello with a release
// stamp, "flaiover <Version> <commit>" signed by a key it generates, which
// Keys verify: a client given Keys as its ReleaseKeys accepts it (S-0237).
type Dashboard struct {
	URL    string
	Proven chan *Conn
	// Keys verify the stamps Sign makes and the one it answers with by default.
	Keys []*ecdsa.PublicKey

	signer  *ecdsa.PrivateKey
	mu      sync.Mutex
	release *channel.Release // what it answers hello with; nil sends none
	refuse  string           // when set, the fault each hello is closed with
}

// Version is what the dashboard's default stamp and its answer to hello name.
const Version = "9.9.9"

// Conn is one proven connection from a flai.
type Conn struct {
	Project string
	Kind    string // what hello said: empty for a project, channel.KindCandidate for one offered for import
	// Release is the stamp flai sent in hello; nil when it sent none.
	Release *channel.Release
	ws      *websocket.Conn
}

// Sign returns statement as a release stamp signed by the dashboard's key,
// as cosign sign-blob signs: ECDSA P-256 over its SHA-256, ASN.1 DER, base64.
func (d *Dashboard) Sign(statement string) channel.Release {
	sum := sha256.Sum256([]byte(statement))
	sig, err := ecdsa.SignASN1(rand.Reader, d.signer, sum[:])
	if err != nil {
		panic("channeltest: sign a release stamp: " + err.Error())
	}
	return channel.Release{Statement: statement, Signature: base64.StdEncoding.EncodeToString(sig)}
}

// Answer sets the stamp the dashboard answers hello with from now on; nil
// answers without one, as a dashboard of protocol 1 does.
func (d *Dashboard) Answer(r *channel.Release) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.release = r
}

// Refuse closes every hello from now on with channel.CloseUnsigned and the
// reason "unsigned peer: <fault>", as a dashboard that refuses an unsigned
// flai does; an empty fault answers hello again.
func (d *Dashboard) Refuse(fault string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.refuse = fault
}

func (d *Dashboard) answer() (*channel.Release, string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.release, d.refuse
}

type msg struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *channel.Error  `json:"error,omitempty"`
}

// New starts a dashboard that holds key and answers with a signed stamp.
func New(t *testing.T, key string) *Dashboard {
	t.Helper()
	signer, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	d := &Dashboard{Proven: make(chan *Conn, 16), Keys: []*ecdsa.PublicKey{&signer.PublicKey}, signer: signer}
	stamp := d.Sign("flaiover " + Version + " 0ef6d26")
	d.release = &stamp
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
			Nonce   string           `json:"nonce"`
			Project channel.Project  `json:"project"`
			Kind    string           `json:"kind"`
			Release *channel.Release `json:"release"`
		}
		_ = json.Unmarshal(data, &hello)
		_ = json.Unmarshal(hello.Params, &hp)
		release, refuse := d.answer()
		if refuse != "" {
			_ = c.Close(channel.CloseUnsigned, channel.UnsignedReason+": "+refuse)
			return
		}
		b := make([]byte, 16)
		_, _ = rand.Read(b)
		mine := hex.EncodeToString(b)
		res, _ := json.Marshal(struct {
			Nonce     string           `json:"nonce"`
			Proof     string           `json:"proof"`
			Dashboard string           `json:"dashboard"`
			Release   *channel.Release `json:"release,omitempty"`
		}{mine, channel.Proof([]byte(key), "dashboard", hp.Nonce, mine), Version, release})
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
		d.Proven <- &Conn{Project: hp.Project.Key, Kind: hp.Kind, Release: hp.Release, ws: c}
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

// Last reads what flai sends until it closes the connection, and returns the
// notifications' methods in order.
func (c *Conn) Last(t *testing.T) []string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var out []string
	for {
		_, data, err := c.ws.Read(ctx)
		if err != nil {
			if ctx.Err() != nil {
				t.Fatal("flai did not close the connection")
			}
			return out
		}
		var m msg
		_ = json.Unmarshal(data, &m)
		out = append(out, m.Method)
	}
}

// Close drops the connection from the dashboard's side.
func (c *Conn) Close() { _ = c.ws.Close(websocket.StatusGoingAway, "gone") }
