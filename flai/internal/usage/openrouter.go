package usage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// keyReads is how many times OpenRouter's /v1/key is read for one reading:
// it answers from copies that lag one another, and a key's usage only grows,
// so the highest of several reads is the newest (S-0351).
const keyReads = 3

// keyReadPause is the wait between two reads of /v1/key, as S-0351's check
// spaced them.
const keyReadPause = time.Second

// maxSpendBody bounds how much of a gateway's answer is decoded.
const maxSpendBody = 1 << 20

// errNotFound is OpenRouter's 404, for a generation it has not recorded.
var errNotFound = errors.New("not found")

// OpenRouter reads what OpenRouter charged for a run, by the run's
// generation IDs through /v1/generation or by the key's usage through
// /v1/key before and after it.
type OpenRouter struct {
	root   string
	key    string
	client *http.Client
	pause  time.Duration
}

// NewOpenRouter is OpenRouter's spend reader for the provider's base URL,
// https://openrouter.ai/api or https://openrouter.ai/api/v1, called with key.
func NewOpenRouter(baseURL, key string, client *http.Client) *OpenRouter {
	root := strings.TrimSuffix(strings.TrimRight(baseURL, "/"), "/v1")
	return &OpenRouter{root: root, key: key, client: spendClient(client), pause: keyReadPause}
}

// Run is what OpenRouter charged for the run: the sum of its generations
// when it has their IDs, else the rise in the key's usage since KeyBefore.
func (o *OpenRouter) Run(ctx context.Context, run RunRef) (*Spent, error) {
	switch {
	case len(run.Generations) > 0:
		return o.byGeneration(ctx, run.Generations)
	case run.KeyBefore != nil:
		return o.byKey(ctx, *run.KeyBefore)
	}
	return nil, fmt.Errorf("OpenRouter: the run has no generation IDs and no reading of the key's usage from its start: %w", ErrNoSpend)
}

// KeyUsage is the key's usage in US dollars, the highest of keyReads reads.
func (o *OpenRouter) KeyUsage(ctx context.Context) (float64, error) {
	var best float64
	for i := range keyReads {
		if i > 0 {
			if err := pause(ctx, o.pause); err != nil {
				return 0, fmt.Errorf("read the key's usage from OpenRouter: %w", err)
			}
		}
		var body struct {
			Data *struct {
				Usage *float64 `json:"usage"`
			} `json:"data"`
		}
		endpoint, err := o.get(ctx, "/v1/key", &body)
		if err != nil {
			return 0, err
		}
		if body.Data == nil || body.Data.Usage == nil {
			return 0, fmt.Errorf("GET %s answered without data.usage; check that the provider's base_url is OpenRouter's API", endpoint)
		}
		best = max(best, *body.Data.Usage)
	}
	return best, nil
}

// byGeneration sums total_cost over the run's generations; a generation not
// recorded yet leaves the run without spend until the next measurement.
func (o *OpenRouter) byGeneration(ctx context.Context, ids []string) (*Spent, error) {
	spent := &Spent{By: SpentByGeneration}
	missing := 0
	for _, id := range ids {
		var body struct {
			Data *struct {
				TotalCost *float64 `json:"total_cost"`
			} `json:"data"`
		}
		endpoint, err := o.get(ctx, "/v1/generation?id="+url.QueryEscape(id), &body)
		if errors.Is(err, errNotFound) {
			missing++
			continue
		}
		if err != nil {
			return nil, err
		}
		if body.Data == nil || body.Data.TotalCost == nil {
			return nil, fmt.Errorf("GET %s answered without data.total_cost; check that the provider's base_url is OpenRouter's API", endpoint)
		}
		spent.Cost += *body.Data.TotalCost
		spent.Calls++
	}
	if missing > 0 {
		return nil, fmt.Errorf("OpenRouter has not recorded %d of the run's %d generations yet; the next measurement reads them again: %w", missing, len(ids), ErrNoSpend)
	}
	return spent, nil
}

// byKey is the rise in the key's usage from before, the reading taken when
// the run started.
func (o *OpenRouter) byKey(ctx context.Context, before float64) (*Spent, error) {
	now, err := o.KeyUsage(ctx)
	if err != nil {
		return nil, err
	}
	if now <= before {
		return nil, fmt.Errorf("OpenRouter's key usage $%.6f has not risen from $%.6f read when the run started; the next measurement reads it again: %w", now, before, ErrNoSpend)
	}
	return &Spent{Cost: now - before, By: SpentByKey}, nil
}

// get GETs path under the API root with the key as a bearer and decodes the
// answer into v. It answers the endpoint it read, errNotFound for a 404,
// and an error naming the endpoint and the status for any other answer but
// a 2xx; the key is in neither.
func (o *OpenRouter) get(ctx context.Context, path string, v any) (string, error) {
	endpoint := o.root + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return endpoint, fmt.Errorf("make the request GET %s: %w; check the provider's base_url", endpoint, err)
	}
	req.Header.Set("Authorization", "Bearer "+o.key)
	resp, err := o.client.Do(req)
	if err != nil {
		return endpoint, fmt.Errorf("read the spend from OpenRouter: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	switch {
	case resp.StatusCode == http.StatusNotFound:
		return endpoint, fmt.Errorf("GET %s answered %d: %w", endpoint, resp.StatusCode, errNotFound)
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return endpoint, fmt.Errorf("GET %s answered %d %s; check the key the provider's key_env names", endpoint, resp.StatusCode, http.StatusText(resp.StatusCode))
	case resp.StatusCode < 200 || resp.StatusCode > 299:
		return endpoint, fmt.Errorf("GET %s answered %d %s", endpoint, resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxSpendBody)).Decode(v); err != nil {
		return endpoint, fmt.Errorf("decode the answer to GET %s: %w", endpoint, err)
	}
	return endpoint, nil
}

// pause waits d, or until ctx is done.
func pause(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
