---
id: ADR-0320
title: "The auto-upgrade drains the verify pass's jobs before it restarts, holds new close-out requests until after the restart for at most 45 minutes, and flai serve asks again once for a close-out or rebase and resolve a restart interrupted"
status: accepted
date: 2026-11-27
supersedes: []
superseded_by: []
refines: [ADR-0158, ADR-0229]
---

# ADR-0320 The auto-upgrade drains the verify pass's jobs before it restarts, holds new close-out requests until after the restart for at most 45 minutes, and flai serve asks again once for a close-out or rebase and resolve a restart interrupted

## Context

[I-0388](../issues/I-0388-flai-serve-s-auto-upgrade-restart-interrupts-the-close-out-jobs-it-was-running.md) counts 11 close-outs in W46 that ended `interrupted, flai serve stopped while it ran`, two more in W45 and two in W47. Each was cut by an auto-upgrade restart. The serve journal holds 16 auto-upgrade restarts in W46. Each interrupted close-out was asked for again from the start by its agent, which ran the integration tier again, about nine minutes, and lost up to twenty minutes of wall clock.

[ADR-0158](0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md) §5 has the upgrade wait only for a dashboard write that `flai serve` runs in its own process. It predates the close-out job of [ADR-0229](0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md). That job runs in a process group serve stops when it stops (`runCloseOut`), and §2 of ADR-0229 has the next serve end it `interrupted`. Nothing asks for it again. The rebase and resolve job ([ADR-0313](0313-rebase-and-resolve-runs-as-a-flai-serve-job-that-rebases-a-review-story-s.md)) and the state job share the verify pass and its fate.

The close-out's p50 was 474 s in W46. The longest job recorded is the 33-minute rebase and resolve of S-0824.

## Decision

The auto-upgrade drains the verify pass's jobs before it restarts, holds the requests that would start a new job until after the restart, gives up waiting after 45 minutes, and the next serve asks again, once, for each close-out or rebase and resolve a restart still interrupted.

1. **The drain.** While any served project's verify pass runs a job, a close-out, a rebase and resolve, or a state job, `flai serve` asks the host for no upgrade, as it waits for a dashboard write (ADR-0158 §5). It logs `flai upgrade waits for a close-out job` (info) once per wait, with the number of jobs running, and asks again on the next minute's look.
2. **No new job meanwhile.** While an upgrade waits on the drain, the verify pass starts no job. A request it would answer by starting one is held instead: the pass takes it as usual, keeps it as `S-nnnn.held.json` beside the request, and answers at once with `held: true` and a reason that names the release and says the job starts once flai serve has restarted on it, with a `verified` event at its end. A request answered from the verify record, or with the id of the story's job that runs, is answered as before. No state job starts; its entry runs once the drain lifts, as an interrupted one does today. The MCP tool `verify` answers a held request as an answer, not an error, and `flai verify --job` prints its reason and exits 0. An MCP server older than this reads the answer as a refusal; its reason says the same.
3. **Held requests are put back.** At each look while no drain holds, at its first look after a start among them, the verify pass renames each held request back to `S-nnnn.request.json`, unless the story has a newer request waiting, and takes it as any request. So a held request starts its job after the restart, or at once when the upgrade the drain waited for was not asked after all.
4. **The bound.** The drain waits at most 45 minutes from the look that first found the upgrade waiting on a job. Then `flai serve` asks for the upgrade with the jobs still running, and logs `flai upgrade stops waiting for close-out jobs` (warn) with their number. No manifest setting is added: the bound covers the longest job recorded with room, and the upgrade's own retry every 15 minutes and its thread after an hour (ADR-0158 §6) stay as they are, counted from the first ask.
5. **Asked again once.** When serve starts and ends a close-out or a rebase and resolve job `interrupted` (ADR-0229 §2), it writes a request for the same story, kind, asker, message, `all`, and `fresh`, marked `again: true`, and the job that request starts carries `again: true` in its state. It writes none when the interrupted job was itself asked again, so a job is asked again at most once; when a request or a held request of the story is already waiting; or, for a rebase and resolve, when the story is no longer in review. The journal's `interrupted` entry says the job was asked again. The agent's `verified` event with the verdict `interrupted` is followed by the new job's own `verified` event, and a request the agent makes meanwhile is answered with that job's id, as for any job that runs. A state job is not asked again: its entry already runs again.

Kluster jobs, the merge queue's batch, and the agent image jobs are outside this decision: each keeps what it does at a restart today.

## Consequences

- An auto-upgrade no longer cuts a close-out that would end within 45 minutes; the upgrade lands at most that much later than it does today.
- A close-out asked for during the drain waits for the restart rather than running and being cut, and its agent waits for the `verified` event as for any job.
- A close-out the restart still cut, past the bound, at an operator's `flai host upgrade`, or after a crash, runs again without its agent, once; a second interruption is left to the agent, as every interruption is today.
- `closeout.Request`, `closeout.Answer`, and `closeout.Job` gain fields, `again` and `held`, that an older flai ignores.

## Alternatives considered

- **Wait for the jobs with no bound.** One job that hangs, or a steady stream of close-outs, would hold the host on an old flai for hours, while ADR-0142 refuses `flai task done` to every agent on it.
- **Keep taking jobs during the drain.** With three stories in progress, close-outs overlap, and the drain would seldom end before the bound.
- **Restart while the jobs run, detached, and have the new serve adopt them.** The job's end, its verdict, the record kept at the commit, and a resolve job's settle all run in serve's process; adopting them needs a job runner outside serve, a larger change than the drain.
- **Resume a job mid-tier.** `flai verify` already reuses the tiers that passed at the same patch, so a job asked again repeats only the tier that was cut; resuming inside a tier would need every tier to checkpoint.
- **Let the agent ask again.** It does today, and loses the time between the interruption and its next turn; the drain removes most interruptions and the ask again the rest.
