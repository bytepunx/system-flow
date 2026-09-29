---
id: T-0559
type: task
nature: feature
title: A dismissible notice component with a right-aligned X
status: done
parent: S-0151
owner: alex
created: 2026-09-29T19:11:39Z
updated: 2026-09-29T19:12:51Z
transitions:
  - to: ready
    at: 2026-09-29T19:11:55Z
    by: agent-S-0151
  - to: in-progress
    at: 2026-09-29T19:11:56Z
    by: agent-S-0151
  - to: done
    at: 2026-09-29T19:12:51Z
    by: agent-S-0151
stream: S-0151
tags: []
usage:
  source: log
  seconds: 55
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 14
      output: 3145
      cache_read: 650458
      cache_write: 15059
      cost: 0.3135
---

# T-0559 A dismissible notice component with a right-aligned X

## Work

Add `flaiover/src/lib/components/DismissibleNotice.svelte`: a banner that takes its classes, role, and test ID from the caller, renders its content, and ends with a right-aligned X button labelled Dismiss that calls `ondismiss`. Add a component test.

## Done when

The component test shows the X is the banner's last element, is labelled Dismiss, calls `ondismiss` when clicked, and leaves the text under the caller's test ID without the X.

## Notes
