package usage

import "testing"

func TestForChoosesTheReaderByHost(t *testing.T) {
	for _, c := range []struct {
		baseURL    string
		openRouter string // the API root of OpenRouter's reader, or "" for another
	}{
		{"https://openrouter.ai/api", "https://openrouter.ai/api"},
		{"https://openrouter.ai/api/v1", "https://openrouter.ai/api"},
		{"https://OpenRouter.ai/api/v1/", "https://OpenRouter.ai/api"},
		{"https://eu.openrouter.ai/api", "https://eu.openrouter.ai/api"},
		{"https://notopenrouter.ai/api", ""},
		{"https://openrouter.ai.example.com/api", ""},
		{"http://localhost:4000", ""},
	} {
		s := For(c.baseURL, testKey, nil)
		o, ok := s.(*OpenRouter)
		switch {
		case c.openRouter != "" && (!ok || o.root != c.openRouter || o.key != testKey):
			t.Errorf("For(%q) = %#v; want OpenRouter's reader at %s", c.baseURL, s, c.openRouter)
		case c.openRouter == "" && ok:
			t.Errorf("For(%q) is OpenRouter's reader; want another", c.baseURL)
		}
	}
	for _, c := range []struct{ baseURL, root string }{
		{"http://127.0.0.1:4000", "http://127.0.0.1:4000"},
		{"http://localhost:4000/v1/", "http://localhost:4000"},
		{"https://notopenrouter.ai/api", "https://notopenrouter.ai/api"},
		{"https://openrouter.ai.example.com/api", "https://openrouter.ai.example.com/api"},
	} {
		l, ok := For(c.baseURL, testKey, nil).(*LiteLLM)
		if !ok || l.root != c.root || l.key != testKey {
			t.Errorf("For(%q) = %#v; want LiteLLM's reader at %s", c.baseURL, l, c.root)
		}
	}
	for _, baseURL := range []string{"https://api.anthropic.com", "https://API.anthropic.com/v1", "::not a url", ""} {
		if s := For(baseURL, testKey, nil); s != nil {
			t.Errorf("For(%q) = %#v; want none", baseURL, s)
		}
	}
}

func TestForGivesAClientWithATimeout(t *testing.T) {
	o, ok := For("https://openrouter.ai/api", testKey, nil).(*OpenRouter)
	if !ok || o.client == nil || o.client.Timeout != spendTimeout {
		t.Errorf("For with no client = %#v; want a client with a %s timeout", o, spendTimeout)
	}
	l, ok := For("http://127.0.0.1:4000", testKey, nil).(*LiteLLM)
	if !ok || l.client == nil || l.client.Timeout != spendTimeout {
		t.Errorf("For with no client = %#v; want a client with a %s timeout", l, spendTimeout)
	}
}
