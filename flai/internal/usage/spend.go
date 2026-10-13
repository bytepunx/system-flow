package usage

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// What a provider's gateway charged for a run, read from its spend log after
// the run (ADR-0132). Each gateway has a reader of its own; a reader is
// handed the key the caller looked up, sends it as a bearer, and never logs
// it or puts it in an error.

// How a reader joined the gateway's spend to a run.
const (
	// SpentByGeneration is the sum of the run's calls, looked up by the
	// gateway's generation IDs.
	SpentByGeneration = "generation"
	// SpentByKey is the rise in the key's usage between the run's start and
	// now.
	SpentByKey = "key"
	// SpentBySession is the gateway's log rows tagged with the run's
	// harness session IDs.
	SpentBySession = "session"
	// SpentByWindow is the gateway's log rows for the key within the run's
	// window.
	SpentByWindow = "window"
)

// ErrNoSpend is what a reader returns, wrapped with why, when the gateway
// answers but holds nothing for the run yet.
var ErrNoSpend = errors.New("the gateway holds no spend for the run")

// spendTimeout bounds one request to a gateway when the caller gives no
// client of its own.
const spendTimeout = 20 * time.Second

// RunRef is what a reader joins a gateway's spend to a run by.
type RunRef struct {
	// Sessions are the run's harness session IDs.
	Sessions []string
	// From and To are the run's window.
	From, To time.Time
	// Generations are the gateway's own IDs for the run's calls, where the
	// harness surfaces them; Claude Code surfaces none.
	Generations []string
	// KeyBefore is the key's usage in US dollars read when the run started,
	// by the reader's KeyReader; nil when none was read.
	KeyBefore *float64
	// KeyHash is the SHA-256 in hex of the key the run's calls used, which
	// LiteLLM files its spend log rows under; never the key itself.
	KeyHash string
}

// Spent is what the gateway charged for a run.
type Spent struct {
	// Cost is in US dollars.
	Cost float64
	// Calls is how many calls the gateway counted, where it says; 0 when it
	// does not.
	Calls int
	// By is how the reader joined the spend to the run: one of the SpentBy
	// constants.
	By string
}

// Spend reads what a gateway charged for a run.
type Spend interface {
	// Run is what the gateway charged for the run, or an error wrapping
	// ErrNoSpend when it holds nothing for it yet.
	Run(ctx context.Context, run RunRef) (*Spent, error)
}

// KeyReader reads a key's usage, which a reader that joins by the key's
// usage before and after a run takes when the run starts.
type KeyReader interface {
	// KeyUsage is the key's usage to date in US dollars.
	KeyUsage(ctx context.Context) (float64, error)
}

// For is the spend reader for a provider whose gateway is at baseURL, called
// with key: OpenRouter's for openrouter.ai and its subdomains, none for
// api.anthropic.com, whose harness figure stands, none for a base URL that
// does not parse, and LiteLLM's for any other host. A nil client is one with
// a timeout per request.
func For(baseURL, key string, client *http.Client) Spend {
	u, err := url.Parse(baseURL)
	if err != nil || u.Hostname() == "" {
		return nil
	}
	host := strings.ToLower(u.Hostname())
	switch {
	case host == "openrouter.ai" || strings.HasSuffix(host, ".openrouter.ai"):
		return NewOpenRouter(baseURL, key, client)
	case host == "api.anthropic.com":
		return nil
	}
	return NewLiteLLM(baseURL, key, client)
}

// spendClient is client, or one with spendTimeout when client is nil.
func spendClient(client *http.Client) *http.Client {
	if client != nil {
		return client
	}
	return &http.Client{Timeout: spendTimeout}
}
