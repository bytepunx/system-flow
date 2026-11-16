---
id: S-0806
type: story
nature: feature
title: The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:03Z
updated: 2026-11-16T23:42:03Z
transitions: []
tags: [flai, template]
topics: [verify, workflow]
after: [S-0805]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §2, §4, §6, §8, §9: the keys and their checks only. Nothing reads them to decide a run yet; the stories after this one do.

## Acceptance criteria

- [ ] The workflow package parses `verify:` on a transition (suites or `all`, with a severity), `runs:` on a state (`suites:` and `kluster:` of `deploy` and `checks`), and top-level `verify_at:` with `batch`, `publish`, and `ci`, and still refuses any other key.
- [ ] The baseline and the template's file give `verify: all` at `warn` on the move to review, `verify: all` on the acceptance, `batch: all`, `ci: all`, no `publish`, and `runs: {kluster: [deploy, checks]}` on `review`, and the two files stay equal.
- [ ] `flai check` reports a suite the file names that no tier has as an error, and a suite of a tier that is not `manual` which no transition, state, or point names as a warning.
- [ ] A flai from before this story refuses a file that uses the keys with words that name the flai that reads them, and `flai upgrade` does not write the keys into a project's own file.
- [ ] The workflow design and the reference describe the keys.

## Tasks

## Notes

- After S-0805: the checks of the file's suites read the suite S-0805 gives each tier.
