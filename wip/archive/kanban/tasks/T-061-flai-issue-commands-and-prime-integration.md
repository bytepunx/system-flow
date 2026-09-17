---
id: T-061
type: task
nature: feature
title: flai issue commands and prime integration
status: done
parent: S-027
owner: alex
created: 2026-09-16T23:54:25Z
updated: 2026-09-16T23:59:09Z
transitions:
  - to: ready
    at: 2026-09-16T23:59:08Z
    by: agent
  - to: in-progress
    at: 2026-09-16T23:59:08Z
    by: agent
  - to: done
    at: 2026-09-16T23:59:09Z
    by: agent
stream: S-027
tags: [cli, issues]
---

# T-061 flai issue commands and prime integration

## Work
cmd/issue.go: new, bump, close, list, summary, each regenerating summary.md; prime lists summary.md after the conventions and --cat appends the table when anything is open.

## Done when
In-process tests cover the lifecycle through the CLI and the prime integration.

## Notes
