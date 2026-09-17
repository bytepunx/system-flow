---
id: T-0052
type: task
nature: feature
title: flai prime command
status: done
parent: S-0025
owner: alex
created: 2026-09-16T23:03:59Z
updated: 2026-09-16T23:07:55Z
transitions:
  - to: ready
    at: 2026-09-16T23:07:54Z
    by: agent
  - to: in-progress
    at: 2026-09-16T23:07:55Z
    by: agent
  - to: done
    at: 2026-09-16T23:07:55Z
    by: agent
stream: S-0025
tags: [cli]
---

# T-0052 flai prime command

## Work
cmd/prime.go: flai prime lists README then files by order as paths; --cat prints contents with headers; --json returns the file metadata; a missing folder is an error with the fix named.

## Done when
Tests for paths, cat order, json, and missing folder.

## Notes
