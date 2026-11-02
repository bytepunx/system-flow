package cmd

import (
	"strings"
	"testing"
)

// servedPublishProject is a git project with a remote and S-0001, a feature
// against cli, accepted and not yet released, so that flai release --pending
// would tag cli/v1.1.0 and push it (S-0423).
func servedPublishProject(t *testing.T) (root, remote string) {
	t.Helper()
	root, remote = researchProject(t, "feature", true)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept S-0001: %s", errOut)
	}
	return root, remote
}

// S-0423: an agent flai serve started, the publish formula's among them
// (ADR-0240), is refused flai release --pending while the project's push host
// action is off, by either variable that says flai serve started it, naming
// flai serve enable push and exiting 4, with nothing tagged or pushed; its
// dry run, which publishes nothing, is not held. With --json the refusal
// still reaches stderr, as the formula's step runs it.
func TestReleasePendingHeldForAServedAgentWhilePushIsOff(t *testing.T) {
	for _, env := range []struct{ name, value string }{{"FLAI_STARTED_BY", "flai-serve"}, {"FLAI_STORY", "S-0009"}} {
		t.Run(env.name, func(t *testing.T) {
			root, remote := servedPublishProject(t)
			before := strings.TrimSpace(gitIn(t, remote, "rev-parse", "main"))
			t.Setenv(env.name, env.value)

			dry, errOut, code := runIn(t, root, "release", "--pending", "--dry-run")
			if code != 0 || !strings.Contains(dry, "cli") || !strings.Contains(dry, "S-0001") || !strings.Contains(dry, "dry run: nothing changed") {
				t.Fatalf("a dry run publishes nothing, so it is not held: %d %s %s", code, dry, errOut)
			}

			for _, args := range [][]string{{"release", "--pending"}, {"release", "--pending", "--json"}} {
				out, errOut, code := runIn(t, root, args...)
				if code != exitReleasePushOff || !strings.Contains(errOut, "the push host action is off for "+root) || !strings.Contains(errOut, "flai serve enable push") || !strings.Contains(errOut, "the operator's to turn on") {
					t.Fatalf("flai %v by a served agent with push off is refused, exit 4, naming flai serve enable push: %d %s %s", args, code, out, errOut)
				}
				if strings.Contains(out, "pushed") {
					t.Errorf("flai %v printed no result: %s", args, out)
				}
			}
			if tags := gitIn(t, root, "tag", "--list", "cli/*"); strings.Contains(tags, "cli/v1.1.0") {
				t.Errorf("nothing tagged here: %s", tags)
			}
			if tags := gitIn(t, remote, "tag", "--list", "cli/*"); strings.Contains(tags, "cli/v1.1.0") {
				t.Errorf("nothing tagged on the remote: %s", tags)
			}
			if after := strings.TrimSpace(gitIn(t, remote, "rev-parse", "main")); after != before {
				t.Errorf("nothing pushed: %s, was %s", after, before)
			}
			if status := gitIn(t, root, "status", "--porcelain"); strings.TrimSpace(status) != "" {
				t.Errorf("nothing applied: %s", status)
			}
		})
	}
}

// S-0423: with the push host action on, an agent flai serve started publishes
// as the operator does.
func TestReleasePendingForAServedAgentWithPushOn(t *testing.T) {
	root, remote := servedPublishProject(t)
	if out, errOut, code := runIn(t, root, "serve", "enable", "push"); code != 0 {
		t.Fatalf("enable push: %d %s %s", code, out, errOut)
	}
	t.Setenv("FLAI_STARTED_BY", "flai-serve")
	t.Setenv("FLAI_STORY", "S-0009")
	out, errOut, code := runIn(t, root, "release", "--pending")
	if code != 0 || !strings.Contains(out, "cli 1.1.0") || !strings.Contains(out, "pushed to origin") {
		t.Fatalf("a served agent with push on publishes: %d %s %s", code, out, errOut)
	}
	if tags := gitIn(t, remote, "tag", "--list", "cli/*"); !strings.Contains(tags, "cli/v1.1.0") {
		t.Errorf("the tag reaches the remote: %s", tags)
	}
}

// S-0423: the operator's own shell, with neither variable set (TestMain
// clears both), is not held while the push host action is off.
func TestReleasePendingForTheOperatorWithPushOff(t *testing.T) {
	t.Parallel()
	root, remote := servedPublishProject(t)
	out, errOut, code := runIn(t, root, "release", "--pending")
	if code != 0 || !strings.Contains(out, "cli 1.1.0") || !strings.Contains(out, "pushed to origin") {
		t.Fatalf("the operator publishes with push off: %d %s %s", code, out, errOut)
	}
	if tags := gitIn(t, remote, "tag", "--list", "cli/*"); !strings.Contains(tags, "cli/v1.1.0") {
		t.Errorf("the tag reaches the remote: %s", tags)
	}
}
