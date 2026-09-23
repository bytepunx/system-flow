package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai agent is the project's default agent (S-0103): the harness, model,
// and options every story created while it is set gets in its front matter.
// It is not flai serve agent, the host's command for starting one (S-0079).
func newAgentCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "agent",
		Short: "The project's default agent: the harness, model, and options every new story gets",
		Long: `Who works a story is an agent: a harness (the command line or API that runs
it, such as claude-code), a model (such as claude-opus-5-5), and options for
that harness (config, key=value). The project's default is kept in
system-flow.yaml; every story created while it is set gets a copy in its front
matter, which flai story new --harness/--model/--agent-config and flai edit
override for that story. Stories created before, or with no default set, carry
none.

flai agent shows the default, set changes it (only what you give; --config
adds or replaces keys), and clear removes it.`,
		Example: `  flai agent
  flai agent set --harness claude-code --model claude-opus-5-5 --config effort=high
  flai agent set --model claude-sonnet-5
  flai agent clear`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showAgent() },
	}
	var harness, model string
	var config, unset []string
	set := &cobra.Command{
		Use:   "set",
		Short: "Set the project's default harness, model, or options",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			cfg, err := manifest.ParseConfig(config)
			if err != nil {
				return err
			}
			if harness == "" && model == "" && len(cfg) == 0 && len(unset) == 0 {
				return fmt.Errorf("nothing to set: give --harness, --model, --config key=value, or --unset key")
			}
			next := repo.Manifest.Agent.With(&manifest.Agent{Harness: harness, Model: model, Config: cfg})
			if next != nil {
				for _, k := range unset {
					delete(next.Config, k)
				}
				if len(next.Config) == 0 {
					next.Config = nil
				}
			}
			if err := manifest.WriteAgent(filepath.Join(repo.Root, manifest.File), next); err != nil {
				return err
			}
			return a.showAgent()
		},
	}
	set.Flags().StringVar(&harness, "harness", "", "the harness that runs the agent, such as claude-code")
	set.Flags().StringVar(&model, "model", "", "the model it runs, such as claude-opus-5-5")
	set.Flags().StringArrayVar(&config, "config", nil, "an option for the harness, key=value (repeatable)")
	set.Flags().StringArrayVar(&unset, "unset", nil, "remove an option by key (repeatable)")
	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Remove the project's default agent; stories keep what they have",
		Args:  cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			if err := manifest.WriteAgent(filepath.Join(repo.Root, manifest.File), nil); err != nil {
				return err
			}
			return a.showAgent()
		},
	}
	c.AddCommand(set, clearCmd)
	return c
}

func (a *app) showAgent() error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	// read again: set and clear have just written it
	m, err := manifest.Load(filepath.Join(repo.Root, manifest.File))
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"agent": m.Agent})
	}
	if m.Agent.IsZero() {
		fmt.Fprintln(a.out, "no default agent; new stories carry none (flai agent set --harness ... --model ...)")
		return nil
	}
	fmt.Fprintf(a.out, "default agent: %s\n  every story created from now on gets it; flai story new and flai edit override it for one story\n", m.Agent)
	return nil
}

// agentFlags is the agent --harness, --model, and --agent-config give, nil
// when none is given.
func agentFlags(harness, model string, config []string) (*manifest.Agent, error) {
	cfg, err := manifest.ParseConfig(config)
	if err != nil {
		return nil, err
	}
	a := &manifest.Agent{Harness: harness, Model: model, Config: cfg}
	if a.IsZero() {
		return nil, nil
	}
	return a, a.Validate()
}

// editedAgent is the agent a story will have after flai edit's agent flags:
// merged into what it has, or, with --clear-agent, exactly what is given; a
// key given with no value (key=) is removed. Nil means none.
func editedAgent(repo *workitem.Repo, id string, replace bool, harness, model string, config []string) (*manifest.Agent, error) {
	cfg, err := manifest.ParseConfig(config)
	if err != nil {
		return nil, err
	}
	var base *manifest.Agent
	if !replace {
		it, err := repo.Get(id)
		if err != nil {
			return nil, err
		}
		base = it.Agent
	}
	set := map[string]string{}
	for k, v := range cfg {
		if v != "" {
			set[k] = v
		}
	}
	next := base.With(&manifest.Agent{Harness: harness, Model: model, Config: set})
	if next != nil {
		for k, v := range cfg {
			if v == "" {
				delete(next.Config, k)
			}
		}
		if len(next.Config) == 0 {
			next.Config = nil
		}
		if next.IsZero() {
			next = nil
		}
	}
	return next, nil
}
