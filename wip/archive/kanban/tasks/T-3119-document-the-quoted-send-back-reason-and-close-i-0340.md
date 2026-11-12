---
id: T-3119
type: task
nature: remediation
title: Document the quoted send-back reason and close I-0340
status: done
parent: S-0711
owner: alex
created: 2026-11-11T09:06:17Z
updated: 2026-11-11T22:03:15Z
transitions:
  - to: ready
    at: 2026-11-11T22:02:39Z
    by: agent-S-0711
  - to: in-progress
    at: 2026-11-11T22:02:47Z
    by: agent-S-0711
  - to: done
    at: 2026-11-11T22:03:15Z
    by: agent-S-0711
stream: S-0711
tags: [docs, issues]
touches: [design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0340-sendbackreason-puts-a-ci-check-name-and-reviewer-login-unquoted-into-the-agent.md, design/issues/summary.md]
after: [T-3118]
usage:
  source: log
  seconds: 28
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 1479
      cache_read: 910156
      cache_write: 23802
      cost: 0.4021
---
# T-3119 Document the quoted send-back reason and close I-0340

## Work

Update the send back's reason example to the quoted form T-3118 writes, in two places:

- `design/system/flai-cli.md`, § The send back: say that each name is quoted in a code span, with its control characters dropped.
- `docs/users/flai.md`, step 2 of the send back.

Then close the issue with `flai issue close I-0340 --reason`, naming the fix and its test. That writes the I-0340 file and `design/issues/summary.md`.

It waits for T-3118 because it quotes the reason's form as T-3118 writes it, and names T-3118's test in the close reason.

## Done when

- Both documents give the quoted reason, as `sendBackReason` writes it.
- I-0340 is closed with a reason that names the fix.
- The markdown lint passes on the two documents.

## Notes
