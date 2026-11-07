---
id: D-0492
title: "A remediation of a test that fails only under load reproduces it with go test directly, GOMAXPROCS=2 -race -count=N -run, beside a CPU load taskset pins to the test's CPUs, since flai test takes none of those flags"
kind: decision
sources: [S-0569, I-0247]
topics: [testing, go, flai]
paths: [design/conventions/delegation.md, system-flow.yaml]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0492 A remediation of a test that fails only under load reproduces it with go test directly, GOMAXPROCS=2 -race -count=N -run, beside a CPU load taskset pins to the test's CPUs, since flai test takes none of those flags

S-0569 remediated I-0247's load-dependent integration failures. Its narrative records an exception to running tests only through `flai test`: "The tasks' stress runs (`GOMAXPROCS=2 go test -race -count=N -run <test>`) are run with `go test` directly ... they are the reproduction and measurement this remediation needs, and `flai test` takes no `-count`, `-run`, or `GOMAXPROCS`. The checks of each change still go through `flai test`."

It also kept the load it added from disturbing the other agents on the host: "The CPU load beside a stress run is confined with `taskset` to the CPUs the test is pinned to, so it contends with the test and not with the other agents on the host (24 CPUs, load average about 40 at the start). Each task of layer 1 takes its own CPU pair." Each fix reproduced its failure this way before anything changed, and all six tests then "passed 60 of 60 together, pinned to CPUs 20-23 beside 6 busy loops" (I-0247, Remediation).

A later story that fixes a flaky test can do the same, and record the exception in its Decisions. A hand `go test` of `flai/internal/serve` must pass `-skip '^TestGateway'`, or it runs the paid gateway checks (D-0451). Reproduce before lengthening a wait (D-0413).

Source: S-0569 (wip/archive/agents/S-0569.md, Decisions); I-0247.
