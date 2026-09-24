---
id: TH-0009
title: Commit trailer names a fixed model
anchor:
  path: design/conventions/git.md
status: resolved
participants: [system-flow, alex]
created: 2026-09-24T06:28:11Z
updated: 2026-09-24T07:49:54Z
---

# TH-0009 Commit trailer names a fixed model

On design/conventions/git.md.

## Entries

### 2026-09-24T06:28:11Z system-flow
The project addition in git.md says commit messages end with `Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>`. S-0112 was worked by Claude Opus 5.5, and its harness's attribution names that model, so its commits carry the Opus 5.5 trailer. Proposal (recommended): change the addition to "the model that authored the change", so that the trailer is never false. Alternative: keep Fable 5.1 fixed, and agents on other models override it. S-0112 is not waiting on this.

### 2026-09-24T07:48:49Z alex
I agree, the model that did the work should be the one listed on the commit message.

### 2026-09-24T07:49:36Z system-flow
Done on story/S-0113 as 43f2765 (T-0396). The project addition in design/conventions/git.md now says the trailer names the model that authored the change, with Opus 5.5 given as an example. It lands on main when S-0113 is accepted. This was S-0113's agent: S-0112's agent opened the thread, and it reached me because both agents share one name (I-0037).

### 2026-09-24T07:49:54Z system-flow
Resolved: The operator decided that the trailer names the model that authored the change. It is applied on story/S-0113 as 43f2765.
