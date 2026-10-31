---
id: ADR-0229
title: "The close-out runs as a flai serve job the story's agent asks for with the verify MCP tool or flai verify --job and waits for by a verified event, and a close-out asked again at a commit the verify record holds a verdict for answers from the record"
status: accepted
date: 2026-10-31
supersedes: []
superseded_by: []
---

# ADR-0229 The close-out runs as a flai serve job the story's agent asks for with the verify MCP tool or flai verify --job and waits for by a verified event, and a close-out asked again at a commit the verify record holds a verdict for answers from the record

## Context

`scripts/close-out.sh` runs `flai verify S-nnnn --record-issues` inside one Bash tool call of the story's agent. Claude Code ends a tool call at 600 s. In the fortnight to 2026-10-29, 64 of 267 close-outs died there and 146 ended with no summary line; the agent then ran the whole close-out again from the start, for 18.5 of the fortnight's 113 session-hours. The timeout is the harness's, not a setting flai can raise, so the close-out has to leave the agent's turn.

flai serve already runs work for other processes: `serve.RequestMerge` hands it a pull request merge through a request file under the project's `.flai-cache`, which serve takes by renaming it and answers in a file beside it. The MCP server is not serve's process, and it does not know serve's state folder. The verify record, `.flai-cache/verify/<story>.json`, keeps only the newest run of each story, and `flai verify` always runs its pre-tier steps, so nothing can answer a second ask without running.

## Decision

The close-out runs as a job of flai serve in the story's worktree, asked for through a request file the `verify` MCP tool or `flai verify --job` writes and answered with the job's id at once, its end announced by a `verified` event in the project's event log and journaled with its duration and verdict, and a close-out asked again at a commit the verify record holds a verdict for answers from the record and runs nothing.

1. **The request.** The story's agent asks for the close-out with the MCP tool `verify`, or `flai verify S-nnnn --job` on the host, each with the commit message the close-out commits with (`message`, `-m`) and `all` and `fresh` as `flai verify` takes them. Either writes `S-nnnn.request.json` under the project's `.flai-cache/verify/jobs/`, as `RequestMerge` writes a merge request, and waits at most ten seconds for serve to take it and answer with the job's id. A request serve does not take in that time is withdrawn, and the asker is told that flai serve is not running and to run the close-out script, or `flai verify` where the project has none, itself. `flai verify` without `--job` runs in its own process as before.
2. **The job.** flai serve looks for requests every second in each served project, beside the merge pass. It starts one job per story at a time, with the id `<story>-<UTC time of its start, yyyymmddThhmmssZ>`; a request for a story whose job is running is answered with that job's id and starts nothing. Stories' jobs run side by side, as their agents' close-outs did. The job runs the project's close-out in the story's worktree: `scripts/close-out.sh S-nnnn -m <message>` where the worktree has that script, and serve's own `flai verify S-nnnn --record-issues` where it has none, with `FLAI_AGENT` set to the agent that asked, at most two hours, its output kept in `S-nnnn.log` beside the request. The job's state, `S-nnnn.job.json`, says its id, story, asker, start, and, once it ended, its end, exit code, verdict, and the head commit it ended at. A job still running when serve stopped is ended failed, with the verdict `interrupted`, when serve starts again.
3. **The verdict.** A job's verdict is `passed` when the close-out exits 0, `stopped` with the step it stopped at when the verify record names one, and `failed` otherwise. A job that answers from the record (5) carries the verdict the record holds.
4. **The event.** When a job ends, serve appends a `verified` event to `.flai-cache/events.jsonl` (ADR-0211) with the story as its ID and fields `job`, `commit`, `verdict`, `stopped_at`, `duration`, and `recorded` (true when it answered from the record). The inbox reports a `verified` event to the agent of its story and to readers with no story of their own, the orchestrator among them, and to no other story's agent. It is not the merge queue's `verified`, which lives in the queue's own log, `.flai-cache/merge-queue/events.jsonl`, and is read only by `flai stats`. `wait_for_events` lists `verified` among its kinds, and does not answer `end` to a story's agent while a job of its story runs.
5. **The answer from the record.** Each run `flai verify` saves is kept per commit as well as newest, under `.flai-cache/verify/S-nnnn/<commit>.json`, when the worktree was clean as it ran; a close-out job whose script committed after a passing verify keeps the verdict under the commit it made too, since that commit holds exactly the files the run verified. A close-out asked for at a commit the record holds, with the worktree clean and the main branch at the record's base commit or ahead of it only by commits the sync step passes over (ADR-0135), answers that verdict, says that it came from the record and when that run was, and runs nothing: the request answers it at once, before any job is started, and `flai verify` answers it the same way. `fresh` always runs. `flai verify --last`, the verifier, and the dashboard's `verify.status` read the newest run as before.
6. **The journal.** serve journals each job's start and end through `Host.Record` under the action `close-out` and the method `verify.job`: the start with the outcome `started`, the end with `done` for a pass and `failed` otherwise. `hostapi.Entry` gains `duration_ms` and `verdict`, set on the end, and `flai serve journal` shows them. A close-out answered from the record is journaled as one end entry with `verdict` and a `detail` saying it came from the record.
7. **The count.** `flai stats` counts close-outs per week from the serve journal's `close-out` end entries of the project, by verdict, with those answered from the record apart, and the median and 85th percentile of their durations, as `design/system/metrics.md` defines.

## Consequences

- No tool call holds the close-out, so the harness's timeout no longer cuts it short, and a close-out asked again at a verified commit costs nothing.
- The story's agent waits for `verified` with `wait_for_events` instead of a Bash call; the prompt flai serve gives it, the conventions, and `scripts/close-out.sh`'s header say so.
- A session with no flai serve keeps the close-out script and `flai verify` in its own process, with the harness's timeout as before.
- The verify record grows by one file per verified clean commit; they are small and `.flai-cache` is the project's cache.
- `flai stats` reads the host's serve journal, which only a host running flai serve has; without it the close-out section is empty.

## Alternatives considered

- **Run the close-out in the background from the agent's Bash call and poll a file.** It leaves the turn, but every poll is a model turn, and nothing records the verdict or the duration for the operator.
- **Run the close-out inside the MCP server's process.** The MCP tool call is held by the same timeout, and an MCP server restarted with the agent kills the run.
- **The dashboard's `verify.run` host method.** It is detached already, but only the dashboard channel reaches it, behind the `checks` host action, and it runs `flai verify`, not the close-out script, which commits.
- **Key the record by the story's patch id instead of the commit.** It would answer across a rebase that changed nothing of the branch's own, but the sync step's verdict depends on the main branch too; the commit and the base say exactly what the verdict covered.
