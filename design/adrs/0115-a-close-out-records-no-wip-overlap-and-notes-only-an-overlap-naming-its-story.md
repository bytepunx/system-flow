---
id: ADR-0115
title: "A close-out records no wip.overlap and notes only an overlap naming its story, and wip.overlap compares the claims of the stories in progress once per pair"
status: proposed
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0096]
---

# ADR-0115 A close-out records no wip.overlap and notes only an overlap naming its story, and wip.overlap compares the claims of the stories in progress once per pair

## Context

[ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) scoped a close-out's `flai check` to its story and recorded each rule's findings outside the story in one issue per rule. It put every `wip.overlap` outside the story, because the pull hold and the other story's agent clear it. So I-0076, ``flai check finds `wip.overlap` outside the story at close-out``, was bumped 14 times between 2026-10-05 and 2026-10-07. Its instances show three causes.

1. **Every close-out records every overlap on the board.** A close-out records the overlaps of every pair of stories in progress, whether or not the closing story is one of them. S-0251 recorded S-0228 against S-0273; S-0264 recorded S-0299 against S-0301; S-0307 recorded S-0269 against S-0294. Both stories of one pair record it again at their own close-outs: S-0253 and S-0260 each recorded the same pair, as did S-0262 and S-0257. An overlap is the state of the board, not a defect. The pull hold ([ADR-0046](0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md)), the `overlapped` change in `inbox`, the report when a write grows a claim (I-0059), and the trial merge at `flai stream sync` already report it to the agents who can act.
2. **A story is reported against another story's task as well as against the story.** `wip.overlap` pairs every item in progress, tasks included, and skips only an item and its own parent. So S-0262 was reported against S-0257 and again against T-0877; S-0273 against S-0228 and again against T-0917 and T-0938.
3. **`wip.overlap` compares raw touches, not the claim.** The pull hold compares claims, in which a story's folder touch is narrowed to the touches its tasks name inside it ([ADR-0096](0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md) §3). `wip.overlap` compares each item's raw touches, so a folder is reported against a file inside it that the folder's tasks may never name: S-0223's `flai/internal/serve` against S-0296's `flai/internal/serve/review_wait_test.go`, S-0228's `flai/internal/hostapi` against three files of S-0273, and S-0286's `flai/internal/harness` against S-0272's files.

The question is what a close-out does with an overlap, and what `wip.overlap` compares.

## Decision

A close-out records no `wip.overlap` in an issue and notes only an overlap that names its story, and `wip.overlap` compares the claims of the stories in progress, one finding per pair of stories, as the pull hold does.

1. **`wip.overlap` compares claims, once per pair of stories.** Each story in progress is compared by its claim, as `Holds.Claim` reads it from the story and its tasks (ADR-0096 §3), with every other story in progress. A pair whose claims overlap where it counts, as `Holds.Overlaps` says, outside the manifest's shared paths, is one warning, on the item file of the story with the lower ID, at its `touches` line. Its message starts with that story's ID and names each overlapping path of both claims. A task is never compared on its own: its touches are part of its story's claim. A story in review, or in any other state, is not compared, as it holds nothing (ADR-0096 §1). The finding carries the two story IDs in the `Finding` field `Stories`, printed in `--json` as `stories`, so that a reader scopes it without parsing the message. `check.Overlaps`, which the designer's inbox lists, and the comparison a write makes before and after it (I-0059) read the same findings.
2. **A check scoped to a story keeps only the overlaps that name it, as notes.** `flai check --story S-nnnn` drops from its result every `wip.overlap` whose `Stories` do not name the story, and keeps one that names it as a note outside the story, as ADR-0085 has every `wip.overlap`: printed with `(outside S-nnnn)`, counted in `outside`, failing neither the run nor `--strict`. Without `--story` every `wip.overlap` is reported, as before.
3. **`--record-issues` records no `wip.overlap`.** The findings of every other rule outside the story are recorded as ADR-0085 says. An overlap is reported where it can be acted on: by the pull hold, the `overlapped` change, the report on a growing claim, and the trial merge. An issue would add a count of how often stories ran side by side, which is not friction.

This refines ADR-0085, whose rule that every `wip.overlap` is outside the story becomes: only one that names the story is kept, and none is recorded. It refines ADR-0096, whose claim `wip.overlap` now reads as the pull hold does.

## Consequences

- A close-out no longer bumps an issue because two other stories run side by side, and I-0076 can be closed: its cause no longer occurs.
- An agent still sees at its close-out, once, each story in progress whose claim overlaps its own, so it can sync and coordinate.
- `wip.overlap` agrees with the pull hold. A folder touch that the story's tasks narrow is no longer reported against files the tasks do not name; a task's touch outside its story's touches is reported as part of its story's claim.
- The designer's inbox lists one overlap per pair of stories instead of one per pair of paths and items. Its entry is keyed by the message, which changes when either claim gains or loses an overlapping path, so the entry then shows again as new.
- A write that grows a claim into a story it already overlapped makes the pair's message new, so the write reports it as an advisory warning, as it reports a new overlap.
- Nothing counts how often stories overlap now that no issue does. The pull hold's figures in `flai stats` ([ADR-0081](0081-flai-stats-reports-forecast-error-cost-of-delay-waiting-holds-touches-drift-and.md)) measure what overlaps cost.
- An open issue for `wip.overlap` recorded by an older flai stays as it is until it is closed.

## Alternatives considered

- **Keep recording, but only the overlaps that name the closing story.** It removes cause 1's instances from stories outside the pair, but each pair is still recorded twice, once at each story's close-out, and the issue still counts concurrency rather than friction: its remediation would have nothing to fix.
- **Drop `wip.overlap` from a scoped check altogether.** The closing story's agent would no longer see, at its close-out, the story in progress it may conflict with when it is accepted.
- **Keep comparing raw touches and only stop pairing a story with another story's task.** It removes cause 2's duplicates, but `wip.overlap` would still disagree with the pull hold about folder touches, which is cause 3.
- **Compare stories in review as well.** A story in review holds nothing (ADR-0096 §1); its conflicts are for the acceptance's trial merge to find.
