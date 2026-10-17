package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// forgeCheckTimeout bounds flai forge check's call to the forge.
const forgeCheckTimeout = 30 * time.Second

// githubAPIBase is the GitHub API flai's forge commands call; tests point it at a stand-in.
var githubAPIBase = forge.GitHubAPI

// newForgeCmd groups the commands about the project's forge (ADR-0165).
func newForgeCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "forge",
		Short: "The project's forge, through which pull-request acceptance merges a story",
		Long: `A project whose system-flow.yaml sets acceptance.mode: pull-request accepts a
story through a pull request on its forge, acceptance.forge (github, the
default), rather than by merging on this machine (ADR-0165). flai reaches the
forge's repository named by the origin remote, with a token it reads from the
environment variable this host names for the forge: GITHUB_TOKEN unless
flai serve forge <forge> --token-env <NAME> names another.`,
		Example: `  flai forge check
  flai forge check --json`,
	}
	c.AddCommand(newForgeCheckCmd(a))
	return c
}

func newForgeCheckCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Say whether the forge token this host names can read and push to the repository behind origin",
		Long: `Print the forge (acceptance.forge, github by default, under either
acceptance mode, so the token can be checked before switching), the
repository the origin remote names, the environment variable the token is
read from, and whether the token can read the repository and push to it. The
variable is named; the token never is.

The exit status is 1, with the reason, when the project has no origin remote,
origin is not a repository on the forge, the variable is unset or empty, the
forge does not accept the token, the repository is not found or not visible to
the token, the forge refuses the token or has rate-limited it, or the token can
read but not push. What was learned before the failure is printed still. With
--json the result is {"forge", "repository", "token_env", "read", "push",
"ok", "reason"}.`,
		Example: `  flai forge check
  flai forge check --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			cfg, _, err := a.loadConfig()
			if err != nil {
				return err
			}
			r := a.checkForge(repo.Root, repo.Manifest.Acceptance.ForgeOrDefault(), cfg)
			if err := a.printForgeCheck(r); err != nil {
				return err
			}
			if !r.OK {
				return errors.New(r.Reason)
			}
			return nil
		},
	}
}

// forgeCheck is what flai forge check found; it never holds the token.
type forgeCheck struct {
	Forge      string `json:"forge"`
	Repository string `json:"repository"`
	TokenEnv   string `json:"token_env"`
	Read       bool   `json:"read"`
	Push       bool   `json:"push"`
	OK         bool   `json:"ok"`
	Reason     string `json:"reason,omitempty"`
	asked      bool   // the forge answered, so Read and Push are its word
}

// checkForge asks name whether the token cfg names for it can read and push
// to the repository origin names in root, stopping at the first failure with
// its reason.
func (a *app) checkForge(root, name string, cfg config.Config) forgeCheck {
	r := forgeCheck{Forge: name, TokenEnv: cfg.ForgeTokenEnv(name)}
	remote, err := a.runner.Run(root, "git", "remote", "get-url", "origin")
	if err != nil {
		r.Reason = fmt.Sprintf("the project has no origin remote to read its repository from (%s): add the forge's repository with git remote add origin <url>", gitSaid(remote, err))
		return r
	}
	token := forge.TokenFunc(cfg.ForgeTokenFunc(name, os.Getenv))
	f, err := newForge(name, strings.TrimSpace(remote), token)
	if err != nil {
		r.Reason = err.Error()
		return r
	}
	r.Repository = f.Repository()
	if _, err := token(); err != nil {
		r.Reason = err.Error()
		return r
	}
	ctx, cancel := context.WithTimeout(context.Background(), forgeCheckTimeout)
	defer cancel()
	access, err := f.Access(ctx)
	if err != nil {
		r.asked = errors.Is(err, forge.ErrUnauthorized) || errors.Is(err, forge.ErrNotFound) || errors.Is(err, forge.ErrForbidden)
		r.Reason = forgeAccessReason(r, err)
		return r
	}
	r.asked, r.Read, r.Push = true, access.Read, access.Push
	if !r.Push {
		r.Reason = fmt.Sprintf("the token in %s can read %s but not push to it: give it push access to the repository (contents: write for a fine-grained token, the repo scope for a classic one)", r.TokenEnv, r.Repository)
		return r
	}
	r.OK = true
	return r
}

// newForge is the adapter for the forge name, bound to the repository remote
// names; a forge flai gains adds its branch here.
func newForge(name, remote string, token forge.TokenFunc) (forge.Forge, error) {
	if name == manifest.ForgeGitHub {
		owner, repo, err := forge.ParseGitHubRemote(remote)
		if err != nil {
			return nil, fmt.Errorf("origin is not a repository on github: %w; point origin at the GitHub repository stories are accepted into", err)
		}
		return forge.NewGitHub(forge.GitHubOptions{BaseURL: githubAPIBase, Owner: owner, Name: repo, Token: token})
	}
	return nil, fmt.Errorf("flai has no adapter for the forge %q; acceptance.forge takes %s", name, strings.Join(manifest.Forges, ", "))
}

// forgeAccessReason says why the forge would not answer whether the token
// can read the repository, naming the variable and never the token.
func forgeAccessReason(r forgeCheck, err error) string {
	switch {
	case errors.Is(err, forge.ErrUnauthorized):
		return fmt.Sprintf("%s did not accept the token in %s (%v): put a valid token in %s, or name the variable that holds one with flai serve forge %s --token-env <NAME>", r.Forge, r.TokenEnv, err, r.TokenEnv, r.Forge)
	case errors.Is(err, forge.ErrNotFound):
		return fmt.Sprintf("%s has no repository %s that the token in %s can see (%v): check that origin names it and that the token is granted access to it", r.Forge, r.Repository, r.TokenEnv, err)
	case errors.Is(err, forge.ErrForbidden):
		return fmt.Sprintf("%s refused the token in %s (%v): check that the token may read %s, or wait for the rate limit to reset", r.Forge, r.TokenEnv, err, r.Repository)
	}
	return fmt.Sprintf("could not ask %s about %s: %v", r.Forge, r.Repository, err)
}

// gitSaid is git's own words for a failed command, or the error when it said nothing.
func gitSaid(out string, err error) string {
	if s := strings.TrimSpace(out); s != "" {
		return s
	}
	return err.Error()
}

// printForgeCheck prints r, one finding a line, or as JSON; the reason goes
// out through the error path.
func (a *app) printForgeCheck(r forgeCheck) error {
	if a.jsonOut {
		return a.printJSON(r)
	}
	fmt.Fprintf(a.out, "forge: %s\n", r.Forge)
	if r.Repository != "" {
		fmt.Fprintf(a.out, "repository: %s (from origin)\n", r.Repository)
	}
	fmt.Fprintf(a.out, "token: read from %s\n", r.TokenEnv)
	if r.asked {
		fmt.Fprintf(a.out, "read: %s\npush: %s\n", yesNo(r.Read), yesNo(r.Push))
	}
	return nil
}

// yesNo is b as yes or no.
func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
