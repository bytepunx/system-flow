---
id: T-1011
type: task
nature: research
title: Choose how flai keeps design/issues/summary.md from stopping a sync, from I-0074's instances, and record it in an ADR
status: done
parent: S-0278
owner: alex
created: 2026-10-06T11:31:58Z
updated: 2026-10-06T19:47:35Z
transitions:
  - to: ready
    at: 2026-10-06T19:46:57Z
    by: agent-S-0278
  - to: in-progress
    at: 2026-10-06T19:46:57Z
    by: agent-S-0278
  - to: done
    at: 2026-10-06T19:47:35Z
    by: agent-S-0278
stream: S-0278
tags: [flai]
touches: [design/adrs]
usage:
  source: log
  seconds: 38
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 10
      output: 3380
      cache_read: 396852
      cache_write: 14703
      cost: 0.2364
---
# T-1011 Choose how flai keeps design/issues/summary.md from stopping a sync, from I-0074's instances, and record it in an ADR

## Work

The story asks for a solution proposed from I-0074's instances before it is built. Read I-0074's three instances, and I-0089, which has the same cause. Then weigh these directions:

- `flai stream sync` regenerates `summary.md` itself when a rebase stops on it alone. `flai accept` gets the same, since `mergeStoryBranch` syncs through `syncStoryBranch` (`flai/cmd/branch.go`). The trial merge in `flai/cmd/stream_sync.go` then stops counting `summary.md` as a conflict. This is the planner's recommendation. The file is derived from the issue files, which merge cleanly when two stories record different issues.
- A git merge driver named in `.gitattributes`. A driver needs `git config` in every clone, and flai installs none (`flai init` and `flai upgrade` write no git config), so this alone fixes nothing for the operator's checkout.
- Dropping the `updated:` line that `issues.Summary` writes (`flai/internal/issues/issues.go`). It narrows the conflict but does not remove it: the first instance's removed rows sat two lines apart, and those hunks meet anyway.
- Not committing `summary.md` and generating it where it is read (I-0089's second direction). `flai prime`, `flai check`'s staleness warning, and the dashboard read it, so this is wider than the story.

Write an ADR under `design/adrs` with `flai adr new`. It records the direction chosen, the paths it covers, and what happens when other paths conflict alongside `summary.md`. It also records what it leaves out: two branches bumping the same issue file, as I-0073 was in the third instance, still conflict in its front matter.

This task waits for none. It comes first because the other tasks build what it decides.

## Done when

- An ADR under `design/adrs` records the chosen fix, the alternatives above and why each was not taken, and what it leaves out.
- `flai check --strict` passes on the ADR.

## Notes
