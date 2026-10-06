---
id: S-0270
type: story
nature: improvement
title: "Verification is a module: flai verify runs the tiers the diff selects and answers structured findings, replacing the verifier sub-agent's run"
status: ready
parent: E-0017
owner: alex
created: 2026-10-05T01:35:29Z
updated: 2026-10-06T23:33:06Z
transitions:
  - to: ready
    at: 2026-10-06T22:48:06Z
    by: alex
tags: [cli, mcp, dashboard]
topics: [automation, mcp, hostapi, code, conventions, template, dashboard]
touches: [flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/root.go, flai/internal/verify, flai/internal/serve/checks.go, flai/internal/serve/checks_config.go, scripts/close-out.sh, scripts/README.md, template/root/scripts/close-out.sh, template/root/scripts/README.md, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/verify_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/verifier.md", template/root/.claude/agents/verifier.md, design/conventions/delegation.md, design/conventions/work-management.md, template/root/design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/server/agent.ts, "flaiover/src/routes/api/items/[id]/verify/+server.ts", "flaiover/src/routes/api/items/[id]/verify/verify.test.ts", design/adrs, design/system/flai-cli.md, design/system/devex.md, design/system/dashboard-host-channel.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md]
after: [S-0273]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 43
  by: planner-E-0017
  at: 2026-10-06T11:36:18Z
forecast:
  duration: 90m
  delivery: 2026-10-07T02:25:00Z
  basis: "Its own forecast of 1h30m; 5th in the pull order with an in-progress limit of 3, behind S-0300, S-0302, S-0273, S-0228 and S-0269."
  by: flai
  at: 2026-10-06T23:33:06Z
finalized:
  by: alex
  at: 2026-10-06T22:48:02Z
---
# S-0270 Verification is a module: flai verify runs the tiers the diff selects and answers structured findings, replacing the verifier sub-agent's run

## Goal

Before review a story agent starts a Sonnet verifier sub-agent to run the close-out and read its log: 88 minutes over 22 stories, and in S-0248 three runs of the same two-minute script. `flai verify S-nnnn`, `verify` over MCP, and `verify.run` on the host channel run what `scripts/close-out.sh` runs for the branch's diff, cheapest tier first, and answer findings as data: step, path, line, message, and whether the step passed, failed, or was not reached, with the duration of each. A finding outside the story (S-0249) is a note, not a failure. The agent reads twenty lines, not a log, and the sub-agent is kept for judgement, checking the diff against the criteria and the conventions, if at all.

## Acceptance criteria
- [ ] `flai verify S-nnnn` runs the project's lint, test tiers, check, and narrative check for the story's worktree and answers one structured result, as text and `--json`, with each step's state and duration
- [ ] The same is `verify` over MCP and `verify.run` on the host channel; the dashboard's story page can show the last result
- [ ] `scripts/close-out.sh` calls it, or shares its implementation, so that the two never disagree
- [ ] `design/conventions/delegation.md`, the harness prompt, and the verifier definition send the agent to `flai verify` for the run and keep the sub-agent for the criteria and conventions review
- [ ] `design/system/flai-cli.md`, `devex.md`, and the user guide describe it

## Tasks
- T-1071 flai/internal/verify runs the close-out's steps for a story's worktree and answers one result with each step's state, duration, and findings
- T-1075 flai verify S-nnnn prints the verification result as text or --json and exits non-zero when a step fails
- T-1076 verify over MCP and verify.run and verify.status on the host channel answer the verification result
- T-1082 Both close-out scripts run flai verify for their checks and keep only the commit around it
- T-1089 The review page shows a story's last verification result, each step with its state, duration, and findings
- T-1096 The conventions, the harness prompts, and the verifier definitions send the agent to flai verify for the run and keep the verifier for the review against the criteria
- T-1108 The CLI design, devex, the host channel design, and the user guides describe flai verify, its MCP tool and host methods, and the review page's result

## Notes

Companion to S-0266 (verifier exit status) and S-0267 (duplicate tier), which stay worth doing until this lands.

### Planning

Revisited by planner-S-0270 on 2026-10-06. Its tasks are drafted in four layers:

1. T-1071.
2. T-1075 and T-1076.
3. T-1082, T-1089, and T-1096.
4. T-1108.

The plan and its assumptions are on the plan thread.

Touches. Every touch planner-E-0017 declared is kept. `flai touches suggest S-0270` listed no further path that bears on the story: its top co-changes, such as `flaiover-dashboard.md` and `docs/operators/index.md`, change with these files for other reasons. Where each touch came from:

- `flai/cmd/verify.go`, `verify_test.go`, `root.go`: declared, layout. A new command (T-1075).
- `flai/internal/verify`: declared, layout. S-0273 adds this package. T-1071 adds `verify.go` and `verify_test.go` to it. **Folder kept**, because S-0273 decides the package's files, and in the story's claim flai replaces the folder with the files the tasks name (ADR-0096).
- `flai/internal/serve/checks.go`, `checks_config.go`: declared, layout. No task plans to change them. `flai verify` runs S-0273's tiers, not the review checks' commands. They are kept because they were declared. The plan thread proposes dropping them.
- `scripts/close-out.sh`, `scripts/README.md`, and their `template/root` copies: declared, co-change (4, 4, and 3 of 7) and criterion 3 (T-1082).
- `flai/internal/mcpserver/folder.go`, `flai/internal/hostapi/hostapi.go`, `writes.go`, `writes_test.go`: declared, layout. The `verify` tool and the `verify.run` and `verify.status` methods (T-1076).
- `flai/internal/mcpserver/verify_test.go`: added, layout. It is the tool's test, beside the other tools' tests (T-1076).
- `flaiover/src/lib/components/Review.svelte`, `flaiover/src/lib/review.ts`: declared, layout. The review panel already shows checks.
- `Review.svelte.test.ts`, `review.test.ts`, `flaiover/src/lib/server/agent.ts`, `flaiover/src/routes/api/items/[id]/verify/+server.ts`, and `verify.test.ts`: added, layout. The panel reads the last result from the host, as the `checks` route does through `agent.ts` (T-1089).
- `flai/internal/harness/harness.go`, `harness_test.go`, both `verifier.md`, `delegation.md`, `work-management.md`, their `template/root` copies, and `template/CHANGELOG.md`: declared, design (criterion 4, T-1096). The orchestrator's `accept_reviews` paragraph in `harness.go` also hands the run to the verifier. From S-0257's release, Claude Code asks the operator on a thread before an agent writes under `.claude/`.
- `design/adrs`: declared, design. The new ADR refines ADR-0059's hand-off. **Folder kept**, because its number is not known until `flai adr new` writes it (T-1096).
- `design/system/dashboard-host-channel.md`: added, design. It lists the host channel's methods (T-1108).
- `design/system/flai-cli.md`, `devex.md`, `docs/users/flai.md`, `flai-reference.md`, `flaiover.md`: declared, criterion 5 (T-1108).

Topics: `conventions` and `template` were added for T-1096, and `dashboard` for T-1089.

Forecast: 90m, against flai's 50m. flai's figure is 83 s per unit over 25 done large improvement stories, times size 36. I raised it, from planner-E-0017's 70m, for three reasons:

- There are seven tasks in four layers.
- Checking T-1082 and the story needs at least two full close-out runs across the flai and flaiover tiers.
- The `.claude` verifier edit waits for the operator's approval.

flai's delivery, 2026-10-07T01:05Z, ignores `after`. S-0270 waits for S-0273, so the delivery is S-0273's forecast delivery, 2026-10-07T08:35Z, plus 90m: 2026-10-07T10:05Z. S-0273 is still a draft, so the date moves with it.

Cost of delay: 43 USD a week, which stands, against flai's 164.92. flai shares E-0017's 900 USD a week by forecast time: this story's 70m of the epic's 6h22m. planner-E-0017 shared it by the turns each story removes instead. The evidence is 88 verifier minutes, about 290 turns at the mean of 18 s a turn. Most of the test-log reading is S-0273's. Turns removed measure the value better than time to build. Raising the forecast does not change that.
