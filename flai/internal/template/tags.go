package template

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// Tag is a release tag of a template: its name as tagged and its version.
type Tag struct {
	Name    string // as tagged, such as v1.0.60
	Version string // X.Y.Z, such as 1.0.60
	parts   [3]int
}

// Remote is what a template's git remote offers: its default branch, its
// branches, and its release tags.
type Remote struct {
	DefaultBranch string
	Branches      []string
	Tags          []Tag // ascending by version
}

// ListRemote reads repo's default branch, branches, and release tags with
// one git ls-remote. A local directory has no remote: it answers an empty
// Remote, whose Latest reports none, so its template.yaml version stands.
func ListRemote(r execx.Runner, repo string) (Remote, error) {
	if IsLocal(repo) {
		return Remote{}, nil
	}
	if err := execx.Require(r, "git", "Install git to list the template's releases, or point --template at a local directory."); err != nil {
		return Remote{}, err
	}
	out, err := r.Run("", "git", "ls-remote", "--symref", repo, "HEAD", "refs/heads/*", "refs/tags/*")
	if err != nil {
		return Remote{}, fmt.Errorf("list the branches and tags of template %s: %w; check the repository URL and your network, or pass --ref to name the version", repo, err)
	}
	return parseLsRemote(out), nil
}

// parseLsRemote reads git ls-remote --symref output: the symbolic HEAD line
// names the default branch; refs/heads lines are branches; refs/tags lines
// whose name is X.Y.Z, with or without a leading v, are release tags.
func parseLsRemote(out string) Remote {
	var rm Remote
	seen := map[string]bool{}
	for line := range strings.SplitSeq(out, "\n") {
		left, ref, ok := strings.Cut(strings.TrimSpace(line), "\t")
		if !ok {
			continue
		}
		if target, sym := strings.CutPrefix(left, "ref: "); sym {
			if ref == "HEAD" {
				rm.DefaultBranch = strings.TrimPrefix(target, "refs/heads/")
			}
			continue
		}
		if branch, isHead := strings.CutPrefix(ref, "refs/heads/"); isHead {
			rm.Branches = append(rm.Branches, branch)
			continue
		}
		name, isTag := strings.CutPrefix(ref, "refs/tags/")
		if !isTag {
			continue
		}
		name = strings.TrimSuffix(name, "^{}")
		if seen[name] {
			continue
		}
		seen[name] = true
		if tag, ok := releaseTag(name); ok {
			rm.Tags = append(rm.Tags, tag)
		}
	}
	sort.SliceStable(rm.Tags, func(i, j int) bool {
		a, b := rm.Tags[i], rm.Tags[j]
		if a.parts != b.parts {
			return less(a.parts, b.parts)
		}
		return a.Name < b.Name
	})
	return rm
}

// releaseTag reads a tag name as a release: X.Y.Z with or without a leading
// v. A pre-release or build suffix is not a release.
func releaseTag(name string) (Tag, bool) {
	if strings.ContainsAny(name, "-+") {
		return Tag{}, false
	}
	parts, ok := buildinfo.Semver(name)
	if !ok {
		return Tag{}, false
	}
	return Tag{Name: name, Version: fmt.Sprintf("%d.%d.%d", parts[0], parts[1], parts[2]), parts: parts}, true
}

func less(a, b [3]int) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

// Latest returns the newest release tag, and false when there is none.
func (rm Remote) Latest() (Tag, bool) {
	if len(rm.Tags) == 0 {
		return Tag{}, false
	}
	return rm.TagFor(rm.Tags[len(rm.Tags)-1].Version)
}

// TagFor returns the release tag of a version X.Y.Z (with or without a
// leading v), preferring v<version>, the name flai template push tags.
func (rm Remote) TagFor(version string) (Tag, bool) {
	want, ok := releaseTag(version)
	if !ok {
		return Tag{}, false
	}
	var found Tag
	for _, t := range rm.Tags {
		if t.parts != want.parts {
			continue
		}
		if t.Name == "v"+want.Version {
			return t, true
		}
		if found.Name == "" {
			found = t
		}
	}
	return found, found.Name != ""
}

// Match returns the release tag a ref the operator gives names: the tag of
// that name, or, when no branch or tag has the name, the tag of the version
// it spells, so 1.0.60 finds v1.0.60.
func (rm Remote) Match(ref string) (Tag, bool) {
	for _, t := range rm.Tags {
		if t.Name == ref {
			return t, true
		}
	}
	if rm.isBranch(ref) {
		return Tag{}, false
	}
	return rm.TagFor(ref)
}

// FollowsReleases reports whether ref follows the template's releases, and
// so resolves to the newest tag: no ref, the default branch, or a release
// tag. Any other branch, or a commit, is a deliberate choice used as given.
func (rm Remote) FollowsReleases(ref string) bool {
	if ref == "" || ref == rm.DefaultBranch {
		return true
	}
	_, ok := rm.Match(ref)
	return ok
}

func (rm Remote) isBranch(name string) bool {
	for _, b := range rm.Branches {
		if b == name {
			return true
		}
	}
	return false
}
