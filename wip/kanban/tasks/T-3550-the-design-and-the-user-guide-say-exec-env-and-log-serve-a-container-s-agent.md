---
id: T-3550
type: task
nature: remediation
title: The design and the user guide say exec, env, and log serve a container's agent, and I-0361 is closed
status: done
parent: S-0834
owner: alex
created: 2026-11-17T08:33:05Z
updated: 2026-11-18T08:07:54Z
transitions:
  - to: ready
    at: 2026-11-18T08:00:28Z
    by: agent-S-0834
  - to: in-progress
    at: 2026-11-18T08:00:28Z
    by: agent-S-0834
  - to: done
    at: 2026-11-18T08:07:54Z
    by: agent-S-0834
stream: S-0834
tags: [docs, kluster]
touches: [design/system/kluster.md, design/system/agent-adapters.md, docs/users/flai.md, design/issues/I-0361-flai-kluster-exec-and-log-read-flai-cache-kluster-which-a-story-s-agent-in-a.md, design/issues/summary.md, design/system/flai-cli.md]
after: [T-3549]
usage:
  source: log
  seconds: 440
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 62
      output: 300
      cache_read: 4789226
      cache_write: 117998
      cost: 2.8801
---
# T-3550 The design and the user guide say exec, env, and log serve a container's agent, and I-0361 is closed

## Work

This task describes what T-3549 built, so it waits for it.

- `design/system/kluster.md`:
  - In § What the agent runs against its namespace, "Reads", say that `flai kluster exec`, `env`, and `log` read the story's copy in a container, and the main checkout's files on the host.
  - In § The namespace's own kubeconfig, add the copy, when it is written and removed, and the server's rewrite.
  - In the file table, add `.flai-cache/agents/<S-nnnn>/kluster/`.
- `design/system/agent-adapters.md`: in the container's Mounts row and the Isolation row, which says only one cache folder is mounted, say that the story's cache folder carries the namespace's copy. Name the ADR.
- `docs/users/flai.md`: in the exec, env, and log sections, say they work for an agent in a container too.
- Close I-0361 with `flai issue close I-0361 --reason`, saying what fixed it and naming the ADR, T-3547, T-3548, and T-3549. Tick both criteria once they are verified.

## Done when

- The three documents say what was built, and the markdown lint passes on them.
- I-0361 is closed, with `design/issues/summary.md` updated by flai.
- Both of S-0834's criteria are ticked.

## Notes
