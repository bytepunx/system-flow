---
id: T-2140
type: task
nature: remediation
title: Close I-0185 with the reason, note why no reproducing test fits, and tick the criteria
status: done
parent: S-0501
owner: alex
created: 2026-10-26T04:31:33Z
updated: 2026-10-26T04:36:16Z
transitions:
  - to: ready
    at: 2026-10-26T04:35:59Z
    by: agent-S-0501
  - to: in-progress
    at: 2026-10-26T04:36:00Z
    by: agent-S-0501
  - to: done
    at: 2026-10-26T04:36:16Z
    by: agent-S-0501
stream: S-0501
tags: [dashboard]
touches: [design/issues/I-0185-search-test-ts-finds-a-body-phrase-across-wip-and-design-times-out-at-5-s-on-the-monorepo-while-other-stories-close-outs-run.md, design/issues/summary.md]
after: [T-2139]
usage:
  source: log
  seconds: 16
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1246
      cache_read: 449122
      cache_write: 13320
      cost: 0.2213
---
# T-2140 Close I-0185 with the reason, note why no reproducing test fits, and tick the criteria

## Work

Waits for T-2139, since the reason names what it measured and changed. Run `flai issue close I-0185 --reason` in the story worktree, the reason saying that the test's one-shot `flai hostapi` builds the whole repository's index on every call, giving the measured figures, and that the `adrs on the monorepo` block now has a 30 s timeout as `writes.test.ts` has had since S-0116; the command writes the issue file and `design/issues/summary.md`. Add a line to the story's `## Notes`, through `flai edit`, saying that no test reproduces a timeout under host load, so the first criterion is met by the cause being removed rather than by a reproducing test. Then tick both criteria with `flai criteria tick S-0501 <n>` once verified.

## Done when

- I-0185 is closed with a reason that names the cause, the figures, and the timeout, and `design/issues/summary.md` no longer lists it as open.
- The story's `## Notes` says why no reproducing test fits.
- Both acceptance criteria of S-0501 are ticked.
- `flai check --strict` is clean for the story.

## Notes

Drafted by the planner. Issue commands write under the current checkout, so run them in the story worktree.
