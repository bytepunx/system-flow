---
id: T-2850
type: task
nature: improvement
title: Close I-0328 naming the lock the checked edit now holds
status: backlog
parent: S-0654
owner: alex
created: 2026-11-10T07:01:30Z
updated: 2026-11-10T07:01:30Z
transitions: []
stream: S-0654
tags: [flai, issues]
touches: [design/issues/I-0328-a-checked-edit-of-a-work-item-writes-checks-and-restores-it-with-no-lock-so-a.md, design/issues/summary.md]
after: [T-2848]
---
# T-2850 Close I-0328 naming the lock the checked edit now holds

## Work

Run `flai issue close I-0328 --reason "<reason>"`, the reason saying that `itemedit.Apply` and `itemedit.ApplyForecasts` now hold the checked-write lock from their first save until they keep or restore the edit, so flai serve's look reads no edit the check refuses, naming the new ADR and the tests T-2844 and T-2847 added. It writes the issue's file and `design/issues/summary.md`.

It waits for T-2848: the reason names the ADR that task writes, by its number.

## Done when

- I-0328 is closed with that reason, and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` passes.

## Notes
