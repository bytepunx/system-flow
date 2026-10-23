---
id: ADR-0188
title: "A task takes no folder touch, the planner's folder touch on a story needs the story's Planning note to name it, and flai check reports a story in progress whose claim keeps a folder and a story three or more open stories name in after"
status: accepted
date: 2026-10-23
supersedes: []
superseded_by: []
refines: [ADR-0046, ADR-0096]
topics: [planning, board, touches]
---

# ADR-0188 A task takes no folder touch, the planner's folder touch on a story needs the story's Planning note to name it, and flai check reports a story in progress whose claim keeps a folder and a story three or more open stories name in after

## Context

[ADR-0046](0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md) holds a ready story while its claim overlaps the claim of a story in progress, or while a story it names in `after` is not done. [ADR-0096](0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md) § 3 narrows a story's folder touch in its claim to the touches of its tasks inside it. So a folder holds only while no task names a file in it. [ADR-0134](0134-flai-serve-asks-the-agent-of-a-story-in-progress-about-each-ready-story-it.md) lets the holding agent share overlapping paths with the story it holds. [I-0087](../issues/I-0087-one-in-progress-story-holds-every-ready-story-by-overlap-through-folder-wide-touches-and-docs-users-flai-md-so-the-board-runs-one-story-at-a-time-under-a-limit-of-three.md), closed, recorded the first form of the problem, through `docs/users/flai.md` and folder touches.

[I-0152](../issues/I-0152-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md) records the form that came after, through the planner. The analysis in [2026-10-15-bottlenecks.md](../analysis/2026-10-15-bottlenecks.md) § Finding 3 found:

- E-0016's planned stories, S-0207 to S-0211 and S-0220 to S-0229, declared folder touches: `flai/cmd`, `flai/internal/{harness,hostapi,mcpserver,serve,guard,manifest}`, `flaiover/src/routes`, and `docs/users`.
- S-0218 was an `after` hub that six stories named. S-0223 was another, for S-0224 and S-0227.
- One story in progress with a folder touch, or one hub awaiting acceptance, held up to ten ready stories a day: 9 on 2026-10-03, 9 on 2026-10-05, and 10 on 2026-10-06, while in progress stood at 0, 1, and 2 of a limit of 5.
- W41 held ready stories for 1,577,991 s (438 h) in all.
- `exact_share` was 0.49 to 0.60. Half the completed stories changed files outside their touches or left declared touches unchanged.

Narrowing works only when tasks name files. The convention and `planPrompt` ask the planner and the story's agent to declare touches file by file, but nothing enforces it. A task may name a folder, and then it narrows nothing. A planner's folder touch on a story stands with no reason given. Nothing reports a story in progress that keeps a folder, or a story that many others wait for in `after`. Shares cut the hold in W42, 179,046 s avoided, but each share costs the holding agent a turn per pair.

The question: what flai refuses and reports so that a claim names the files a story changes, and `after` names only a change a story needs.

## Decision

**A task takes no folder touch. The planner's folder touch on a story needs the story's `### Planning` note to name it. `flai check` warns on a story in progress whose claim keeps a folder, and on a story that three or more open stories name in `after`.** This refines ADR-0046 and ADR-0096. The claim, the hold, its reasons, and shares stay as they are.

1. **A folder touch.** A touches entry is read as a path, as a claim reads it: a component's name or tag means its sub-project's path. It is a folder touch when it names a folder of the repository tree holding at least one file. A path absent from the tree is a file to be made. A folder lying wholly inside a `claims.shared` pattern, by ADR-0096's rule, does not count: it holds nothing.
2. **A task takes no folder touch.** `flai task new --touches`, `flai edit --touches`, `flai touches`, and the MCP tools `item_new` and `item_edit` refuse one on a task, writing nothing. The refusal names the folder and how many files it holds. It says to name the files the task changes, or to leave the folder to the story's touches with a `### Planning` note. `flai task done` widens a task's touches with the files its commit changed ([ADR-0144](0144-flai-task-done-widens-touches-only-with-the-committed-paths-the-story-s-branch.md)), never a folder, and is unchanged.
3. **The planner's folder touch on a story needs its note.** The writer is the planner when `FLAI_ROLE=plan` is in the environment, or else when the author's name starts with `planner-`. A folder touch it writes on a story, through `flai edit --touches`, `flai touches`, `item_new`, or `item_edit`, is refused, writing nothing, unless the story's Notes name the folder's path under `### Planning`. The note is read as the body will be after the write, so one call may add both. The operator's and the story's agent's story touches are not refused.
4. **`touches.folder`.** `flai check` gives an advisory warning, which `--strict` passes over, on the story's file, on a story in progress whose claim, ADR-0096's with its tasks' narrowing, keeps a folder touch. The finding names the folder, its file count, and the ready stories it holds. A check scoped to a story keeps the finding on that story and drops it on another story, recording no issue, as a scoped check drops a `wip.overlap` between two other stories ([ADR-0115](0115-a-close-out-records-no-wip-overlap-and-notes-only-an-overlap-naming-its-story.md)).
5. **`wip.after-hub`.** `flai check` gives an advisory warning, which `--strict` passes over, on the hub's file, on a story not done or cancelled that three or more stories in backlog, ready, or in progress name in `after`, naming them. A check scoped to a story drops it when the hub is not that story, recording no issue.
6. **What `after` means.** `after` on a story names only a story whose accepted change the story needs on main. It never states an order of preference: the pull order does that.

## Consequences

- A task's touches name files, so ADR-0096's narrowing replaces a story's folder touch with them as soon as the tasks are written. A story in progress holds only the files its tasks change, and the folders it keeps on purpose.
- A planner that keeps a folder says why, in the note the operator reads before finalizing the story.
- The two findings show a hold while it lasts, naming the stories held, rather than in a later analysis.
- Both findings are advisory, as `narrative.state` is: a story keeps its folder touch from the start until its tasks name files, a planner may keep one with its note, and a hub's own agent cannot drop the `after` entries of the stories that name it. Failing `--strict` would stop `flai task done`, the close-out, and the main checkout's check over a hold nobody can clear at that moment; the finding informs, the refusals enforce.
- Reading the tree makes a touch's kind depend on the checkout: an entry is a folder once files exist under it. A task that named a path to be made keeps it when the path later becomes a folder, since only writes are refused.
- The planner's prompt, its definition in `.claude/agents/planner.md`, the conventions, and the user guide must say the new rules. Until they do, the planner meets them as refusals.
- An older flai neither refuses nor reports: the rules add no front matter key, so they raise no `flai.minimum`.

## Alternatives considered

- **Warn instead of refusing a task's folder touch.** The convention already asks for files, and the planner did not follow it; a warning adds one more thing not followed.
- **Refuse every writer's folder touch on a story.** The operator's words are theirs, and a story's agent may need a folder where files no task can name yet will be added.
- **Release an `after` hold when the hub reaches review.** The hub's change is not on main until it is accepted, so the dependent story's branch cannot build on it.
- **Findings that fail `--strict`.** A story keeps its folder touch until its tasks are written, and only the stories that name a hub can drop their `after`; a failing check would stop task closes and close-outs that cannot clear it.
- **A manifest setting for the hub's threshold.** Three is stated; a setting adds a settings row and a test for something nobody has asked for.
