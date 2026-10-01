---
id: TH-0043
title: "S-0180's \"pass in CI\": open a pull request to run the flai workflow, or count the Docker run?"
anchor:
  path: wip/kanban/stories/S-0180-flai-s-tests-and-install-test-pass-on-a-macos-host.md
  item: S-0180
status: resolved
participants: [agent-S-0180, alex]
created: 2026-10-01T08:26:52Z
updated: 2026-10-01T08:30:01Z
---

# TH-0043 S-0180's "pass in CI": open a pull request to run the flai workflow, or count the Docker run?

On wip/kanban/stories/S-0180-flai-s-tests-and-install-test-pass-on-a-macos-host.md.

## Entries

### 2026-10-01T08:26:52Z agent-S-0180
All of S-0180's work is committed on story/S-0180 (worktree clean), and four of the five criteria are ticked. The fourth says `scripts/flai-test.sh` and `scripts/install-test.sh` pass on a macOS host and in CI.

- On this Mac: flai-test.sh passes gofmt, vet, golangci-lint, behavior, integration, and the template smoke. Its repository check (`flai check --strict`) stops on six warnings about the main checkout's board that this story doesn't touch: E-0003, E-0010, E-0011, and E-0012 are done but not archived, and TH-0026 and TH-0032 are on archived stories. main's system-flow-check CI has been red for a while too. The markdown lint and install-test.sh pass on their own.
- On Linux: `go test -race -count=1 ./...` passes in `golang:1.26` in Docker, the same command as the CI test job.
- In GitHub Actions: nothing from this branch has run. flai.yml runs only on a pull request or a push to main, and the branch isn't pushed.

Which do you want?
(a) I push story/S-0180 and open a draft pull request so flai.yml runs on ubuntu-latest and the new macos-latest job, and tick the criterion when both are green. This is my recommendation: it is the only way to see the new macOS job run before acceptance.
(b) Count this Mac plus the Docker Linux run, and let CI confirm on the push to main after acceptance.

### 2026-10-01T08:29:47Z alex
b

### 2026-10-01T08:30:01Z agent-S-0180
Resolved: Option b: this Mac plus the Docker Linux run count for "in CI"; CI confirms on the push to main after acceptance.
