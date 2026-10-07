---
id: T-1229
type: task
nature: feature
title: The Updates page shows a flai release this binary cannot verify as not installable and offers no way to install it
status: backlog
parent: S-0233
owner: alex
created: 2026-10-07T22:40:40Z
updated: 2026-10-07T22:40:40Z
transitions: []
stream: S-0233
tags: [dashboard]
touches: [flaiover/src/lib/components/HostProcesses.svelte, flaiover/src/lib/components/HostProcesses.svelte.test.ts, docs/users/flaiover.md]
after: [T-1228]
---
# T-1229 The Updates page shows a flai release this binary cannot verify as not installable and offers no way to install it

## Work

In `flaiover/src/lib/components/HostProcesses.svelte`, add the mark T-1228 gives each release to the inline `Release` type. Show a release with it as not installable, with the reason, and give it no install control, so the page cannot ask the host to install it. Cover both kinds of release, one verifiable and one not, in `HostProcesses.svelte.test.ts`. Say in the Updates section of `docs/users/flaiover.md` what the mark means.

Waits for T-1228: the page reads the field name and the reason that task gives the JSON of `flai host versions`.

## Done when

- A release marked as not verifiable shows as not installable, with its reason, and has no install control.
- Other releases install as before.
- The component test covers both, and `flai test` passes on `flaiover`.
- `docs/users/flaiover.md` describes the mark.

## Notes
