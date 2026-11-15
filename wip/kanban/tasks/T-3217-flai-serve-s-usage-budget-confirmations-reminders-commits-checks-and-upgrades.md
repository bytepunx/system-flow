---
id: T-3217
type: task
nature: feature
title: flai serve's usage, budget, confirmations, reminders, commits, checks, and upgrades read a state's kind
status: done
parent: S-0742
owner: alex
created: 2026-11-12T16:07:09Z
updated: 2026-11-15T01:55:14Z
transitions:
  - to: ready
    at: 2026-11-15T01:27:04Z
    by: agent-S-0742
  - to: in-progress
    at: 2026-11-15T01:27:04Z
    by: agent-S-0742
  - to: done
    at: 2026-11-15T01:55:14Z
    by: agent-S-0742
stream: S-0742
tags: [cli]
touches: [flai/internal/serve/usage.go, flai/internal/serve/confirm.go, flai/internal/serve/remind.go, flai/internal/serve/remind_test.go, flai/internal/serve/commit.go, flai/internal/serve/claudecheck.go, flai/internal/serve/components.go, flai/internal/serve/imagejob.go]
after: [T-3370]
usage:
  source: log
  seconds: 1690
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 103
      output: 479
      cache_read: 5372548
      cache_write: 237026
      cost: 2.1147
---
# T-3217 flai serve's usage, budget, confirmations, reminders, commits, checks, and upgrades read a state's kind

## Work

Replace the one to three references to a state's name in each of these files with a read of the kind: usage is measured when a story reaches an accepted or cancelled state, a reminder and a confirmation skip a closed item, and a check, a verify job, and an upgrade find the stories in active-work and wait states. Test the reminder on a workflow with a state added.

## Done when

These files name no state, and the existing `serve` tests pass unchanged.

## Notes
