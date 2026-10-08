---
id: S-0346
type: story
nature: improvement
title: "flai check finds `markdown.MD034` outside the story at close-out"
status: backlog
owner: alex
created: 2026-10-08T08:08:21Z
updated: 2026-10-08T08:53:22Z
transitions: []
tags: [flai, mdlint, serve]
topics: [markdown, planning]
touches: [flai/internal/mdlint/inline.go, flai/internal/mdlint/quote.go, flai/internal/mdlint/quote_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, design/system/agent-narrative.md, design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md, design/issues/summary.md]
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
      seconds: 83
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 24
          output: 448
          cache_read: 6476975
          cache_write: 7796
          cost: 1.5978
draft: true
forecast:
  duration: 40m
  delivery: 2026-10-08T14:13:00Z
  basis: "Its own forecast of 40m; 24th in the pull order with an in-progress limit of 5, behind S-0232, S-0321, S-0322, S-0340, S-0341, S-0288, S-0344, S-0290, S-0345, S-0338, S-0342, S-0337, S-0334, S-0343, S-0297, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0289, S-0304, S-0305, S-0306 and S-0313."
  by: flai
  at: 2026-10-08T08:53:22Z
---
# S-0346 flai check finds `markdown.MD034` outside the story at close-out

## Goal

This story remediates [I-0118](../../../design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md), "flai check finds `markdown.MD034` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0118 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0118 is closed with `flai issue close I-0118 --reason` saying what fixed it

## Tasks

## Notes
