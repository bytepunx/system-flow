---
id: TH-0468
title: "E-0023 plan: seven draft stories, S-0406 to S-0412, with 38 tasks, in the epic's order; values wait on TH-0467"
anchor:
  path: wip/kanban/epics/E-0023-an-opencode-adapter-runs-a-story-s-agent-on-opencode-with-the-guard-the-permission-prompt-the-roles-and-usage-measured-through-its-plugin-hooks-and-its-server-s-events-so-a-non-claude-model-works-a-story-end-to-end.md
  item: E-0023
status: answered
participants: [planner-E-0023, orchestrator, flai serve]
created: 2026-10-16T05:24:38Z
updated: 2026-10-28T21:06:53Z
---

# TH-0468 E-0023 plan: seven draft stories, S-0406 to S-0412, with 38 tasks, in the epic's order; values wait on TH-0467

On wip/kanban/epics/E-0023-an-opencode-adapter-runs-a-story-s-agent-on-opencode-with-the-guard-the-permission-prompt-the-roles-and-usage-measured-through-its-plugin-hooks-and-its-server-s-events-so-a-non-claude-model-works-a-story-end-to-end.md.

## Entries

### 2026-10-16T05:24:38Z planner-E-0023
I drafted seven stories for E-0023, all in the backlog, each with its tasks, touches, and forecast. The cost of delay values wait on TH-0467, question 1.

## Stories and order

| Story | What | After | Forecast | flai's |
|-------|------|-------|----------|--------|
| S-0406 | `opencode` adapter, headless run, MCP entry, protected list, instruction-file check; the four open findings checked against the binary | none | 1h15m | 51m |
| S-0407 | guard plugin on `tool.execute.before`; `Guard: true` only where a sub-agent's call can be told apart | S-0406 | 45m | 34m |
| S-0408 | protected write asked through hold-and-ask, or `--auto` with deny rules, by an ADR | S-0407 | 50m | 39m |
| S-0409 | explorer, verifier, and strategic agents translated from `.claude/agents` | S-0408 | 40m | 31m |
| S-0410 | usage reader, `priced_by`, provider entry from the host's map | S-0409 | 50m | 37m |
| S-0411 | card and page status through the run's adapter; ADR on `--format json` or `opencode serve` | S-0410 | 35m | 24m |
| S-0412 | experiment: a non-Claude model works one story end to end through OpenRouter | S-0411 | 30m | 17m |

The chain is the epic's own order. Every story carries its own docs and template changes, so there is no docs story. Each story's `### Planning` notes give the reason for each raise.

## Tasks and layers

- **S-0406:** T-1682, T-1683, T-1684, T-1685, T-1686; then T-1687; then T-1688; then T-1689.
- **S-0407:** T-1690; then T-1691, T-1692; then T-1693; then T-1694.
- **S-0408:** T-1695; then T-1696, T-1697; then T-1698; then T-1699; then T-1700.
- **S-0409:** T-1701; then T-1702, T-1703; then T-1704; then T-1705.
- **S-0410:** T-1706, T-1707; then T-1708; then T-1709; then T-1710; then T-1711.
- **S-0411:** T-1712, T-1713; then T-1714; then T-1715.
- **S-0412:** T-1716; then T-1717; then T-1718, T-1719.

## Assumptions

1. TH-0467's recommended answers hold: OpenCode installed on this host before S-0406 starts, and an OpenRouter key with a limit for S-0412.
2. A story's OpenCode configuration must leave its worktree clean, or the move to review is refused. S-0406's ADR (T-1686) chooses between a file under `.flai-cache` that `OPENCODE_CONFIG` names and a git-ignored path. S-0409's agents land in the same place, so "generated into `.opencode/agents/`" in the epic may become "generated where the configuration lives".
3. The template keeps `CLAUDE.md` as the map, which OpenCode reads as its fallback, as `agent-adapters.md` § Definitions recommends. No `AGENTS.md` is shipped.
4. `.claude/agents/` stays the source of the role definitions (ADR-0130), translated at start through a tool-name table.
5. S-0407's capability follows S-0406's finding. If a plugin cannot tell a sub-agent's call, `Guard` stays false, and S-0409's roles stay refused under ADR-0131. S-0412 then runs without roles.
6. S-0412's target story and model are the operator's choice, on a thread on S-0412. The operator sets the target's `agent`.

## Proposals

- **Split S-0406 (optional, not drafted):** the four findings (T-1682) as a `research` story of their own first. The findings would land and be reviewed before any adapter code. I kept them in one story because the epic names them together, and the adapter needs the binary anyway.
- **Split S-0411 if its ADR chooses `opencode serve`:** a long-running OpenCode server under `flai serve` would be a story of its own. S-0411's Notes say so.
- **Not drafted:** LiteLLM through OpenCode. The epic names OpenRouter only, and S-0410's provider entry covers any `openai-chat` gateway.

## Overlap with other epics

These stories share files with E-0021's and E-0022's drafts, so the holds run them one after the other:

- `flai/internal/guard/guard.go`: S-0407 and S-0393.
- `flai/internal/serve/agents.go` and `agents_test.go`: S-0406 and S-0411, with S-0394.

Nothing else overlaps outside `design/adrs` and the shared docs.

### 2026-10-21T11:04:58Z orchestrator (recommendation)
Recommended: resolve this thread, approving the plan as drafted.

- Each story has its value from the 500 USD a week penalty you set on E-0023 (TH-0467, question 1).
- S-0407 to S-0411 are finalized and wait on S-0406 by `after`.
- S-0406 and S-0412 stay drafts. They need only TH-0467's questions 2 and 3: OpenCode on the host, and an OpenRouter key and model.

Nothing on this thread needs a decision beyond those two.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-28T21:06:53Z flai serve
Confirmed the recommendation of 2026-10-21T11:04:58Z orchestrator after 48h with no answer from the operator (orchestration.permissions.confirm_after).

Source: design/conventions/strategic-agents.md § As the planner
