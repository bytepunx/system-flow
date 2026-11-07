---
id: T-2457
type: task
nature: remediation
title: Both copies of safety.md name flai serve's force-with-lease push of a story's branch as the one exception git.md makes, and the template's changelog records it
status: done
parent: S-0558
owner: alex
created: 2026-11-05T01:34:11Z
updated: 2026-11-07T04:38:32Z
transitions:
  - to: ready
    at: 2026-11-07T04:37:59Z
    by: agent-S-0558
  - to: in-progress
    at: 2026-11-07T04:38:09Z
    by: agent-S-0558
  - to: done
    at: 2026-11-07T04:38:32Z
    by: agent-S-0558
stream: S-0558
tags: [conventions, template]
touches: [design/conventions/safety.md, template/root/design/conventions/safety.md, template/CHANGELOG.md]
usage:
  source: log
  seconds: 23
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 1235
      cache_read: 386569
      cache_write: 11768
      cost: 0.1962
---
# T-2457 Both copies of safety.md name flai serve's force-with-lease push of a story's branch as the one exception git.md makes, and the template's changelog records it

## Work

In `design/conventions/safety.md` and `template/root/design/conventions/safety.md`, § Rules, the sub-item under "Confirm before any action that is hard to reverse" ends "force pushes always do". `git.md` § Rules, in both copies since S-0393 (ADR-0166), makes one exception: under pull-request acceptance, after `flai stream sync` rebases `story/S-nnnn`, `flai serve`, never an agent, pushes it with `--force-with-lease` on the head it last pushed. An agent that reads both meets a contradiction (I-0189).

Rewrite that clause in the baseline of both copies, above the `system-flow:end-of-baseline` marker, word for word the same, so that it:

- keeps "an agent's force push always needs confirmation";
- names `flai serve`'s `--force-with-lease` push of a story's own branch under pull-request acceptance as the one exception, and links `git.md` for it rather than restating the rule.

Bump `updated` in the front matter of both copies. Change nothing else in either file, and leave `git.md` as it is: it already states the exception.

Add one line under `## Unreleased` in `template/CHANGELOG.md`, in the form the other entries take: `S-0558` and what changed in `safety.md`, marked `(patch)`.

This task waits for none: it is the first layer.

## Done when

- Both copies of `safety.md` carry the same new clause, and `diff` of their baselines shows no difference.
- Neither copy says that every force push needs confirmation without naming the exception `git.md` makes.
- `template/CHANGELOG.md` has the `S-0558` line under `## Unreleased`.
- `flai test` on the three paths passes, and the markdown lint is clean on them.

## Notes
