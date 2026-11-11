---
id: T-2993
type: task
nature: remediation
title: Make pnpm install both the glibc and the musl native bindings for flaiover
status: backlog
parent: S-0706
owner: alex
created: 2026-11-10T10:30:28Z
updated: 2026-11-11T08:54:15Z
transitions: []
stream: S-0706
tags: [flaiover]
touches: [flaiover/package.json, flaiover/pnpm-workspace.yaml]
---
# T-2993 Make pnpm install both the glibc and the musl native bindings for flaiover

## Work

It waits for no task: it is the first layer, beside T-2994, which changes other files.

Set pnpm's `supportedArchitectures` for flaiover so that one install puts in the native bindings for both libcs: `os: [current]`, `cpu: [current]`, `libc: [glibc, musl]`. Put it in `flaiover/pnpm-workspace.yaml`, which already holds flaiover's other pnpm setting, `onlyBuiltDependencies`, so that its pnpm settings stay in one file. Use the `pnpm` field of `flaiover/package.json` only if the pinned pnpm, 10.34.5, does not read `supportedArchitectures` from the workspace file, and say why in this task's Notes. T-2994 adds `pnpm-workspace.yaml` to the paths of flaiover's accept `install` in `system-flow.yaml`, so this task leaves the manifest alone.

`flaiover/pnpm-lock.yaml` already records `@rolldown/binding-linux-x64-gnu` and `-musl`, and the musl variants of rollup, lightningcss, and tailwind oxide. Check that the setting leaves the lockfile as it is: `pnpm install --frozen-lockfile` must still pass.

## Done when

- `flaiover/pnpm-workspace.yaml` (or `flaiover/package.json`) declares `supportedArchitectures` with `libc` holding both `glibc` and `musl`.
- A fresh `pnpm install --frozen-lockfile` in `flaiover` leaves both `@rolldown/binding-linux-x64-gnu` and `@rolldown/binding-linux-x64-musl` under `node_modules/.pnpm`. The lockfile is unchanged.
- `flai test` on the changed file passes.

## Notes
