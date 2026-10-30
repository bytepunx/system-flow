---
id: T-2193
type: task
nature: remediation
title: search.test.ts reads the ADR statuses it allows from the template's status comment instead of its own list
status: done
parent: S-0512
owner: alex
created: 2026-10-28T08:28:10Z
updated: 2026-10-30T14:15:44Z
transitions:
  - to: ready
    at: 2026-10-30T14:13:46Z
    by: agent-S-0512
  - to: in-progress
    at: 2026-10-30T14:13:47Z
    by: agent-S-0512
  - to: done
    at: 2026-10-30T14:15:44Z
    by: agent-S-0512
stream: S-0512
tags: [dashboard]
touches: [flaiover/src/lib/server/search.test.ts]
after: [T-2191]
usage:
  source: log
  seconds: 117
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 2779
      cache_read: 623513
      cache_write: 27393
      cost: 0.3653
---
# T-2193 search.test.ts reads the ADR statuses it allows from the template's status comment instead of its own list

## Work

In `flaiover/src/lib/server/search.test.ts`, the test `lists ADRs with the supersession chain` keeps its own list, `['proposed', 'accepted', 'superseded', 'deprecated', 'rejected']`, which is what failed when ADR-0202 became `rejected` (I-0208). Replace the hard-coded list by one read from the monorepo's `design/adrs/0000-template.md`: the `status:` line's comment, `# proposed | accepted | superseded | deprecated | rejected`, split on `|` and trimmed. Assert that the list read has at least the four the test had, so that a broken read fails loudly rather than allowing everything, and keep the comment that names ADR-0202 as the rejected one. Everything else in the test stays.

This task waits for T-2191, which fixes the comment line's shape in both template copies.

## Done when

- The test takes its allowed statuses from `design/adrs/0000-template.md`'s status comment and has no list of its own.
- The vitest tier passes on the monorepo (`flai test flaiover/src/lib/server/search.test.ts`), with `rejected` among the statuses read.
- prettier and eslint pass on the file.

## Notes

Drafted by planner-S-0512. The test runs against the monorepo (`monorepo` in the file), so the path is the repository's own `design/adrs/0000-template.md`, resolved from the test's location as the file already resolves `monorepo`.
