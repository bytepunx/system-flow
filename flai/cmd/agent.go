package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

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

An agent may also name roles: the agents for the work a story's agent hands
to sub-agents, explore and verify (S-0189), each with a harness, model, and
config of its own (--role-harness, --role-model, --role-config). flai serve
starts claude-code with each role's model over the project's sub-agent
definition for it.

flai agent shows the default, set changes it (only what you give; --config
adds or replaces keys, --role-config role.key= removes one, --unset-role
removes a role), and clear removes it.`,
		Example: `  flai agent
  flai agent set --harness claude-code --model claude-opus-5-5 --config effort=high
  flai agent set --model claude-sonnet-5
  flai agent set --role-model explore=haiku --role-model verify=sonnet
  flai agent clear`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showAgent("") },
	}
	var harness, model string
	var config, unset, trailers []string
	var replace, autocommit bool
	var rf roleFlags
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
			roles, err := rf.roles()
			if err != nil {
				return err
			}
			if harness == "" && model == "" && len(cfg) == 0 && len(unset) == 0 && !rf.given() {
				return fmt.Errorf("nothing to set: give --harness, --model, --config key=value, --unset key, or a --role- flag")
			}
			base := repo.Manifest.Agent
			if replace {
				base = nil
			}
			next := base.With(&manifest.Agent{Harness: harness, Model: model, Config: cfg, Roles: roles})
			if next != nil {
				for _, k := range unset {
					delete(next.Config, k)
				}
				if len(next.Config) == 0 {
					next.Config = nil
				}
				rf.prune(next)
				if next.IsZero() {
					next = nil
				}
			}
			if err := next.Validate(); err != nil {
				return err
			}
			if err := manifest.WriteAgent(filepath.Join(repo.Root, manifest.File), next); err != nil {
				return err
			}
			return a.showAgent(a.commitManifest(repo, autocommit, "chore: set the project's default agent", trailers))
		},
	}
	set.Flags().StringVar(&harness, "harness", "", "the harness that runs the agent, such as claude-code")
	set.Flags().StringVar(&model, "model", "", "the model it runs, such as claude-opus-5-5")
	set.Flags().StringArrayVar(&config, "config", nil, "an option for the harness, key=value (repeatable)")
	set.Flags().StringArrayVar(&unset, "unset", nil, "remove an option by key (repeatable)")
	rf.register(set, true)
	set.Flags().BoolVar(&replace, "replace", false, "start from nothing: the default becomes exactly what is given")
	set.Flags().BoolVar(&autocommit, "autocommit", false, "commit system-flow.yaml on its own, unless dashboard.autocommit is false")
	set.Flags().StringArrayVar(&trailers, "trailer", nil, "a trailer line for the commit (repeatable)")
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
			return a.showAgent(a.commitManifest(repo, autocommit, "chore: clear the project's default agent", trailers))
		},
	}
	clearCmd.Flags().BoolVar(&autocommit, "autocommit", false, "commit system-flow.yaml on its own, unless dashboard.autocommit is false")
	clearCmd.Flags().StringArrayVar(&trailers, "trailer", nil, "a trailer line for the commit (repeatable)")
	c.AddCommand(set, clearCmd)
	return c
}

// commitManifest commits system-flow.yaml alone, when asked and the project
// commits what flai writes for the dashboard; it says the commit, or nothing.
func (a *app) commitManifest(repo *workitem.Repo, autocommit bool, msg string, trailers []string) string {
	if !autocommit || !repo.Manifest.Autocommit() {
		return ""
	}
	if len(trailers) > 0 {
		msg += "\n\n" + strings.Join(trailers, "\n")
	}
	if _, err := a.runner.Run(repo.Root, "git", "add", "--", manifest.File); err != nil {
		return ""
	}
	if _, err := a.runner.Run(repo.Root, "git", "commit", "-q", "-m", msg, "--", manifest.File); err != nil {
		return ""
	}
	sha, err := a.runner.Run(repo.Root, "git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(sha)
}

func (a *app) showAgent(commit string) error {
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
		out := map[string]any{"agent": m.Agent}
		if commit != "" {
			out["commit"] = commit
		}
		return a.printJSON(out)
	}
	if m.Agent.IsZero() {
		fmt.Fprintln(a.out, "no default agent; new stories carry none (flai agent set --harness ... --model ...)")
		return nil
	}
	fmt.Fprintf(a.out, "default agent: %s\n  every story created from now on gets it; flai story new and flai edit override it for one story\n", m.Agent)
	return nil
}

// roleFlags are the agent's roles as flags give them (S-0189): the harness
// and model of a role (role=value), its config (role.key=value, removed when
// given no value), and roles to remove whole.
type roleFlags struct {
	harness, model, config, unset []string
}

// register adds the role flags to c; --unset-role only where the agent is
// laid over one that exists.
func (rf *roleFlags) register(c *cobra.Command, unset bool) {
	c.Flags().StringArrayVar(&rf.harness, "role-harness", nil, "the harness of a sub-agent role, role=harness, such as verify=claude-code (repeatable)")
	c.Flags().StringArrayVar(&rf.model, "role-model", nil, "the model a sub-agent role runs, role=model, such as verify=sonnet (repeatable)")
	c.Flags().StringArrayVar(&rf.config, "role-config", nil, "an option for a sub-agent role, role.key=value (repeatable)")
	if unset {
		c.Flags().StringArrayVar(&rf.unset, "unset-role", nil, "remove a sub-agent role (repeatable)")
	}
}

func (rf roleFlags) given() bool {
	return len(rf.harness)+len(rf.model)+len(rf.config)+len(rf.unset) > 0
}

// roles are the roles to lay over an agent: what the flags set, without the
// config keys given no value, which prune removes.
func (rf roleFlags) roles() (map[string]manifest.Role, error) {
	roles, err := manifest.ParseRoles(rf.harness, rf.model, rf.config)
	if err != nil {
		return nil, err
	}
	for n, r := range roles {
		for k, v := range r.Config {
			if v == "" {
				delete(r.Config, k)
			}
		}
		if len(r.Config) == 0 {
			r.Config = nil
		}
		roles[n] = r
	}
	return roles, nil
}

// prune removes from next the role config keys given no value, the roles
// unset, and a role left setting nothing.
func (rf roleFlags) prune(next *manifest.Agent) {
	if next == nil {
		return
	}
	given, _ := manifest.ParseRoles(nil, nil, rf.config)
	for n, r := range given {
		o, ok := next.Roles[n]
		if !ok {
			continue
		}
		for k, v := range r.Config {
			if v == "" {
				delete(o.Config, k)
			}
		}
		if len(o.Config) == 0 {
			o.Config = nil
		}
		next.Roles[n] = o
	}
	for _, n := range rf.unset {
		delete(next.Roles, n)
	}
	for n, r := range next.Roles {
		if r.IsZero() {
			delete(next.Roles, n)
		}
	}
	if len(next.Roles) == 0 {
		next.Roles = nil
	}
}

// agentFlags is the agent --harness, --model, --agent-config, and the role
// flags give, nil when none is given.
func agentFlags(harness, model string, config []string, rf roleFlags) (*manifest.Agent, error) {
	cfg, err := manifest.ParseConfig(config)
	if err != nil {
		return nil, err
	}
	roles, err := rf.roles()
	if err != nil {
		return nil, err
	}
	a := &manifest.Agent{Harness: harness, Model: model, Config: cfg, Roles: roles}
	if a.IsZero() {
		return nil, nil
	}
	return a, a.Validate()
}

// editedAgent is the agent a story will have after flai edit's agent flags:
// merged into what it has, or, with --clear-agent, exactly what is given; a
// key given with no value (key=, role.key=) is removed, and so is a role
// --unset-role names. Nil means none.
func editedAgent(repo *workitem.Repo, id string, replace bool, harness, model string, config []string, rf roleFlags) (*manifest.Agent, error) {
	cfg, err := manifest.ParseConfig(config)
	if err != nil {
		return nil, err
	}
	roles, err := rf.roles()
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
	next := base.With(&manifest.Agent{Harness: harness, Model: model, Config: set, Roles: roles})
	if next != nil {
		for k, v := range cfg {
			if v == "" {
				delete(next.Config, k)
			}
		}
		if len(next.Config) == 0 {
			next.Config = nil
		}
		rf.prune(next)
		if next.IsZero() {
			next = nil
		}
	}
	return next, nil
}
