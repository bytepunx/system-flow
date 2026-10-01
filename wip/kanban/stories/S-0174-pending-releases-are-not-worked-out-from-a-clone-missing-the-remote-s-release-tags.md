---
id: S-0174
type: story
nature: remediation
title: Pending releases are not worked out from a clone missing the remote's release tags
status: ready
owner: alex
created: 2026-10-01T06:55:51Z
updated: 2026-10-01T07:39:51Z
transitions:
  - to: ready
    at: 2026-10-01T07:39:51Z
    by: alex
tags: [flai, dashboard]
topics: [release]
touches: [flai/internal/release, flai/cmd/release.go, flai/internal/hostapi, flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0174 Pending releases are not worked out from a clone missing the remote's release tags

## Goal

flai works out what is waiting to be published from local tags only (`release.CurrentVersion` runs `git tag --list <name>/v*`), and flai never fetches. When a clone's `main` is fast-forwarded from `origin/main` without its tags, as happens after publishing from another clone, the dashboard's done lane and publish banner show every story accepted since the last local tag as unpublished, and propose versions computed from that stale tag (seen 2026-10-01: local `flai/v1.18.0` and `flaiover/v0.27.8` against origin's `flai/v1.26.3` and `flaiover/v0.32.1`, with flai offering `1.18.0 -> 1.19.0` across 33+ stories). Publishing in that state would apply and tag versions that already exist on the remote before the push fails. A clone missing the remote's release tags should be detected and said, never offered as a pending release.

## Acceptance criteria
- [ ] `publish.preview` and `flai release --pending` detect when origin has a `<name>/v*` tag newer than the highest local one for a component (`git ls-remote --tags`), and report it with what fixes it (`git fetch --tags origin`) instead of a plan built on the stale tag
- [ ] `publish.run` and `flai release --pending` refuse to apply, tag, or push while origin has a newer release tag than the clone, before changing anything locally
- [ ] When origin cannot be reached, the pending plan is still shown, with a warning that it could not be checked against the remote
- [ ] The dashboard's done lane and publish banner say the clone is missing published tags and how to fetch them, rather than listing published stories as waiting
- [ ] Tests cover a clone whose tags lag the remote, one in step with it, and an unreachable remote
- [ ] The design (`design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`) and the user guide describe the check

## Tasks

## Notes

Alternative considered: have publish fetch tags itself before planning. Detecting and saying keeps flai from changing refs the operator did not ask for; revisit if the check proves noisy.
