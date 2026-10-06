package template

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

// fakeGit answers git commands from a table keyed by the first argument and
// records every call.
type fakeGit struct {
	out   map[string]string
	fail  map[string]error
	calls [][]string
}

func (f *fakeGit) Run(dir, name string, args ...string) (string, error) {
	f.calls = append(f.calls, append([]string{name}, args...))
	if len(args) == 0 {
		return "", nil
	}
	return f.out[args[0]], f.fail[args[0]]
}

func (f *fakeGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return f.Run(dir, name, args...)
}

func (f *fakeGit) LookPath(name string) (string, error) { return "/usr/bin/" + name, nil }

func (f *fakeGit) ran(sub string) bool {
	return slices.ContainsFunc(f.calls, func(c []string) bool { return len(c) > 1 && c[1] == sub })
}

const lsRemote = "ref: refs/heads/main\tHEAD\n" +
	"aaa\tHEAD\n" +
	"aaa\trefs/heads/main\n" +
	"bbb\trefs/heads/next\n" +
	"ccc\trefs/tags/v1.0.9\n" +
	"ddd\trefs/tags/v1.0.10\n" +
	"eee\trefs/tags/v1.0.10^{}\n" +
	"fff\trefs/tags/1.0.18\n" +
	"ggg\trefs/tags/v1.1.0-rc.1\n" +
	"hhh\trefs/tags/before-split\n" +
	"iii\trefs/tags/v1.0.60\n" +
	"jjj\trefs/tags/v1.0.60^{}\n"

func listFake(t *testing.T, out string) Remote {
	t.Helper()
	f := &fakeGit{out: map[string]string{"ls-remote": out}}
	rm, err := ListRemote(f, "https://example.com/t.git")
	if err != nil {
		t.Fatal(err)
	}
	if len(f.calls) != 1 || !f.ran("ls-remote") {
		t.Fatalf("want one git ls-remote, ran %v", f.calls)
	}
	return rm
}

func tagNames(rm Remote) []string {
	var names []string
	for _, tag := range rm.Tags {
		names = append(names, tag.Name)
	}
	return names
}

func TestListRemoteReadsDefaultBranchAndReleaseTags(t *testing.T) {
	rm := listFake(t, lsRemote)
	if rm.DefaultBranch != "main" {
		t.Fatalf("default branch %q, want main", rm.DefaultBranch)
	}
	if want := []string{"main", "next"}; !slices.Equal(rm.Branches, want) {
		t.Fatalf("branches %v, want %v", rm.Branches, want)
	}
	if want := []string{"v1.0.9", "v1.0.10", "1.0.18", "v1.0.60"}; !slices.Equal(tagNames(rm), want) {
		t.Fatalf("tags %v, want %v (numeric order, peeled lines folded, pre-release and non-version skipped)", tagNames(rm), want)
	}
}

func TestLatestIsTheHighestVersion(t *testing.T) {
	rm := listFake(t, "ref: refs/heads/main\tHEAD\n"+
		"a\trefs/tags/v1.0.9\n"+
		"b\trefs/tags/v1.0.10\n"+
		"c\trefs/tags/1.0.18\n"+
		"d\trefs/tags/v2.0.0-rc.1\n"+
		"e\trefs/tags/latest\n")
	tag, ok := rm.Latest()
	if !ok || tag.Name != "1.0.18" || tag.Version != "1.0.18" {
		t.Fatalf("latest %+v %v, want 1.0.18", tag, ok)
	}
	tag, ok = listFake(t, lsRemote).Latest()
	if !ok || tag.Name != "v1.0.60" || tag.Version != "1.0.60" {
		t.Fatalf("latest %+v %v, want v1.0.60", tag, ok)
	}
}

func TestLatestReportsNoneWithoutReleaseTags(t *testing.T) {
	rm := listFake(t, "ref: refs/heads/main\tHEAD\na\tHEAD\na\trefs/heads/main\nb\trefs/tags/v1.1.0-rc.1\n")
	if tag, ok := rm.Latest(); ok {
		t.Fatalf("latest %+v, want none", tag)
	}
}

func TestLatestPrefersTheVTagOfAVersionTaggedTwice(t *testing.T) {
	rm := listFake(t, "a\trefs/tags/v1.0.18\nb\trefs/tags/1.0.18\n")
	if tag, _ := rm.Latest(); tag.Name != "v1.0.18" {
		t.Fatalf("latest %q, want v1.0.18", tag.Name)
	}
}

func TestMatchAndTagFor(t *testing.T) {
	rm := listFake(t, lsRemote)
	cases := []struct {
		ref, want string
	}{
		{"1.0.60", "v1.0.60"},
		{"v1.0.60", "v1.0.60"},
		{"1.0.18", "1.0.18"},
		{"v1.0.18", "1.0.18"},
		{"1.0.61", ""},
		{"main", ""},
		{"v1.1.0-rc.1", ""},
	}
	for _, c := range cases {
		tag, ok := rm.Match(c.ref)
		if ok != (c.want != "") || tag.Name != c.want {
			t.Errorf("Match(%q) = %q %v, want %q", c.ref, tag.Name, ok, c.want)
		}
	}
	if tag, ok := rm.TagFor("1.0.10"); !ok || tag.Name != "v1.0.10" {
		t.Fatalf("TagFor(1.0.10) = %+v %v", tag, ok)
	}
}

func TestMatchLeavesABranchNamedLikeAVersion(t *testing.T) {
	rm := listFake(t, "a\trefs/heads/1.0.60\nb\trefs/tags/v1.0.60\n")
	if tag, ok := rm.Match("1.0.60"); ok {
		t.Fatalf("Match(1.0.60) = %q, want none: a branch has that name", tag.Name)
	}
}

func TestFollowsReleases(t *testing.T) {
	rm := listFake(t, lsRemote)
	for _, ref := range []string{"", "main", "v1.0.18", "1.0.18", "1.0.60"} {
		if !rm.FollowsReleases(ref) {
			t.Errorf("FollowsReleases(%q) = false, want true", ref)
		}
	}
	for _, ref := range []string{"next", "0123456789abcdef0123456789abcdef01234567", "v1.1.0-rc.1"} {
		if rm.FollowsReleases(ref) {
			t.Errorf("FollowsReleases(%q) = true, want false", ref)
		}
	}
}

func TestListRemoteOfALocalDirectoryIsEmpty(t *testing.T) {
	f := &fakeGit{}
	rm, err := ListRemote(f, mini)
	if err != nil || len(f.calls) != 0 {
		t.Fatalf("local: %v, ran %v", err, f.calls)
	}
	if _, ok := rm.Latest(); ok {
		t.Fatal("a local template has no release tags")
	}
}

func TestListRemoteErrorNamesTheRepo(t *testing.T) {
	f := &fakeGit{fail: map[string]error{"ls-remote": errors.New("could not resolve host")}}
	_, err := ListRemote(f, "https://example.com/t.git")
	if err == nil || !strings.Contains(err.Error(), "https://example.com/t.git") || !strings.Contains(err.Error(), "could not resolve host") {
		t.Fatalf("error %v should name the repo and the cause", err)
	}
}
