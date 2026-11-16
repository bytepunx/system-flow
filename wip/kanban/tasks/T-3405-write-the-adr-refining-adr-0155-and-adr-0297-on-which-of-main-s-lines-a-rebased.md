---
id: T-3405
type: task
nature: remediation
title: Write the ADR refining ADR-0155 and ADR-0297 on which of main's lines a rebased branch undoes
status: done
parent: S-0800
owner: alex
created: 2026-11-16T01:44:25Z
updated: 2026-11-16T01:48:22Z
transitions:
  - to: ready
    at: 2026-11-16T01:47:31Z
    by: agent-S-0800
  - to: in-progress
    at: 2026-11-16T01:47:32Z
    by: agent-S-0800
  - to: done
    at: 2026-11-16T01:48:22Z
    by: agent-S-0800
stream: S-0800
tags: [adr, storygit]
touches: [design/adrs/README.md, design/adrs/0309-a-line-main-added-is-undone-only-when-the-rebased-branch-removes-it-at-the-line.md]
usage:
  source: log
  seconds: 50
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 12
      output: 104
      cache_read: 1372441
      cache_write: 6843
      cost: 0.5226
---
# T-3405 Write the ADR refining ADR-0155 and ADR-0297 on which of main's lines a rebased branch undoes

## Work

Settle the rule I-0385 asks for in a new ADR, with `flai adr new`, refining ADR-0155 and ADR-0297. Choose between the issue's two options and record the other under Alternatives considered:

- By position (recommended): a line main added since the branch's old base is undone only when the branch's net change against main removes it at a line main added in main's file. Read the positions from the `-U0` hunk headers of `git diff oldBase onto` and `git diff onto HEAD`, rather than counting the file's lines as one multiset of text. It also settles other lines with no change of their own, such as `return nil` or `default:`.
- Narrow: a line that holds only brackets, braces, or white space is never counted alone.

Say how the rule applies to `RecordedUndone`, which checks the record against main now, past `onto`. Keep the I-0338 case, a line only a commit after `onto` adds, out of the count. Say how the rule sits beside ADR-0297's markdown words rule. Show that the revert of I-0124, replayed over main's identical change, is still caught. Name the file the ADR is written to in the story's `### Planning` note, narrowing the story's `design/adrs/` touch to it with `flai touches`.

This task waits for nothing: it is the first layer, and every other task reads its decision.

## Done when

- The ADR is accepted, refines ADR-0155 and ADR-0297, and states the rule, its reach over `FindUndone` and `RecordedUndone`, and why I-0124's revert is still caught.
- `design/adrs/README.md` lists it, and `flai check --strict` is clean.

## Notes
