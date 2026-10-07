---
id: T-1122
type: task
nature: improvement
title: The design, the users' and operators' guides, and flai serve actions say what auto-approve allows, what the review shows, and what the Claude Code check does
status: done
parent: S-0286
owner: alex
created: 2026-10-06T22:54:18Z
updated: 2026-10-07T01:36:17Z
transitions:
  - to: ready
    at: 2026-10-07T01:28:32Z
    by: agent-S-0286
  - to: in-progress
    at: 2026-10-07T01:28:32Z
    by: agent-S-0286
  - to: done
    at: 2026-10-07T01:36:17Z
    by: agent-S-0286
stream: S-0286
tags: [flai]
touches: [design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, docs/operators/settings.md, flai/internal/hostapi/writes.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, docs/operators/index.md, docs/operators/runbooks/backup.md]
after: [T-1114, T-1120]
usage:
  source: log
  seconds: 465
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 111
      output: 44849
      cache_read: 7680172
      cache_write: 194330
      cost: 3.6121
---
# T-1122 The design, the users' and operators' guides, and flai serve actions say what auto-approve allows, what the review shows, and what the Claude Code check does

## Work

Write down what T-1110, T-1114, and T-1120 built:

- **`flai/internal/hostapi/writes.go`:** `ActionAutoApprove`'s description, which `flai serve actions` prints. It says auto-approve allows writes to every protected path in the agent's own in-progress story's worktree but `.git`, and that a story changing one is accepted by the operator only.
- **`design/system/flai-cli.md`:** the `flai serve` and `flai mcp` rows and `flai accept`. Cover the wider scope, the acceptance gate, and the check of each new Claude Code version.
- **`design/system/flaiover-dashboard.md`:** the review page's list of protected files.
- **`docs/users/flai.md`:** the section on writes under `.claude/`, renamed if it now covers more. Say what auto-approve allows, that acceptance is where the operator checks these changes, and what the check does and when it opens a thread.
- **`docs/users/flai-reference.md`:** regenerate it for any change to `flai accept`'s help.
- **`docs/users/flaiover.md`:** under "Reviewing a story", the protected files and who accepts.
- **`docs/operators/settings.md`:** the `auto-approve` host action's row, and where the check records its outcome.

Waits for T-1114 and T-1120, whose behaviour it describes. It needs T-1110 too, which T-1120 already waits for.

## Done when

- Each file above says what auto-approve allows, what the review shows, and what the check does, as built.
- `flai serve actions` prints the new description.
- The markdown lint and `flai check --strict` are clean, and `scripts/flai-test.sh` passes.

## Notes

Drafted by planner-S-0286. `docs/operators/index.md` was left out because it does not describe auto-approve today. The story's `docs/operators` folder touch covers it if the agent finds it should.
