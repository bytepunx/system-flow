---
id: T-3072
type: task
nature: remediation
title: agent.test.ts's GET /api/agent describe gets a 30 s timeout, as search.test.ts's monorepo describe has, and I-0341 is closed
status: done
parent: S-0669
owner: alex
created: 2026-11-10T10:36:23Z
updated: 2026-11-10T10:38:50Z
transitions:
  - to: ready
    at: 2026-11-10T10:36:28Z
    by: agent-S-0669
  - to: in-progress
    at: 2026-11-10T10:36:30Z
    by: agent-S-0669
  - to: done
    at: 2026-11-10T10:38:50Z
    by: agent-S-0669
stream: S-0669
tags: [flaiover, test]
touches: [flaiover/src/routes/api/agent/agent.test.ts, design/issues/I-0341-agent-test-ts-s-does-not-ask-a-candidate-for-project-info-times-out-at-vitest-s.md, design/issues/summary.md]
after: [T-2912]
usage:
  source: log
  seconds: 140
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 354
      cache_read: 829432
      cache_write: 2524
      cost: 0.1928
---
# T-3072 agent.test.ts's GET /api/agent describe gets a 30 s timeout, as search.test.ts's monorepo describe has, and I-0341 is closed

## Work

This task waits for T-2912 only because it came after it: S-0669's close-out stopped twice on I-0341, a test the story does not change.

`flaiover/src/routes/api/agent/agent.test.ts`'s first test pays the cold `await import('./+server')` after `vi.resetModules()`, with `importOriginal` of `$lib/server/agent`. Under a load average of 60 to 90 that import alone outlasts vitest's 5 s default. Give the `GET /api/agent` describe `{ timeout: 30000 }`, as `search.test.ts`'s `adrs on the monorepo` describe has, with a comment naming I-0341. Then close I-0341 naming the fix.

No deterministic test reproduces the load (D-0550). The two close-outs' timeouts are the evidence.

## Done when

- [ ] The describe carries a 30 s timeout with a comment naming I-0341.
- [ ] `flai test` passes on the file.
- [ ] I-0341 is closed with the fix as its reason.

## Notes
