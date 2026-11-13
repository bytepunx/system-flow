---
id: T-3233
type: task
nature: feature
title: A conversation between two stories closes when either reaches an accepted or cancelled state, and a share lasts while both are open
status: backlog
parent: S-0749
owner: alex
created: 2026-11-12T16:08:23Z
updated: 2026-11-12T16:08:23Z
transitions: []
stream: S-0749
tags: [cli]
touches: [flai/internal/messages/messages.go, flai/internal/messages/share.go, flai/internal/messages/messages_test.go, flai/internal/messages/share_test.go]
---
# T-3233 A conversation between two stories closes when either reaches an accepted or cancelled state, and a share lasts while both are open

## Work

Replace the 14 references to a state's name in `messages.go` and `share.go` with reads of kinds: a conversation closes when either story reaches an accepted or a cancelled state (ADR-0120), and a share holds while both stories are in active-work or wait states. Test a close and a share on a workflow with a state added.

## Done when

Both files name no state, and their tests pass on both workflows.

## Notes
