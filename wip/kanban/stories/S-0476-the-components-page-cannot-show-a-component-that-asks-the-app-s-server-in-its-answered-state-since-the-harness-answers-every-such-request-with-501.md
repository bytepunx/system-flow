---
id: S-0476
type: story
nature: improvement
title: The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501
status: backlog
owner: alex
created: 2026-10-21T11:59:18Z
updated: 2026-10-22T20:52:34Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
finalized:
  by: alex
  at: 2026-10-22T20:52:34Z
---
# S-0476 The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501

## Goal

This story remediates [I-0164](../../../design/issues/I-0164-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md), "The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0164 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0164 is closed with `flai issue close I-0164 --reason` saying what fixed it

## Tasks

## Notes

I-0164's cost "0m" is not a duration longer than zero, so no time_lost_per_cycle was derived from it.
