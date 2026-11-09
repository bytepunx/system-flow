---
id: T-2736
type: task
nature: remediation
title: The operators guide, the runbooks, the user guide, and the CLI and Docker design describe the refusal, and I-0254 is closed
status: done
parent: S-0619
owner: alex
created: 2026-11-09T17:03:41Z
updated: 2026-11-09T17:21:44Z
transitions:
  - to: ready
    at: 2026-11-09T17:20:44Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T17:20:44Z
    by: system-flow
  - to: done
    at: 2026-11-09T17:21:44Z
    by: agent-S-0619
stream: S-0619
tags: [docs, flai]
touches: [docs/operators/index.md, docs/operators/settings.md, docs/operators/runbooks/install.md, docs/operators/runbooks/agent-container.md, docs/users/flai.md, design/system/flai-cli.md, design/tech/docker.md, design/issues/I-0254-an-agent-container-on-a-native-linux-docker-engine-cannot-reach-flai-s-http-mcp.md, design/issues/summary.md]
after: [T-2735]
usage:
  source: log
  seconds: 60
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 2534
      cache_read: 1502238
      cache_write: 29485
      cost: 0.5812
---
# T-2736 The operators guide, the runbooks, the user guide, and the CLI and Docker design describe the refusal, and I-0254 is closed

## Work

This task waits for T-2735, because the documents quote the refusal T-2735 writes and the conditions it checks.

- In `docs/operators/index.md` § "MCP over HTTP", in the paragraph on `flai serve mcp bind`, say that `flai serve` now refuses a container start on an engine that is not Docker Desktop while `mcp.bind` is loopback. Give the refusal's words and the fix it names.
- Say the same, shorter, in:
  - the `mcp.bind` row of `docs/operators/settings.md`;
  - `docs/operators/runbooks/install.md`, where it tells a native engine's operator to set the bind;
  - `docs/operators/runbooks/agent-container.md`;
  - `docs/users/flai.md` § "Serving agents over MCP".
- In `design/system/flai-cli.md`, where it says a native Linux engine needs `mcp.bind`, describe the preflight: how it finds the engine's kind, when it refuses, and that a `docker info` that cannot answer does not refuse.
- In `design/tech/docker.md`'s "Agent harness" row, add the two `docker` calls the preflight makes.
- Close the issue once T-2735 is done: `flai issue close I-0254 --reason "<what fixed it>"`. Name T-2735's refusal and its test, and say that no native engine was measured on this host. flai rewrites `design/issues/summary.md`.

## Done when

- Each document above describes the refusal, and none still says that the first MCP call fails with nothing to say why.
- `flai issue close I-0254` has run with a reason that names the refusal and its test.
- `flai test` on the changed paths passes, the markdown lint and `settings-index` among them.

## Notes
