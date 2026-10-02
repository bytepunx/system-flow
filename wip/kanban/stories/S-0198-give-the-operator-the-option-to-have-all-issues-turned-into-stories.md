---
id: S-0198
type: story
nature: feature
title: Give the operator the option to have all issues turned into stories
status: backlog
owner: arobson
created: 2026-10-02T09:46:13Z
updated: 2026-10-02T10:03:23Z
transitions: []
tags: [flai]
touches: [flai/internal/issues, flai/cmd/issue.go, flai/internal/check, flai/internal/mcpserver, flai/internal/harness]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0198 Give the operator the option to have all issues turned into stories

## Goal

In a story's review page, provide the operator with the opportunity to create stories for any open issues that do not already have one. The planning agent creates a `remediation` or `improvement` story in the backlog where none exists, linking the issue and recommending a solution. This replaces presenting the summary at an epic's end. Ensure that `design/conventions/continuous-improvement.md` and `design/system/continuous-improvement.md` say so.

## Acceptance criteria
- [ ] flai records which story recorded or bumped each issue instance (the story the agent works, from `FLAI_AGENT` or the stream), so it can be listed per story
- [ ] Viewing a story's review page, the operator has the option (via check box) to create remediation or improvement story for each open issue that does not have a story.
- [ ] `flai issue` can create the remediation or improvement story for an issue (`flai issue story I-nnnn`, or the MCP equivalent), with the issue linked, its nature from the issue's class, and its last criterion closing the issue
- [ ] `flai check` warns about an open issue with no open story that links it, older than a threshold the project sets, so issues from before this change are found too
- [ ] `harness.Prompt`, the user guide, and `design/system/flai-cli.md` describe it
- [ ] Tests cover the option at review, the story created from an issue, and the check

## Tasks

## Notes

Decided by the designer on 2026-10-02 while reviewing their convention edits.
