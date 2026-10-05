---
id: TH-0136
title: "S-0220 plan: eight tasks in four layers, and the assumptions behind them"
anchor:
  path: wip/kanban/stories/S-0220-the-orchestrator-answers-threads-or-recommends-an-answer-as-its-permission-allows.md
  item: S-0220
status: resolved
participants: [planner-S-0220, alex]
created: 2026-10-05T04:49:27Z
updated: 2026-10-05T05:03:35Z
---

# TH-0136 S-0220 plan: eight tasks in four layers, and the assumptions behind them

On wip/kanban/stories/S-0220-the-orchestrator-answers-threads-or-recommends-an-answer-as-its-permission-allows.md.

## Entries

### 2026-10-05T04:49:27Z planner-S-0220
I planned S-0220 with eight tasks in the backlog, in four layers. No two tasks in a layer share a path.

| Layer | Task | What it delivers | Waits for |
|-------|------|------------------|-----------|
| 1 | T-0894 | A thread entry can be a recommendation with a source; `threads.Confirm`; ADR | nothing |
| 2 | T-0902 | `flai thread reply --recommend --source`, `flai thread confirm`, the MCP `thread_reply` parameters, host API `thread.confirm`, decision log entry; docs | T-0894 |
| 2 | T-0906 | `flai stats` counts waits the orchestrator ended, and confirmed ones, apart; a recommendation ends no wait; ADR for `metrics.md` | T-0894 |
| 3 | T-0908 | A pending recommendation is no answer to the asking agent (MCP `awaiting`, flai serve's `answered`), and `inbox.designer` lists it | T-0902 (both change `mcpserver/server.go`) |
| 3 | T-0910 | `flai guard` applies `answer_threads` to the orchestrator's thread calls | T-0902 |
| 3 | T-0912 | The orchestrator's prompt, definition, and convention: recommend, answer, or escalate, citing a source | T-0902 |
| 4 | T-0913 | Dashboard: recommendation mark, source link, and **Confirm** in Threads and the inbox | T-0908 |
| 4 | T-0914 | End-to-end test of a recommendation, an autonomous answer, and an escalation | T-0906, T-0908, T-0910, T-0912 |

Figures: touches widened from 15 to 33 paths. The forecast is 1h15m (flai's history gives 45m to 1h24m depending on the touch count; six comparable cross-cutting stories took 53 to 72 min). The cost of delay value is 76.01 USD a week, S-0220's share of E-0016's 1500 USD. The reasons are under `### Planning` in the story's Notes.

Assumptions to confirm or correct:

1. S-0218 lands first and gives the orchestrator its role in `FLAI_ROLE`, its agent name, `orchestration.permissions.answer_threads`, guard enforcement by permission, the decision log, and `.claude/agents/orchestrator.md` with a template copy. S-0220 refines these and creates none of them.
2. A recommendation is a mark on an ordinary entry, such as a heading suffix plus a `Source:` line, so older flai reads it as a plain entry. T-0894 settles the form in an ADR.
3. Confirming is the operator's alone. It is a CLI command and a dashboard action, not an MCP tool. It appends an entry by the operator and sets `answered`.
4. With `recommend`, a recommendation does not wake the agent that asked; the operator's confirmation does. With `autonomous`, a sourced answer wakes it at once.
5. flai enforces escalation in one way only: in `autonomous`, an answer without a source is refused, and the guard says to post it as a recommendation. Whether a question names the operator's judgement (a decision, scope, money) is the orchestrator's call, made in its prompt; flai does not classify questions.
6. "Never answers a thread opened by itself" means its replies on its own threads are follow-ups, with no recommendation or source. "Never resolves a thread it did not open" holds whatever the permission is.
7. `thread_reply` in role `orchestrate` writes the decision log entry with the source itself. Logging then does not depend on the model calling `activity_log`.
8. The metrics add fields and change none, so S-0215's waiting chart keeps working. A confirmed recommendation counts as `confirmed`, not as the orchestrator's.
9. The orchestrator's thread rules go into the baseline convention (`design/conventions/strategic-agents.md`, `## As the orchestrator`) and its template copy, with a CHANGELOG entry. If you would rather keep them only in the definition and prompt, T-0912 drops those three paths.

Split, merge, or drop: none proposed. T-0914 could fold into the earlier tasks' tests. I kept it separate because criterion 5 asks for the three scenarios end to end, across the guard, the MCP server, serve, and stats.

A flai defect found while planning: the MCP `item_edit` refuses `.claude/agents/orchestrator.md` as a touch ("letters, digits, and _ . / @ + - only, not starting with a dash"), but `item_new` and `flai touches` accept it. I set the story's touches with `flai touches`.

### 2026-10-05T05:03:35Z alex
Resolved.
