---
id: S-0360
type: story
nature: improvement
title: flai serve checks each new claude version through the provider a project's agents name, so a host that reaches Claude only through a gateway is not reported failing
status: done
parent: E-0019
owner: alex
created: 2026-10-08T08:59:39Z
updated: 2026-10-12T21:06:32Z
transitions:
  - to: ready
    at: 2026-10-12T20:03:43Z
    by: alex
  - to: in-progress
    at: 2026-10-12T20:40:20Z
    by: agent-S-0360
  - to: review
    at: 2026-10-12T21:05:07Z
    by: agent-S-0360
  - to: done
    at: 2026-10-12T21:06:32Z
    by: orchestrator
tags: [cli]
topics: [agents]
touches: [flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/system/agent-adapters.md, design/issues/I-0143-flai-check-finds-board-order-outside-the-story-at-close-out.md, design/issues/summary.md, design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md]
after: [S-0351]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1508
  turns:
    - day: 2026-10-12
      ceremony: 3
      hand_edits: 2
      work: 38
  models:
    - model: claude-opus-5-5
      input: 134
      output: 60077
      cache_read: 9218609
      cache_write: 335637
      cost: 5.4015
  strategic:
    - kind: orchestrator
      seconds: 189
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 45
          output: 396
          cache_read: 6901850
          cache_write: 89302
          cost: 1.7502
cost_of_delay:
  value: 22.12
  by: planner-E-0019
  at: 2026-10-08T09:00:28Z
forecast:
  duration: 20m
  delivery: 2026-10-12T21:02:00Z
  basis: "Its own forecast of 20m; 1st in the pull order with an in-progress limit of 5, behind S-0238 and S-0370."
  by: flai
  at: 2026-10-12T20:40:11Z
finalized:
  by: alex
  at: 2026-10-12T20:03:40Z
---
# S-0360 flai serve checks each new claude version through the provider a project's agents name, so a host that reaches Claude only through a gateway is not reported failing

## Goal

`flai serve` checks `permission_prompt` against each `claude` version it has not seen. It runs one headless `claude -p` on `haiku` with the adapter's permission arguments, and opens a thread on every project whose agents run `claude-code` when the check fails ([ADR-0106](../../../design/adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md)). The check always calls Anthropic directly. On a host whose agents reach Claude only through a gateway, because the operator's subscription is capped and the key is OpenRouter's or a LiteLLM proxy's, the check fails for want of Anthropic access, not because `permission_prompt` broke. This story runs the check through the provider the project's default agent names, with the environment S-0350 derives, and records which provider the check ran through.

## Acceptance criteria

- [x] When a served project's default agent names a provider, the check for that project runs through it: the `claude-code` adapter's provider environment, the `haiku` alias resolved through the entry's `models`, and the key copied from its variable as S-0350 copies it.
- [x] `claude-checks.json` records, per version, the provider each check ran through, `anthropic` when none; a version passed through one provider is checked again through another the first time a project names it.
- [x] A provider whose key's variable is unset is not checked and opens no thread; `flai serve`'s log says the check was skipped and why.
- [x] Tests cover a check with a provider, without one, a second provider for a version already passed, and an unset key, with a fake `claude`.
- [x] `design/system/flai-cli.md` describes the check per provider.

## Tasks

- T-1426 The claude version check runs through the project's provider, records the provider per version, and skips a provider whose key is unset
- T-1427 flai-cli.md describes the claude version check per provider

## Notes

- Proposed on the plan's thread on E-0019 as an addition: the operator's reason for the epic's cost of delay is a capped subscription, which is the case where today's check fails.
- It comes after S-0351, which factors `claudecheck.go`'s scratch project for the gateway test.

### Planning

Planned by planner-E-0019 on 2026-10-08. Every touch is a file; no folder touch is kept. Two layers: the check (T-1426), then the document (T-1427).

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/serve/claudecheck.go`, `claudecheck_test.go` | layout | The check and its record (T-1426) |
| `design/system/flai-cli.md` | design | Where the check is described (T-1427) |

`touches suggest` listed the user and operator guides and the dashboard's documents. None is taken: the check's thread and its text are unchanged.

Forecast: 20m. `flai forecast` gave 12m, 85.5 s per unit over 11 small improvement stories, size 8. It is raised because the check moves from once per host to once per provider, with a fake `claude` in its tests.

Cost of delay: 22.12 USD a week, as `flai cod` works it out: 20m of the 7h32m forecast over E-0019's 12 stories, of the operator's 500 USD a week penalty. It stands.

### Accepted by the orchestrator

- Verified: d8cd9b6dd0edd68d1376b72b906cf669ce9df361
- At: 2026-10-12T21:06:32Z

Verdict: pass at d8cd9b6dd0edd68d1376b72b906cf669ce9df361; flai verify passed every step at that commit, the verifier matched all five criteria to the diff, and no key or secret value is in it.

- 1: flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go
- 2: flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go
- 3: flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go
- 4: flai/internal/serve/claudecheck_test.go
- 5: design/system/flai-cli.md, design/system/agent-adapters.md, docs/users/flai.md, docs/operators/index.md
