package manifest

import (
	"maps"
	"reflect"
	"strings"
	"testing"
)

// S-0349: an entry's values are checked only where they are set, each refusal
// naming the key and what to write.
func TestProviderProblems(t *testing.T) {
	good := Provider{API: APIAnthropicMessages, BaseURL: "http://127.0.0.1:4000", KeyEnv: "LITELLM_API_KEY", SpendKeyEnv: "LITELLM_MASTER_KEY", Models: map[string]string{"opus": "anthropic/claude-opus-5-5", "haiku": "claude-haiku-4-5@2026"}}
	if errs := good.Problems("providers.litellm"); len(errs) != 0 {
		t.Errorf("good: %v", errs)
	}
	if errs := (Provider{}).Problems("providers.litellm"); len(errs) != 0 {
		t.Errorf("nothing set: %v", errs)
	}
	for _, c := range []struct {
		p    Provider
		want string
	}{
		{Provider{API: "openai"}, `providers.litellm.api "openai" is not an API flai speaks; write anthropic-messages, openai-chat, openai-responses`},
		{Provider{BaseURL: "litellm.example.com"}, `providers.litellm.base_url "litellm.example.com" is not an absolute http or https URL`},
		{Provider{BaseURL: "/v1"}, `providers.litellm.base_url "/v1" is not an absolute http or https URL`},
		{Provider{BaseURL: "ftp://litellm.example.com"}, `providers.litellm.base_url "ftp://litellm.example.com" is not an absolute http or https URL`},
		{Provider{BaseURL: "https://"}, `providers.litellm.base_url "https://" is not an absolute http or https URL`},
		{Provider{BaseURL: "https://me:secret@litellm.example.com"}, `providers.litellm.base_url carries a user or password`},
		{Provider{KeyEnv: "1KEY"}, `providers.litellm.key_env is not the name of an environment variable`},
		{Provider{KeyEnv: "LITELLM-KEY"}, `providers.litellm.key_env is not the name of an environment variable`},
		{Provider{SpendKeyEnv: "1KEY"}, `providers.litellm.spend_key_env is not the name of an environment variable; write one such as LITELLM_MASTER_KEY`},
		{Provider{SpendKeyEnv: "LITELLM MASTER KEY"}, `providers.litellm.spend_key_env is not the name of an environment variable`},
		{Provider{Models: map[string]string{"": "claude-opus-5-5"}}, `providers.litellm.models alias "" is not a name`},
		{Provider{Models: map[string]string{"opus": ""}}, `providers.litellm.models.opus "" is not a model ID`},
		{Provider{Models: map[string]string{"opus": "has space"}}, `providers.litellm.models.opus "has space" is not a model ID`},
	} {
		errs := c.p.Problems("providers.litellm")
		if len(errs) != 1 || !strings.Contains(errs[0], c.want) {
			t.Errorf("%+v: got %q, want %q", c.p, errs, c.want)
		}
	}
	// A key pasted where its variable's name goes, or a password in the URL,
	// is never repeated in the message.
	for _, p := range []Provider{{KeyEnv: "sk-or-v1-0123456789"}, {SpendKeyEnv: "sk-or-v1-0123456789"}, {BaseURL: "https://me:sk-or-v1-0123456789@litellm.example.com"}} {
		if errs := p.Problems("providers.litellm"); len(errs) != 1 || strings.Contains(errs[0], "sk-or-v1-0123456789") {
			t.Errorf("%+v: got %q, want one problem without the secret", p, errs)
		}
	}
}

// S-0349: a host's entry is whole, each part it lacks named with the flag
// that sets it, and its name has a harness's shape.
func TestHostProviderProblems(t *testing.T) {
	whole := Provider{API: APIOpenAIChat, BaseURL: "https://openrouter.ai/api/v1", KeyEnv: "OPENROUTER_API_KEY"}
	if errs := HostProviderProblems(map[string]Provider{"openrouter": whole}); len(errs) != 0 {
		t.Errorf("whole: %v", errs)
	}
	if errs := HostProviderProblems(nil); len(errs) != 0 {
		t.Errorf("none: %v", errs)
	}
	withSpend := whole
	withSpend.SpendKeyEnv = "OPENROUTER_MANAGEMENT_KEY"
	if errs := HostProviderProblems(map[string]Provider{"openrouter": withSpend}); len(errs) != 0 {
		t.Errorf("whole, with the optional spend_key_env: %v", errs)
	}
	got := strings.Join(HostProviderProblems(map[string]Provider{"litellm": {Models: map[string]string{"opus": "x y"}}}), "\n")
	for _, w := range []string{
		"agent.providers.litellm has no api; set it with flai serve agent provider litellm --api",
		"agent.providers.litellm has no base_url; set it with flai serve agent provider litellm --base-url",
		"agent.providers.litellm has no key_env; set it with flai serve agent provider litellm --key-env",
		`agent.providers.litellm.models.opus "x y" is not a model ID`,
	} {
		if !strings.Contains(got, w) {
			t.Errorf("got:\n%s\nwant %q", got, w)
		}
	}
	errs := HostProviderProblems(map[string]Provider{"Open Router": whole})
	if len(errs) != 1 || !strings.Contains(errs[0], `agent.providers "Open Router" is not a provider name`) {
		t.Errorf("a name that is not one: %q", errs)
	}
}

// S-0349: a provider resolves to the host's entry with the project's override
// laid over it, api and base_url where set and models alias by alias, without
// changing either map; no name is no provider, and a name the host lacks, or
// an entry that merged is not valid, is refused.
func TestResolveProvider(t *testing.T) {
	host := map[string]Provider{
		"litellm":    {API: APIAnthropicMessages, BaseURL: "http://127.0.0.1:4000", KeyEnv: "LITELLM_API_KEY", Models: map[string]string{"opus": "claude-opus-5-5", "haiku": "claude-haiku-4-5"}},
		"openrouter": {API: APIOpenAIChat, BaseURL: "https://openrouter.ai/api/v1", KeyEnv: "OPENROUTER_API_KEY"},
	}
	project := map[string]Provider{
		"litellm":    {BaseURL: "https://litellm.example.com", Models: map[string]string{"opus": "anthropic/claude-opus-5-5", "sonnet": "anthropic/claude-sonnet-5"}},
		"openrouter": {API: APIOpenAIResponses, Models: map[string]string{"opus": "anthropic/claude-opus-5-5"}},
		"anthropic":  {API: APIAnthropicMessages},
	}
	hostBefore, projectBefore := cloneProviders(host), cloneProviders(project)

	got, err := ResolveProvider("litellm", host, project)
	if err != nil {
		t.Fatal(err)
	}
	want := Provider{API: APIAnthropicMessages, BaseURL: "https://litellm.example.com", KeyEnv: "LITELLM_API_KEY", Models: map[string]string{"opus": "anthropic/claude-opus-5-5", "haiku": "claude-haiku-4-5", "sonnet": "anthropic/claude-sonnet-5"}}
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("litellm:\n got %+v\nwant %+v", *got, want)
	}
	got.Models["opus"] = "changed"
	got, err = ResolveProvider("openrouter", host, project)
	if err != nil {
		t.Fatal(err)
	}
	want = Provider{API: APIOpenAIResponses, BaseURL: "https://openrouter.ai/api/v1", KeyEnv: "OPENROUTER_API_KEY", Models: map[string]string{"opus": "anthropic/claude-opus-5-5"}}
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("openrouter, models only in the project's:\n got %+v\nwant %+v", *got, want)
	}
	got.Models["opus"] = "changed"
	if got, err := ResolveProvider("openrouter", host, nil); err != nil || !reflect.DeepEqual(*got, host["openrouter"]) {
		t.Errorf("no override: %+v, %v", got, err)
	}
	if !reflect.DeepEqual(host, hostBefore) || !reflect.DeepEqual(project, projectBefore) {
		t.Errorf("an input was changed:\n host %+v\n project %+v", host, project)
	}

	if got, err := ResolveProvider("", host, project); got != nil || err != nil {
		t.Errorf("no name: %+v, %v", got, err)
	}
	if _, err := ResolveProvider("anthropic", host, project); err == nil || !strings.Contains(err.Error(), "flai serve agent provider anthropic") {
		t.Errorf("a name the host lacks, though the project overrides it: %v", err)
	}
	if _, err := ResolveProvider("litellm", nil, nil); err == nil || !strings.Contains(err.Error(), "flai serve agent provider litellm") {
		t.Errorf("no host entries: %v", err)
	}
	if _, err := ResolveProvider("litellm", host, map[string]Provider{"litellm": {BaseURL: "litellm.example.com"}}); err == nil || !strings.Contains(err.Error(), `base_url "litellm.example.com"`) {
		t.Errorf("an override that is not valid: %v", err)
	}
	if _, err := ResolveProvider("bare", map[string]Provider{"bare": {API: APIOpenAIChat}}, nil); err == nil || !strings.Contains(err.Error(), "has no key_env") {
		t.Errorf("a host entry that is not whole: %v", err)
	}
}

// S-0357, ADR-0132: spend_key_env names a variable of the host's environment,
// so a manifest entry that sets it is refused with the command that sets it
// on the host, and the variable's value, which may be the key pasted in its
// place, is never repeated; a resolved provider keeps the host's.
func TestSpendKeyEnvIsTheHosts(t *testing.T) {
	const pasted = "sk-litellm-master-0123456789"
	for _, v := range []string{"LITELLM_MASTER_KEY", pasted} {
		errs := (Manifest{Providers: map[string]Provider{"litellm": {SpendKeyEnv: v}}}).ProviderErrors()
		if len(errs) != 1 || !strings.Contains(errs[0], "providers.litellm.spend_key_env is the host's") ||
			!strings.Contains(errs[0], "flai serve agent provider litellm --spend-key-env <variable>") || strings.Contains(errs[0], pasted) {
			t.Errorf("spend_key_env %q in a manifest: %q", v, errs)
		}
	}
	host := map[string]Provider{"litellm": {API: APIAnthropicMessages, BaseURL: "http://127.0.0.1:4000", KeyEnv: "LITELLM_API_KEY", SpendKeyEnv: "LITELLM_MASTER_KEY"}}
	for _, project := range []map[string]Provider{nil, {"litellm": {BaseURL: "https://litellm.example.com", SpendKeyEnv: "PROJECT_KEY"}}} {
		got, err := ResolveProvider("litellm", host, project)
		if err != nil || got.SpendKeyEnv != "LITELLM_MASTER_KEY" {
			t.Errorf("project %+v: %+v, %v; want the host's spend_key_env", project, got, err)
		}
	}
}

func cloneProviders(in map[string]Provider) map[string]Provider {
	out := map[string]Provider{}
	for n, p := range in {
		p.Models = maps.Clone(p.Models)
		out[n] = p
	}
	return out
}
