---
id: I-0259
title: A story worktree's flaiover/node_modules installed inside the Alpine agent container holds only musl native bindings, so the host's flaiover tier then fails, and the reverse
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-11-02T19:22:35Z
last_reported: 2026-11-02T19:22:35Z
updated: 2026-11-11T09:24:20Z
---

# I-0259 A story worktree's flaiover/node_modules installed inside the Alpine agent container holds only musl native bindings, so the host's flaiover tier then fails, and the reverse

## Description
A story worktree's flaiover/node_modules installed inside the Alpine agent container holds only musl native bindings, so the host's flaiover tier then fails, and the reverse

## Instances

### 2026-11-02T19:22:35Z
Story: S-0510.
S-0510 ran `scripts/flai.sh test` inside flai-agent-sf-claude-code:local on its worktree, and the flaiover tiers passed there. That run installed `flaiover/node_modules` in the worktree, which the host and the container share, and pnpm on Alpine put in only `@rolldown/binding-linux-x64-musl`. The host's close-out then stopped at the flaiover tier: rolldown could not load a binding on glibc. Removing `flaiover/node_modules` let the host install again. A host install would leave the container without the musl binding the same way. A story whose agent runs in a container and whose close-out runs on the host meets this whenever a flaiover tier runs in both. A remedy could install both libc variants (pnpm's `supportedArchitectures` with `libc: [glibc, musl]`), or keep node_modules apart per platform.

## Remediation

Story S-0706 remediates this issue, created from it at 2026-11-10T09:11:37Z.
Closed 2026-11-11T09:24:20Z: S-0706: flaiover/pnpm-workspace.yaml sets supportedArchitectures with libc glibc and musl (T-2993), so one node_modules holds both native bindings for the host and the Alpine agent container; scripts/flaiover-install.sh --if-needed installs again when pnpm-lock.yaml, package.json, pnpm-workspace.yaml, or .npmrc is newer than the last install, and flai accept's install watches them (T-2994); scripts/flaiover-install-test.sh reproduces both, offline (T-2997).
