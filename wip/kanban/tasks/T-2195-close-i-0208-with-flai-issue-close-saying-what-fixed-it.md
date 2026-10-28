---
id: T-2195
type: task
nature: remediation
title: Close I-0208 with flai issue close saying what fixed it
status: backlog
parent: S-0512
owner: alex
created: 2026-10-28T08:28:41Z
updated: 2026-10-28T08:28:41Z
transitions: []
stream: S-0512
tags: [cli]
touches: [design/issues/I-0208-search-test-ts-lists-adrs-with-the-supersession-chain-fails-on-main-since-adr-0202-s-status-became-rejected-a-status-outside-the-four-the-test-and-the-adr-template-name.md, design/issues/summary.md]
after: [T-2192, T-2193, T-2194]
---
# T-2195 Close I-0208 with flai issue close saying what fixed it

## Work

From the story's worktree, run `flai issue close I-0208 --reason "..."` with a reason that says what fixed it: `rejected` is the fifth ADR status, named in both ADR template copies and the documentation standard (T-2191), immutable to `flai doc` as an accepted one (T-2192), `flai check` reports `adr.status` on a status outside the five (T-2190), and `search.test.ts` reads its allowed statuses from the template's comment (T-2193), with the ADR T-2189 recorded. The command writes the issue file and regenerates `design/issues/summary.md`; commit both with `flai task done`.

Then verify each acceptance criterion of S-0512 and tick it with `flai criteria tick`: criterion 1 by the check test of T-2190 and the vitest run of T-2193, criterion 2 by the closed issue.

This task waits for T-2192, T-2193, and T-2194: the issue is closed only once every part of the fix is on the branch.

## Done when

- I-0208's status is `closed` with the reason recorded, and `design/issues/summary.md` no longer lists it as open.
- Both criteria of S-0512 are ticked through flai.
- `flai check --strict` reports nothing new.

## Notes

Drafted by planner-S-0512. `design/issues` is in the manifest's `claims.shared`, so the two files hold no other story; `flai stream sync` regenerates `summary.md` when a rebase stops on it alone (ADR-0098).
