---
id: T-0074
type: task
nature: feature
title: "flai upgrade command: dry-run, dirty-tree guard, interactive keep/replace/diff, keep-all, replace-all"
status: done
parent: S-0020
owner: alex
created: 2026-09-17T03:19:34Z
updated: 2026-09-17T03:25:22Z
transitions:
  - to: ready
    at: 2026-09-17T03:25:21Z
    by: agent
  - to: in-progress
    at: 2026-09-17T03:25:22Z
    by: agent
  - to: done
    at: 2026-09-17T03:25:22Z
    by: agent
stream: S-0020
tags: [cli, template]
---

# T-0074 flai upgrade command: dry-run, dirty-tree guard, interactive keep/replace/diff, keep-all, replace-all

## Work
cmd/upgrade.go: fetch the template (config or flags), compare versions (no-op unless --force), refuse a dirty git tree unless --force, --dry-run prints the plan, interactive conflicts offer keep, replace, or show diff (git diff --no-index), non-interactive needs --keep-all or --replace-all else exits 1 with the list and changes nothing.

## Done when
In-process tests on a project rendered from the fixture then upgraded from a modified copy.

## Notes
