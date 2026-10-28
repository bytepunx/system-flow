---
id: TH-0473
title: "E-0024 plan: twelve draft stories, S-0414 to S-0425, with 66 tasks, in the epic's order; values wait on TH-0471"
anchor:
  path: wip/kanban/epics/E-0024-workflows-as-data-a-project-defines-a-strategic-role-or-a-repeatable-formula-of-steps-in-files-flai-runs-and-the-built-in-roles-are-expressed-the-same-way-so-a-new-role-or-a-release-process-needs-no-flai-release.md
  item: E-0024
status: resolved
participants: [planner-E-0024, orchestrator, alex]
created: 2026-10-16T05:54:33Z
updated: 2026-10-28T12:02:59Z
---

# TH-0473 E-0024 plan: twelve draft stories, S-0414 to S-0425, with 66 tasks, in the epic's order; values wait on TH-0471

On wip/kanban/epics/E-0024-workflows-as-data-a-project-defines-a-strategic-role-or-a-repeatable-formula-of-steps-in-files-flai-runs-and-the-built-in-roles-are-expressed-the-same-way-so-a-new-role-or-a-release-process-needs-no-flai-release.md.

## Entries

### 2026-10-16T05:54:33Z planner-E-0024
I drafted twelve stories for E-0024, all in the backlog as drafts. Each has its tasks, file-level touches, and a forecast. The cost of delay values wait on TH-0471, question 1. Every story carries its own docs and template changes, so there is no docs story. Each story's `### Planning` notes say where every touch came from and why each figure stands.

## Stories and order

| Story | What | After | Forecast | flai's |
|-------|------|-------|----------|--------|
| S-0414 | Role file format, ADR, `flai/internal/roles`, `flai role list/show`, `roles.*` check findings, the three baseline role files, `flai upgrade` | none | 1h1m | 1h1m |
| S-0415 | `flai prime --role` for a defined role; the conventions' `roles` take defined roles | S-0414 | 29m | 29m |
| S-0416 | `flai guard` reads each role's rules from its file; the three roles' rules move there, tests unchanged | S-0414 | 1h | 43m |
| S-0417 | `flai role run`, `role_run`, `role.run`; activity document per role; `roles.<name>.agent` | S-0415, S-0416 | 1h15m | 1h15m |
| S-0418 | Triggers: events, schedule, `continuous`, with hold | S-0417 | 51m | 51m |
| S-0419 | Planner, orchestrator, and analyzer start through the role path, with prompts from their files and an embedded baseline | S-0418 | 1h | 44m |
| S-0420 | One role page per role, settings panel, `agent.status` `roles` | S-0419 | 1h4m | 1h4m |
| S-0421 | Formula file format, ADR, check, `flai formula run` makes a story | S-0420 | 51m | 51m |
| S-0422 | `flai serve` works an instance; failure and resume; acceptance with no branch; starters | S-0421 | 1h15m | 59m |
| S-0423 | Baseline `publish` formula, the orchestrator starts it, run once here | S-0422 | 50m | 36m |
| S-0424 | Baseline `import` and `analysis-cadence` formulas, the board's Import starts its formula, each run once here | S-0422 | 50m | 36m |
| S-0425 | Strategic usage per role by an ADR, `flai stats`, charts | S-0423, S-0424 | 50m | 50m |

The total is 11h16m. The last delivery flai plays out is 2026-10-17T01:57Z, behind the 23 stories of E-0021 to E-0023 ahead of these in the pull order.

## Tasks and layers

- **S-0414:** T-1720; then T-1721, T-1722; then T-1723; then T-1724.
- **S-0415:** T-1725, T-1726; then T-1727; then T-1728.
- **S-0416:** T-1729; then T-1730, T-1732, T-1733; then T-1731; then T-1734.
- **S-0417:** T-1735, T-1736, T-1737; then T-1738; then T-1739; then T-1740.
- **S-0418:** T-1741; then T-1743, T-1744, T-1746; then T-1748; then T-1751.
- **S-0419:** T-1753; then T-1754, T-1756; then T-1757, T-1759; then T-1761.
- **S-0420:** T-1762, T-1763; then T-1764; then T-1765; then T-1766; then T-1767.
- **S-0421:** T-1768; then T-1769, T-1770; then T-1771, T-1772; then T-1773.
- **S-0422:** T-1774; then T-1775, T-1776, T-1777, T-1778; then T-1779.
- **S-0423:** T-1780; then T-1781, T-1782; then T-1783, T-1784.
- **S-0424:** T-1785, T-1786; then T-1787, T-1788, T-1789.
- **S-0425:** T-1790; then T-1791, T-1792; then T-1793; then T-1794.

No two tasks of one layer share a path.

## Assumptions

1. TH-0471's recommended answers hold. Questions 2 to 6 shape S-0414, S-0417, S-0421, and S-0422; each task that depends on an answer says what changes with the other.
2. Rules that are not plain allow-lists stay in Go as named conditions the role file references. Examples are the planner's draft-only `item_new`, the analyzer's reports folder, and the orchestrator's thread and permission rules. A project's own role can use them but not add new ones without a flai release.
3. "Behave exactly as today" means the existing tests for the three roles pass with no changed assertion, and their records in `serve/agents.json` keep their fields beside the new `roles`.
4. A project with no `design/roles/` keeps running the three roles from a baseline embedded in flai, with a `roles.baseline` warning (S-0419).
5. `flai release --evaluate` exits 0 whether or not the policy is met. S-0423 adds `--require-met` so that a formula step can stop on it.
6. The live runs in S-0423 and S-0424 need the operator's `formula`, `push`, and `analyze` host actions on, and a scratch repository under a folder named for import.
7. A defined role's activity charges no item until S-0425. It stays in the role's project strategic total.

## Proposals

- **Relax the order (not applied):** S-0421 reads only S-0414's loader, so it could follow S-0414 and run beside the role chain. S-0425 could start after S-0420. Holds on `check.go` and the shared docs would still serialize parts, but the formula chain would finish about four hours sooner. I kept the epic's order.
- **Split S-0417 (optional, not drafted):** its 41 touches are the serve path plus three entry points. The entry points (T-1739) could be a story of their own after the serve path. I kept them together because a role no one can start is not a demonstrable increment.
- **Not drafted:** a fourth baseline role, such as a reviewer. The epic names it as an example, and the fixture role `reviewer` in the tests shows the path. A baseline one would be a story of its own after S-0419.

## Overlap with other epics

These stories share files with E-0021's, E-0022's, and E-0023's drafts, so the holds run them one after the other:

- `flai/internal/guard/guard.go`: S-0416 with S-0393 and S-0407.
- `flai/internal/harness/harness.go`: S-0417 and S-0419 with S-0406.
- `flai/internal/serve/agents.go` and `serve.go`: S-0417, S-0418, S-0420, and S-0422 with S-0393, S-0394, S-0399, and S-0411.
- `flai/cmd/accept.go` and `flai/internal/preview/accept.go`: S-0422 with S-0395, S-0397, S-0399, S-0400, and S-0403.
- `flai/internal/manifest/manifest.go` and `settings.go`, and `flai/cmd/serve_actions.go`: S-0417, S-0420, S-0422, and S-0424 with S-0392, S-0399, and S-0405.
- `flai/internal/hostapi/writes.go` and `hostapi.go`: S-0417, S-0418, S-0420, and S-0422 with S-0393, S-0395, S-0397, S-0398, and S-0399.

S-0409 (E-0023) translates the strategic agents from `.claude/agents` for OpenCode. With role files, it should read each role's `definition` field rather than a fixed list.

### 2026-10-16T05:57:27Z planner-E-0024
One addition for the order between epics. TH-0476 recommends that E-0025 build a fourth built-in role, the curator, in Go on the analyzer's machinery, without waiting for E-0024.

- If you take that answer and E-0025 lands first, S-0419 should also move the curator onto the role path: a fourth role file, `curator.md`, and its starter rewired. That adds about one task to S-0414 and one to S-0419. I have not changed either, since TH-0476 is open.
- If E-0025's curator waits for E-0024, it becomes a role file after S-0419 with no Go of its own, and nothing here changes.

### 2026-10-21T11:05:02Z orchestrator (recommendation)
Recommended: approve the plan, with your TH-0471 answers applied story by story through the story planner. Take the planner's addition about the curator.

- Your TH-0471 answers change the drafts:
  - `design/workflows/roles` and `design/workflows/formulas`, in place of `design/roles` and `design/formulas`;
  - a formula instance gets a branch or worktree, as a story does (question 5);
  - a generated role file, with the Claude agent file symlinked where Claude allows it (question 6).
- I started the planner on S-0414 at 11:03Z under plan_backlog_stories, and S-0415 to S-0425 follow one at a time. Each run brings its story's touches, words, and tasks in line with those answers.
- Each story's value needs the 4h on E-0024 that you chose on TH-0471, written on the epic: `flai edit E-0024` with its cost of delay flags.
- The curator: TH-0476 recommends building it now on the analyzer machinery. So if E-0025 lands first, add one task to S-0414 and one to S-0419 to move it onto the role path, as this thread proposes.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T12:58:20Z alex
Confirmed the recommendation of 2026-10-21T11:05:02Z orchestrator.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-28T12:02:59Z alex
Resolved.
