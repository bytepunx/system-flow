---
id: TH-0009
title: Commit trailer names a fixed model
anchor:
  path: design/conventions/git.md
status: open
participants: [system-flow]
created: 2026-09-24T06:28:11Z
updated: 2026-09-24T06:28:11Z
---

# TH-0009 Commit trailer names a fixed model

On design/conventions/git.md.

## Entries

### 2026-09-24T06:28:11Z system-flow
The project addition in git.md says commit messages end with `Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>`. S-0112 was worked by Claude Opus 5.5, and its harness's attribution names that model, so its commits carry the Opus 5.5 trailer. Proposal (recommended): change the addition to "the model that authored the change", so that the trailer is never false. Alternative: keep Fable 5.1 fixed, and agents on other models override it. S-0112 is not waiting on this.
