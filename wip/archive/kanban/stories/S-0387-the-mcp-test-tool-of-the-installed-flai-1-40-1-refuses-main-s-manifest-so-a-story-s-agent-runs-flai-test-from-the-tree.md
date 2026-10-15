---
id: S-0387
type: story
nature: improvement
title: The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree
status: done
owner: alex
created: 2026-10-14T20:47:00Z
updated: 2026-10-15T02:28:33Z
transitions:
  - to: ready
    at: 2026-10-15T01:11:34Z
    by: alex
  - to: in-progress
    at: 2026-10-15T01:32:11Z
    by: agent-S-0387
  - to: review
    at: 2026-10-15T02:20:43Z
    by: agent-S-0387
  - to: done
    at: 2026-10-15T02:28:33Z
    by: orchestrator
tags: []
topics: [release, cli, mcp, host]
touches: [design/adrs/0153-flai-serve-upgrades-the-host-s-flai-when-the-project-s-history-has-a-newer-published-release-behind-auto-upgrade.md, design/adrs/README.md, design/system/flai-cli.md, flai/internal/serve/upgrade.go, flai/internal/serve/upgrade_test.go, flai/internal/serve/serve.go, flai/internal/serve/outdated_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, flai/internal/mcpserver/test.go, flai/internal/mcpserver/test_test.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, docs/operators/index.md, docs/operators/settings.md, docs/operators/runbooks/update.md, docs/users/flai.md, design/issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md, design/issues/summary.md, design/adrs/0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md, flai/cmd/test.go, flai/cmd/test_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/serve.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2936
  priced_by: harness
  turns:
    - day: 2026-10-15
      ceremony: 2
      hand_edits: 1
      work: 45
  models:
    - model: claude-haiku-5-5
      input: 160
      output: 54996
      cache_read: 11528857
      cache_write: 238837
      cost: 0.7525
    - model: claude-opus-5-5
      input: 268
      output: 105894
      cache_read: 15494746
      cache_write: 438127
      cost: 8.0138
  strategic:
    - kind: orchestrator
      seconds: 840
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 92
          output: 1534
          cache_read: 24314642
          cache_write: 52731
          cost: 6.1087
        - model: claude-sonnet-5-5
          input: 16
          output: 83
          cache_read: 308349
          cache_write: 61139
          cost: 0.3024
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4m
    by: flai
    at: 2026-10-14T20:47:00Z
  value: 10
  by: planner-S-0387
  at: 2026-10-14T20:53:56Z
forecast:
  duration: 1h15m
  delivery: 2026-10-15T03:12:00Z
  basis: "Its own forecast of 1h15m; 7th in the pull order with an in-progress limit of 5, behind S-0381, S-0383, S-0382, S-0241, S-0384, S-0385 and S-0386."
  by: flai
  at: 2026-10-15T01:10:41Z
finalized:
  by: alex
  at: 2026-10-15T00:59:41Z
---
# S-0387 The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree

## Goal

This story remediates [I-0141](../../../design/issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md), "The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0141 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0141 is closed with `flai issue close I-0141 --reason` saying what fixed it

## Tasks
- T-1567 Record the remedy for an outdated host flai in a new ADR and in flai-cli.md § Versions
- T-1568 flai serve checks whether the host's flai is behind each served project on its reconcile loop, upgrades it with auto-upgrade on, and asks the operator otherwise
- T-1569 flai test and the MCP tools test and verify say the running flai is older than the project when the manifest's tiers refuse
- T-1570 Document auto-upgrade, the periodic check, and the thread for operators and users
- T-1571 Close I-0141 with the remedy that fixed it

## Notes

Cost of delay inputs set by flai from I-0141. time_lost_per_cycle 4m: 2m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-12T17:07:12Z, 2.2 days before this story; under one cycle counts as one).

### Planning

The planner's proposed remedy, for the story's agent to confirm or change in the first task's ADR (see the plan's thread on S-0387):

- **The cause.** Both instances are the host's installed flai lagging main after a flai release: 1.40.1 did not know the `{importers}` placeholder S-0364 added to the manifest's tests, and 1.48.0 was a minor behind 1.49.1, so ADR-0142 refused `task_done`. Nothing upgrades the host, and `flai serve` warns only when it starts serving a project (`flai/internal/serve/serve.go`, the `flai is older than the project it serves` log), so a host that runs on through a release never says so again.
- **The remedy.** `flai serve` checks `release.FlaiOutdated` for each served project on its reconcile loop, not only at start. With a new host action `auto-upgrade` on, turned on only in a shell like `auto-approve`, and the newer release published and verifiable, it asks the host to upgrade itself (the path `flai host upgrade` takes). Story agents are detached from the host, so they run on, and their next session's `flai mcp` is the new flai. With it off, it opens one thread asking the operator to run `flai host upgrade`, as ADR-0108's restart limit asks. The MCP tools `test` and `verify`, and `flai test`, say the running flai is older than the project and name the upgrade when the manifest's tiers refuse and `FlaiOutdated` answers, instead of telling the agent to fix the manifest.
- **Not proposed.** Raising `flai.minimum` for a change to the manifest's test vocabulary would make an outdated host stop serving the project at all, which is worse than either instance.

Where each touch came from:

| Touch | Source |
|-------|--------|
| `design/adrs/0153-…-behind-auto-upgrade.md`, `design/adrs/README.md` | design: the remedy is a decision; the ADR's number and slug are predicted, and `flai task done` widens the touches with the real path |
| `design/system/flai-cli.md` | design: § Versions: the host's flai and the tree |
| `flai/internal/serve/upgrade.go`, `upgrade_test.go` | layout: new files for the periodic check, the upgrade, and the thread |
| `flai/internal/serve/serve.go`, `outdated_test.go` | layout: the reconcile loop and today's outdated warning |
| `flai/cmd/serve_actions.go`, `flai/cmd/serve_actions_test.go` | layout and co-change (40 commits with `serve_actions.go`): `serve.Options` wiring |
| `flai/internal/hostapi/writes.go`, `writes_test.go` | layout: the host actions and `shellOnly` |
| `flai/internal/verify/manifest.go`, `manifest_test.go` | layout: `CheckoutTiers`, the refusal both `flai test` and the MCP tools give |
| `flai/internal/mcpserver/test.go`, `test_test.go`, `verify.go`, `verify_test.go` | layout: the MCP tools that refused in the first instance |
| `docs/operators/index.md`, `docs/operators/settings.md`, `docs/operators/runbooks/update.md`, `docs/users/flai.md` | design: § Keeping the host's flai current, the settings index, the update runbook, the host actions row |
| `design/issues/I-0141-…md`, `design/issues/summary.md` | goal: the second criterion closes I-0141 |

No folder touch is kept. Co-change's `docs/users/flai-reference.md` is left out because no flag is added; `template/CHANGELOG.md` because the template does not change.

Figures:

- **Forecast 1h15m**, raised from flai's 34m: five tasks with an ADR, a host action, and a serve loop that drives the host's upgrade, as S-0372 and S-0374 were raised to 1h for comparable plans. Delivery moves by the 41m added, to 2026-10-15T00:08:00Z.
- **Cost of delay 10 USD a week**, flai cod's figure from the 4m time_lost_per_cycle input at 150 USD an hour, kept: both instances were worked around in about 2m each.

### Accepted by the orchestrator

- Verified: 626b4f55f67531a0f7ff536a9c5792ccf41b0d16
- At: 2026-10-15T02:28:33Z

Verdict: meets both criteria; flai verify passed every step at 626b4f55; every changed path is within the story's touches, none under .claude/; auto-upgrade is off by default and shell-only, as the operator approved on TH-0447 (ADR-0158).

- 1: design/adrs/0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md, design/adrs/README.md, design/system/flai-cli.md, flai/internal/serve/upgrade.go, flai/internal/serve/upgrade_test.go, flai/internal/serve/serve.go, flai/internal/serve/outdated_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/serve.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, flai/internal/mcpserver/test.go, flai/internal/mcpserver/test_test.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, flai/cmd/test.go, flai/cmd/test_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, docs/operators/index.md, docs/operators/settings.md, docs/operators/runbooks/update.md, docs/users/flai.md
- 2: design/issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md, design/issues/summary.md
