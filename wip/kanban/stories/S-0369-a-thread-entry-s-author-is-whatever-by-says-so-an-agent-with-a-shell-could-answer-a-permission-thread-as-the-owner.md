---
id: S-0369
type: story
nature: improvement
title: A thread entry's author is whatever --by says, so an agent with a shell could answer a permission thread as the owner
status: backlog
owner: alex
created: 2026-10-12T16:52:42Z
updated: 2026-10-12T16:52:42Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0369 A thread entry's author is whatever --by says, so an agent with a shell could answer a permission thread as the owner

## Goal

This story remediates [I-0091](../../../design/issues/I-0091-a-thread-entry-s-author-is-whatever-by-says-so-an-agent-with-a-shell-could-answer-a-permission-thread-as-the-owner.md), "A thread entry's author is whatever --by says, so an agent with a shell could answer a permission thread as the owner". The issue recommends this solution:

Directions to weigh: flai refuses a `--by` that is not the session's own name when `FLAI_AGENT` is set, so that an agent's entries always carry its name; and `flai guard` refuses a story's agent `flai thread reply`, `resolve`, and `confirm` with a `--by` at all. Either leaves the operator's own session, where `FLAI_AGENT` is unset, free to name itself. The dashboard already writes the operator's name itself.

## Acceptance criteria
- [ ] The cause I-0091 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0091 is closed with `flai issue close I-0091 --reason` saying what fixed it

## Tasks

## Notes
