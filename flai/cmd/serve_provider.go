package cmd

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// flai serve agent provider: where the model calls of a story's agent that
// names a provider go, on this host (ADR-0129). The flags given change those
// fields of the entry, --model given at all replaces its aliases, and --reset
// removes it. No flag takes a key, and nothing here reads the variable that
// holds it.
func newServeAgentProviderCmd(a *app) *cobra.Command {
	var p manifest.Provider
	var models []string
	var reset bool
	c := &cobra.Command{
		Use:   "provider [<name>] [--api <api>] [--base-url <url>] [--key-env <variable>] [--spend-key-env <variable>] [--model <alias>=<name>]... [--reset]",
		Short: "Set where the model calls of a story's agent that names a provider go: its API, its URL, and the variable that holds its key",
		Long: `A story's agent may name a provider, such as litellm or openrouter, beside
its harness and model; you say here what that is on this host. --api is what
the provider speaks (` + strings.Join(manifest.ProviderAPIs, ", ") + `), --base-url
where it answers, and --key-env the name of the environment variable that
holds its key. No flag takes a key: the key stays in that variable, which
flai serve reads when it starts the agent, and is never written to the
configuration. --spend-key-env, optional, names the variable of a key that
can read the gateway's spend log where the agent's own key cannot, such as
LiteLLM's master key; an empty value clears it. --model <alias>=<name>,
repeated, maps the harness's model aliases to names on the provider; given
at all, it replaces the entry's aliases.

A new entry needs --api, --base-url, and --key-env; on an entry that is
set, the flags given change those fields and leave the rest. --reset removes
the entry. With a name and no flag, the entry is shown; with no name, every
entry. A project's manifest may override an entry's api, base_url, and
models under providers.<name>, never its key_env or spend_key_env.`,
		Example: `  flai serve agent provider
  flai serve agent provider litellm --api anthropic-messages --base-url http://127.0.0.1:4000 --key-env LITELLM_API_KEY
  flai serve agent provider openrouter --api openai-chat --base-url https://openrouter.ai/api/v1 --key-env OPENROUTER_API_KEY --model opus=anthropic/claude-opus-5-5 --model sonnet=anthropic/claude-sonnet-5
  flai serve agent provider litellm --base-url https://litellm.example.com
  flai serve agent provider litellm --spend-key-env LITELLM_MASTER_KEY
  flai serve agent provider litellm --reset`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			changed := func(flags ...string) bool {
				return slices.ContainsFunc(flags, cmd.Flags().Changed)
			}
			given := changed("api", "base-url", "key-env", "spend-key-env", "model")
			if len(args) == 0 {
				if given || reset {
					return fmt.Errorf("name the provider, such as litellm: flai serve agent provider <name> --api <api> --base-url <url> --key-env <variable>")
				}
				return a.showProviders("")
			}
			name := args[0]
			if !given && !reset {
				return a.showProviders(name)
			}
			if reset && given {
				return fmt.Errorf("--reset takes no other flag: it removes provider %s whole", name)
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if reset {
				cfg.Agent = cfg.Agent.WithoutProvider(name)
			} else {
				entry := cfg.Agent.Providers[name]
				if changed("api") {
					entry.API = p.API
				}
				if changed("base-url") {
					entry.BaseURL = p.BaseURL
				}
				if changed("key-env") {
					entry.KeyEnv = p.KeyEnv
				}
				if changed("spend-key-env") {
					entry.SpendKeyEnv = p.SpendKeyEnv
				}
				if changed("model") {
					if entry.Models, err = providerModels(models); err != nil {
						return err
					}
				}
				if cfg.Agent, err = cfg.Agent.WithProvider(name, entry); err != nil {
					return err
				}
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showProviders(name)
		},
	}
	c.Flags().StringVar(&p.API, "api", "", "what the provider speaks: "+strings.Join(manifest.ProviderAPIs, ", "))
	c.Flags().StringVar(&p.BaseURL, "base-url", "", "the provider's absolute http or https URL, with no user or password in it")
	c.Flags().StringVar(&p.KeyEnv, "key-env", "", "the name of the environment variable that holds the provider's key, such as LITELLM_API_KEY; never the key")
	c.Flags().StringVar(&p.SpendKeyEnv, "spend-key-env", "", "the name of the environment variable that holds a key that can read the provider's spend log where the agent's key cannot, such as LITELLM_MASTER_KEY; never the key; empty clears it")
	c.Flags().StringArrayVar(&models, "model", nil, "a model alias and its name on the provider, alias=name (repeatable; replaces the entry's aliases)")
	c.Flags().BoolVar(&reset, "reset", false, "remove the provider's entry")
	return c
}

// providerModels are the model aliases --model gives, alias=name each; the
// names themselves are checked with the rest of the entry.
func providerModels(pairs []string) (map[string]string, error) {
	out := map[string]string{}
	for _, pair := range pairs {
		alias, model, ok := strings.Cut(pair, "=")
		alias = strings.TrimSpace(alias)
		if !ok || alias == "" {
			return nil, fmt.Errorf("--model %q is not alias=name; write one such as opus=anthropic/claude-opus-5-5", pair)
		}
		out[alias] = strings.TrimSpace(model)
	}
	return out, nil
}

// showProviders prints the host's provider name, or every one when name is
// empty: what each speaks, where, and the name of the variable that holds its
// key and, where set, of the one that holds its spend log's key, never a key.
func (a *app) showProviders(name string) error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	providers := cfg.Agent.Providers
	if name != "" {
		providers = map[string]manifest.Provider{}
		if p, ok := cfg.Agent.Providers[name]; ok {
			providers[name] = p
		}
	}
	if a.jsonOut {
		if providers == nil {
			providers = map[string]manifest.Provider{}
		}
		return a.printJSON(map[string]any{"providers": providers})
	}
	const add = "--api <api> --base-url <url> --key-env <variable>"
	switch {
	case len(providers) > 0:
	case name == "":
		fmt.Fprintf(a.out, "no provider is set on this host; add one with flai serve agent provider <name> %s\n", add)
		return nil
	default:
		fmt.Fprintf(a.out, "provider %s has no entry on this host; add it with flai serve agent provider %s %s\n", name, name, add)
		return nil
	}
	for _, n := range slices.Sorted(maps.Keys(providers)) {
		p := providers[n]
		fmt.Fprintf(a.out, "provider %s: %s at %s, its key in the variable %s\n", n, p.API, p.BaseURL, p.KeyEnv)
		if p.SpendKeyEnv != "" {
			fmt.Fprintf(a.out, "  spend log read with the key in the variable %s\n", p.SpendKeyEnv)
		}
		for _, alias := range slices.Sorted(maps.Keys(p.Models)) {
			fmt.Fprintf(a.out, "  model %s: %s\n", alias, p.Models[alias])
		}
	}
	return nil
}

// flai serve forge: the name of the environment variable a forge's token is
// read from, on this host (S-0392), kept beside the providers as their keys'
// variables are (ADR-0129). No flag takes a token, and the command reads the
// variable only to say whether it is set, never printing its value.
func newServeForgeCmd(a *app) *cobra.Command {
	var tokenEnv string
	var reset bool
	c := &cobra.Command{
		Use:   "forge [<forge>] [--token-env <NAME>] [--reset]",
		Short: "Set the environment variable a forge's token is read from, such as GITHUB_TOKEN for github",
		Long: `A project that accepts by pull request reaches its forge, such as github,
with a token; you say here which environment variable holds it on this host.
--token-env names the variable, ` + config.DefaultGitHubTokenEnv + ` for github when you name none, and
--reset goes back to that default. No flag takes a token: it stays in the
variable, which flai reads only when it calls the forge, and is never written
to the configuration, a log, a work item, or a start record.

With no forge, every forge flai calls is shown; with a forge and no flag,
that forge: the variable its token is read from, and whether that variable is
set in this environment, never its value. flai serve reads its own
environment, so export the variable where flai serve starts too.`,
		Example: `  flai serve forge
  flai serve forge github --token-env FLAI_GITHUB_TOKEN
  flai serve forge github --reset
  flai serve forge --json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			given := cmd.Flags().Changed("token-env")
			if len(args) == 0 {
				if given || reset {
					return fmt.Errorf("name the forge, such as github: flai serve forge <forge> --token-env <NAME>")
				}
				return a.showForges("")
			}
			name := args[0]
			if !slices.Contains(config.KnownForges(), name) {
				return fmt.Errorf("flai calls no forge %q; it calls %s", name, strings.Join(config.KnownForges(), ", "))
			}
			if !given && !reset {
				return a.showForges(name)
			}
			if reset && given {
				return fmt.Errorf("--reset takes no --token-env: it goes back to the variable %s", config.DefaultForgeTokenEnv(name))
			}
			if given && strings.TrimSpace(tokenEnv) == "" {
				return fmt.Errorf("--token-env needs the name of a variable, such as %s; --reset goes back to it", config.DefaultForgeTokenEnv(name))
			}
			cfg, path, err := a.loadConfig()
			if err != nil {
				return err
			}
			if reset {
				cfg = cfg.WithoutForge(name)
			} else if cfg, err = cfg.WithForge(name, config.ForgeHost{TokenEnv: tokenEnv}); err != nil {
				return err
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			return a.showForges(name)
		},
	}
	c.Flags().StringVar(&tokenEnv, "token-env", "", "the name of the environment variable that holds the forge's token, such as "+config.DefaultGitHubTokenEnv+"; never the token")
	c.Flags().BoolVar(&reset, "reset", false, "read the forge's token from its default variable again")
	return c
}

// showForges prints forge, or every forge flai calls when forge is empty: the
// variable its token is read from, whether that is the default, and whether
// it is set and not empty in this environment, never its value.
func (a *app) showForges(forge string) error {
	cfg, _, err := a.loadConfig()
	if err != nil {
		return err
	}
	names := config.KnownForges()
	if forge != "" {
		names = []string{forge}
	}
	out := map[string]any{}
	for _, n := range names {
		env := cfg.ForgeTokenEnv(n)
		value, ok := os.LookupEnv(env)
		set := ok && strings.TrimSpace(value) != ""
		isDefault := env == config.DefaultForgeTokenEnv(n)
		if a.jsonOut {
			out[n] = map[string]any{"token_env": env, "default": isDefault, "token_set": set}
			continue
		}
		note := ""
		if isDefault {
			note = " (default)"
		}
		state := "set in this environment"
		if !set {
			state = "not set in this environment"
		}
		fmt.Fprintf(a.out, "forge %s: its token in the variable %s%s, %s\n", n, env, note, state)
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"forges": out})
	}
	return nil
}
