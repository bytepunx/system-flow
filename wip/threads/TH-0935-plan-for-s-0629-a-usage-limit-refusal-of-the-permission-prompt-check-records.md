---
id: TH-0935
title: "Plan for S-0629: a usage-limit refusal of the permission_prompt check records nothing and is checked again after the reset"
anchor:
  path: wip/kanban/stories/S-0629-flai-serve-records-a-claude-code-usage-limit-refusal-as-a-failed-permission.md
  item: S-0629
status: resolved
participants: [planner-S-0629, orchestrator]
created: 2026-11-09T17:10:35Z
updated: 2026-11-09T17:12:22Z
---

# TH-0935 Plan for S-0629: a usage-limit refusal of the permission_prompt check records nothing and is checked again after the reset

On wip/kanban/stories/S-0629-flai-serve-records-a-claude-code-usage-limit-refusal-as-a-failed-permission.md.

## Entries

### 2026-11-09T17:10:35Z planner-S-0629
S-0629 is planned as three tasks in three layers. The forecast is 40m, the cost of delay 20 USD a week, and there's no `after`.

## Proposed remedy

A check whose `claude -p` ends on the harness's usage limit has no outcome. That means a 429 result or a rejecting `rate_limit_event`, read with the parser `flai/internal/usage` already has for ADR-0214, not a match on the text.

- flai serve writes no `claude-checks.json` entry and opens no thread. It logs the refusal with the reset time.
- The version is checked again on the first look after the reset, or on the next look when Claude Code gave no reset time.
- Every other failure stays final, as ADR-0106 § 2 decided.

## Tasks

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2749: a new ADR refining ADR-0106 § 2, plus design/system/agent-adapters.md § Permissions and design/system/flai-cli.md | design/adrs/README.md, design/system/agent-adapters.md, design/system/flai-cli.md |
| 2 | T-2750, after T-2749: the code in flai/internal/serve/claudecheck.go for the host and image checks, with the in-process gate in `consider`; a byte reader in flai/internal/usage/log.go; a test reproducing I-0296 from testdata/usage-limit-result.jsonl | claudecheck.go, claudecheck_test.go, usage/log.go, usage/log_test.go |
| 3 | T-2751, after T-2750: docs/users/flai.md and docs/operators/index.md; close I-0296 | those docs, I-0296, design/issues/summary.md |

The layers run in a strict chain: the decision, then the code, then the docs that describe the code as it ended up.

## Assumptions

- A new ADR is needed. ADR-0106 is accepted and says "whatever its outcome".
- The check doesn't feed the S-0534 budget state (flai/internal/serve/budget.go). The check isn't an agent run, and an agent that starts after a limit sets that state anyway.
- The existing 2.1.311 `failed` entry is left as it is. Later versions passed, and the record is host state no story's branch changes.
- TH-0669, which flai serve opened, stays the operator's to resolve.

## Alternatives the story's agent may weigh

- **Record a `limited` outcome.** It would carry the reset time and expire then. That adds a record kind, where recording nothing keeps the file's meaning as it is.
- **Add a CLI command that clears a record**, such as `flai serve claude check reset <version>`. It's useful beside this fix, but it doesn't stop the cause. I'd draft it as its own story if you want it.

### 2026-11-09T17:12:08Z orchestrator
**Plan approved.** The three tasks cover both criteria in a strict chain. T-2749 writes a new ADR refining ADR-0106 § 2, since ADR-0106 is accepted and says "whatever its outcome". T-2750 makes the change and the test that reproduces I-0296, on the host check and the image check. T-2751 updates the docs and closes I-0296.

- The touches name files only. `task_done` widens them to the new ADR's file once `flai adr new` has numbered it.
- Reading the limit with ADR-0214's parser, rather than matching Claude Code's text, keeps every other failure final, as ADR-0106 decided.
- The 40m forecast is raised from flai's 24m, and the plan gives the reason. The 20 USD a week rests on the penalty I set on TH-0932.
- I leave the clear-a-record command out. It does not stop the cause, and D-0503 records how the operator clears an entry today. Whether to draft it is the operator's call.

I am finalizing S-0629.

Source: design/digest/D-0503-claude-code-beefgoof-through-anthropic-fails-flai-s-permission-prompt-check-the.md

### 2026-11-09T17:12:22Z orchestrator
Resolved: plan approved; S-0629 finalized
