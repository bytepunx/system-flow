---
id: T-2881
type: task
nature: feature
title: The close-out records the namespace's last deploy and checks run with their commits and outcomes, and flai verify --last prints them
status: done
parent: S-0662
owner: alex
created: 2026-11-10T08:27:22Z
updated: 2026-11-13T14:27:26Z
transitions:
  - to: ready
    at: 2026-11-13T13:56:31Z
    by: agent-S-0662
  - to: in-progress
    at: 2026-11-13T13:56:32Z
    by: agent-S-0662
  - to: done
    at: 2026-11-13T14:27:26Z
    by: agent-S-0662
stream: S-0662
tags: [cli]
touches: [flai/internal/verify/story.go, flai/internal/verify/kluster.go, flai/internal/verify/kluster_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go]
after: [T-2875]
usage:
  source: log
  seconds: 1841
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 114
      output: 600
      cache_read: 6207448
      cache_write: 683557
      cost: 2.5759
---
# T-2881 The close-out records the namespace's last deploy and checks run with their commits and outcomes, and flai verify --last prints them

## Work

Give the verify report a `kluster` part, filled in `flai/internal/verify/kluster.go` from T-2875's evidence when `Verify` saves the report (`SaveReport` in `story.go`): the last deploy and the last checks run, each with its commit, outcome, who asked, when it ended, and its log paths, and whether each passed at the head verified. It is absent while the action is off or no `deploy` is set. The close-out, `scripts/close-out.sh`, runs `flai verify --record-issues`, so it records it with no change to the script.

The part comes from a verify step of its own, `kluster`, as ADR-0292 §6 decides: it runs nothing, passes when the action is off, the block has no `deploy`, or `kluster.agent_namespaces` is 0, otherwise reports the record's commits, outcomes, and whether each is at the head, and fails only on a deploy or checks run that failed at the head. Not deploying is no failure.

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
