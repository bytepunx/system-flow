---
id: D-0570
title: A flaiover route test whose first case pays the route's cold import after vi.resetModules() gets a 30 s timeout on its describe block, since the import, not the test's own steps, outlasts vitest's 5 s default under load (I-0341)
kind: cause
sources: [I-0341, S-0669, T-3072]
topics: [testing, flaiover, sveltekit, code]
paths: [flaiover/src/routes/api/agent/agent.test.ts, design/issues/I-0341-agent-test-ts-s-does-not-ask-a-candidate-for-project-info-times-out-at-vitest-s.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0570 A flaiover route test whose first case pays the route's cold import after vi.resetModules() gets a 30 s timeout on its describe block, since the import, not the test's own steps, outlasts vitest's 5 s default under load (I-0341)

I-0341 was closed on 2026-11-10 by S-0669 (T-3072). The test "does not ask a candidate for project.info", in `flaiover/src/routes/api/agent/agent.test.ts`, timed out at vitest's 5,000 ms default and stopped two of S-0669's close-outs at the flaiover tier, alone of 2,432 tests, though S-0669 changes neither the route nor its test.

The cause, from I-0341's second instance: "The test is the file's first: after vi.resetModules() it pays the cold import of ./+server and $lib/server/agent's importOriginal, which alone outlasts 5 s under that load." The host's load average was 61 to 88, with seven stories' agents and close-outs at work.

The fix, from I-0341's Remediation: "The GET /api/agent describe in flaiover/src/routes/api/agent/agent.test.ts carries { timeout: 30000 }, as search.test.ts's monorepo describe does." S-0669's Decisions (`wip/archive/agents/S-0669.md`) call it a five-minute fix, so it was added as a task on the story that met it rather than left to a story of its own.

This is a third case beside D-0544 and D-0552. D-0544 batches settles a test can cut; D-0552 raises a block's timeout where the time is the page's own mount. Here the time is the module import the test cannot cut, so the describe's timeout is raised. No deterministic test reproduces the load, so the two stopped close-outs are the evidence (D-0550). The same shape stopped close-outs on `search.test.ts` (I-0185) and `repo.test.ts` (I-0269, S-0695).

Sources: I-0341; S-0669 (`wip/archive/agents/S-0669.md`, Decisions); T-3072.
