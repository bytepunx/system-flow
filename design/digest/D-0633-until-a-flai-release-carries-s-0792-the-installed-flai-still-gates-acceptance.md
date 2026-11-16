---
id: D-0633
title: "Until a flai release carries S-0792, the installed flai still gates acceptance by ADR-0291, so a digest or docs commit on main sends a story in review back; its agent syncs and closes out again, changing no code"
kind: cause
sources: [S-0792, TH-1220, I-0378]
topics: [acceptance, verify, flai]
paths: [flai/internal/preview/accept.go, flai/cmd/accept.go, design/adrs/0303-an-acceptance-s-verify-behind-main-sends-back-only-the-tiers-selected-by-the.md]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: [TH-1237]
---

# D-0633 Until a flai release carries S-0792, the installed flai still gates acceptance by ADR-0291, so a digest or docs commit on main sends a story in review back; its agent syncs and closes out again, changing no code

S-0792 narrowed `verify_behind_main` to the main paths that overlap a story's change set ([ADR-0303](../adrs/0303-an-acceptance-s-verify-behind-main-sends-back-only-the-tiers-selected-by-the.md)). But acceptances run the installed flai, and S-0792 was accepted on 2026-11-15 without a release. So S-0792 itself was sent back on TH-1220 for two digest files main brought in, `D-0627` and `D-0632`, which select the markdown tier. Its narrative's Decisions record the outcome:

> TH-1220: no change to the code. The installed flai that ran the acceptance predates ADR-0303, so its gate sent the markdown tier back for two digest files the story does not change. Until a flai release carries S-0792, acceptances keep the old gate. The fix is to sync and close out again.

So a story sent back like this before that release does what TH-1220 asks: run `flai stream sync`, ask for the close-out again with the MCP tool `verify`, and move back to review. It does not change its code, and it does not ask the operator. Once a release carries S-0792 and the host runs it, this entry no longer holds.

Source: S-0792's narrative, `wip/archive/agents/S-0792.md`, § Decisions; TH-1220; I-0378.
