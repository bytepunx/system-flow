---
id: T-0841
type: task
nature: improvement
title: flai check --story --record-issues opens or bumps one issue per rule for the findings outside the story
status: done
parent: S-0249
owner: alex
created: 2026-10-05T00:18:03Z
updated: 2026-10-05T00:51:43Z
transitions:
  - to: ready
    at: 2026-10-05T00:44:47Z
    by: agent-S-0249
  - to: in-progress
    at: 2026-10-05T00:44:47Z
    by: agent-S-0249
  - to: done
    at: 2026-10-05T00:51:43Z
    by: agent-S-0249
stream: S-0249
tags: [flai, issues]
touches: [flai/internal/issues, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-0840]
usage:
  source: log
  seconds: 416
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 69
      output: 22849
      cache_read: 2980345
      cache_write: 89772
      cost: 1.5807
---
# T-0841 flai check --story --record-issues opens or bumps one issue per rule for the findings outside the story

## Work

TH-0056's answer asks flai to record each finding outside the story as an issue, opening one when none exists and bumping its count otherwise, so the operator sees where these findings occur.

Add `--record-issues` to `flai check`. It is valid only with `--story`. For each rule with findings outside the story, find the open issue that records that rule. If there is none, open one with `issues.New`, titled for instance "flai check finds `story.unaccepted` outside the story at close-out", with class `efficiency`. Then bump it with `issues.Bump`, giving the story and an instance that names each finding's path and message.

Add a lookup to `flai/internal/issues` that finds an open issue by its title, or by a key kept in its front matter, so that a later run finds the issue again. An instance for the same story and the same findings is written once, so that running the close-out again does not inflate the count.

The issue files are written in the checkout `flai check` runs in, which is the story's worktree at close-out. The close-out then commits them with the story, as `flai issue new` would. Print each issue opened or bumped.

It waits for T-0840, whose outside marker it reads, and because both change `flai/cmd/check.go`.

## Done when

- A test in `flai/cmd/check_test.go` covers three runs. A first run of `flai check --strict --story S-0001 --record-issues`, with a finding outside the story, opens an issue with count 1 and an instance naming S-0001 and the finding. A second identical run leaves the count at 1. A run for S-0003 with the same finding bumps it to 2.
- `--record-issues` without `--story` is refused with a message that says why.
- A check with no findings outside the story writes nothing under `design/issues`.
- `docs/users/flai-reference.md` matches the help text, and `scripts/flai-test.sh` passes.

## Notes
