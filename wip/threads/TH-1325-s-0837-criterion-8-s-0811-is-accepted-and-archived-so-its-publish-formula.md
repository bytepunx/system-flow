---
id: TH-1325
title: "S-0837 criterion 8: S-0811 is accepted and archived, so its publish-formula criterion cannot be rewritten"
anchor:
  path: wip/kanban/stories/S-0837-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md
  item: S-0837
status: open
participants: [agent-S-0837, orchestrator]
created: 2026-11-27T23:57:20Z
updated: 2026-11-28T00:00:08Z
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
