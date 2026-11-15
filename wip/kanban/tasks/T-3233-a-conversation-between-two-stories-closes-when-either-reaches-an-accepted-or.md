---
id: T-3233
type: task
nature: feature
title: A conversation between two stories closes when either reaches an accepted or cancelled state, and a share lasts while both are open
status: done
parent: S-0749
owner: alex
created: 2026-11-12T16:08:23Z
updated: 2026-11-15T19:57:30Z
transitions:
  - to: ready
    at: 2026-11-15T19:45:26Z
    by: agent-S-0749
  - to: in-progress
    at: 2026-11-15T19:45:26Z
    by: agent-S-0749
  - to: done
    at: 2026-11-15T19:57:30Z
    by: agent-S-0749
stream: S-0749
tags: [cli]
touches: [flai/internal/messages/messages.go, flai/internal/messages/share.go, flai/internal/messages/messages_test.go, flai/internal/messages/share_test.go]
usage:
  source: log
  seconds: 707
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 52
      output: 247
      cache_read: 1701037
      cache_write: 87548
      cost: 0.6758
---
# T-3233 A conversation between two stories closes when either reaches an accepted or cancelled state, and a share lasts while both are open

## Work

Replace the 14 references to a state's name in `messages.go` and `share.go` with reads of kinds: a conversation closes when either story reaches an accepted or a cancelled state (ADR-0120), and a share holds while both stories are in active-work or wait states. Test a close and a share on a workflow with a state added.

## Done when

Both files name no state, and their tests pass on both workflows.

## Notes
