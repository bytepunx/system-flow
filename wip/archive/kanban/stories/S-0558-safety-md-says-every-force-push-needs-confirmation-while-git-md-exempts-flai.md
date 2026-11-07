---
id: S-0558
type: story
nature: remediation
title: safety.md says every force push needs confirmation, while git.md exempts flai serve's force-with-lease push of a story's branch
status: done
owner: alex
created: 2026-11-02T00:04:34Z
updated: 2026-11-07T04:44:25Z
transitions:
  - to: ready
    at: 2026-11-07T04:21:54Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T04:37:10Z
    by: agent-S-0558
  - to: review
    at: 2026-11-07T04:42:37Z
    by: agent-S-0558
  - to: done
    at: 2026-11-07T04:44:25Z
    by: orchestrator
tags: [conventions, template]
topics: [git, safety]
touches: [design/conventions/safety.md, template/root/design/conventions/safety.md, template/CHANGELOG.md, design/issues/I-0189-safety-md-says-every-force-push-needs-confirmation-while-git-md-exempts-flai-serve-s-force-with-lease-push-of-a-story-s-branch.md, design/issues/summary.md, design/issues/I-0270-the-integration-and-smoke-tiers-when-names-template-but-their-paths-do-not.md, design/issues/I-0272-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md]
agent:
  harness: container
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 376
  priced_by: harness
  turns:
    - day: 2026-11-07
      test_runs: 1
      work: 38
  models:
    - model: claude-opus-5-5
      input: 84
      output: 25930
      cache_read: 8116449
      cache_write: 247089
      cost: 4.1189
  strategic:
    - kind: orchestrator
      seconds: 270
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 61
          output: 390
          cache_read: 2351469
          cache_write: 284199
          cost: 0.6335
        - model: claude-sonnet-5-5
          input: 6
          output: 769
          cache_read: 39930
          cache_write: 36733
          cost: 0.0553
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: orchestrator
    at: 2026-11-05T01:35:25Z
  value: 7.5
  by: planner-S-0558
  at: 2026-11-05T01:35:53Z
forecast:
  duration: 12m
  delivery: 2026-11-07T09:35:00Z
  basis: "Its own forecast of 12m; 11th in the pull order with an in-progress limit of 5, behind S-0555, S-0565, S-0567, S-0569, S-0570, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0241, S-0543 and S-0556."
  by: flai
  at: 2026-11-07T04:18:28Z
finalized:
  by: orchestrator
  at: 2026-11-05T01:36:22Z
---
# S-0558 safety.md says every force push needs confirmation, while git.md exempts flai serve's force-with-lease push of a story's branch

## Goal

This story remediates [I-0189](../../../design/issues/I-0189-safety-md-says-every-force-push-needs-confirmation-while-git-md-exempts-flai-serve-s-force-with-lease-push-of-a-story-s-branch.md), "safety.md says every force push needs confirmation, while git.md exempts flai serve's force-with-lease push of a story's branch". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0189 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0189 is closed with `flai issue close I-0189 --reason` saying what fixed it

## Tasks
- T-2457 Both copies of safety.md name flai serve's force-with-lease push of a story's branch as the one exception git.md makes, and the template's changelog records it
- T-2458 I-0189 is closed with a reason naming the safety.md clause that now names git.md's force-with-lease exception

## Notes

### Planning

Proposed fix, from I-0189's one instance: rewrite the clause "force pushes always do" in both copies of `safety.md` § Rules so that an agent's force push still always needs confirmation and `flai serve`'s `--force-with-lease` push of a story's own branch under pull-request acceptance is named as the one exception, linking `git.md`, which states it (S-0393, ADR-0166). `git.md` stays as it is.

Touches, each a file:

| Path | Source | Why |
|------|--------|-----|
| `design/conventions/safety.md` | design | I-0189 names it as the file to fix |
| `template/root/design/conventions/safety.md` | design | the template's copy; conventions land in both in one story |
| `template/CHANGELOG.md` | co-change, layout | S-0393 recorded its `git.md` change there; a template change gets an Unreleased line |
| `design/issues/I-0189-…md` | design | criterion 2 closes it |
| `design/issues/summary.md` | layout | `flai issue close` rewrites the summary |

`flai touches suggest` with the two `safety.md` copies listed both `git.md` copies at 75% and other conventions at 50 to 62%. They are left out: the fix needs no change to them. `template/CHANGELOG.md` and `design/issues` are in `claims.shared`.

Forecast: 12m, flai's figure, kept. It rests on 13 done small remediation stories on the same model, times size 6. Each touch is a line of prose, a changelog line, or a generated close, so nothing argues for more.

Cost of delay: 7.50 USD a week, flai's figure, kept. The input, `time_lost_per_cycle: 3m`, was asked on TH-0824 and set by the orchestrator. At 150 USD an hour and one 168h cycle a week it gives 7.50. Nothing argues for another figure: the fix changes no behaviour.

No story in `after`: the story reads only `git.md`'s exception, which S-0393 wrote and is on main.

Test: none is planned. The cause is a contradiction between two sentences of prose. The template tier renders the template and runs `flai check --strict` on it, which is the check that applies.

### Working

Plan: the planner's two tasks kept unchanged, in two layers. T-2457 is the first. T-2458 waits for it, because its reason names the clause T-2457 writes. The story's agent worked both, without sub-agents: each is one line of prose or one call.

No automated test reproduces the cause: it is a contradiction between two sentences of prose, which no test can judge. The template tier's render and `flai check --strict` of the rendered project is the check that applies, and a `diff` of the two baselines showed them identical.

These decisions are recorded here, not in the narrative's `## Decisions`: this agent runs in a container that mounts `wip/` read-only, and `stream_state` writes only `## Current state` and `## Next steps`.

### Accepted by the orchestrator

- Verified: 6e4f917cdcd4afe0ec2e2dfbbd999b5aaeeff841
- At: 2026-11-07T04:44:25Z

Verdict: pass — both safety.md copies now name flai serve's --force-with-lease push of a story's own branch as the one exception git.md makes, and I-0189 is closed with a reason.
- 1: design/conventions/safety.md, template/root/design/conventions/safety.md, template/CHANGELOG.md
- 2: design/issues/I-0189-safety-md-says-every-force-push-needs-confirmation-while-git-md-exempts-flai-serve-s-force-with-lease-push-of-a-story-s-branch.md, design/issues/summary.md
