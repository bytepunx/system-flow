---
id: I-0341
title: "agent.test.ts's \"does not ask a candidate for project.info\" times out at vitest's 5 s default under load and stops a close-out that does not touch it"
class: defect
status: open
count: 1
cost: 21m
first_reported: 2026-11-10T10:16:25Z
last_reported: 2026-11-10T10:16:25Z
updated: 2026-11-10T10:16:25Z
---

# I-0341 agent.test.ts's "does not ask a candidate for project.info" times out at vitest's 5 s default under load and stops a close-out that does not touch it

## Description
agent.test.ts's "does not ask a candidate for project.info" times out at vitest's 5 s default under load and stops a close-out that does not touch it

## Instances

### 2026-11-10T10:16:25Z
Story: S-0669.
S-0669's close-out (job S-0669-20261110T095525Z, head e6fe6df) stopped at the flaiover tier after 20m44s. One test of 2432 failed: flaiover/src/routes/api/agent/agent.test.ts > GET /api/agent > does not ask a candidate for project.info, "Test timed out in 5000ms". S-0669 changes neither that route nor its test. The run's sums show a loaded host: import 1332.93s and transform 721.97s, for a 106.50s run, with seven stories' agents and close-outs at work. This is the same shape as I-0269 (repo.test.ts) and I-0185 (search.test.ts): a test at vitest's 5 s default stops an unrelated close-out under load.

## Remediation
