---
id: ADR-0313
title: "Rebase and resolve runs as a flai serve job that rebases a review story's branch and closes it out, and sends the story back to its agent only on conflicts or a failed step"
status: accepted
date: 2026-11-17
supersedes: [ADR-0302]
superseded_by: []
refines: [ADR-0229, ADR-0248]
topics: [flai, cli, dashboard, acceptance, verify]
---

# ADR-0313 Rebase and resolve runs as a flai serve job that rebases a review story's branch and closes it out, and sends the story back to its agent only on conflicts or a failed step

## Context

[ADR-0302](0302-the-operator-sends-a-review-story-that-cannot-be-accepted-back-to-its-agent-to.md) made `flai accept S-nnnn --resolve`, and the dashboard's Rebase and resolve button, send a review story that cannot be accepted straight back to in-progress, and start its agent. Most such stories need only a rebase onto main and a close-out over the new main, which no agent has to do. On 2026-11-15 five stories were sent back so in one day, each to an agent that rebased, closed out, and moved the story back to review.

flai serve already runs a story's close-out as a job with no model ([ADR-0229](0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md)), and an acceptance already sends a story back to its agent through a thread when its rebase stops on conflicts ([ADR-0248](0248-every-acceptance-commits-the-main-checkout-s-uncommitted-files-under-wip-design.md)). The close-out job's log said only that `flai verify` ran, because `flai verify` printed its steps once it had ended, so nothing could show what a long close-out was doing.

## Decision

Rebase and resolve runs as a flai serve job that rebases a review story's branch and closes it out, and sends the story back to its agent only on conflicts or a failed step.

1. **The request.** `flai accept S-nnnn --resolve`, and the host method `accept.resolve` with `{id}`, which runs `flai accept <id> --resolve --by=<owner> --json`, run the acceptance preview as before. When it lists a blocker, they ask flai serve for a job of kind `resolve`, through a request under `.flai-cache/verify/jobs` as `flai verify --job` asks for a close-out, and answer at once. With `--json` the answer is `{id, job, started, running, log, reason}`. The request moves nothing, opens no thread, and starts no agent.
2. **One job per story.** flai serve refuses a story that is not in review, or that has no worktree. A request while a job of the story runs, a close-out or a rebase and resolve, is answered with that job and starts nothing. The verify record never answers it, since the job always rebases first.
3. **The rebase.** The job runs this flai's `flai stream sync S-nnnn` in the story's worktree, on the host. A rebase that stops on conflicts is aborted, after the job reads the paths in conflict, so the branch is as it was. A rebase that was already waiting before the sync is left as it was found. A sync that fails with no rebase waiting, such as on uncommitted changes, ends the job too. Either way the job is `stopped` at `rebase` and the close-out does not run.
4. **The close-out.** After a clean rebase the job runs the story's close-out at the rebased head, exactly as a close-out job does: `scripts/close-out.sh`, or for a story whose agent runs in a container, its steps there and on the host ([ADR-0274](0274-the-close-out-of-a-story-whose-agent-runs-in-a-container-runs-its-verify-in.md)). Its commit message is `chore: [S-nnnn] close out again after rebase and resolve`.
5. **The progress.** Every close-out and resolve job runs with `FLAI_VERIFY_PROGRESS=1`, and `flai verify` then writes a line to stderr as each of its steps and tiers starts and ends: `verify: S-nnnn <time> step|tier started <name>` and `verify: S-nnnn <time> step|tier ended <outcome> <name>`. The resolve job writes the same lines, from `resolve:`, for its steps `rebase` and `close-out`. All of it goes to the job's log. `closeout.Status` reads the job and its log into the step, the tier running, and every step and tier with its start, end, and outcome. The host read `resolve.status` with `{id}` answers it, and the dashboard follows a running job through `resolve.status` and `log.read`.
6. **The end.** A clean rebase and a passing close-out leave the story in review, with nothing sent back and no agent started. Only a rebase that stopped, or a close-out that stopped at a step or failed, sends the story back through ADR-0248's send-back: the job moves the story to in-progress as `flai`, and opens the thread whose title begins `Acceptance sent back:`, or writes on the one still open. The thread names the paths in conflict, or the step, at most ten of its findings with the rest counted, and the job's log. flai serve starts the story's agent on it and resolves it once the story is back in review, as for every acceptance's send-back. A story that left review while the job ran is left as it is.
7. **The record.** The job's state, `S-nnnn.job.json`, gains `kind`, and for a resolve job its `outcome` (`stayed in review` or `sent back`), `thread`, `conflicts`, and `reason`. Its start and end are journal entries with action `close-out` and method `resolve.job`, which `flai stats` does not count among the close-outs. Its end announces a `verified` event, as a close-out job's does. A job a restart of flai serve interrupts is ended `interrupted` and sends nothing back: the story stays in review.

What of ADR-0302 stays: `--resolve` accepts nothing and merges nothing; it is the operator's alone, refused with `--dry-run`, `--yes`, `--batch`, the flags of an acceptance, the orchestrator's `--by`, and under `FLAI_ROLE=orchestrate`; a story not in review, or one nothing holds, is refused with exit 4, under `--json` as `{"refused": {id, reason}}`; `accept.resolve` is gated by no host action; the dashboard's Rebase and resolve button runs it; a send-back goes through ADR-0248's thread. What goes: the send-back on every request, the thread that lists each blocker for the agent to take, and the send-back object `{id, reason, sent_back, branch, base, conflicts, thread, blockers}`. A request flai serve does not take, or that an older flai serve answers with no job, is now refused with exit 4 too, and nothing runs.

This supersedes ADR-0302. It refines ADR-0229, whose close-out job gains a kind and its progress lines, and ADR-0248, whose send-back it uses from flai serve.

## Consequences

- A story held only by a branch behind main, or a verify that no longer stands, is made acceptable with no agent and no model. The operator sees the job run and accepts the story once it stays in review.
- An agent is started only for what needs one: a conflict, or a step that fails at the rebased head. Its thread says which, and where the output is.
- Rebase and resolve now needs flai serve running for the project. Without it the request is refused and nothing changes.
- A blocker the rebase and the close-out do not settle, such as uncommitted files in the main checkout, is not sent back: the job passes, the story stays in review, and the preview still lists it. The operator settles it, or sends the story back by hand.
- Every close-out's log now holds a line per step and tier, which `closeout.Running` and `resolve.status` read, so a reader can say which tier runs, not only that `flai verify` does. Stdout, whose last line the close-out script reads, is unchanged.
- The close-out's rules are not copied: the job runs the same close-out as a close-out job, so a change to one changes both.

## Alternatives considered

- **Keep sending every story back, and let the agent rebase.** It spends an agent's session on work with no judgement in it, which is what the operator asked to stop.
- **Rebase in `flai accept --resolve` itself and run the close-out there.** A close-out runs for minutes. A CLI call or a host method that held it would block the dashboard and be cut short, as tool calls were before ADR-0229.
- **A job of its own beside the close-out job.** Two runners for the same steps would drift apart. One job with a kind shares the request files, the log, the sweep, and the journal.
- **Answer from the verify record when it holds a pass.** The record is for the branch before the rebase. The job exists to verify over the new main, so it always rebases first.
