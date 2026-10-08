---
id: S-0360
type: story
nature: improvement
title: flai serve checks each new claude version through the provider a project's agents name, so a host that reaches Claude only through a gateway is not reported failing
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:59:39Z
updated: 2026-10-08T10:29:12Z
transitions: []
tags: [cli]
topics: [agents]
touches: [flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go, design/system/flai-cli.md]
after: [S-0351]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 12
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 34
          cache_read: 2042826
          cache_write: 8532
          cost: 0.5054
draft: true
cost_of_delay:
  value: 22.12
  by: planner-E-0019
  at: 2026-10-08T09:00:28Z
forecast:
  duration: 20m
  delivery: 2026-10-08T22:25:00Z
  basis: "Its own forecast of 20m; 29th in the pull order with an in-progress limit of 5, behind S-0232, S-0338, S-0342, S-0337, S-0347, S-0343, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0289, S-0304, S-0305, S-0306, S-0313, S-0349, S-0350, S-0351, S-0352, S-0353, S-0354, S-0355, S-0356, S-0357, S-0358 and S-0359."
  by: flai
  at: 2026-10-08T10:29:12Z
---
# S-0360 flai serve checks each new claude version through the provider a project's agents name, so a host that reaches Claude only through a gateway is not reported failing

## Goal

`flai serve` checks `permission_prompt` against each `claude` version it has not seen. It runs one headless `claude -p` on `haiku` with the adapter's permission arguments, and opens a thread on every project whose agents run `claude-code` when the check fails ([ADR-0106](../../../design/adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md)). The check always calls Anthropic directly. On a host whose agents reach Claude only through a gateway, because the operator's subscription is capped and the key is OpenRouter's or a LiteLLM proxy's, the check fails for want of Anthropic access, not because `permission_prompt` broke. This story runs the check through the provider the project's default agent names, with the environment S-0350 derives, and records which provider the check ran through.

## Acceptance criteria

- [ ] When a served project's default agent names a provider, the check for that project runs through it: the `claude-code` adapter's provider environment, the `haiku` alias resolved through the entry's `models`, and the key copied from its variable as S-0350 copies it.
- [ ] `claude-checks.json` records, per version, the provider each check ran through, `anthropic` when none; a version passed through one provider is checked again through another the first time a project names it.
- [ ] A provider whose key's variable is unset is not checked and opens no thread; `flai serve`'s log says the check was skipped and why.
- [ ] Tests cover a check with a provider, without one, a second provider for a version already passed, and an unset key, with a fake `claude`.
- [ ] `design/system/flai-cli.md` describes the check per provider.

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
