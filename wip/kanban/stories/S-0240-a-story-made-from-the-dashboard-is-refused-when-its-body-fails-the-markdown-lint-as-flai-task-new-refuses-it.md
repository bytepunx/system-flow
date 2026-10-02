---
id: S-0240
type: story
nature: remediation
title: A story made from the dashboard is refused when its body fails the markdown lint, as flai task new refuses it
status: in-progress
owner: arobson
created: 2026-10-02T12:42:51Z
updated: 2026-10-02T16:50:12Z
transitions:
  - to: ready
    at: 2026-10-02T16:24:25Z
    by: alex
  - to: in-progress
    at: 2026-10-02T16:48:17Z
    by: agent-S-0240
tags: [dashboard, cli]
topics: [dashboard, markdown]
touches: [flai/internal/mdlint, flai/internal/itemnew, flai/internal/hostapi, flaiover/src/lib/server, design/issues, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0240 A story made from the dashboard is refused when its body fails the markdown lint, as flai task new refuses it

## Goal

[I-0055](../../../design/issues/I-0055-a-story-made-from-the-dashboard-reached-main-with-a-markdown-lint-error-flai-task-new-would-have-refused.md): S-0231, made from the dashboard on 2026-10-02, reached `main` with its criteria list indented one space (MD007), which stopped S-0193's close-out at the markdown lint. `flai task new` and `flai story new` refuse such a body with the rule and the line (ADR-0061), so the dashboard's path for a new item, or the edit that follows it, writes a body without that check. Find where, and close the gap.

## Acceptance criteria
- [ ] The path the dashboard's new-item form takes to write a story's or epic's body is traced, and the step that writes without the wip markdown lint is named in the story's notes.
- [ ] Making or editing an item from the dashboard with a body the lint rejects is refused, and the form shows the rule and the line as the CLI does.
- [ ] A test on the side that was missing the check reproduces S-0231's body and fails without the fix.
- [ ] I-0055 is closed with `flai issue close`.

## Tasks
- T-0708 mdlint reports MD007 ul-indent as markdownlint does
- T-0709 The dashboard's new-item and save paths refuse S-0231's body naming MD007 and its line
- T-0710 The design names MD007 among the rules flai lints and I-0055 is closed

## Notes
