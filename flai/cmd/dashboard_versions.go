package cmd

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/dashboard"
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
	// Unverifiable marks a release with no digest list this flai can
	// verify, which flai dashboard upgrade --published refuses (ADR-0070);
	// Signature says why (missing, bad, or unknown-key) and Reason in words.
	// All three are left out otherwise.
	Unverifiable bool                      `json:"unverifiable,omitempty"`
	Signature    selfupgrade.SignatureKind `json:"signature,omitempty"`
	Reason       string                    `json:"reason,omitempty"`
}

func newDashboardVersionsCmd(a *app) *cobra.Command {
	var src releaseSource
	c := &cobra.Command{
		Use:   "versions",
		Short: "List the published dashboard releases, marking the running, the configured, and the not deployable ones",
		Long: `Lists the published dashboard releases, the repository's flaiover/vX.Y.Z tags,
newest first, and installs nothing. Each release's image tag is its bare
version, X.Y.Z. It marks the newest, the release the running container's image
carries (its FLAIOVER_VERSION, else its image tag), and the configured tag
(flags, then the dashboard section of system-flow.yaml, then config). It marks
"not deployable", with the reason, each release with no digest list this flai
can verify, such as one from before signing (ADR-0070): flai dashboard
refuses to run it. To check, it downloads each release's digest list and its
signature; a download that fails fails the list. Docker is not needed to
list: without it, or with no container running, nothing is marked running.
flai dashboard upgrade --published --tag X.Y.Z deploys one of them, an
earlier one included, for that container once; dashboard.tag is not changed.
With --json it prints a list of version, tag, running, configured, and
latest, and, for a release that is not deployable, unverifiable, signature
(missing, bad, or unknown-key), and reason.`,
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
	refusals, err := digestListRefusals(ctx, opt, published)
	if err != nil {
		return err
	}
	running := a.runningDashboardVersion(s.Name)
	listed := make([]listedDashboard, len(published))
	runningListed, configuredListed := false, false
	for i, r := range published {
		listed[i] = listedDashboard{Version: r.Version, Tag: r.Tag, Running: r.Version == running, Configured: r.Version == s.Tag, Latest: i == 0}
		if refused := refusals[i]; refused != nil {
			listed[i].Unverifiable, listed[i].Signature, listed[i].Reason = true, refused.Kind, refused.Reason
		}
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
			if r.Unverifiable {
				marks = append(marks, "not deployable: "+r.Reason)
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
// release, naming the published ones (ADR-0117 §3), and one with no digest
// list this flai can verify, naming why and the deployable ones (ADR-0070).
func requirePublishedDashboard(ctx context.Context, opt selfupgrade.Options, tag string) error {
	published, err := selfupgrade.ListTags(ctx, opt, selfupgrade.DashboardTagPrefix)
	if err != nil {
		return err
	}
	if len(published) == 0 {
		return fmt.Errorf("upgrade the dashboard to %s: %s has no published dashboard release (no tag %sX.Y.Z); check the repository the releases come from, or give --image and --tag without --published to run another image's tag as it is, unsigned", tag, opt.Repo, selfupgrade.DashboardTagPrefix)
	}
	i := slices.IndexFunc(published, func(r selfupgrade.Release) bool { return r.Version == tag })
	if i < 0 {
		list := selfupgrade.VersionList(published)
		return fmt.Errorf("upgrade the dashboard to %s: it is not a published dashboard release of %s (a tag %sX.Y.Z, whose image is tagged X.Y.Z); published are %s: give one of them as the bare X.Y.Z, or give --image and --tag without --published to run another image's tag as it is, unsigned", tag, opt.Repo, selfupgrade.DashboardTagPrefix, list)
	}
	chosen, err := digestListRefusals(ctx, opt, published[i:i+1])
	if err != nil || chosen[0] == nil {
		return err
	}
	// only a refusal costs the other releases' lists, to name those to deploy
	refusals, err := digestListRefusals(ctx, opt, published)
	if err != nil {
		return err
	}
	var deployable []selfupgrade.Release
	for j, r := range published {
		if refusals[j] == nil {
			deployable = append(deployable, r)
		}
	}
	head := fmt.Sprintf("upgrade the dashboard to %s: it is not deployable: %s, so this flai cannot verify the image it names (ADR-0070)", tag, chosen[0].Reason)
	if len(deployable) == 0 {
		return fmt.Errorf("%s; no published dashboard release of %s has a digest list this flai can verify: upgrade flai (flai self-upgrade), or give --image and --tag without --published to run another image's tag as it is, unsigned", head, opt.Repo)
	}
	return fmt.Errorf("%s; deployable are %s: give one of them as the bare X.Y.Z", head, selfupgrade.VersionList(deployable))
}

// digestListRefusals are, for each of tags, the published dashboard releases,
// why this flai cannot verify its release's digest list, nil for one it can.
// The kind is as SignatureError's, and Reason is worded for a list that marks
// the release. A tag with no GitHub release costs no download; one with a
// list costs the list and its signature, and its cosign.pub when no trusted
// key verifies it. A download that fails fails them all: a release that could
// not be checked is not one that cannot be verified.
func digestListRefusals(ctx context.Context, opt selfupgrade.Options, tags []selfupgrade.Release) ([]*selfupgrade.SignatureError, error) {
	refusals := make([]*selfupgrade.SignatureError, len(tags))
	if len(tags) == 0 {
		return refusals, nil
	}
	keys, err := selfupgrade.TrustedKeys(opt)
	if err != nil {
		return nil, fmt.Errorf("check the dashboard releases' digest lists: %w", err)
	}
	rels, err := selfupgrade.List(ctx, opt, selfupgrade.DashboardTagPrefix)
	if err != nil {
		return nil, err
	}
	byVersion := make(map[string]selfupgrade.Release, len(rels))
	for _, r := range rels {
		byVersion[r.Version] = r
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errs := make([]error, len(tags))
	slots := make(chan struct{}, verifyAtOnce)
	var wg sync.WaitGroup
	for i, t := range tags {
		rel, ok := byVersion[t.Version]
		if !ok {
			refusals[i] = &selfupgrade.SignatureError{Tag: t.Tag, Kind: selfupgrade.SignatureMissing, Reason: "it has no GitHub release with a signed digest list"}
			continue
		}
		wg.Go(func() {
			slots <- struct{}{}
			defer func() { <-slots }()
			if refusals[i], errs[i] = digestListRefusal(ctx, opt, rel, keys); errs[i] != nil {
				cancel()
			}
		})
	}
	wg.Wait()
	// the failure that cancelled the others, rather than one it caused
	for _, err := range errs {
		if err != nil && !errors.Is(err, context.Canceled) {
			return nil, err
		}
	}
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return refusals, nil
}

// digestListRefusal is why this flai cannot verify rel's digest list, nil
// when it can, or the failure to check it, as dashboard.FetchList refuses
// it, worded for a list that marks the release: a list that does not parse
// or is for another version is bad, as a signature that does not verify is.
func digestListRefusal(ctx context.Context, opt selfupgrade.Options, rel selfupgrade.Release, keys []*ecdsa.PublicKey) (*selfupgrade.SignatureError, error) {
	_, _, _, err := dashboard.FetchList(ctx, opt, rel, keys)
	var refused *dashboard.Refusal
	switch {
	case errors.As(err, &refused):
	case err != nil:
		return nil, fmt.Errorf("check whether this flai can verify the digest list of dashboard release %s: %w", rel.Tag, err)
	default:
		return nil, nil
	}
	refusal := func(kind selfupgrade.SignatureKind, reason string) (*selfupgrade.SignatureError, error) {
		return &selfupgrade.SignatureError{Tag: rel.Tag, Kind: kind, Reason: reason}, nil
	}
	switch refused.Kind {
	case dashboard.RefusalMissing:
		return refusal(selfupgrade.SignatureMissing, "it publishes no signed digest list")
	case dashboard.RefusalUnknownKey:
		return refusal(selfupgrade.SignatureUnknownKey, "its digest list is signed by a key this flai does not trust")
	case dashboard.RefusalMalformed:
		return refusal(selfupgrade.SignatureBad, "its digest list is not in the format release-flaiover.yml writes")
	case dashboard.RefusalVersion:
		return refusal(selfupgrade.SignatureBad, "its digest list is for another version")
	}
	return refusal(selfupgrade.SignatureBad, "its digest list's signature does not verify, so it may have been altered")
}
