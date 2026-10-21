---
id: T-1953
type: task
nature: feature
title: flaiover has the icon library TH-0553 settles, recorded in design/tech
status: backlog
parent: S-0458
owner: alex
created: 2026-10-21T07:01:37Z
updated: 2026-10-21T07:01:37Z
transitions: []
stream: S-0458
tags: [dashboard]
touches: [flaiover/package.json, flaiover/pnpm-lock.yaml, design/tech/icons.md, design/tech/README.md]
---
# T-1953 flaiover has the icon library TH-0553 settles, recorded in design/tech

## Work

- With `@lucide/svelte`, the recommendation on TH-0553 question 2, add it to `flaiover/package.json` with `pnpm add`, so that `pnpm-lock.yaml` follows.
- Write `design/tech/icons.md`: the library, its version, its licence, where it is used (the nav rail), and why. List it in `design/tech/README.md`.
- If the operator chose hand-written SVG instead, add no dependency. Record in `design/tech/icons.md` that the icons are inline components under `flaiover/src/lib/icons/`, and widen this task's touches to that folder.
- Waits for nothing. It runs beside the rail's data.

## Done when

The library resolves in a Svelte component under vitest, `pnpm install --frozen-lockfile` passes, and `design/tech` names it.

## Notes
