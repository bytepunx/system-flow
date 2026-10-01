---
id: T-0648
type: task
nature: feature
title: The dashboard image has a HEALTHCHECK and flai dashboard status says running, not answering, or gone
status: done
parent: S-0184
owner: arobson
created: 2026-10-01T08:55:54Z
updated: 2026-10-01T08:58:25Z
transitions:
  - to: ready
    at: 2026-10-01T08:56:22Z
    by: agent-S-0184
  - to: in-progress
    at: 2026-10-01T08:56:22Z
    by: agent-S-0184
  - to: done
    at: 2026-10-01T08:58:25Z
    by: agent-S-0184
stream: S-0184
tags: []
usage:
  source: log
  seconds: 123
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 36
      output: 279
      cache_read: 2568177
      cache_write: 30160
      cost: 1.0278
---

# T-0648 The dashboard image has a HEALTHCHECK and flai dashboard status says running, not answering, or gone

## Work

- `flaiover/Dockerfile`: a `HEALTHCHECK` that asks `/_health` on the container's own port.
- `flai/cmd`: a probe of the running container's published address (`docker inspect`'s host IP and port, loopback for `0.0.0.0`), and `flai dashboard status` saying `running`, `not answering`, or `gone`, with `state` in `--json` and Docker's own health beside it.

## Done when

- The image's HEALTHCHECK is in the Dockerfile.
- `flai dashboard status` and its `--json` report the three states, and behaviour tests in `flai/cmd` cover each with a fake runner and probe.

## Notes
