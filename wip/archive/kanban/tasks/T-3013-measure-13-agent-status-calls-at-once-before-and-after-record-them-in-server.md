---
id: T-3013
type: task
nature: improvement
title: Measure 13 agent.status calls at once before and after, record them in server-performance.md, and close I-0313
status: done
parent: S-0678
owner: alex
created: 2026-11-10T10:31:37Z
updated: 2026-11-11T06:08:19Z
transitions:
  - to: ready
    at: 2026-11-11T05:50:09Z
    by: agent-S-0678
  - to: in-progress
    at: 2026-11-11T05:50:09Z
    by: agent-S-0678
  - to: done
    at: 2026-11-11T06:08:19Z
    by: agent-S-0678
stream: S-0678
tags: [cli, docs]
touches: [design/system/server-performance.md, design/issues/I-0313-agent-status-with-13-projects-at-once-is-not-under-200-ms.md, design/issues/summary.md]
after: [T-3004, T-3008]
usage:
  source: log
  seconds: 1090
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 80
      output: 31341
      cache_read: 5109079
      cache_write: 228622
      cost: 3.0321
---
# T-3013 Measure 13 agent.status calls at once before and after, record them in server-performance.md, and close I-0313

## Work

It waits for T-3004, whose phases say where the time goes, and T-3008, the fix it measures.

Measure as S-0603 measured cause 12 (`design/system/server-performance.md`, § Cause 12 fixed). Build an in-process Go program under `/tmp` against each tree, which does per call what the host's `Agent` does, with a request timer: before is `main` at the story's base, after is the story's head. Read the host's `~/.flai/serve/agents.json` and the projects it names that exist. Run a warm-up round, one sequential call per project, then five rounds of 13 calls released together. Run the builds in turn, twice each, and record the load average. Write nothing of the program's into the repository.

Add a section to `design/system/server-performance.md` after cause 21's: how it was measured, a table of the calls before and after (p50, p90, max, calls under 200 ms, this repository's p50), and the phases inside `agent.activity` from T-3004. Say what is left in the time and why.

When every call answers under 200 ms, or what is left is the wait for a core on a loaded host, which the section shows with the phases, close I-0313 with `flai issue close I-0313 --reason` naming what fixed it. When a part of `Activity` still takes most of a call, do not close it: add a task for that part, which this task then waits for.

## Done when

- The section records both builds' figures with the date, the commits, and the load average.
- I-0313 is closed with its reason, or a task for what is left is open on S-0678 and named in the section.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes
