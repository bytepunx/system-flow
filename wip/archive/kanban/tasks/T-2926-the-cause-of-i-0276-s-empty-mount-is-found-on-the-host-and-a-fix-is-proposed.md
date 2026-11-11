---
id: T-2926
type: task
nature: research
title: The cause of I-0276's empty mount is found on the host and a fix is proposed
status: done
parent: S-0667
owner: alex
created: 2026-11-10T08:48:42Z
updated: 2026-11-10T11:02:56Z
transitions:
  - to: ready
    at: 2026-11-10T10:59:27Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T10:59:28Z
    by: system-flow
  - to: done
    at: 2026-11-10T11:02:56Z
    by: agent-S-0667
stream: S-0667
tags: [docker, agent-container]
touches: [design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md]
usage:
  source: log
  seconds: 208
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 43
      output: 8207
      cache_read: 4219862
      cache_write: 124773
      cost: 2.0065
---
# T-2926 The cause of I-0276's empty mount is found on the host and a fix is proposed

## Work

Waits for no task: it finds what the other two build on.

Reproduce I-0276's instance on the host, from this story's agent's shell, inside and outside Claude Code's sandbox, and from the operator's shell where the operator can run it:

```sh
docker run --rm -v "$ROOT:$ROOT:ro" flai-agent-sf-claude-code:local ls -la "$ROOT"
```

`$ROOT` is the main checkout, `/home/alex/git/bytepunx/system-flow`. For each shell, record what differs between them and `flai serve`'s container start, which mounts the same path with its files (ADR-0247): `docker context show`, `DOCKER_HOST` and `DOCKER_CONTEXT`, `command -v docker`, `docker info --format '{{.Name}} {{.OperatingSystem}}'`, and the mount's `Source` in `docker inspect` of a container left running. Read how `flai serve` runs `docker` in `flai/internal/serve/container.go`: the program, its environment, and the mount arguments.

Name the cause in the story's Notes under `### Cause`, with the commands and what each printed, and the fix it calls for. A fix in `scripts/agent-container.sh` or the runbook is T2's and T3's. When the fix is outside this repository, such as a Docker Desktop setting or the host's docker configuration, ask the operator on a thread on the story, the recommended fix first, before T2 starts.

Add the finding to I-0276 as an instance with `flai issue bump I-0276`.

## Done when

- The story's Notes name the cause of the empty mount, with the commands that show it and their output.
- The fix is named there, and a fix outside the repository is asked on a thread on the story.

## Notes
