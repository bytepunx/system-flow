---
id: T-2751
type: task
nature: remediation
title: The user and operator docs say a usage-limit refusal is checked again, and I-0296 is closed
status: done
parent: S-0629
owner: alex
created: 2026-11-09T17:09:13Z
updated: 2026-11-09T17:50:32Z
transitions:
  - to: ready
    at: 2026-11-09T17:46:25Z
    by: agent-S-0629
  - to: in-progress
    at: 2026-11-09T17:46:26Z
    by: agent-S-0629
  - to: done
    at: 2026-11-09T17:50:32Z
    by: agent-S-0629
stream: S-0629
tags: [flai, docs, issues]
touches: [docs/users/flai.md, docs/operators/index.md, design/issues/I-0296-flai-serve-records-a-claude-code-usage-limit-refusal-as-a-failed-permission.md, design/issues/summary.md]
after: [T-2750]
usage:
  source: log
  seconds: 246
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 3480
      cache_read: 1903383
      cache_write: 53750
      cost: 0.7931
---
# T-2751 The user and operator docs say a usage-limit refusal is checked again, and I-0296 is closed

## Work

Bring the outward docs in line with the behaviour T-2750 ships, then close the issue. This task waits for T-2750 because the docs and the close reason describe the code as it ends up.

- `docs/users/flai.md`
  - § Writes to protected paths says "flai never checks a version again through the same provider, passed or failed". Say that a check refused on the usage limit records nothing and is checked again after the reset, and link T-2749's ADR.
  - Make the same change in § In a container: the container harness.
- `docs/operators/index.md` § Starting an agent when a story becomes ready: say the same where it describes `claude-checks.json`.
- Close I-0296 with `flai issue close I-0296 --reason "<what fixed it>"`. Name the ADR, the rule, and the test that reproduces the 2026-10-26 refusal. That updates `design/issues/summary.md`.

## Done when

- Neither doc says a usage-limit refusal is final, and both link the new ADR.
- I-0296 is closed with a reason naming the fix and its test.
- `flai check --strict` and the markdown lint pass.

## Notes
