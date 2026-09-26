---
id: S-0126
type: story
nature: remediation
title: "Threads don’t correctly render markdown"
status: done
parent: E-0003
owner: alex
created: 2026-09-26T07:30:00Z
updated: 2026-09-26T07:40:57Z
transitions:
  - to: ready
    at: 2026-09-26T07:30:10Z
    by: alex
  - to: in-progress
    at: 2026-09-26T07:30:44Z
    by: agent-S-0126
  - to: review
    at: 2026-09-26T07:34:19Z
    by: agent-S-0126
  - to: done
    at: 2026-09-26T07:40:57Z
    by: alex
tags: [dashboard]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0126 Threads don’t correctly render markdown

## Goal

When viewing any workflow item with threads, they should render the markdown.

## Acceptance criteria
- [x] Thread text is rendered markdown, not raw markdown.

## Tasks
- T-0459 Thread entries render as markdown on every page that shows threads

## Notes

Verified by `flaiover/src/lib/components/Threads.svelte.test.ts`, which renders the real markdown module in jsdom; not viewed in a running dashboard.
