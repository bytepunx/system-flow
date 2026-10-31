---
id: T-2396
type: task
nature: improvement
title: "An ADR settles the close-out job: how the CLI and MCP reach flai serve, what the job runs, the verified event, the per-commit verdict, and the journal's fields, and metrics.md defines the close-out count"
status: backlog
parent: S-0546
owner: alex
created: 2026-10-31T07:43:28Z
updated: 2026-10-31T07:43:28Z
transitions: []
stream: S-0546
tags: [flai]
touches: [design/adrs/0224-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-waits-on-by-event.md, design/adrs/README.md, design/system/metrics.md, design/system/flai-cli.md]
---
# T-2396 An ADR settles the close-out job: how the CLI and MCP reach flai serve, what the job runs, the verified event, the per-commit verdict, and the journal's fields, and metrics.md defines the close-out count

## Work

Write the ADR with `flai adr new` (it gives the number; rename the touch to it with `flai touches`) and decide, from the code the planner's explorer read:

- How the job is asked for. `flai verify --job` writes a request under the serve directory and serve takes it, as `serve.RequestMerge` (`flai/internal/serve/forgemerge.go`) does for merges, polling every second. The MCP `verify` tool writes the same request, since the MCP server is not serve's process. A session with no serve running is told so, and `flai verify` without `--job` runs in process as today.
- What the job runs. The project's close-out (`scripts/close-out.sh` here, `flai verify` directly in the template's copy), in the story's worktree, with the commit options the agent gives. Say whether the job runs the script or the steps in Go. If it runs the script, the script stops calling itself the agent's turn.
- The event. `verified`, appended to `.flai-cache/events.jsonl` (`flai/internal/inbox/events.go`), with the story, the job id, the commit, and the verdict. Say how it differs from `mergequeue.EventVerified`, which has the same name and lives in the merge queue's own log.
- The answer from the record. Today the record at `.flai-cache/verify/<story>.json` is keyed by story and overwritten by each run (`verify.SaveReport`, `flai/internal/verify/story.go`). Decide how the record keeps a verdict per commit. Today the close-out verifies at the head before its own commit, so decide which commit a second ask is matched against: before the commit, after it, or both.
- The journal. `hostapi.Entry` has no duration or verdict. Decide whether to add fields or use `detail`, and what `flai stats` counts.

Then add the close-out count to `design/system/metrics.md`, citing the ADR, since a metric changes only with an ADR. Also update the `flai verify` row and the event-kind table in `design/system/flai-cli.md`. Waits for no other task: every other task builds on what this one decides.

## Done when

- [ ] The ADR is accepted in `design/adrs/` and indexed in `design/adrs/README.md`
- [ ] `design/system/metrics.md` defines what `flai stats` counts of close-outs, and cites the ADR
- [ ] `design/system/flai-cli.md` lists `verified` among the event kinds and gives `flai verify --job`
- [ ] `flai check --strict` and the markdown lint pass on the changed files

## Notes
