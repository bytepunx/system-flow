---
id: T-1141
type: task
nature: improvement
title: An ADR refining ADR-0085 and ADR-0096 records the remedy for I-0076, proposed from its instances
status: backlog
parent: S-0279
owner: alex
created: 2026-10-07T01:20:06Z
updated: 2026-10-07T01:20:06Z
transitions: []
stream: S-0279
tags: [flai]
touches: [design/adrs]
---
# T-1141 An ADR refining ADR-0085 and ADR-0096 records the remedy for I-0076, proposed from its instances

## Work

The story asks for a remedy proposed from I-0076's instances before it is built. This task proposes it and records it as an ADR with `flai adr new`, refining ADR-0085 and ADR-0096. It waits for nothing. The code tasks build what this ADR decides, so they wait for it.

The instances show three causes:

1. **Every close-out records every overlap.** `ScopeToStory` in `flai/internal/check/scope.go` marks each `wip.overlap` outside the story, as ADR-0085 says. `--record-issues` then bumps I-0076 at every close-out while any two stories in progress overlap. That happens even when the closing story is in neither pair: S-0251 recorded S-0228 against S-0273. An overlap is the state of the board, not a defect. The pull hold, the `overlapped` change, and the trial merge at `flai stream sync` already report it.
2. **A story is reported against another story's task as well as against the story.** `overlap()` in `flai/internal/check/check.go` pairs every in-progress item. It skips only an item and its own parent. So "S-0262 against T-0877" and "S-0273 against T-0917" repeat the story-to-story pair.
3. **wip.overlap reads raw touches, not the claim.** The hold reads the claim (`Holds.Claim`), where the tasks narrow a folder touch (ADR-0096 §3). `overlap()` reads each item's raw touches. So a folder such as `flai/internal/serve`, `flai/internal/hostapi`, or `flai/internal/harness` is reported against files that the folder's tasks may not name.

The recommended remedy:

- A close-out records no `wip.overlap` in an issue.
- A scoped run lists as notes only the overlaps that name the story.
- `wip.overlap` compares the claims of the stories in progress, one finding per pair of stories, as the pull hold does. It keeps the shared-path exception.

Weigh it against the alternative: keep recording, and record only the overlaps that name the closing story. Say in the ADR why you chose one over the other. Open a thread on S-0279 for the operator to accept the ADR, and go on while it is open.

## Done when

- A new ADR under `design/adrs/` names the three causes, with the instances that show each. It states the remedy and refines ADR-0085 and ADR-0096. It is indexed in `design/adrs/README.md`, as `flai adr new` writes it.
- A thread on S-0279 asks the operator to accept it.
- `flai check --strict` and the markdown lint pass on the ADR.

## Notes
