---
id: S-0174
type: story
nature: remediation
title: Pending releases are not worked out from a clone missing the remote's release tags
status: done
owner: alex
created: 2026-10-01T06:55:51Z
updated: 2026-10-01T08:18:49Z
transitions:
  - to: ready
    at: 2026-10-01T07:39:51Z
    by: alex
  - to: in-progress
    at: 2026-10-01T07:48:30Z
    by: agent-S-0174
  - to: review
    at: 2026-10-01T08:10:35Z
    by: agent-S-0174
  - to: done
    at: 2026-10-01T08:18:49Z
    by: alex
tags: [flai, dashboard]
topics: [release]
touches: [flai/internal/release, flai/cmd/release.go, flai/internal/hostapi, flai/internal/check, flaiover/src, flai/cmd/hostapi_reads_test.go, flai/cmd/push_test.go, flai/cmd/release_pending_test.go, flai/internal/pending/pending.go, flai/internal/preview/push.go, flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users, design/issues]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1361
  models:
    - model: claude-opus-5-5
      input: 286
      output: 95435
      cache_read: 25254557
      cache_write: 278174
      cost: 9.1861
---
# S-0174 Pending releases are not worked out from a clone missing the remote's release tags

## Goal

flai works out what is waiting to be published from local tags only (`release.CurrentVersion` runs `git tag --list <name>/v*`), and flai never fetches. When a clone's `main` is fast-forwarded from `origin/main` without its tags, as happens after publishing from another clone, the dashboard's done lane and publish banner show every story accepted since the last local tag as unpublished, and propose versions computed from that stale tag (seen 2026-10-01: local `flai/v1.18.0` and `flaiover/v0.27.8` against origin's `flai/v1.26.3` and `flaiover/v0.32.1`, with flai offering `1.18.0 -> 1.19.0` across 33+ stories). Publishing in that state would apply and tag versions that already exist on the remote before the push fails. A clone missing the remote's release tags should be detected and said, never offered as a pending release.

Added 2026-10-01 by the designer (I-0024): the pending plan also drops work silently. An item whose plan cannot be recomputed, such as a story with no component tag that touches two components ("no tag says which it delivers to", `release.go` ~L278), is skipped with `continue` (~L550-566), so its work can miss its bump, and nothing warns before then.

## Acceptance criteria
- [x] `publish.preview` and `flai release --pending` detect when origin has a `<name>/v*` tag newer than the highest local one for a component (`git ls-remote --tags`), and report it with what fixes it (`git fetch --tags origin`) instead of a plan built on the stale tag
- [x] `publish.run` and `flai release --pending` refuse to apply, tag, or push while origin has a newer release tag than the clone, before changing anything locally
- [x] When origin cannot be reached, the pending plan is still shown, with a warning that it could not be checked against the remote
- [x] The dashboard's done lane and publish banner say the clone is missing published tags and how to fetch them, rather than listing published stories as waiting
- [x] Tests cover a clone whose tags lag the remote, one in step with it, and an unreachable remote
- [x] (I-0024) The pending plan lists every accepted item it could not plan, with the reason, in `flai release --pending`, `publish.preview`, and the dashboard's publish banner, instead of skipping it silently
- [x] (I-0024) `flai check` warns on a story with no component tag, its own or its epic's, whose touches reach two or more components, saying which tag to add
- [x] The design (`design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`) and the user guide describe the check
- [x] I-0024 is closed with what fixed it

## Tasks
- T-0615 flai detects a clone whose release tags lag the remote and says so in the pending plan
- T-0616 flai release --pending and publish.run refuse to publish while the remote has newer release tags
- T-0617 The dashboard says the clone is missing published tags instead of listing published stories as waiting
- T-0618 Design and user guide describe the remote release-tag check
- T-0625 The pending plan lists every accepted item it could not plan, with the reason
- T-0626 flai check warns on a story with no component tag whose touches reach two components

## Notes

Alternative considered: have publish fetch tags itself before planning. Detecting and saying keeps flai from changing refs the operator did not ask for; revisit if the check proves noisy.

The I-0024 criteria were added while the story was in progress: write a task for them.

Decisions made in the story (narrative `wip/agents/S-0174.md` has the reasons): while any component lags, the preview offers no plan for any component, since publishing is one run; a remote that cannot be asked refuses the publish too, before anything changes, because the push needs it; a clone with no remote is not checked; the remote's tags are kept a minute per root for the board's preview; the done lane hides archived cards while lagging in the dashboard rather than flai asking the network on every board read.
