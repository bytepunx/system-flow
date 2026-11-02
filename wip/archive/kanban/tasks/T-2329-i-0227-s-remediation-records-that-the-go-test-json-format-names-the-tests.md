---
id: T-2329
type: task
nature: remediation
title: I-0227's Remediation records that the go-test-json format names the tests running at a timeout, and the mcpserver package passes flai test here
status: done
parent: S-0540
owner: alex
created: 2026-10-29T07:37:46Z
updated: 2026-11-02T06:46:03Z
transitions:
  - to: ready
    at: 2026-11-02T06:42:48Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T06:42:49Z
    by: system-flow
  - to: done
    at: 2026-11-02T06:46:03Z
    by: agent-S-0540
stream: S-0540
tags: [flai, issues]
touches: [design/issues/I-0227-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md]
after: [T-2326]
usage:
  source: log
  seconds: 194
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1603
      cache_read: 822565
      cache_write: 20536
      cost: 0.3609
---
# T-2329 I-0227's Remediation records that the go-test-json format names the tests running at a timeout, and the mcpserver package passes flai test here

## Work

Criterion 2. S-0548 closed I-0227 on 2026-10-31 with `flai issue close`, its reason naming the tier's `-timeout 20m`, the `-short` gate on the two prime tests, and the parallel package. That close left out the issue's second cause: a timed-out package's finding names only the package. T-2326 fixes that. `flai issue close` refuses an issue already closed, so record the fix in the issue's `## Remediation`, whose text is edited freely (`continuous-improvement.md`), as one line after S-0548's close line.

Then run `flai test flai/internal/mcpserver` once in the worktree and record its time and result in the narrative. That shows the cause no longer occurs on this host, for criterion 1.

It waits for T-2326: the line names what T-2326 changed.

## Done when

- I-0227 stays closed, and its `## Remediation` names S-0540's T-2326 and what it changed.
- `flai test flai/internal/mcpserver` passed in the worktree, and its time is in the narrative.

## Notes
