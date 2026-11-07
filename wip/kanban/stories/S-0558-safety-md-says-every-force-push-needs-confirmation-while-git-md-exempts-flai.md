---
id: S-0558
type: story
nature: remediation
title: safety.md says every force push needs confirmation, while git.md exempts flai serve's force-with-lease push of a story's branch
status: backlog
owner: alex
created: 2026-11-02T00:04:34Z
updated: 2026-11-07T01:34:32Z
transitions: []
tags: [conventions, template]
topics: [git, safety]
touches: [design/conventions/safety.md, template/root/design/conventions/safety.md, template/CHANGELOG.md, design/issues/I-0189-safety-md-says-every-force-push-needs-confirmation-while-git-md-exempts-flai-serve-s-force-with-lease-push-of-a-story-s-branch.md, design/issues/summary.md]
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
      seconds: 148
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 37
          output: 225
          cache_read: 1396365
          cache_write: 141478
          cost: 0.3661
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
  delivery: 2026-11-07T06:55:00Z
  basis: "Its own forecast of 12m; 15th in the pull order with an in-progress limit of 5, behind S-0578, S-0580, S-0581, S-0579, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0241, S-0543, S-0555, S-0556 and S-0557."
  by: flai
  at: 2026-11-07T01:34:32Z
finalized:
  by: orchestrator
  at: 2026-11-05T01:36:22Z
---
# S-0558 safety.md says every force push needs confirmation, while git.md exempts flai serve's force-with-lease push of a story's branch

## Goal

This story remediates [I-0189](../../../design/issues/I-0189-safety-md-says-every-force-push-needs-confirmation-while-git-md-exempts-flai-serve-s-force-with-lease-push-of-a-story-s-branch.md), "safety.md says every force push needs confirmation, while git.md exempts flai serve's force-with-lease push of a story's branch". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0189 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0189 is closed with `flai issue close I-0189 --reason` saying what fixed it

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
