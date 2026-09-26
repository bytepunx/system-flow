---
title: Coordinating stories that run in parallel
updated: 2026-09-26
status: active
---

# Coordinating stories that run in parallel

The finding of S-0124, for E-0009. The operator drags stories into `ready`, and flai starts an agent for each one while the in-progress limit has room. Nothing checks whether two of them will change the same files or depend on each other's work. This document says what flai does today, surveys how others decide whether two pieces of work may run at once, and ranks the designs that fit flai. The goal is a mechanism that holds a story back when it would collide with one already in progress, shows it yellow on the board, and says why and what clears it. The decision and the stories that follow are recorded at the end once the designer has chosen.

## What flai does today

| Mechanism | Where | Effect |
|-----------|-------|--------|
| `touches` on stories and tasks | Front matter, set with `--touches`, `flai touches`, `flai edit`, MCP and hostapi writes ([ADR-0019](../adrs/0019-story-branches-and-touches.md)) | Free strings, paths or component names. Advisory only. |
| Overlap test | `pathsOverlap` in `flai/internal/check/check.go`: equal, or one is a `/`-bounded prefix of the other | `flai check` warns `wip.overlap` for two in-progress items; the designer's inbox lists it; the MCP `who_touches` tool answers who covers a path. |
| Document editor | flaiover `DocEditor.svelte` | The designer's Save is held until they tick "I know" when an in-progress item touches the document. The only gate anywhere, and not on agents. |
| In-progress limit | `CanPull` in `workitem/boardview.go`, the launcher's `free` in `serve/agents.go` | The only thing that holds a ready story back. `flai move` only warns past it. |
| Pull order | `order:` in `board.md`, `flai order` (S-0057) | `wait_for_work` offers `Ready[0]`; the launcher starts ready stories in this order. |
| Blocked flag | `flai block` ([ADR-0004](../adrs/0004-workflow-states-and-transitions.md)) | The launcher skips a blocked ready story without saying so; `wait_for_work` does not skip it. Red ring and `BLOCKED` on the card. |
| Story branches | `story/S-nnnn` in `.flai-cache/worktrees/` (ADR-0019) | A conflict between two open stories surfaces only after one is accepted, as a rebase conflict at the other's `flai stream sync` or `flai accept`. Nothing compares two open branches. |
| Dependencies | none | Only `parent`. No `depends_on` or `blocked_by`. |
| Yellow on the board | `AgentDot` in flaiover, fed by `serve.Activity()` | Yellow means the agent is waiting: a question is open, the story is blocked, or a retry is queued for room. A story the launcher skips gets no per-story reason; the reasons are joined into one `Waiting` string in `serve/agents.json`. |

ADR-0019 considered locks instead of `touches` and rejected them: "a lock nobody releases blocks work; a warning plus a visible badge keeps humans in charge." Any hold chosen here refines that decision, so it has to answer the same worry: a hold must be released by something that always happens (acceptance, cancellation, a return to backlog), be visible with its reason, and be overridable by the operator.

Where a hold plugs in: the launcher's classification of ready stories (`launcher.look` in `serve/agents.go`), `wait_for_work`'s choice of story (`mcpserver/work.go`) and `CanPull`, the MCP `inbox` ready list, `flai board`, and `Activity()` for the yellow dot and its reason.

## What is at stake

A July 2026 study of 33,596 agent pull requests found textual conflicts in 19.8% of pairs from the same agent and 41.7% of pairs from different agents ([arXiv 2607.04697](https://arxiv.org/html/2607.04697v2)). Uber measured a 40% conflict probability at 16 concurrent changes ([SubmitQueue](https://blog.acolyer.org/2019/04/18/keeping-master-green-at-scale/)). Text is not the only failure: two patches can each pass and fail together with no textual conflict when one changes an interface the other uses; telling the second agent what the first changed recovered 82% of those failures in constructed tasks ([STALE, arXiv 2609.25396](https://arxiv.org/abs/2609.25396), a preprint).

Hard locks do badly with language-model agents. In Cursor's experiment, 20 agents sharing a locked coordination file slowed to the throughput of two or three, and agents forgot to release locks ([Cursor](https://cursor.com/blog/scaling-agents)). CoAgent measured two-phase locking at 0.81 deadlocks per trial and optimistic concurrency at 1.83 times the token cost; a fixed order decided at launch plus advisory notices gave a 1.4 times speedup ([arXiv 2606.15376](https://arxiv.org/html/2606.15376)). Coarse, ordered, advisory mechanisms do best.

## How others decide

| | Mechanism | How it works | Fit for flai | Cost |
|-|-----------|--------------|--------------|------|
| A | Isolate, then merge | One worktree, sandbox, or VM per agent; conflicts surface at merge. Claude Code worktrees, Cursor parallel agents, Codex cloud, Copilot coding agent, MultiDevin. | flai has it (ADR-0019). It finds conflicts late. | none |
| B | Disjoint scopes from the planner | The planner gives each agent its own files. Claude Code agent teams, Copilot `/fleet` ("no file locking… assign each agent distinct files"), Cursor planners and workers, OpenHands by directory, CODEOWNERS. | This is `touches`. What is missing is acting on it. | low |
| C | Advisory path leases | An agent takes an exclusive or shared lease on a glob, with a time to live; the grant lists conflicts; an optional pre-commit guard. MCP Agent Mail, Chubby. | Good. A story's lease is its `touches`, held from in progress to accept or cancel; no time to live is needed because acceptance and cancellation always end it. | medium |
| D | Mandatory locks | One holder may edit a file. Perforce `+l`, Git LFS locking. | Poor as a default; possible for a few named hot files. | medium to high |
| E | Hierarchical intent locks | Locks over a tree: a lock on `flai/` conflicts with one on `flai/cmd/serve`; two subtrees coexist. Gray's granularity locking. | Excellent as the overlap rule. It is the prefix test flai already has, applied at pull time. | low |
| F | Early conflict detection | Merge open work speculatively and report conflicts before anyone finishes. Palantír, Syde, WeCode, Crystal; today `git merge-tree --write-tree`, a three-way merge that touches no worktree. | Excellent as a safety net under declared `touches`: it finds what the declaration missed. Textual only unless builds and tests run too. | low |
| G | Merge queues | Test the merge result, not the branch. Bors, GitHub merge queue, GitLab merge trains, Zuul; Uber's conflict analyzer lets changes that share no build target land in parallel. | `flai accept` already serialises rebase and fast-forward. Uber's "disjoint targets run in parallel" is the precedent for pull-time holds. | low |
| H | Change-impact graphs | Map changed files to projects and walk reverse dependencies. Nx affected, Bazel target-determinator. | Moderate. flai is language-neutral; a few component edges declared in the manifest would do. | medium |
| I | Co-change mining | Files often changed in the same commit are coupled. Zimmermann's ROSE found a correct location in its top three suggestions over 70% of the time. | Good for widening `touches` that were declared too narrowly, as a warning. Needs flai's own bookkeeping commits filtered. | low to medium |
| J | Trained conflict prediction | A classifier on git features; good at predicting safe merges, weak at predicting conflicts. | Low; flai has no training data. | high |
| K | Order plus notices | A total order fixed at launch avoids deadlock; when earlier work writes what later work reads, the later agent is told. CoAgent, STALE. | High and cheap. The pull order is already a total order; the inbox can tell a story what an accepted story changed. | low |
| L | Shared live workspace | Every agent edits one workspace backed by conflict-free replicated data types. AgentRoom. | None; it contradicts one branch per story. | high |
| M | Explicit dependencies | "Blocked by" edges: GitHub issue dependencies, Claude agent-team tasks that cannot be claimed until their dependencies are done, Zuul `Depends-On`. | Good. flai has no such edge; one field and a check would add it. | low |
| N | Explaining a wait | A reason code, the evidence, and what clears it. GitHub merge queue timeline reasons, GitLab train notes, Kubernetes `0/3 nodes are available: 3 Insufficient cpu`, Bazel `somepath`. | The shape for the yellow card. | low |

## Candidate designs

### 1. Declared-touch holds, the pull order, and a trial-merge safety net (recommended)

- **Hold at pull.** A story's `touches` are its claim from `in-progress` until it is accepted, cancelled, or sent back. A ready story whose touches overlap an in-progress story's (the prefix rule, E) is held: the launcher does not start it, `wait_for_work` does not offer it, `inbox` and `flai board` list it as held, and the card turns yellow with the reason, the evidence, and what clears it: `held: touches flai/cmd/serve, inside flai/cmd which S-0128 (in progress) touches; starts when S-0128 is accepted or cancelled`.
- **The pull order is the only order.** A story is only ever held by one already in progress, so there is no cycle and no deadlock. Whether flai starts a later ready story that does not overlap while an earlier one is held is a policy question below.
- **Explicit dependencies.** An `after: [S-nnnn]` field (M) holds a story until those stories are done, for dependencies that are not about files.
- **Safety net at sync.** `flai stream sync` trial-merges the story branch with every other open story branch using `git merge-tree` (F) and reports a conflict to both stories' inboxes before either is accepted. It also compares the paths the branch actually changed with its `touches` and says when they drift, so a declaration that was too narrow is corrected while the story runs.
- **Notice at accept.** Accepting a story tells every open story that touches what it changed (K).
- **The operator overrides.** Start agent on a held story starts it anyway, as it already does past the in-progress limit.

Cheap: Go and git, no new service. Explainable, because every hold names a path and a story. It is only as good as the declarations; the drift check and trial merges limit the damage, and a story that declares a whole component serialises everything under it, which the board should make visible.

### 2. Design 1 with widening from history and component edges

Adds co-change mining (I) and component edges declared in the manifest (H). A direct overlap holds the story; a widened one (the files are usually changed together, or one component depends on the other) shows amber: the story starts, flagged, and is checked again at sync. Catches narrow declarations and interface coupling, at the cost of thresholds to tune and statistical reasons that are harder to explain. Worth doing after design 1 has run for a while.

### 3. Optimistic only

Ignore `touches` for scheduling. Run up to the in-progress limit, trial-merge at every sync, open a thread on a conflict, and let acceptance be the queue. Cheapest and never holds a story needlessly; this is how Codex, Cursor, and Copilot work today. It wastes agent work on the conflicts it finds late and gives the operator nothing to see before work starts.

### 4. Enforced claims

Design 1 plus a guard at `flai stream sync` or a pre-commit hook that refuses a commit to a path outside the story's claim or inside another story's, until the claim is widened, which may itself have to wait. The strongest guarantee, and it keeps declarations honest; it adds friction mid-story and risks the head-of-line blocking that made Cursor and Agent Mail back away from hard locks. Better kept for a few paths marked exclusive in the manifest.

## Recommendation

Design 1, built in this order: the hold with its reason (the part the operator sees), the explicit dependency, the trial merge and drift check at sync, the notice at accept. Design 2 later, if drift reports show declarations are often too narrow.

## Sources

- Claude Code [agent teams](https://code.claude.com/docs/en/agent-teams) and [worktrees](https://code.claude.com/docs/en/worktrees)
- Cursor, [scaling agents](https://cursor.com/blog/scaling-agents) and [worktrees](https://cursor.com/docs/configuration/worktrees)
- OpenAI [Codex cloud](https://developers.openai.com/codex/cloud), [issue 1351](https://github.com/openai/codex/issues/1351)
- GitHub, [Copilot `/fleet`](https://github.blog/ai-and-ml/github-copilot/run-multiple-agents-at-once-with-fleet-in-copilot-cli/), [Copilot coding agent](https://docs.github.com/copilot/concepts/agents/coding-agent/about-coding-agent), [merge queue](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/managing-a-merge-queue), [issue dependencies](https://docs.github.com/en/issues/tracking-your-work-with-issues/using-issues/creating-issue-dependencies), [CODEOWNERS](https://docs.github.com/articles/about-code-owners)
- Cognition, [Devin managing Devins](https://cognition.ai/blog/devin-can-now-manage-devins)
- OpenHands, [parallel agents for refactors](https://www.openhands.dev/blog/automating-massive-refactors-with-parallel-agents)
- [MCP Agent Mail](https://github.com/dicklesworthstone/mcp_agent_mail)
- Preprints: [CoAgent](https://arxiv.org/html/2606.15376), [STALE](https://arxiv.org/abs/2609.25396), [agent PR conflict rates](https://arxiv.org/html/2607.04697v2), [AgentRoom](https://arxiv.org/pdf/2608.23740), [S-Bus](https://arxiv.org/pdf/2605.17076)
- Early conflict detection: [Crystal](https://www.cs.ubc.ca/~rtholmes/papers/fse_2011_brun.pdf), [Palantír](https://web.engr.oregonstate.edu/~sarmaa/wp-content/uploads/2020/08/05928359.pdf), [Syde](https://dl.acm.org/doi/10.1145/1810295.1810339), [git merge-tree](https://git-scm.com/docs/git-merge-tree)
- Conflict prediction: [Owhadi-Kareshk, Nadi, Rubin](https://arxiv.org/abs/1907.06274); co-change: [Zimmermann et al., ROSE](https://thomas-zimmermann.com/publications/files/zimmermann-tse-2005.pdf)
- Merge queues: [Uber SubmitQueue](https://github.com/uber/submitqueue), [Zuul gating](https://zuul-ci.org/docs/zuul/latest/gating.html), [GitLab merge trains](https://github.com/gitlabhq/gitlabhq/blob/master/doc/ci/pipelines/merge_trains.md), [Bors](https://bors.tech/newsletter/2023/04/30/tmib-76/)
- Locking: [Gray, granularity of locks](https://mwhittaker.github.io/papers/html/gray1976granularity.html), [Chubby](https://research.google.com/archive/chubby.html), [Git LFS locking](https://github.com/git-lfs/git-lfs/wiki/File-Locking), [Perforce exclusive checkout](https://portal.perforce.com/s/article/3114)
- Impact graphs: [Nx affected](https://nx.dev/docs/features/ci-features/affected), [target-determinator](https://github.com/bazel-contrib/target-determinator)
- [Trunk-based development, short-lived branches](https://trunkbaseddevelopment.com/short-lived-feature-branches/)

The September 2026 preprints are not peer reviewed; only their central ideas are relied on here.

## Decision

[ADR-0046](../adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md), 2026-09-26. The designer chose design 1 on TH-0018 and took every policy recommendation on TH-0019:

- A story's claim is its `touches` and its open tasks' `touches`, with sub-project names and tags read as their paths. The claim holds while the story is in progress or in review.
- A ready story is held when its claim overlaps an open story's claim, when its claim is empty or an open story's is, or when it names a story in `after:` that is not done.
- The launcher and `wait_for_work` take the first ready story that is not held, in pull order. A held story keeps its place.
- `flai move` warns on a held story and does not refuse it. Start agent overrides the hold.
- The card is yellow and says the reason and what clears it.
- `flai stream sync` trial-merges open branches and reports paths changed outside the claim. Acceptance tells overlapping open stories what changed.

The stories that build it are under E-0009. Until they are accepted, flai behaves as described in "What flai does today".
