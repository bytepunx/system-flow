---
id: T-3098
type: task
nature: remediation
title: "Clear a budget state with a reason: journaled under serve.budget and sent as a notice, and budgetHeldWhy names flai serve budget clear"
status: backlog
parent: S-0720
owner: alex
created: 2026-11-11T08:54:48Z
updated: 2026-11-11T08:54:48Z
transitions: []
stream: S-0720
tags: [flai, serve, budget]
touches: [flai/internal/serve/budget.go, flai/internal/serve/budget_test.go, flai/internal/serve/notify.go, flai/internal/serve/notify_test.go]
after: [T-3095]
---
# T-3098 Clear a budget state with a reason: journaled under serve.budget and sent as a notice, and budgetHeldWhy names flai serve budget clear

## Work

- Add to `flai/internal/serve/budget.go` an exported clear on the serve folder: by key, or every state whose harness, program, or login matches what a harness change names. It removes the state under the file's lock, journals one `cleared` entry under `serve.budget` with why (`by hand`, or the harness change and what it changed), and sends the notice `budgetNotice` sends for any other clear.
- Carry why in the notice where `flai/internal/serve/notify.go`'s `BudgetNotice` has no field for it, and cover its JSON in `flai/internal/serve/notify_test.go`.
- Make `budgetHeldWhy` and the error in `Dir.budgets` name `flai serve budget clear <key>` in place of removing the file.
- Waits for T-3095: it clears states under the key T-3095 writes, and both change `budget.go` and `budget_test.go`.

## Done when

- A test clears a held state by key and by a harness's identity, and finds the journal entry and the notice with why.
- A cleared state lets the next start through as the probe, in a test.
- `budgetHeldWhy` names the command, in a test.
- `flai test` on the changed paths passes.

## Notes
