---
id: D-0392
title: flai ties a child process to its host with --exit-with the host's PID, never Linux's Pdeathsig, because Go forks from any OS thread and that thread's exit would kill the child
kind: decision
sources: [S-0106]
topics: [go, serve, code]
paths: [flai/internal/host]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0392 flai ties a child process to its host with --exit-with the host's PID, never Linux's Pdeathsig, because Go forks from any OS thread and that thread's exit would kill the child

S-0106 chose `--exit-with <host pid>`, the mechanism `flai mcp http` already had, over `Pdeathsig`: "Go forks from any OS thread, and that thread's exit would kill the child." The host also stops its children when it returns, so a killed host leaves nothing running within seconds.

Any new child process flai starts follows the same rule.

Source: S-0106.
