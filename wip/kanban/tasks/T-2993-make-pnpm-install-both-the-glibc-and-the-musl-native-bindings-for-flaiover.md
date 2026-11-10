---
id: T-2993
type: task
nature: remediation
title: Make pnpm install both the glibc and the musl native bindings for flaiover
status: backlog
parent: S-0706
owner: alex
created: 2026-11-10T10:30:28Z
updated: 2026-11-10T10:30:46Z
transitions: []
stream: S-0706
tags: [flaiover]
touches: [flaiover/package.json]
---
# T-2993 Make pnpm install both the glibc and the musl native bindings for flaiover

## Work

It waits for no task: it is the first layer, beside T-2994, which changes other files.

Set pnpm's `supportedArchitectures` for flaiover so that one install puts in the native bindings for both libcs: `os: [current]`, `cpu: [current]`, `libc: [glibc, musl]`. Put it in the `pnpm` field of `flaiover/package.json`. The pinned pnpm, 10.34.5, reads `supportedArchitectures` from that field. `system-flow.yaml`'s accept `install` already watches `package.json`. If pnpm 10.34.5 also takes the setting from `flaiover/pnpm-workspace.yaml` and prefers it there, it may go there instead. In that case, add `flaiover/pnpm-workspace.yaml` to this task's touches. Also add it to the `paths` of the accept `install` and of the `svelte-check` tier in `system-flow.yaml`.

`flaiover/pnpm-lock.yaml` already records `@rolldown/binding-linux-x64-gnu` and `-musl`, and the musl variants of rollup, lightningcss, and tailwind oxide. Check that the setting leaves the lockfile as it is: `pnpm install --frozen-lockfile` must still pass.

## Done when

- `flaiover/package.json` (or `flaiover/pnpm-workspace.yaml`) declares `supportedArchitectures` with `libc` holding both `glibc` and `musl`.
- A fresh `pnpm install --frozen-lockfile` in `flaiover` leaves both `@rolldown/binding-linux-x64-gnu` and `@rolldown/binding-linux-x64-musl` under `node_modules/.pnpm`. The lockfile is unchanged.
- `flai test flaiover/package.json` passes.

## Notes
