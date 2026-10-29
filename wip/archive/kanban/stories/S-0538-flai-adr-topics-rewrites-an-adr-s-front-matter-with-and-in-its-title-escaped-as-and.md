---
id: S-0538
type: story
nature: remediation
title: "flai adr topics rewrites an ADR's front matter with < and > in its title escaped as < and >"
status: done
owner: alex
created: 2026-10-29T02:01:36Z
updated: 2026-10-29T04:03:42Z
transitions:
  - to: ready
    at: 2026-10-29T02:59:03Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-29T02:59:08Z
    by: agent-S-0538
  - to: review
    at: 2026-10-29T03:18:22Z
    by: system-flow
  - to: in-progress
    at: 2026-10-29T03:35:46Z
    by: flai
  - to: review
    at: 2026-10-29T03:47:10Z
    by: system-flow
  - to: done
    at: 2026-10-29T04:03:42Z
    by: alex
tags: [flai]
touches: [flai/internal/adr/adr.go, flai/internal/adr/adr_test.go, flai/cmd/adr_test.go, design/issues/I-0222-flai-adr-topics-rewrites-an-adr-s-front-matter-with-and-in-its-title-escaped-as-and.md, design/issues/summary.md, design/issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md, flai/internal/guard/roles_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1982
  priced_by: harness
  turns:
    - day: 2026-10-29
      ceremony: 3
      hand_edits: 4
      work: 40
  models:
    - model: claude-opus-5-5
      input: 106
      output: 21200
      cache_read: 7786820
      cache_write: 207159
      cost: 3.6391
  strategic:
    - kind: orchestrator
      seconds: 1916
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 174
          output: 2782
          cache_read: 26378773
          cache_write: 60644
          cost: 5.8306
        - model: claude-sonnet-5-5
          input: 6
          output: 49
          cache_read: 34277
          cache_write: 29448
          cost: 0.044
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-29T02:01:36Z
  value: 25
  by: planner-S-0538
  at: 2026-10-29T02:56:45Z
forecast:
  duration: 30m
  delivery: 2026-10-29T03:41:00Z
  basis: "Its own forecast of 30m; 1st in the pull order with an in-progress limit of 5, behind S-0532."
  by: flai
  at: 2026-10-29T02:59:06Z
finalized:
  by: orchestrator
  at: 2026-10-29T02:58:55Z
---
# S-0538 flai adr topics rewrites an ADR's front matter with < and > in its title escaped as < and >

## Goal

This story remediates [I-0222](../../../design/issues/I-0222-flai-adr-topics-rewrites-an-adr-s-front-matter-with-and-in-its-title-escaped-as-and.md), "flai adr topics rewrites an ADR's front matter with < and > in its title escaped as < and >". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0222 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0222 is closed with `flai issue close I-0222 --reason` saying what fixed it

## Tasks
- T-2318 flai adr new writes <, >, and & in an ADR's title as themselves, and the title flai prints is the title
- T-2319 Close I-0222 with the reason that flai adr new now quotes a title without HTML escapes
- T-2324 TestADefinedRoleIsHeldToItsFile expects the to-backlog condition's wording since S-0535, so main's integration tier passes again

## Notes

Cost of delay inputs set by flai from I-0222. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T21:41:05Z, 0.2 days before this story; under one cycle counts as one).

### Planning

Cause, read from the code: `flai adr new` (and the MCP tool `adr_new`, which calls the same `adr.New`) writes the title through `yamlString` in `flai/internal/adr/adr.go`. A title holding `>`, `&`, or any other character of its quote set is quoted with `json.Marshal`, which escapes `<`, `>`, and `&` as `\u003c`, `\u003e`, and `\u0026`. YAML reads those escapes back as the characters, so the value is right but the file's text is not. `flai adr topics` only edits the `topics:` line (`topics.SetInFrontMatter`), so it did not make the escape; it showed it, because `titleOf` prints the title as written, quotes trimmed and escapes kept. No ADR on main holds such an escape today (searched `design/`, `docs/`, `wip/`), so no file needs repair. The issue's title reads "escaped as < and >" because its escapes were written through a JSON tool argument, which decodes `\u003c` to `<`: T-2319 gives its reason through the shell.

Proposed fix, for the story's agent to confirm before building it: quote with a JSON encoder whose HTML escaping is off, so `<`, `>`, and `&` are written as themselves, and make `titleOf` read a double-quoted title as YAML does, so the title flai prints is the title. T-2318 makes the fix and its tests; T-2319, after it, closes I-0222.

Touches:

| Path | From | Why |
| --- | --- | --- |
| `flai/internal/adr/adr.go` | design (the cause above) | `yamlString` and `titleOf` change |
| `flai/internal/adr/adr_test.go` | layout | new unit test of the quoting; the package has no test file yet |
| `flai/cmd/adr_test.go` | co-change (`flai touches suggest`) and layout | the command-level reproduction: `adr new` then `adr topics` |
| `design/issues/I-0222-…md` | criterion 2 | `flai issue close` writes it |
| `design/issues/summary.md` | criterion 2 | `flai issue close` regenerates it |

Left out from `flai touches suggest`: `flai/cmd/adr.go` (57%), `flai/internal/check/check.go` and its test, `flai/internal/docedit/docedit.go`, `design/system/flai-cli.md`, and `design/system/flaiover-dashboard.md`. The fix changes no command, flag, output shape, or check, so none of them changes; a story's agent that finds otherwise widens the touches. No folder touch is kept.

Figures:

- Forecast: flai forecast gave 14m from 34 done medium remediation stories. Raised to 30m: the change is one function, but each task close and the close-out run the `flai/cmd` test package and the lint tiers, which alone take longer than 14m here. Delivery is flai's (36th in the pull order), moved by the 16m added.
- Cost of delay: 25.00 USD a week, flai cod's figure from the input flai set from I-0222 (10m lost per 168h cycle at 150 USD an hour). It stands: one occurrence, and the workaround is to leave angle brackets out of a title.
- 2026-10-29T03:35:46Z: moved to in-progress: merge queue batch 20261029-032429 fails with S-0538 at 2a955340bc1e: integration failed: flai/internal/guard/roles_test.go:139 TestADefinedRoleIsHeldToItsFile: {Session: Subagent:false Agent: Role:review Story: Kind:3 Name: Command: Path: Tool:item_move Args:{ID:S-0001 Story: To:review Type: Draft:false Commit:false Recommendation:false Source: Fields:[] CostOfDelay:map[]}}: {Outcome:1 Reason:the role reviewer cannot move S-0001 to review: it may make it only under to-backlog (on a call, a move to backlog and nowhere else; on a trigger, only a move to backlog), or to-ready (a move to ready and nowhere else). Ask the operator with thread_open on the item if it …
