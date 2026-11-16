---
id: S-0807
type: story
nature: feature
title: A story's move needs the suites its transition's verify lists to have passed at its head, and the close-out runs the suites its next transitions need
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:04Z
updated: 2026-11-16T23:42:04Z
transitions: []
tags: [flai]
topics: [verify, workflow]
after: [S-0806]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0807 A story's move needs the suites its transition's verify lists to have passed at its head, and the close-out runs the suites its next transitions need

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §2 and §7. Today no rule reads the verify report: a close-out before review is convention, and it runs every tier the story's paths select.

## Acceptance criteria

- [ ] The rule `verified` is in the vocabulary: a story's move by a transition with `verify:` is refused, or warned of at `warn`, unless its last `flai verify` passed at its branch's head every tier of those suites its paths select.
- [ ] `flai verify <S-nnnn>` and the close-out job run the suites the transitions out of the story's state need, and `--all` every suite.
- [ ] The verify report records each suite's tiers and the head and base they passed at, and a later run reuses them as ADR-0156 reuses tiers.
- [ ] On the baseline, a story's close-out and its move to review run and require what they did before this story.
- [ ] The conventions, the workflow design, and the user guide say what a close-out runs.

## Tasks

## Notes
