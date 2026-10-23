---
id: T-1813
type: task
nature: feature
title: One rule in flai/internal/digest turns a thread into its operator's answers and overruled recommendations, each a candidate entry with its topics and paths
status: backlog
parent: S-0429
owner: alex
created: 2026-10-16T06:03:49Z
updated: 2026-10-23T16:42:53Z
transitions: []
stream: S-0429
tags: [cli]
touches: [flai/internal/digest/answer.go, flai/internal/digest/answer_test.go]
---
# T-1813 One rule in flai/internal/digest turns a thread into its operator's answers and overruled recommendations, each a candidate entry with its topics and paths

## Work

Add `digest.Answers(thread, owner, items)` in `answer.go`: for a thread someone other than the operator (TH-0476 question 5) opened, each entry by the operator that answers it becomes a candidate.

- An `answer`: the operator's reply, or their confirmation of a recommendation, which takes the recommendation's text as the answer.
- A `recommendation`: an operator's reply that follows a pending recommendation without confirming it, carrying both texts, so the entry says what was overruled and why.
- Each candidate's title is the thread's question (its title) and the answer's first sentence; its source is the thread and the entry's time; its topics and paths are those of the thread's anchored item (its topics, tags through the manifest, and touches), none when it is anchored on a path.

S-0433 calls the same function at the moment the operator answers, so keep it free of I/O beyond what it is given.

First layer: runs beside the narrative and issue sources task, which touches other files.

## Done when

`answer_test.go` covers a reply, a confirmation, an overruled recommendation, a thread resolved with no answer, a thread whose only replies are agents', and the operator's own thread; `flai test flai/internal/digest` passes.

## Notes

`threads.PendingRecommendation` and the confirmation line ADR-0090 writes tell a confirmation from a reply.

Planner-S-0429, 2026-10-23: in `flai/internal/threads/threads.go`, `Thread.Entries()` gives each entry's `At`, `Author`, `Recommendation`, and `Source`; `Entry.Confirms()` tells a confirmation, written as `Confirmed the recommendation of <At> <author>.`; `Thread.Opener()` names the opener, and `Thread.Anchor` holds `Item`, `Path`, and `Heading`. An overruled recommendation is an entry by the operator, after a recommendation by someone else, that does not confirm it; `PendingRecommendation` returns nil then, so read the entries in order. The operator is `manifest.Manifest.Owner`. A story's topics come from `topics.OfStory` in `flai/internal/topics/story.go`, with its epic and tasks; `threads.StoryOf` gives the story of a task anchor; an epic anchor takes the epic's own topics and its tags through the manifest.
