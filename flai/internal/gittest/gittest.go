// Package gittest gives a test the git configuration it makes and none of the
// host's, so that what git does in a test does not depend on whose machine
// runs it (I-0045).
package gittest

import (
	"os"
	"path/filepath"
	"testing"
)

// identityVars are the variables git takes an identity from before its
// configuration.
var identityVars = []string{"GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL", "GIT_COMMITTER_NAME", "GIT_COMMITTER_EMAIL", "EMAIL"}

// Isolate makes git, for the rest of t, read no system configuration, read a
// file of t's own as its global configuration, and take no identity from the
// environment. It returns the path of that file, which starts empty.
func Isolate(t testing.TB) string {
	t.Helper()
	global := filepath.Join(t.TempDir(), "gitconfig")
	if err := os.WriteFile(global, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", global)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	for _, k := range identityVars {
		t.Setenv(k, "x") // registers the restore; an empty value would still override git config
		if err := os.Unsetenv(k); err != nil {
			t.Fatal(err)
		}
	}
	return global
}

// Identity isolates git for the rest of t and gives it name and email to
// author and commit with.
func Identity(t testing.TB, name, email string) {
	t.Helper()
	Isolate(t)
	t.Setenv("GIT_AUTHOR_NAME", name)
	t.Setenv("GIT_AUTHOR_EMAIL", email)
	t.Setenv("GIT_COMMITTER_NAME", name)
	t.Setenv("GIT_COMMITTER_EMAIL", email)
}

// NoIdentity isolates git for the rest of t and leaves it without an
// identity: git refuses to commit, instead of working one out from the
// user's account and the host's name, as it does on macOS.
func NoIdentity(t testing.TB) {
	t.Helper()
	global := Isolate(t)
	if err := os.WriteFile(global, []byte("[user]\n\tuseConfigOnly = true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}
