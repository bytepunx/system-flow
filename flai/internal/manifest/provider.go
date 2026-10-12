package manifest

import (
	"fmt"
	"maps"
	"net/url"
	"regexp"
	"slices"
	"strings"
)

// Provider is where a story's agent's model calls go (ADR-0129): the API a
// gateway speaks, its base URL, the name of the host's environment variable
// that holds its key, and the harness's model aliases on it. The host's
// configuration holds each entry whole, under agent.providers; a project's
// manifest, under providers, may override an entry's api, base_url, and
// models, never key_env. It never holds a key.
type Provider struct {
	API     string            `yaml:"api,omitempty" json:"api,omitempty"`
	BaseURL string            `yaml:"base_url,omitempty" json:"base_url,omitempty"`
	KeyEnv  string            `yaml:"key_env,omitempty" json:"key_env,omitempty"`
	Models  map[string]string `yaml:"models,omitempty" json:"models,omitempty"`
}

// The APIs a provider may speak.
const (
	APIAnthropicMessages = "anthropic-messages"
	APIOpenAIChat        = "openai-chat"
	APIOpenAIResponses   = "openai-responses"
)

// ProviderAPIs are the values of a provider's api, in the order they are
// listed to the operator.
var ProviderAPIs = []string{APIAnthropicMessages, APIOpenAIChat, APIOpenAIResponses}

var envNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Problems are what is wrong with the values p sets, one phrase each, each
// beginning with what, the entry's key such as providers.litellm: an api flai
// does not know, a base_url that is not an absolute http or https URL, a
// key_env that is not a variable's name, and a model alias or name that is
// not one. What p leaves out is not checked.
func (p Provider) Problems(what string) []string {
	var errs []string
	if p.API != "" && !slices.Contains(ProviderAPIs, p.API) {
		errs = append(errs, fmt.Sprintf("%s.api %q is not an API flai speaks; write %s", what, p.API, strings.Join(ProviderAPIs, ", ")))
	}
	if p.BaseURL != "" {
		if r := baseURLProblem(p.BaseURL); r != "" {
			errs = append(errs, fmt.Sprintf("%s.base_url %s", what, r))
		}
	}
	// The value is not quoted: one that is not a variable's name may be the
	// key itself, pasted in its place.
	if p.KeyEnv != "" && !envNamePattern.MatchString(p.KeyEnv) {
		errs = append(errs, fmt.Sprintf("%s.key_env is not the name of an environment variable; write one such as LITELLM_API_KEY (letters, digits, _, not beginning with a digit), never the key", what))
	}
	for _, alias := range sortedKeys(p.Models) {
		switch model := p.Models[alias]; {
		case !configKeyPattern.MatchString(alias):
			errs = append(errs, fmt.Sprintf("%s.models alias %q is not a name such as opus or sonnet (lower case letters, digits, . _ -)", what, alias))
		case !modelPattern.MatchString(model):
			errs = append(errs, fmt.Sprintf("%s.models.%s %q is not a model ID on the provider, such as anthropic/claude-opus-5-5", what, alias, model))
		}
	}
	return errs
}

// baseURLProblem says what is wrong with a base URL, or "" when it is an
// absolute http or https URL with a host and no user or password in it. It
// quotes the URL only when it carries no user or password.
func baseURLProblem(s string) string {
	u, err := url.Parse(s)
	switch {
	case err == nil && u.User != nil:
		return "carries a user or password; write the URL without it, and put the key in the variable key_env names"
	case err != nil || u.Host == "" || u.Scheme != "http" && u.Scheme != "https":
		return fmt.Sprintf("%q is not an absolute http or https URL; write one such as https://litellm.example.com or http://127.0.0.1:4000", s)
	}
	return ""
}

// providerNameProblem is what is wrong with a provider's name, under what such
// as providers, or "" when it is a name of a harness's shape.
func providerNameProblem(what, name string) string {
	if harnessPattern.MatchString(name) {
		return ""
	}
	return fmt.Sprintf("%s %q is not a provider name such as litellm or openrouter (lower case letters, digits, . _ -)", what, name)
}

// ProviderErrors are what is wrong with the manifest's providers, one sentence
// each: a name that is not one, a value Problems refuses, and a key_env, which
// is the host's to set.
func (m Manifest) ProviderErrors() []string {
	var errs []string
	for _, name := range sortedKeys(m.Providers) {
		if r := providerNameProblem("providers", name); r != "" {
			errs = append(errs, r)
			continue
		}
		p := m.Providers[name]
		what := "providers." + name
		if p.KeyEnv != "" {
			errs = append(errs, fmt.Sprintf("%s.key_env is the host's: it names a variable of the host's environment, which a committed manifest cannot know; remove it here and set it with flai serve agent provider %s --key-env <variable>", what, name))
			p.KeyEnv = ""
		}
		errs = append(errs, p.Problems(what)...)
	}
	return errs
}

// HostProviderProblems are what is wrong with the host's providers, the
// configuration's agent.providers, one sentence each: a name that is not one,
// a value Problems refuses, and an entry without the api, base_url, or key_env
// every host entry sets.
func HostProviderProblems(providers map[string]Provider) []string {
	var errs []string
	for _, name := range sortedKeys(providers) {
		if r := providerNameProblem("agent.providers", name); r != "" {
			errs = append(errs, r)
			continue
		}
		what := "agent.providers." + name
		errs = append(errs, incomplete(what, name, providers[name])...)
		errs = append(errs, providers[name].Problems(what)...)
	}
	return errs
}

// incomplete are the fields a whole entry sets and p leaves out, each naming
// the flag of flai serve agent provider that sets it.
func incomplete(what, name string, p Provider) []string {
	var errs []string
	for _, f := range []struct{ key, flag, value, example string }{
		{"api", "api", p.API, strings.Join(ProviderAPIs, ", ")},
		{"base_url", "base-url", p.BaseURL, "such as https://litellm.example.com"},
		{"key_env", "key-env", p.KeyEnv, "the name of the variable that holds the key, such as LITELLM_API_KEY"},
	} {
		if strings.TrimSpace(f.value) == "" {
			errs = append(errs, fmt.Sprintf("%s has no %s; set it with flai serve agent provider %s --%s <%s> (%s)", what, f.key, name, f.flag, f.key, f.example))
		}
	}
	return errs
}

// ResolveProvider is the provider name stands for: a copy of the host's entry
// with the project's override for it laid over, its api and base_url where
// the project sets them and its models merged alias by alias, the project's
// winning. Nil for no name. It refuses a name the host has no entry for, and
// an entry that, merged, lacks what a whole entry sets or sets a value
// Problems refuses. Neither map is changed.
func ResolveProvider(name string, host, project map[string]Provider) (*Provider, error) {
	if name == "" {
		return nil, nil
	}
	h, ok := host[name]
	if !ok {
		return nil, fmt.Errorf("provider %q has no entry on this host: add it with flai serve agent provider %s --api <api> --base-url <url> --key-env <variable>, or name another provider in the agent", name, name)
	}
	out := h
	out.Models = maps.Clone(h.Models)
	if over, ok := project[name]; ok {
		if over.API != "" {
			out.API = over.API
		}
		if over.BaseURL != "" {
			out.BaseURL = over.BaseURL
		}
		for alias, model := range over.Models {
			if out.Models == nil {
				out.Models = map[string]string{}
			}
			out.Models[alias] = model
		}
	}
	what := "provider " + name
	if errs := append(incomplete(what, name, out), out.Problems(what)...); len(errs) > 0 {
		return nil, fmt.Errorf("start with provider %s: %s", name, strings.Join(errs, "; "))
	}
	return &out, nil
}
