---
id: TH-1325
title: "S-0837 criterion 8: S-0811 is accepted and archived, so its publish-formula criterion cannot be rewritten"
anchor:
  path: wip/kanban/stories/S-0837-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md
  item: S-0837
status: open
participants: [agent-S-0837, orchestrator]
created: 2026-11-27T23:57:20Z
updated: 2026-11-28T01:06:41Z
---

# TH-1325 S-0837 criterion 8: S-0811 is accepted and archived, so its publish-formula criterion cannot be rewritten

On wip/kanban/stories/S-0837-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md.

## Entries

### 2026-11-27T23:57:20Z agent-S-0837
S-0837's criterion 8 asks that S-0811's criterion naming the publish formula be rewritten to speak of `release_publish`. S-0811 was accepted and is archived (`wip/archive/kanban/stories/S-0811-…md`), with its criterion 2 ticked: "The publish formula stops on that code and its story says which suite failed." The archive is not edited, and S-0837's Notes say archived publish stories stay as they are.

Recommendation: leave S-0811 as history, and tick criterion 8 as met by this answer. What S-0811's criterion described moves to `release_publish`, which already answers exit 5 as a `Refused` error naming the tier and suite (S-0811 built that); ADR-0316's Consequences record the change. I will tick 8 once you confirm.

Alternative: you edit the archived S-0811 yourself to read "`release_publish` and the dashboard's Publish refuse on that code, naming the suite and tier that failed".

The rest of S-0837 does not wait on this.

### 2026-11-28T00:00:08Z orchestrator (recommendation)
**Recommended:** leave the archived S-0811 as history, and tick S-0837's criterion 8 once ADR-0316 records that choice.

- **The archive is not edited.** S-0811 is under `wip/archive`. No agent, me included, rewrites an archived item's words.
- **ADR-0316 says otherwise today.** Its Consequences say "S-0811's criterion about the publish formula's step is rewritten to speak of `release_publish`". The ADR is still `proposed`, and T-3600 accepts it on S-0837's branch. Before it is accepted, change that line to say that S-0811 stays as archived and that `release_publish` takes over the refusal S-0811's criterion 2 described. The same section already says the 51 archived publish stories "stay as they are".
- **Then criterion 8 is met.** The ADR is where the change is recorded, as the criterion's "by a thread" path allows.

**Why I recommend rather than answer:** criterion 8 asks for a rewrite, and this replaces it with a record in the ADR. That changes how the operator's criterion is met, so the call is yours.

**Alternative:** you edit the archived S-0811's criterion 2 yourself, as agent-S-0837's entry words it.

Source: design/adrs/0316-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md § Consequences

### 2026-11-28T01:05:32Z agent-S-0837
A correction to the orchestrator's recommendation: T-3600 already accepted ADR-0316 on S-0837's branch, so its Consequences line can no longer be changed. An accepted ADR changes only in `superseded_by` and `topics`.

The record is already made in the living design. On S-0837's branch, `design/system/strategic-agents.md` § Publishing a release says that `release_publish` refuses on exit 5, naming the suite and tier. That is the refusal S-0811's criterion 2 described.

Answer with one of:

1. **Leave S-0811 as history** (recommended). I tick criterion 8 on this answer and move S-0837 to review. Its close-out passed at d99e197, and criteria 1 to 7 are ticked.
2. **You edit the archived S-0811's criterion 2** to read "`release_publish` and the dashboard's Publish refuse on that code, naming the suite and tier that failed". I tick 8 once that is done.

No other work on S-0837 waits on this.

### 2026-11-28T01:06:41Z orchestrator (recommendation)
**Recommended:** option 1. Leave the archived S-0811 as history. agent-S-0837 then ticks criterion 8 on your answer and moves S-0837 to review.

This replaces my entry of 00:00:08Z. agent-S-0837's correction is right. ADR-0316 is `accepted` on `story/S-0837`, whose head is d99e197, so its Consequences line can no longer be changed.

**Why option 1:**

- **Nobody may rewrite S-0811.** It is under `wip/archive`, and archived items are not edited. I edit no documents, so option 2 is yours alone, and it would be a hand edit to the archive.
- **ADR-0316 already covers S-0811's case.** The same Consequences section says the 51 archived publish stories "stay as they are". S-0811 is archived work of the same kind.
- **Nothing that S-0811's criterion 2 described is lost.** On S-0837's branch, `design/system/strategic-agents.md` § Publishing a release says that `release_publish` returns exit 5 from a failed `verify_at.publish` suite as flai's message. That refusal now lives in the living design.
- **The ADR's line becomes history, and nothing more.** The line "S-0811's criterion … is rewritten" will not be carried out. S-0837's Notes should say so, with this thread as the reason, so a reader of ADR-0316 can see why.

**Why I recommend rather than answer:** criterion 8 asks for a rewrite, and this drops the rewrite. That changes how your criterion is met, so the call is yours.

**Alternative:** option 2. You edit the archived S-0811's criterion 2 yourself, in the words agent-S-0837 gave.

Source: design/adrs/0316-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md § Consequences
