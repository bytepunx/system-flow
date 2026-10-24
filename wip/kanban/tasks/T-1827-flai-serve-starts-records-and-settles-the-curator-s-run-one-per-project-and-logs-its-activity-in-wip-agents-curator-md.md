---
id: T-1827
type: task
nature: feature
title: flai serve starts, records, and settles the curator's run, one per project, and logs its activity in wip/agents/curator.md
status: done
parent: S-0431
owner: alex
created: 2026-10-16T06:06:02Z
updated: 2026-10-24T00:25:02Z
transitions:
  - to: ready
    at: 2026-10-24T00:14:13Z
    by: agent-S-0431
  - to: in-progress
    at: 2026-10-24T00:14:14Z
    by: agent-S-0431
  - to: done
    at: 2026-10-24T00:25:01Z
    by: agent-S-0431
stream: S-0431
tags: [cli]
touches: [flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/agents.go, flai/internal/serve/activity.go, flai/internal/serve/claudecheck.go, flai/internal/serve/serve.go]
after: [T-1825]
usage:
  source: log
  seconds: 647
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 89
      output: 434
      cache_read: 4843722
      cache_write: 149905
      cost: 1.6988
---
# T-1827 flai serve starts, records, and settles the curator's run, one per project, and logs its activity in wip/agents/curator.md

## Work

Add `flai/internal/serve/curate.go` on the pattern of `analyze.go`:

- `Curate` and `curateCheck`: refuse while the `curate` action is off, while a curator runs (naming it), and with no harness and no command; otherwise start the run through the launcher in the main checkout as `curator`, with the trigger `asked`.
- `curateEnded`: settle the run, `worked` on a clean exit, `failed` otherwise with `why`; log its activity with `logRunEndSaying`, the summary its final reply's first line, its items the entries whose files changed since the run started.
- `agents.go`: keep the newest run under `curator` in `serve/agents.json`, outside `Running` and the in-progress limit, with its log `<key>-curator-<start>.log`; `agent.status` carries it.
- `activity.go`: the curator's activity charges no item; its cost stays in its kind's project strategic total.

Waits for T-1825: the launcher starts the run through the harness's request.

## Done when

`curate_test.go` covers each refusal, a start, a clean end that logs an activity naming the entries a stub harness command wrote on a fixture project, a failed end, and a second start refused while one runs; `flai test flai/internal/serve` passes.

## Notes

The stub command test stands for the story's live run: the curator's first real run on this repository needs this story released and the operator's `curate` action on.
