---
id: I-0276
title: A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project
class: efficiency
status: closed
count: 2
cost: 20m
first_reported: 2026-11-07T05:43:57Z
last_reported: 2026-11-10T11:02:42Z
updated: 2026-11-10T11:03:31Z
---

# I-0276 A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project

## Description
A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project

## Instances

### 2026-11-07T05:43:57Z
Story: S-0567.
S-0567, T-2500: a docker run from agent-S-0567's shell with -v /home/alex/git/bytepunx/system-flow:/home/alex/git/bytepunx/system-flow:ro mounted an empty folder, inside and outside Claude Code's sandbox, on Docker Desktop (docker-desktop context), while flai serve's planner container mounted the same path with its files. So the project's own guard hook could not be replayed by hand in the container; only the image's flai guard, which needs no checkout, could.

### 2026-11-10T11:02:42Z
Story: S-0667.
S-0667, T-2926: reproduced from agent-S-0667's zsh shell. The cause is zsh, not Docker: in zsh, `$ROOT:r` is a history modifier even inside double quotes, so `-v "$ROOT:$ROOT:ro"` reaches docker as `<ROOT>:<ROOT>o`. docker inspect shows Binds `[".../system-flow:.../system-flowo"]`, a read-write mount at `<ROOT>o` and nothing at `<ROOT>`, and `-w "$ROOT"` makes `<ROOT>` an empty working directory. `-v "${ROOT}:${ROOT}:ro"` lists the checkout. The docker context, the daemon (Docker Desktop 29.7.2), and DOCKER_HOST change nothing. flai serve and scripts/agent-container.sh build docker's arguments in Go and in sh, so they are not affected. S-0667's Notes, under Cause, hold the commands and what they printed.

## Remediation

Story S-0667 remediates this issue, created from it at 2026-11-10T08:34:40Z.
Closed 2026-11-10T11:03:31Z: S-0667 found the cause in the shell, not Docker. zsh reads `$ROOT:r` in an unbraced `-v "$ROOT:$ROOT:ro"` as its modifier, even inside double quotes, so docker mounted the checkout at `<ROOT>o` and left `<ROOT>` empty. Fixed by T-3078's rule in design/conventions/tooling.md's project additions, which says to write `${name}` before a colon, and by T-2929's note in docs/operators/runbooks/agent-container.md. scripts/agent-container.sh and flai serve build docker's arguments in sh and Go, so they were never affected.
