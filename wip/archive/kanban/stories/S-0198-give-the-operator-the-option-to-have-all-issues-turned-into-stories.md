---
id: S-0198
type: story
nature: feature
title: Give the operator the option to have all issues turned into stories
status: done
owner: arobson
created: 2026-10-02T09:46:13Z
updated: 2026-10-03T02:43:34Z
transitions:
  - to: ready
    at: 2026-10-03T01:23:03Z
    by: alex
  - to: in-progress
    at: 2026-10-03T01:23:30Z
    by: agent-S-0198
  - to: review
    at: 2026-10-03T02:42:14Z
    by: agent-S-0198
  - to: done
    at: 2026-10-03T02:43:34Z
    by: alex
tags: [flai]
touches: [flai/internal/issues, flai/cmd/issue.go, flai/cmd/issue_test.go, flai/internal/check, flai/internal/manifest, flai/internal/mcpserver, flai/internal/harness, flai/internal/hostapi, flaiover/src/lib, flaiover/src/routes/api/issues, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, docs/operators/settings.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/system/project-manifest.md, design/system/flaiover-dashboard.md, design/conventions/continuous-improvement.md, template/root/design/conventions/continuous-improvement.md, design/system/conventions.md, docs/users/conventions.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4756
  models:
    - model: claude-haiku-4-5-20251001
      input: 330
      output: 8887
      cache_read: 2259722
      cache_write: 91509
      cost: 0.3851
    - model: claude-opus-5-5
      input: 646
      output: 217014
      cache_read: 37086845
      cache_write: 760146
      cost: 16.2862
    - model: claude-sonnet-5-5
      input: 64
      output: 11139
      cache_read: 1139406
      cache_write: 146939
      cost: 0.7067
---
# S-0198 Give the operator the option to have all issues turned into stories

## Goal

In a story's review page, provide the operator with the opportunity to create stories for any open issues that do not already have one. The planning agent creates a `remediation` or `improvement` story in the backlog where none exists, linking the issue and recommending a solution. This replaces presenting the summary at an epic's end. Ensure that `design/conventions/continuous-improvement.md` and `design/system/continuous-improvement.md` say so.

## Acceptance criteria
- [x] flai records which story recorded or bumped each issue instance (the story the agent works, from `FLAI_AGENT` or the stream), so it can be listed per story
- [x] Viewing a story's review page, the operator has the option (via check box) to create remediation or improvement story for each open issue that does not have a story.
- [x] `flai issue` can create the remediation or improvement story for an issue (`flai issue story I-nnnn`, or the MCP equivalent), with the issue linked, its nature from the issue's class, and its last criterion closing the issue
- [x] `flai check` warns about an open issue with no open story that links it, older than a threshold the project sets, so issues from before this change are found too
- [x] `harness.Prompt`, the user guide, and `design/system/flai-cli.md` describe it
- [x] Tests cover the option at review, the story created from an issue, and the check

## Tasks
- T-0730 Each issue instance records its story, and an issue can be turned into a story's title, nature, and body
- T-0731 flai issue records the story it runs for, lists issues per story, and makes a story from an issue
- T-0732 flai check warns about an open issue older than the project's threshold that no open story links
- T-0733 The MCP tool issue_story makes a story from an issue
- T-0734 A story's review page offers a story for each open issue no open story links, by checkbox
- T-0735 The prompt, the continuous-improvement convention, and its design say how issues become stories

## Notes

Decided by the designer on 2026-10-02 while reviewing their convention edits.

The designer refined the review page on TH-0074. The story's own open issues come first, checked, then every other open issue with no open story, unchecked. With any checked, Accept reads `Accept and Create Stories`, and it makes the stories only once the acceptance succeeds. The check's threshold is `issues.story_after`, `168h` when unset. S-0243 was made from I-0007 so that this repository's check passes. On TH-0075 the designer confirmed that the story's agent no longer makes stories for its own issues before review: the operator chooses at acceptance. `flai issue story` writes nothing to the issue file. The story's link to the issue is the record, so nothing outside `wip/` is left uncommitted.
