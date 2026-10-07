package cmd

import (
	"context"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

// listedDashboard is one published dashboard release as dashboard versions
// prints it.
type listedDashboard struct {
	Version    string `json:"version"`
	Tag        string `json:"tag"`
	Running    bool   `json:"running"`
	Configured bool   `json:"configured"`
	Latest     bool   `json:"latest"`
}

func newDashboardVersionsCmd(a *app) *cobra.Command {
	var src releaseSource
	c := &cobra.Command{
		Use:   "versions",
		Short: "List the published dashboard releases, marking the running and the configured one",
		Long: `Lists the published dashboard releases, the repository's flaiover/vX.Y.Z tags,
newest first, and installs nothing. Each release's image tag is its bare
version, X.Y.Z. It marks the newest, the release the running container's image
carries (its FLAIOVER_VERSION, else its image tag), and the configured tag
(flags, then the dashboard section of system-flow.yaml, then config). Docker
is not needed to list: without it, or with no container running, nothing is
marked running. flai dashboard upgrade --published --tag X.Y.Z deploys one of
them, an earlier one included, for that container once; dashboard.tag is not
changed. With --json it prints a list of version, tag, running, configured,
and latest.`,
		Example: `  flai dashboard versions
  flai dashboard versions --json
  flai dashboard upgrade --published --tag 0.4.0`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runDashboardVersions(cmd.Context(), src.options(a, ""))
		},
	}
	src.register(c, "GitHub repository whose flaiover/vX.Y.Z tags are the published dashboard releases")
	return c
}

func (a *app) runDashboardVersions(ctx context.Context, opt selfupgrade.Options) error {
	repo, err := a.projectOrNone()
	if err != nil {
		return err
	}
	s, err := a.dashboardSettings(repo, "", "", 0, "")
	if err != nil {
		return err
	}
	published, err := selfupgrade.ListTags(ctx, opt, selfupgrade.DashboardTagPrefix)
	if err != nil {
		return err
	}
	running := a.runningDashboardVersion(s.Name)
	listed := make([]listedDashboard, len(published))
	runningListed, configuredListed := false, false
	for i, r := range published {
		listed[i] = listedDashboard{Version: r.Version, Tag: r.Tag, Running: r.Version == running, Configured: r.Version == s.Tag, Latest: i == 0}
		runningListed = runningListed || listed[i].Running
		configuredListed = configuredListed || listed[i].Configured
	}
	if a.jsonOut {
		return a.printJSON(listed)
	}
	if len(listed) == 0 {
		fmt.Fprintf(a.out, "no published dashboard release in %s\n", opt.Repo)
	} else {
		fmt.Fprintf(a.out, "published dashboard releases in %s, newest first (image %s):\n", opt.Repo, s.Image)
		w := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)
		for _, r := range listed {
			var marks []string
			if r.Running {
				marks = append(marks, "running")
			}
			if r.Configured {
				marks = append(marks, "configured")
			}
			if r.Latest {
				marks = append(marks, "latest")
			}
			fmt.Fprintf(w, "  %s\t%s\n", r.Version, strings.Join(marks, ", "))
		}
		if err := w.Flush(); err != nil {
			return err
		}
	}
	switch {
	case running == "":
		fmt.Fprintf(a.out, "running: no %s container found (none runs, or Docker is not available)\n", s.Name)
	case !runningListed:
		fmt.Fprintf(a.out, "running: %s, which is not a published release\n", running)
	}
	if !configuredListed {
		fmt.Fprintf(a.out, "configured: %s, whose tag names no published release\n", s.ref())
	}
	return nil
}

// runningDashboardVersion is the release the running dashboard container's
// image carries: its FLAIOVER_VERSION, else its image tag when that is a
// version; "" when Docker is absent, no container runs, or neither says.
func (a *app) runningDashboardVersion(name string) string {
	if _, err := a.runner.LookPath("docker"); err != nil {
		a.logger().Debug("docker not found: no running dashboard marked", "component", "dashboard")
		return ""
	}
	up, err := a.containerRunning(name)
	if err != nil {
		a.logger().Warn("running dashboard not read", "component", "dashboard", "container", name, "err", err.Error())
		return ""
	}
	if !up {
		return ""
	}
	if out, err := a.runner.Run("", "docker", "inspect", "--format", "{{range .Config.Env}}{{println .}}{{end}}", name); err == nil {
		for _, l := range strings.Split(out, "\n") {
			if v, ok := strings.CutPrefix(strings.TrimSpace(l), "FLAIOVER_VERSION="); ok {
				if _, ok := buildinfo.Semver(v); ok {
					return v
				}
			}
		}
	}
	ref, _ := a.containerInfo(name)
	if i := strings.LastIndex(ref, ":"); i >= 0 {
		if tag := ref[i+1:]; !strings.Contains(tag, "/") {
			if _, ok := buildinfo.Semver(tag); ok {
				return tag
			}
		}
	}
	return ""
}

// requirePublishedDashboard refuses a tag that is not a published dashboard
// release, naming the published ones (ADR-0117 §3).
func requirePublishedDashboard(ctx context.Context, opt selfupgrade.Options, tag string) error {
	published, err := selfupgrade.ListTags(ctx, opt, selfupgrade.DashboardTagPrefix)
	if err != nil {
		return err
	}
	if len(published) == 0 {
		return fmt.Errorf("upgrade the dashboard to %s: %s has no published dashboard release (no tag %sX.Y.Z); check the repository the releases come from, or give --tag without --published to use a tag as it is", tag, opt.Repo, selfupgrade.DashboardTagPrefix)
	}
	for _, r := range published {
		if r.Version == tag {
			return nil
		}
	}
	list := selfupgrade.VersionList(published)
	return fmt.Errorf("upgrade the dashboard to %s: it is not a published dashboard release of %s (a tag %sX.Y.Z, whose image is tagged X.Y.Z); published are %s: give one of them as the bare X.Y.Z, or give --tag without --published to use a tag as it is", tag, opt.Repo, selfupgrade.DashboardTagPrefix, list)
}
