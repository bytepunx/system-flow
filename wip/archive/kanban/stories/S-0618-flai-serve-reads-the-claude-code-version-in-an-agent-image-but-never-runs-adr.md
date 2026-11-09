---
id: S-0618
type: story
nature: remediation
title: flai serve reads the Claude Code version in an agent image but never runs ADR-0106's permission_prompt check inside it, so a new image's version is only warned of
status: done
owner: alex
created: 2026-11-07T23:39:32Z
updated: 2026-11-09T11:24:26Z
transitions:
  - to: ready
    at: 2026-11-09T10:51:10Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T10:51:23Z
    by: agent-S-0618
  - to: review
    at: 2026-11-09T11:22:56Z
    by: agent-S-0618
  - to: done
    at: 2026-11-09T11:24:26Z
    by: orchestrator
tags: [flai, serve, container, claude-code]
touches: [flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/agents.go, flai/internal/serve/serve.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go, design/system/flai-cli.md, design/system/agent-adapters.md, docs/users/flai.md, docs/operators/index.md, design/issues/I-0253-flai-serve-reads-the-claude-code-version-in-an-agent-image-but-never-runs-adr.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1910
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 2
      hand_edits: 1
      work: 41
  models:
    - model: claude-opus-5-5
      input: 196
      output: 70145
      cache_read: 26042492
      cache_write: 425880
      cost: 9.6716
  strategic:
    - kind: planner
      seconds: 249
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 64
          output: 21810
          cache_read: 4738972
          cache_write: 175934
          cost: 2.7917
    - kind: orchestrator
      seconds: 126
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 38
          output: 343
          cache_read: 1378724
          cache_write: 179354
          cost: 0.4091
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-07T23:39:32Z
  value: 25
  by: planner-S-0618
  at: 2026-11-09T10:49:13Z
forecast:
  duration: 2h
  delivery: 2026-11-09T13:18:00Z
  basis: "Its own forecast of 2h; 1st in the pull order with an in-progress limit of 5, behind S-0617."
  by: flai
  at: 2026-11-09T10:51:14Z
finalized:
  by: orchestrator
  at: 2026-11-09T10:50:49Z
---
# S-0618 flai serve reads the Claude Code version in an agent image but never runs ADR-0106's permission_prompt check inside it, so a new image's version is only warned of

## Goal

This story remediates [I-0253](../../../design/issues/I-0253-flai-serve-reads-the-claude-code-version-in-an-agent-image-but-never-runs-adr.md), "flai serve reads the Claude Code version in an agent image but never runs ADR-0106's permission_prompt check inside it, so a new image's version is only warned of". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0253 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0253 is closed with `flai issue close I-0253 --reason` saying what fixed it

## Tasks
- T-2726 Run ADR-0106's permission_prompt check inside an agent image through docker
- T-2727 Start the in-image check from a container start whose Claude Code version has no record
- T-2728 Document the in-image check and close I-0253

## Notes

Cost of delay inputs set by flai from I-0253. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T07:30:18Z, 5.7 days before this story; under one cycle counts as one).

### Planning

The proposed solution, from I-0253's one instance (S-0508, T-2162, TH-0666): when a container start reads a Claude Code version from its image that `claude-checks.json` has no record of through the project's provider, flai serve runs ADR-0106's check inside that image, beside the start, and records it with the image. The check's claude is the image's, and its MCP server is the image's own `flai` on stdio against the scratch project's configuration, with auto-approve on. Stdio, because over HTTP permission_prompt auto-approves only the agent flai serve runs for the story (ADR-0243). The plan's thread names the alternatives.

Touches, all files, no folder:

| Touch | From |
|-------|------|
| `flai/internal/serve/claudecheck.go` | layout: the check lives here |
| `flai/internal/serve/claudecheck_test.go` | co-change, 33% with the two above |
| `flai/internal/serve/container.go` | layout: `containerVersion` warns today |
| `flai/internal/serve/container_test.go` | co-change, 44%; holds the tests of the warning |
| `flai/internal/serve/agents.go` | co-change, 67%: the launcher's start path |
| `flai/internal/serve/serve.go` | co-change, 22%: makes the `claudeChecker` |
| `flai/internal/harness/container.go` | layout: the docker run's mounts, network, and provider refusal |
| `flai/internal/harness/container_test.go` | layout: the tests of the helper above |
| `design/system/flai-cli.md` | design: § The container harness (S-0508) says the check is not run |
| `design/system/agent-adapters.md` | design: the check's paragraph |
| `docs/users/flai.md` | design: names the warning |
| `docs/operators/index.md` | design: the `claude-checks.json` entry |
| `design/issues/I-0253-flai-serve-reads-the-claude-code-version-in-an-agent-image-but-never-runs-adr.md` | criterion 2 |

`flai touches suggest` also listed the strategic launchers (`curate.go`, `analyze.go`, `orchestrate.go`, `plan.go` and their tests), at 22–33%. They co-changed when S-0509 brought strategic agents into containers. All of them start through `startOf` and `containerReady`, so this story need not change them; they are left out.

Tasks, three layers: T-2726 (the in-image check), then T-2727 (start it from a container start), then T-2728 (docs, close I-0253).

Figures:

- Forecast: flai forecast said 26m, from 104 s per unit over 58 done large remediation stories, size 15. Raised to 2h. The story adds a new docker path, fake-docker tests in two packages, and four documents, and its close-out alone has taken about 8 to 10 minutes today (S-0608, S-0610, S-0616).
- Cost of delay: 25 USD a week, as flai cod works it out from flai's input, 10m lost per 168h cycle at 150 USD an hour. Left as it is: one occurrence, and no container agent hung on it yet.
- `after`: none. The story reads no output of another open story. S-0629 (I-0296, a usage-limit refusal recorded as a failed check) may change `claudecheck.go` too; a trial merge judges that overlap, not an `after`.

### Accepted by the orchestrator

- Verified: 2c4fd3f9aeec78b40ba294a00b0b39baf4df94d5
- At: 2026-11-09T11:24:26Z

Verdict: met — an unrecorded Claude Code version in a container image is now checked inside the image, with a reproducing test, and I-0253 is closed with a reason.
- 1: flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/agents.go, flai/internal/serve/serve.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go
- 2: design/issues/I-0253-flai-serve-reads-the-claude-code-version-in-an-agent-image-but-never-runs-adr.md, design/issues/summary.md
