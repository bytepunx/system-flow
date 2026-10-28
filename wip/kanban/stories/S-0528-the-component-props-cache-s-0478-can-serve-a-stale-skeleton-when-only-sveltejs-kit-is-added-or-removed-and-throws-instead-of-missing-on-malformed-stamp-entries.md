---
id: S-0528
type: story
nature: remediation
title: The component props cache (S-0478) can serve a stale skeleton when only @sveltejs/kit is added or removed, and throws instead of missing on malformed stamp entries
status: backlog
owner: alex
created: 2026-10-28T12:54:18Z
updated: 2026-10-28T12:57:41Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
finalized:
  by: alex
  at: 2026-10-28T12:57:41Z
---
# S-0528 The component props cache (S-0478) can serve a stale skeleton when only @sveltejs/kit is added or removed, and throws instead of missing on malformed stamp entries

## Goal

This story remediates [I-0198](../../../design/issues/I-0198-the-component-props-cache-s-0478-can-serve-a-stale-skeleton-when-only-sveltejs-kit-is-added-or-removed-and-throws-instead-of-missing-on-malformed-stamp-entries.md), "The component props cache (S-0478) can serve a stale skeleton when only @sveltejs/kit is added or removed, and throws instead of missing on malformed stamp entries". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0198 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0198 is closed with `flai issue close I-0198 --reason` saying what fixed it

## Tasks

## Notes
