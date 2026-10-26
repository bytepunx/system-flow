---
id: T-2139
type: task
nature: remediation
title: Measure the one-shot search.query on the monorepo and give the monorepo search tests a timeout that covers a loaded host
status: done
parent: S-0501
owner: alex
created: 2026-10-26T04:31:17Z
updated: 2026-10-26T04:35:40Z
transitions:
  - to: ready
    at: 2026-10-26T04:33:06Z
    by: agent-S-0501
  - to: in-progress
    at: 2026-10-26T04:33:07Z
    by: agent-S-0501
  - to: done
    at: 2026-10-26T04:35:40Z
    by: agent-S-0501
stream: S-0501
tags: [dashboard]
touches: [flaiover/src/lib/server/search.test.ts]
usage:
  source: log
  seconds: 153
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 39
      output: 6869
      cache_read: 2476618
      cache_write: 73453
      cost: 1.2205
---
# T-2139 Measure the one-shot search.query on the monorepo and give the monorepo search tests a timeout that covers a loaded host

## Work

Confirm the cause before changing the test. Time `bin/flai hostapi search.query '{"q":"front matter","docs":false,"limit":30}'` in the repository root, three runs idle and three while another vitest or Go test run loads the host, and time `adrs.list` the same way; record the figures in the narrative's `## Decisions`. The test's `flaiAsk` runs one such process per call, and a one-shot process builds the index over every markdown file under `design/`, `wip/` with its archive, and `docs/` with no cache across processes (`flai/internal/hostapi/searching.go`), which is what crosses vitest's 5 s default under other stories' close-outs.

Then give the `adrs on the monorepo` block in `flaiover/src/lib/server/search.test.ts` a timeout that covers a loaded host, 30 s, as `writes.test.ts` does for its twenty flai processes (S-0116), with a comment naming I-0185, the measured figures, and why: the one-shot flai builds the whole repository's index. Prefer the block's options so both monorepo tests get it. Do not change `searching.go` or the index: if the idle measurement shows the build itself takes seconds, say so in the narrative and ask on a thread on S-0501 whether the index should be made cheaper, and leave that to another story.

Run `flai test flaiover/src/lib/server/search.test.ts` and the file alone with vitest a few times to see the block pass.

## Done when

- The narrative's `## Decisions` records the idle and loaded timings of `search.query` and `adrs.list` on the monorepo.
- `flaiover/src/lib/server/search.test.ts` gives the `adrs on the monorepo` block a 30 s timeout with a comment naming I-0185 and the reason.
- The file's tests pass under `flai test` on the paths changed, and the flaiover-lint tier passes on it.

## Notes

Drafted by the planner. A test that reproduces a timeout under host load does not fit a unit test; the story's notes say so for the first criterion.
