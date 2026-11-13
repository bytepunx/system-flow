---
id: T-3170
type: task
nature: feature
title: flai builds the harness images on the host from its embedded Dockerfiles and its release's Linux binaries checked against the signed checksums, or from a development tree's own
status: done
parent: S-0740
owner: alex
created: 2026-11-12T15:59:35Z
updated: 2026-11-13T12:25:58Z
transitions:
  - to: ready
    at: 2026-11-13T11:55:44Z
    by: agent-S-0740
  - to: in-progress
    at: 2026-11-13T11:55:44Z
    by: agent-S-0740
  - to: done
    at: 2026-11-13T12:25:58Z
    by: agent-S-0740
stream: S-0740
tags: [cli]
touches: [flai/internal/agentimage/build.go, flai/internal/agentimage/build_test.go]
after: [T-3168, T-3287]
usage:
  source: log
  seconds: 1800
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 43
      output: 245
      cache_read: 2937846
      cache_write: 235746
      cost: 1.1846
---
# T-3170 flai builds the harness images on the host from its embedded Dockerfiles and its release's Linux binaries checked against the signed checksums, or from a development tree's own

## Work

`flai/internal/agentimage/build.go` builds bedrock, then claude-code over it:

- with `docker build` and a context folder under `.flai-cache` that holds T-3168's embedded Dockerfile;
- for an installed flai, with the running release's Linux `flai` and `pilot`, fetched with `selfupgrade.FetchSigned` and checked against `checksums.txt`;
- in a development tree, as S-0733's ADR names one, from the tree's own Dockerfiles and source, as `scripts/agent-image.sh` builds today;
- with the label's build arguments set.

Run docker through a runner that:

- stops a step with no output for the ADR's limit, and reports it;
- tries again, up to the ADR's count, a step whose output shows a network failure: a connection reset, a timeout, a TLS handshake failure, or an HTTP 5xx;
- writes every line to standard output and to a log file under `.flai-cache`, whose path it returns.

The rebuild of 2026-11-11 (E-0036 § What the rebuild showed) hung four hours in `apk add` and `npm install`, and failed on resets from `proxy.golang.org`: the tests replay both over a fake docker.

A story's container has no Docker socket (ADR-0203), so a real build is the operator's to run on the host. Name the command in the task's Notes for them.

It waits for T-3168, whose embedded Dockerfiles it builds.

## Done when

- `go test ./internal/agentimage/ -run Build` passes, with a silent step stopped, a reset retried and then passed, a reset retried until given up, and the log file written.
- The docker arguments a test records for each image match what `scripts/agent-image.sh` passes today, but for the binaries' source and the label.

## Notes
