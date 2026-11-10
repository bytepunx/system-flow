---
id: T-2957
type: task
nature: remediation
title: The agent-container runbook and scripts/README.md say what make agent-images does when the credential helper fails
status: done
parent: S-0713
owner: alex
created: 2026-11-10T10:24:03Z
updated: 2026-11-10T10:56:33Z
transitions:
  - to: ready
    at: 2026-11-10T10:54:28Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T10:54:29Z
    by: system-flow
  - to: done
    at: 2026-11-10T10:56:33Z
    by: agent-S-0713
stream: S-0713
tags: [docs, scripts]
touches: [docs/operators/runbooks/agent-container.md, scripts/README.md]
after: [T-2955]
usage:
  source: log
  seconds: 124
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 3014
      cache_read: 1693145
      cache_write: 30411
      cost: 0.6423
---
# T-2957 The agent-container runbook and scripts/README.md say what make agent-images does when the credential helper fails

## Work

In `docs/operators/runbooks/agent-container.md`, under the section on building the images, say what happens when the Docker configuration names a `credsStore` whose helper does not answer, as Docker Desktop's `desktop.exe` may not from WSL. The build goes on with no login, since the base images are public, prints the line saying so, and leaves `~/.docker/config.json` as it was. Also say how to build with the login once the helper answers again, which is to run `make agent-images` again. In `scripts/README.md`, add the fallback to the `agent-image.sh` row, add a row for `agent-image-test.sh`, and name it in the `smoke.sh` row. It waits for T-2955, whose behaviour it describes. It shares no path with T-2956, so the two run together.

## Done when

- The runbook says when the build pulls with no login and what it prints.
- `scripts/README.md` describes the fallback, the new test, and that smoke runs it.
- The markdown lint passes on both files.

## Notes
