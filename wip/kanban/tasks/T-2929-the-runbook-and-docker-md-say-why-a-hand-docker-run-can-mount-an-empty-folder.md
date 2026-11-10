---
id: T-2929
type: task
nature: improvement
title: The runbook and docker.md say why a hand docker run can mount an empty folder, and I-0276 is closed
status: backlog
parent: S-0667
owner: alex
created: 2026-11-10T08:49:05Z
updated: 2026-11-10T08:49:05Z
transitions: []
stream: S-0667
tags: [docker, agent-container]
touches: [docs/operators/runbooks/agent-container.md, design/tech/docker.md, design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md, design/issues/summary.md]
after: [T-2927]
---
# T-2929 The runbook and docker.md say why a hand docker run can mount an empty folder, and I-0276 is closed

## Work

Waits for T-2927: it documents the check and the message T-2927 wrote, and closes I-0276 on what T-2927 fixed.

In `docs/operators/runbooks/agent-container.md`, say what makes a hand `docker run` mount the main checkout as an empty folder, as T-2926 found it, how `scripts/agent-container.sh` now refuses it, and what the operator does when it does. In `design/tech/docker.md`, record the fact about the host's Docker that the cause rests on, such as the context or engine whose daemon resolves a bind mount's source, and bump each document's `updated`.

Close the issue with `flai issue close I-0276 --reason "<what fixed it>"`, naming T-2927's check and the fix, which updates `design/issues/summary.md`.

## Done when

- The runbook names the cause, the script's refusal, and what to do about it.
- `design/tech/docker.md` records the fact the cause rests on.
- I-0276 is closed with a reason that says what fixed it.

## Notes
