---
id: T-0859
type: task
nature: remediation
title: flai accept refuses a story branch that adds a conflict marker
status: done
parent: S-0253
owner: alex
created: 2026-10-05T03:13:24Z
updated: 2026-10-05T03:20:32Z
transitions:
  - to: ready
    at: 2026-10-05T03:13:37Z
    by: agent-S-0253
  - to: in-progress
    at: 2026-10-05T03:16:51Z
    by: agent-S-0253
  - to: done
    at: 2026-10-05T03:20:32Z
    by: agent-S-0253
stream: S-0253
tags: []
touches: [flai/cmd/branch.go, flai/cmd/accept_conflict_test.go, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md]
after: [T-0858]
usage:
  source: log
  seconds: 221
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 60
      output: 16666
      cache_read: 2092747
      cache_write: 72915
      cost: 1.2197
---
# T-0859 flai accept refuses a story branch that adds a conflict marker

## Work

Acceptance is the gate onto main that I-0066 names. In `mergeStoryBranch` (`flai/cmd/branch.go`), after the sync and before the fast-forward, read each file the branch adds or changes against the main branch (`git diff --name-only --diff-filter=AMR <base>...<branch>`, then `git show <branch>:<path>`), skip binary content, and scan it with `conflictmark`. When any line is a marker, refuse the acceptance, merging nothing, with an error naming each `path:line` and saying to resolve the conflict in the story's worktree, commit it, and accept again. Every file, not only markdown, as a backstop: a marker in code or YAML usually fails elsewhere, but nothing should reach main with one. The accept, the operator's move to done, and the dashboard share this path. A test in `flai/cmd` reproduces the cause: a story branch whose commit carries the markers is refused, and main is unchanged. Describe the refusal in `design/system/flai-cli.md`'s `flai accept` row, the acceptance in `design/system/workflow.md`, and `docs/users/flai.md`.

Waits for T-0858: it uses `conflictmark`, and both edit `design/system/flai-cli.md` and `docs/users/flai.md`.

## Done when

- `go test ./cmd/ -run Accept` passes, with a test that fails without the refusal.
- The design and the user guide describe the refusal.

## Notes
