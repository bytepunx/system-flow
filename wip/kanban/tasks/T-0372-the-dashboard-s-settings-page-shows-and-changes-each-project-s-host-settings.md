---
id: T-0372
type: task
nature: feature
title: The dashboard's settings page shows and changes each project's host settings
status: done
parent: S-0105
owner: alex
created: 2026-09-23T20:36:55Z
updated: 2026-09-23T20:50:04Z
transitions:
  - to: ready
    at: 2026-09-23T20:44:05Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T20:44:05Z
    by: system-flow
  - to: done
    at: 2026-09-23T20:50:04Z
    by: system-flow
stream: S-0105
tags: []
---

# T-0372 The dashboard's settings page shows and changes each project's host settings

## Work
A Settings page per project, reached from the header. It shows every section read-only while `settings` is off, and says the shell command that turns it on. With it on, each section changes through flai and shows what flai answered: host actions, default agent, agent command and harnesses, checks, import folders, MCP server, and dashboard token. Host-wide sections say that they apply to every project.

## Done when
vitest covers the page read-only and writable, each section's request, and a refusal.

## Notes
