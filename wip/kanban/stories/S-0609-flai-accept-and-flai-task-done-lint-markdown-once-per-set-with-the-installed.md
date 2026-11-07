---
id: S-0609
type: story
nature: improvement
title: flai accept and flai task done lint markdown once per set with the installed markdownlint, not five times through npx, saving about 4 s an acceptance
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:41:56Z
updated: 2026-11-07T22:41:56Z
transitions: []
tags: [cli]
touches: [flai/cmd/accept_markdown.go, flai/cmd/accept.go, scripts/lint-md.sh]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0609 flai accept and flai task done lint markdown once per set with the installed markdownlint, not five times through npx, saving about 4 s an acceptance

## Goal

Cause 19 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `flai accept` runs `scripts/lint-md.sh` five times, one after another, each through `npx --yes markdownlint-cli2` (`flai/cmd/accept_markdown.go`), and checks again after each `--fix`. Each call takes 1.0 to 1.2 s, of which markdownlint itself is 0.27 to 0.36 s. `flai task done` runs it through `npx` twice. `accept` also starts 58 git processes, `git status --porcelain` six times, and lists every item for `EpicWalk`. The serve log has `accept.run` at p50 16.6 s, up to 42 s.

## Acceptance criteria

- [ ] `flai accept` lints the uncommitted set and the branch's set in one call each, runs the installed `markdownlint-cli2` without `npx` where it is installed, and runs no check after a `--fix` whose exit already says what remains.
- [ ] `flai accept` reads `git status --porcelain` once per state of the checkout it needs, not six times.
- [ ] An acceptance in a scratch copy of this repository takes at least 4 s less than before, refusing the same unfixed markdown as ADR-0140 says, and server-performance.md records the before and after.

## Tasks

## Notes

- `releaseDue`, `tellAccepted`, and `installAccepted`, the acceptance's last steps, are independent and could run at once; the saving is small but larger when an install runs.
