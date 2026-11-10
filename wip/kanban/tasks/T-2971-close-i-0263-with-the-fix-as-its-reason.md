---
id: T-2971
type: task
nature: remediation
title: Close I-0263 with the fix as its reason
status: backlog
parent: S-0700
owner: alex
created: 2026-11-10T10:28:28Z
updated: 2026-11-10T10:28:28Z
transitions: []
stream: S-0700
tags: [dashboard]
touches: [design/issues/I-0263-the-new-item-form-s-lint-refusal-is-easy-to-miss-so-a-refused-epic-looks-like-a.md, design/issues/summary.md]
after: [T-2966, T-2968]
---
# T-2971 Close I-0263 with the fix as its reason

## Work

Run `flai issue close I-0263 --reason "<reason>"`, the reason naming S-0700 and what fixed it: the new item form scrolls to and focuses its refusal and marks the Create button refused, with the test that reproduces the missed refusal in `flaiover/src/lib/components/NewItemForm.svelte.test.ts`. flai regenerates `design/issues/summary.md`.

Waits for T-2966 and T-2968: the issue closes only once the fix and its documentation are in.

## Done when

- I-0263 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes
