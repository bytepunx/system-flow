---
id: S-0144
type: story
nature: remediation
title: Automatic publish/push should be a setting
status: done
owner: alex
created: 2026-09-29T01:04:15Z
updated: 2026-09-29T02:07:25Z
transitions:
  - to: ready
    at: 2026-09-29T01:10:17Z
    by: alex
  - to: in-progress
    at: 2026-09-29T01:13:58Z
    by: agent-S-0144
  - to: review
    at: 2026-09-29T02:06:59Z
    by: agent-S-0144
  - to: done
    at: 2026-09-29T02:07:25Z
    by: alex
tags: [cli]
touches: [flai/cmd/push.go, flai/cmd/push_test.go, flai/cmd/release.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/writes.go, flai/internal/pending, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/conventions/git.md, design/adrs, design/issues, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0144 Automatic publish/push should be a setting

## Goal

At some point moving stories from review to done started automatically publishing without user involvement, breaking the operator's ability to batch stories for release. This needs to be an explicit setting that the operator enables, not something it defaults to.

## Acceptance criteria
- [x] By default, enabling push in the server settings does _not_ mean that stories are automatically published/pushed
- [x] A new setting that enables automatic publishing should be added so that controlling that behavior is explicit

## Tasks
- T-0514 flai push --pending tags a release only when the publish host action is enabled
- T-0515 Record the decision and document the publish setting, then run every check
- T-0516 Rename the publish host action to auto-publish, as TH-0031 decided

## Notes

- Cause: since S-0094 `flai push --pending` tagged the pending release before every push, and agents push after each acceptance, so each acceptance was released on its own.
- The setting is the `auto-publish` host action, off by default (`flai serve enable auto-publish`, or its toggle in the dashboard's settings): ADR-0048, the name and option A decided by the operator in TH-0031. Agents still push merged commits after an acceptance; that cuts no release.
- Criterion 1 verified by `TestHostActionPush` (push action on, Push now pushes and tags nothing) and `TestPushPendingReleasesNothingByDefault`; criterion 2 by `TestPushPendingTagsWhatAcceptLeftUnreleased` and `TestPushPendingAutoPublishesAnAcceptanceAlreadyPushed`.
- Found and fixed on the way (TH-0031's question): `flai release --pending` after the acceptances were pushed tagged a Go component and left the tag local, since nothing on the branch was ahead. It now pushes such tags on their own.
- `flai check --strict` on the repository reports warnings this story does not clear: `wip.overlap` with S-0138 (docs/users, design/system/flai-cli.md; clears when S-0138 is accepted) and `threads.archived` for TH-0029 on S-0137.
