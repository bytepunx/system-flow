---
id: T-3217
type: task
nature: feature
title: flai serve's usage, budget, confirmations, reminders, commits, checks, and upgrades read a state's kind
status: backlog
parent: S-0742
owner: alex
created: 2026-11-12T16:07:09Z
updated: 2026-11-12T16:07:09Z
transitions: []
stream: S-0742
tags: [cli]
touches: [flai/internal/serve/usage.go, flai/internal/serve/budget.go, flai/internal/serve/confirm.go, flai/internal/serve/remind.go, flai/internal/serve/remind_test.go, flai/internal/serve/commit.go, flai/internal/serve/claudecheck.go, flai/internal/serve/components.go, flai/internal/serve/checks.go, flai/internal/serve/verifyjob.go, flai/internal/serve/upgrade.go]
---
# T-3217 flai serve's usage, budget, confirmations, reminders, commits, checks, and upgrades read a state's kind

## Work

Replace the one to three references to a state's name in each of these files with a read of the kind: usage is measured when a story reaches an accepted or cancelled state, a reminder and a confirmation skip a closed item, and a check, a verify job, and an upgrade find the stories in active-work and wait states. Test the reminder on a workflow with a state added.

## Done when

These files name no state, and the existing `serve` tests pass unchanged.

## Notes
