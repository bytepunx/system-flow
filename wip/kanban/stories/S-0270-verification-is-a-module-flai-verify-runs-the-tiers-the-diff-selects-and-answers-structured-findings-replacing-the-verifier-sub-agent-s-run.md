---
id: S-0270
type: story
nature: improvement
title: "Verification is a module: flai verify runs the tiers the diff selects and answers structured findings, replacing the verifier sub-agent's run"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:29Z
updated: 2026-10-06T20:20:49Z
transitions: []
tags: [cli, mcp, dashboard]
topics: [automation, mcp, hostapi, code]
touches: [flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/root.go, flai/internal/verify, flai/internal/serve/checks.go, flai/internal/serve/checks_config.go, scripts/close-out.sh, scripts/README.md, template/root/scripts/close-out.sh, template/root/scripts/README.md, flai/internal/mcpserver/folder.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/verifier.md", template/root/.claude/agents/verifier.md, design/conventions/delegation.md, design/conventions/work-management.md, template/root/design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/review.ts, design/adrs, design/system/flai-cli.md, design/system/devex.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md]
after: [S-0273]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  value: 43
  by: planner-E-0017
  at: 2026-10-06T11:36:18Z
forecast:
  duration: 70m
  delivery: 2026-10-07T08:50:00Z
  basis: "Its own forecast of 1h10m; 27th in the pull order with an in-progress limit of 3, behind S-0223, S-0224, S-0227, S-0229, S-0212, S-0213, S-0214, S-0215, S-0216, S-0228, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254, S-0261, S-0264, S-0265 and S-0269."
  by: flai
  at: 2026-10-06T20:20:49Z
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

## Notes

Companion to S-0266 (verifier exit status) and S-0267 (duplicate tier), which stay worth doing until this lands.

### Planning

Touches, none declared before. `flai touches suggest S-0270` was seeded with `scripts/close-out.sh` and `flai/cmd/checks.go`, which 7 commits changed:

- `flai/cmd/verify.go`, `flai/cmd/verify_test.go`, `flai/cmd/root.go`: layout. A new command.
- `flai/internal/verify`: layout. The tier runner S-0273 adds, extended with the close-out's steps.
- `flai/internal/serve/checks.go`, `checks_config.go`: layout. `flai checks` already runs a review's steps with state and duration per step. Its runner may be shared rather than duplicated.
- `scripts/close-out.sh`, `scripts/README.md`, `template/root/scripts/close-out.sh`, `template/root/scripts/README.md`: co-change (4, 4, and 3 of 7) and criterion 3.
- `flai/internal/mcpserver/folder.go`, `flai/internal/hostapi/hostapi.go`, `writes.go`, `writes_test.go`: layout. The `verify` tool and the `verify.run` method.
- `flaiover/src/lib/components/Review.svelte`, `flaiover/src/lib/review.ts`, `docs/users/flaiover.md`: layout. The review panel already shows checks, and criterion 2 shows the last result there.
- `flai/internal/harness/harness.go`, `harness_test.go`, `.claude/agents/verifier.md`, `template/root/.claude/agents/verifier.md`, `design/conventions/delegation.md`, `work-management.md`, their `template/root` copies, and `template/CHANGELOG.md`: design (criterion 4). The close-out rule in `work-management.md` names the script. `.claude/` writes need the operator to paste them, so the story's agent will ask on a thread.
- `design/adrs`: design. Moving the verifier's run into flai changes ADR-0059's hand-off.
- `design/system/flai-cli.md`, `devex.md`, `docs/users/flai.md`, `flai-reference.md`: criterion 5.

Forecast: flai gave 54m (89 s per unit over 21 done large improvement stories, times size 36). I raised it to 70m for the dashboard panel and the operator's copy of the verifier definition. flai's delivery, 2026-10-07T01:07Z, counts from 32nd place in the pull order and ignores `after`. S-0270 waits for S-0273, so I set the delivery to S-0273's, 2026-10-07T02:01Z, plus 70m: 2026-10-07T03:11Z.

Cost of delay: 43 USD a week, against flai's 184.21. This is E-0017's 900 USD a week shared by the turns each story removes. The evidence is 88 verifier minutes, about 290 turns at the mean of 18 s a turn. Most of the test-log reading this story replaces is S-0273's.
