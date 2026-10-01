---
id: S-0181
type: story
nature: improvement
title: The host's flai says when it is older than the tree it serves, and an older flai reads newer work items
status: done
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T11:19:10Z
transitions:
  - to: ready
    at: 2026-10-01T08:31:46Z
    by: alex
  - to: in-progress
    at: 2026-10-01T10:42:32Z
    by: agent-S-0181
  - to: review
    at: 2026-10-01T11:05:17Z
    by: agent-S-0181
  - to: done
    at: 2026-10-01T11:19:10Z
    by: alex
tags: [flai]
touches: [flai/internal/workitem, flai/internal/threads, flai/internal/issues, flai/internal/manifest, flai/internal/serve, flai/internal/mcpserver, flai/internal/check, flai/internal/hostapi, flai/internal/release, flai/internal/buildinfo, flai/cmd, flaiover/src, design/system/flai-cli.md, design/system/work-hierarchy.md, design/system/project-manifest.md, design/issues, docs/operators, docs/users/flai.md, system-flow.yaml]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1393
  models:
    - model: claude-opus-5-5
      input: 452
      output: 128220
      cache_read: 31423861
      cache_write: 536153
      cost: 12.3815
---
# S-0181 The host's flai says when it is older than the tree it serves, and an older flai reads newer work items

## Goal

The flai installed on the host serves MCP and starts agents, and it lags the tree it serves whenever a release is accepted and not installed. Two things go wrong:

- **It acts on rules it does not have.** I-0049: a `flai serve` that predated holds started three ready stories that open stories held. Nothing says the host's flai is older than the tree, or than the latest release; `flai host upgrade` exists but is manual.
- **It refuses the tree outright.** I-0051: work items, threads, and issues are decoded with `yaml.Strict()` (`workitem/item.go` ~124, `threads.go` ~85, `issues.go` ~60), so once any item carries a front-matter field an older flai does not know, that flai refuses every listing that reads it: board, inbox, `flai serve`, MCP. Every new field (`after`, `topics`, `usage`) has required "upgrade the host first" in its design.

## Acceptance criteria
- [x] `flai serve`, the MCP `inbox`, and the dashboard's host badge say when the running flai is older than the newest `flai/v*` tag reachable from the project's HEAD, with the command that upgrades it
- [x] `system-flow.yaml` can name a minimum flai version; a flai below it says so, naming the version needed, before reading any item, instead of failing on an unknown field. A release that adds a front-matter field raises it
- [x] The paths that read work items, threads, and issues to serve or list them tolerate unknown front-matter fields with a warning naming the field and the item; `flai check` still reports unknown fields strictly, and no write drops a field it did not understand
- [x] Tests cover an item with an unknown field read by the listing paths and by `flai check`, a manifest minimum above the running version, and the version warning
- [x] The design (`design/system/flai-cli.md`, `design/system/work-hierarchy.md`, `design/system/project-manifest.md`) and the operator guide describe the version check and the minimum
- [x] I-0049 and I-0051 are closed with what fixed them

## Tasks
- T-0667 Listing paths read work items, threads, and issues with unknown front-matter fields, warn, and keep them on write; flai check reports them
- T-0668 system-flow.yaml names a minimum flai version, and a flai below it says so before reading any item
- T-0669 flai serve, the MCP inbox, and the dashboard's host badge say when the running flai is older than the newest flai tag reachable from HEAD
- T-0670 Design, operator guide, and release rule describe the version check and the minimum; I-0049 and I-0051 closed

## Notes

- Upgrading the host automatically after an acceptance that publishes a release was suggested in triage; it is out of scope unless the story finds it cheap and safe.
- I-0050 (a host flai that had holds still started a held story) is a separate investigation.
