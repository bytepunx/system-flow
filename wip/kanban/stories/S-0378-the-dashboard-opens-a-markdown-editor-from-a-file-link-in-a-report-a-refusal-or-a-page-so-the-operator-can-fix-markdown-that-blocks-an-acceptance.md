---
id: S-0378
type: story
nature: feature
title: The dashboard opens a markdown editor from a file link in a report, a refusal, or a page, so the operator can fix markdown that blocks an acceptance
status: backlog
owner: alex
created: 2026-10-12T17:07:24Z
updated: 2026-10-12T17:07:24Z
transitions: []
tags: [dashboard]
topics: [markdown]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0378 The dashboard opens a markdown editor from a file link in a report, a refusal, or a page, so the operator can fix markdown that blocks an acceptance

## Goal

When an acceptance is refused because markdown the lint fix cannot fix still has findings (S-0367, TH-0413), the operator opens the file from the link the refusal names, corrects it in a markdown editor in the dashboard, saves it through flai, and accepts again.

## Acceptance criteria

- [ ] A repository file link shown inline in a report, a refusal, a thread, or a page opens the file in a markdown editor in the dashboard, at the line the link names when it names one.
- [ ] The editor shows the markdown lint findings of the file as it is edited, with their lines and rules.
- [ ] Saving writes the file through flai on the host (ADR-0023), on the branch or checkout the link names, and refuses a save that conflicts with a change made meanwhile.
- [ ] After the findings that blocked an acceptance are corrected and saved, the acceptance can be made again from the same page.

## Tasks

## Notes

- Asked for by the operator on TH-0413, as a story of its own for the planner, beside S-0367, which makes `flai accept` refuse with the files, lines, and errors of markdown the fix cannot fix.
