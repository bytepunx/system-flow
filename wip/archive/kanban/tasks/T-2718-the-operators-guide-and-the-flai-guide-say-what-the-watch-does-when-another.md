---
id: T-2718
type: task
nature: remediation
title: The operators guide and the flai guide say what the watch does when another process holds the dashboard's port
status: done
parent: S-0615
owner: alex
created: 2026-11-09T09:52:25Z
updated: 2026-11-09T10:30:05Z
transitions:
  - to: ready
    at: 2026-11-09T10:23:33Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T10:23:33Z
    by: system-flow
  - to: done
    at: 2026-11-09T10:30:05Z
    by: agent-S-0615
stream: S-0615
tags: [flai, dashboard, docs]
touches: [docs/operators/index.md, docs/users/flai.md]
after: [T-2714, T-2717]
usage:
  source: log
  seconds: 392
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 1513
      cache_read: 1306254
      cache_write: 18003
      cost: 0.4071
---
# T-2718 The operators guide and the flai guide say what the watch does when another process holds the dashboard's port

## Work

- `docs/operators/index.md` § The dashboard's watch: in **How it looks**, the fourth answer, port held; in **When it restarts it**, that it does not when another process holds the recorded port, and logs it once; add a bullet, **When another process holds the port**, saying what `flai dashboard status` and `flai host status` show, how to free the port (stop the `docker-proxy` they name and the orphan's `containerd-shim`, or restart the Docker daemon), and that the watch then starts the dashboard on its own. Name I-0243's case: a container Docker no longer lists, even with `docker ps --all`, whose processes still run. Link the ADR T-2714 wrote. Keep the "Whether everything runs" row (line 17 or so) in step: running, not answering, port held, or gone.
- `docs/users/flai.md`: the `flai dashboard status` line in the command listing says the fourth state.

It waits for T-2714, whose ADR it links, and T-2717, whose messages it quotes.

## Done when

- Both files say the `port-held` state in the words the commands print, link the new ADR from the operators guide, and bump their `updated`.
- `flai test docs/operators/index.md docs/users/flai.md` passes.

## Notes
