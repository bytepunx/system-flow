---
id: S-0181
type: story
nature: improvement
title: The host's flai says when it is older than the tree it serves, and an older flai reads newer work items
status: ready
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:31:46Z
transitions:
  - to: ready
    at: 2026-10-01T08:31:46Z
    by: alex
tags: [flai]
touches: [flai/internal/workitem, flai/internal/threads, flai/internal/issues, flai/internal/manifest, flai/internal/serve, flai/internal/mcpserver, flai/internal/check, flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0181 The host's flai says when it is older than the tree it serves, and an older flai reads newer work items

## Goal

The flai installed on the host serves MCP and starts agents, and it lags the tree it serves whenever a release is accepted and not installed. Two things go wrong:

- **It acts on rules it does not have.** I-0049: a `flai serve` that predated holds started three ready stories that open stories held. Nothing says the host's flai is older than the tree, or than the latest release; `flai host upgrade` exists but is manual.
- **It refuses the tree outright.** I-0051: work items, threads, and issues are decoded with `yaml.Strict()` (`workitem/item.go` ~124, `threads.go` ~85, `issues.go` ~60), so once any item carries a front-matter field an older flai does not know, that flai refuses every listing that reads it: board, inbox, `flai serve`, MCP. Every new field (`after`, `topics`, `usage`) has required "upgrade the host first" in its design.

## Acceptance criteria
- [ ] `flai serve`, the MCP `inbox`, and the dashboard's host badge say when the running flai is older than the newest `flai/v*` tag reachable from the project's HEAD, with the command that upgrades it
- [ ] `system-flow.yaml` can name a minimum flai version; a flai below it says so, naming the version needed, before reading any item, instead of failing on an unknown field. A release that adds a front-matter field raises it
- [ ] The paths that read work items, threads, and issues to serve or list them tolerate unknown front-matter fields with a warning naming the field and the item; `flai check` still reports unknown fields strictly, and no write drops a field it did not understand
- [ ] Tests cover an item with an unknown field read by the listing paths and by `flai check`, a manifest minimum above the running version, and the version warning
- [ ] The design (`design/system/flai-cli.md`, `design/system/work-hierarchy.md`, `design/system/project-manifest.md`) and the operator guide describe the version check and the minimum
- [ ] I-0049 and I-0051 are closed with what fixed them

## Tasks

## Notes

- Upgrading the host automatically after an acceptance that publishes a release was suggested in triage; it is out of scope unless the story finds it cheap and safe.
- I-0050 (a host flai that had holds still started a held story) is a separate investigation.
