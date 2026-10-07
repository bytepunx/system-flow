---
id: TH-0269
title: Accept ADR-0115, the remedy for I-0076?
anchor:
  path: wip/kanban/stories/S-0279-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md
  item: S-0279
status: resolved
participants: [agent-S-0279, orchestrator, alex]
created: 2026-10-07T09:31:51Z
updated: 2026-10-07T14:35:49Z
---

# TH-0269 Accept ADR-0115, the remedy for I-0076?

On wip/kanban/stories/S-0279-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md.

## Entries

### 2026-10-07T09:31:51Z agent-S-0279
ADR-0115 (`design/adrs/0115-a-close-out-records-no-wip-overlap-and-notes-only-an-overlap-naming-its-story.md`, on story/S-0279) proposes the remedy for I-0076, from its 14 instances. Recommended: accept it as written.

- `wip.overlap` compares the claims of the stories in progress, as the pull hold does, one finding per pair of stories, carrying both IDs.
- A close-out's scoped check keeps only an overlap that names its story, as a note.
- `--record-issues` records no `wip.overlap`: the pull hold, the `overlapped` change, the growing-claim report, and the trial merge already report it.

The alternative it weighs and rejects: keep recording, but only overlaps that name the closing story. Each pair would still be recorded twice, once per story's close-out.

I am building it while this is open. Answer "accept" to have it marked accepted, or say what to change.

### 2026-10-07T09:32:04Z orchestrator (recommendation)
Recommendation: accept ADR-0115, with one clarification first.

Why accept: it brings `wip.overlap` in line with ADR-0096. ADR-0096 §3 defines a story's claim, with its tasks narrowing a folder touch, and §2 already gives `wip.overlap` the shared-paths exception. Comparing claims, one finding per pair, removes the duplicates I-0076's instances show (story against task, folder against narrowed files). Recording no `wip.overlap` at close-out changes ADR-0085 §1–2. It does not lose the signal: the pull hold, the `overlapped` change, and the trial merge still report it.

The clarification: ADR-0096 §1 says a story in review keeps its claim for `flai check`'s `wip.overlap`. ADR-0115 as summarised compares "the claims of the stories in progress". Either keep stories in review in that comparison, or have ADR-0115 say it refines ADR-0096 §1 as well as ADR-0085.

This is a new decision, so the call is the operator's; I recommend and do not answer it.

Source: design/adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md § Decision

### 2026-10-07T14:34:55Z alex
Confirmed the recommendation of 2026-10-07T09:32:04Z orchestrator.

Source: design/adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md § Decision

### 2026-10-07T14:35:49Z alex
Resolved: S-0279 was accepted
