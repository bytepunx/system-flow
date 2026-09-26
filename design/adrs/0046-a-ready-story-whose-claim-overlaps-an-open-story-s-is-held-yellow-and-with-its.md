---
id: ADR-0046
title: "A ready story whose claim overlaps an open story's is held, yellow and with its reason, until that story is accepted"
status: accepted
date: 2026-09-26
supersedes: []
superseded_by: []
refines: [ADR-0019]
---

# ADR-0046 A ready story whose claim overlaps an open story's is held, yellow and with its reason, until that story is accepted

## Context

flai serve starts an agent for every ready story while the in-progress limit has room ([ADR-0043](0043-flai-serve-starts-a-ready-story-s-agent-whenever-the-in-progress-limit-has-room.md)), and `wait_for_work` offers the first ready story to an idle agent. Nothing checks whether two of these stories will change the same files. [ADR-0019](0019-story-branches-and-touches.md) gave stories and tasks `touches`, but only as a warning: `flai check` reports `wip.overlap`, and the dashboard badges documents. A conflict between two open story branches shows up only after one of them is accepted, as a rebase conflict in the other story's worktree. ADR-0019 rejected locks because "a lock nobody releases blocks work; a warning plus a visible badge keeps humans in charge."

E-0009 asks for mechanisms that avoid parallel changes to the same files. S-0124 asks that the operator be able to see why a story is flagged yellow and waits for the stories before it. S-0124 surveyed agent tools, research, merge queues, conflict prediction, and locking; the findings are in [agent-coordination.md](../system/agent-coordination.md). Every agent product isolates work and merges it later, and none predicts overlap from declared paths. The preprints report high conflict rates between parallel agent changes. Hard locks and optimistic retries both did badly with language-model agents. Coarse, advisory claims taken in a fixed order did well.

On TH-0018 (2026-09-26) the designer chose design 1 of the four offered: "go with design 1". On TH-0019 the designer took every policy recommendation: "continue with recommendations".

## Decision

**A story's claim is its `touches` together with those of its open tasks. The claim holds from the moment the story goes in progress until it is accepted, cancelled, or sent back to backlog. A ready story whose claim overlaps an open story's claim is held. So is a ready story that names, in `after:`, a story that is not done. flai serve does not start a held story's agent, and `wait_for_work` does not offer it. The board shows it yellow, with the reason and what clears it. `flai stream sync` trial-merges open branches to catch what the claims missed, and acceptance tells overlapping stories what changed.** This refines ADR-0019. `touches` becomes a claim that scheduling acts on, not a lock on editing. The warning, the badge, and the operator's control all stay.

### The claim and the overlap

- **Open stories.** A story is open while it is in progress or in review. Its branch is not merged until acceptance, so its claim holds through review.
- **The claim.** The story's own `touches` plus the `touches` of its tasks that are not done or cancelled.
- **Component names.** Before comparison, a `touches` entry that is a sub-project's name or one of its tags in `system-flow.yaml` becomes that sub-project's path. So `cli` and `flai` both become `flai`. An entry that matches nothing is compared as it is.
- **Overlap.** Two claims overlap when an entry of one equals an entry of the other, or is a `/`-bounded prefix of it. This is the rule `flai check` already applies (`flai/cmd` overlaps `flai/cmd/serve`; `flai/cmd` does not overlap `flaiover`).
- **No touches.** A story whose claim is empty overlaps every story. While any story is open, it is held with the reason that it declares no touches. While it is open itself, it holds every ready story.

### The hold

A ready story is held when any of these is true:

- `overlap`: its claim overlaps an open story's claim;
- `no-touches`: its claim is empty and some story is open, or some open story's claim is empty;
- `after`: its new front matter field `after: [S-nnnn, …]` names a story that is not done.

A cancelled story named in `after:` keeps the hold. The reason says it was cancelled, so that the operator decides whether the dependent story still makes sense. `flai check` reports an `after:` entry that names no story, names the story itself, or forms a cycle.

- **Pull order with skip-ahead.** The launcher and `wait_for_work` walk ready stories in the pull order and take the first that is not held, within the in-progress limit. A held story keeps its place and is taken first once it is clear. A story can only be held by a story already open, or by one it names, so the pull order never deadlocks.
- **Hand moves warn.** `flai move … in-progress` and the MCP `item_move` warn about a held story, as they do past the in-progress limit, and do not refuse. Start agent (`flai serve agent start`) starts a held story's agent on the operator's word, as it does past the limit.
- **Declaring touches.** Whoever creates a story declares its `touches`, as today; they are not required for ready. The agent that pulls a story may widen them when it writes the tasks. Widening an open story's claim holds only ready stories. It never stops a story that is already open.

### Saying why

A held story carries a reason code (`overlap`, `no-touches`, `after`), the evidence, and what clears it. For example:

`held (overlap): touches flai/cmd/serve, inside flai/cmd which S-0128 (in progress) touches; starts when S-0128 is accepted, cancelled, or sent back`

- **Dashboard.** The card's dot is yellow and its line is the reason, through `agent.status`, for a story that has never had an agent as well as for one that has.
- **Board and inbox.** `flai board`, the MCP `board` and `inbox` ready lists, and `wait_for_work`'s answer when it has nothing to offer all carry `held` and the reason.

### The safety net

- **At sync.** `flai stream sync` does two checks against every other open story branch:
  - It trial-merges the branch with `git merge-tree --write-tree`, which writes nothing to any worktree. A textual conflict is reported in sync's output and to both stories.
  - It compares the paths the branch has changed since main with its claim. Paths outside the claim are reported, so that the agent widens `touches`.
- **At acceptance.** Accepting a story reports, to every open story whose claim overlaps the accepted change, which paths changed. Those agents sync and check their work against it.

## Consequences

- Two stories that declare overlapping paths no longer run at once without anyone deciding it. The operator sees which one waits and why before any agent work is wasted.
- `touches` now changes what runs. Coarse declarations serialise work: `flai`, or no touches at all, holds every other story that touches `flai`. Narrow declarations let conflicts through, which the sync checks then catch.
- Stories gain an optional `after:` field. An older flai ignores it and holds nothing.
- `flai stream sync` runs one `git merge-tree` for each other open branch, which takes milliseconds per branch at the in-progress limits used here.
- Semantic conflicts, where two changes each work alone and fail together without a textual conflict, are caught only as well as `touches` and the notice at acceptance catch them. Tests after a sync remain the agent's job.
- The implementation is split into stories under E-0009, in this order: the hold with its reason, `after:`, the trial merge and drift check at sync, and the notice at acceptance.

## Alternatives considered

- **Design 2: widen claims from co-change history and component edges** (amber, starts flagged). This was offered on TH-0018 and deferred. The drift reports from sync will show whether declarations are too narrow often enough to need it.
- **Design 3: optimistic only.** Never hold, trial-merge and report. It is cheapest, and it is how Codex, Cursor, and Copilot work. It wastes agent work on conflicts found late, and it gives the operator nothing to see before work starts.
- **Design 4: enforced claims.** Refuse commits outside a story's claim. This is the strongest guarantee, but it is the head-of-line blocking that led Cursor and MCP Agent Mail away from hard locks. It may return later for a few paths marked exclusive in the manifest.
- **A strict pull order.** Nothing behind a held story starts. This was offered on TH-0019 and declined, because it idles capacity that a story with no overlap could use.
- **Stories without touches run freely**, flagged amber. This was declined on TH-0019. The safe reading of "declares nothing" is "may change anything".
- **Refusing hand moves of a held story.** This was declined on TH-0019. The operator keeps the last word, as with the in-progress limit.
