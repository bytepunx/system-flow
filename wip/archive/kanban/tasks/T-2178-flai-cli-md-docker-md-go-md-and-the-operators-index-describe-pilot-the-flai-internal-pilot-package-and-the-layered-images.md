---
id: T-2178
type: task
nature: feature
title: flai-cli.md, docker.md, go.md, and the operators' index describe pilot, the flai/internal/pilot package, and the layered images
status: done
parent: S-0507
owner: alex
created: 2026-10-26T04:52:20Z
updated: 2026-10-28T23:14:41Z
transitions:
  - to: ready
    at: 2026-10-28T23:11:17Z
    by: agent-S-0507
  - to: in-progress
    at: 2026-10-28T23:11:18Z
    by: agent-S-0507
  - to: done
    at: 2026-10-28T23:14:41Z
    by: agent-S-0507
stream: S-0507
tags: [cli, docker]
touches: [design/system/flai-cli.md, design/tech/docker.md, design/tech/go.md, docs/operators/index.md, design/system/agent-adapters.md]
after: [T-2153, T-2174]
usage:
  source: log
  seconds: 203
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 39
      output: 15911
      cache_read: 2895828
      cache_write: 80304
      cost: 1.4207
---
# T-2178 flai-cli.md, docker.md, go.md, and the operators' index describe pilot, the flai/internal/pilot package, and the layered images

## Work

Add `pilot` and `flai/internal/pilot` to the package map in `design/system/flai-cli.md`, with a short section on the start request, the relay, and the resume; check `design/tech/docker.md`'s rows from T-2153 against what T-2174 built and correct them; name the second binary and its build in `design/tech/go.md`; and say in `docs/operators/index.md` what pilot is, that it is the entrypoint of flai's agent images, and that it reads nothing of the operator's but the start request and the variables it is given. Waits for T-2153, whose ADR it points to, and T-2174, whose images it describes.

## Done when

The four documents describe pilot as the tests of T-2168 check it, the markdown lint passes, and `flai check --strict` passes.

## Notes

`docs/users/flai.md` gains its container section in S-0508, where the user-facing command exists.
