---
id: T-2881
type: task
nature: feature
title: The close-out records the namespace's last deploy and checks run with their commits and outcomes, and flai verify --last prints them
status: backlog
parent: S-0662
owner: alex
created: 2026-11-10T08:27:22Z
updated: 2026-11-10T08:28:09Z
transitions: []
stream: S-0662
tags: [cli]
touches: [flai/internal/verify/story.go, flai/internal/verify/kluster.go, flai/internal/verify/kluster_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go]
after: [T-2875]
---
# T-2881 The close-out records the namespace's last deploy and checks run with their commits and outcomes, and flai verify --last prints them

## Work

Give the verify report a `kluster` part, filled in `flai/internal/verify/kluster.go` from T-2875's evidence when `Verify` saves the report (`SaveReport` in `story.go`): the last deploy and the last checks run, each with its commit, outcome, who asked, when it ended, and its log paths, and whether each passed at the head verified. It is absent while the action is off or no `deploy` is set. The close-out, `scripts/close-out.sh`, runs `flai verify --record-issues`, so it records it with no change to the script.

`printLastVerify` in `flai/cmd/verify.go` prints it after the verdict, and `--json` carries it.

A story whose agent runs in a container verifies there, where the main checkout's `.flai-cache/kluster` is not mounted (ADR-0274). Read the records from where the close-out can reach them, or fill the part on the host; say which in the story's narrative.

It waits for T-2875, whose evidence it reads.

## Done when

- `kluster_test.go` covers the part filled, absent, and at an older commit.
- `verify_test.go` shows `flai verify --last` and `--last --json` with the deploy and checks run.
- A report saved before this change still reads.
- `flai test` passes on the five files.

## Notes

- Criterion 3 of S-0662, its record and `--last`; T-2885's `item_get` shows it.
