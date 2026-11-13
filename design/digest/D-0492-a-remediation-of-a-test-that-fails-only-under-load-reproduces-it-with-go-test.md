---
id: D-0492
title: "A remediation of a test that fails only under load reproduces it with go test directly, GOMAXPROCS=2 -race -count=N -run, beside a CPU load taskset pins to the test's CPUs, since flai test takes none of those flags"
kind: decision
sources: [S-0569, I-0247, S-0560, I-0246, S-0543, S-0617, I-0252, S-0653, I-0314, S-0651, I-0309, S-0777, S-0779, I-0372]
topics: [testing, go, flai]
paths: [design/conventions/delegation.md, system-flow.yaml, flai/internal/mcpserver, flai/internal/serve, flai/internal/release, flai/internal/channel, flai/cmd]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [S-0617, S-0653, S-0777, S-0779]
---

# D-0492 A remediation of a test that fails only under load reproduces it with go test directly, GOMAXPROCS=2 -race -count=N -run, beside a CPU load taskset pins to the test's CPUs, since flai test takes none of those flags

S-0569 remediated I-0247's load-dependent integration failures. Its narrative records an exception to running tests only through `flai test`: "The tasks' stress runs (`GOMAXPROCS=2 go test -race -count=N -run <test>`) are run with `go test` directly ... they are the reproduction and measurement this remediation needs, and `flai test` takes no `-count`, `-run`, or `GOMAXPROCS`. The checks of each change still go through `flai test`."

It also kept the load it added from disturbing the other agents on the host: "The CPU load beside a stress run is confined with `taskset` to the CPUs the test is pinned to, so it contends with the test and not with the other agents on the host (24 CPUs, load average about 40 at the start). Each task of layer 1 takes its own CPU pair." Each fix reproduced its failure this way before anything changed, and all six tests then "passed 60 of 60 together, pinned to CPUs 20-23 beside 6 busy loops" (I-0247, Remediation).

S-0560, which remediated I-0246, made the same exception in its Decisions: "The stress runs T-2556, T-2557, and T-2558 name (`GOMAXPROCS=2 go test -race -count=N`) are run by hand, against the convention that tests run through `flai test`: `flai test` has no way to repeat a test or bound GOMAXPROCS, and the tasks' done-when needs those runs as measurements. Pass or fail of the change is still judged with `flai test`." Its log adds a measurement for a stress run of a whole package: `GOMAXPROCS=2 go test -race -count=5 -timeout 60m ./internal/mcpserver` "passed in 1926 s (about 385 s a pass on 2 CPUs, against S-0548's 136 s on 24 CPUs at count 1)". So a whole-package run under `GOMAXPROCS=2` needs a `-timeout` well past go test's 10-minute default; a `-run` of the failing tests alone takes seconds.

S-0543, which remediated I-0233, made the exception a third time: "The stress runs call `go test` by hand. `flai test` takes no `-count` or `GOMAXPROCS`, and a repeated race run under load is the task's deliverable. The routine checks stay with `flai test` and the close-out." (wip/archive/agents/S-0543.md, Decisions)

S-0617, which remediated I-0252, made it a fourth time, citing this entry: "T-2724's stress runs call `go test` by hand with `GOMAXPROCS=2 -race -count=N -run`, beside a CPU load pinned with `taskset`, and `-skip '^TestGateway'` ... Pass or fail of each change is still judged with `flai test`." (wip/archive/agents/S-0617.md, Decisions)

S-0653, which remediated I-0314, made it a fifth time, citing this entry, to show its new tests fail without the fix: "The reproducing tests were run with `go test -race -run` directly, as D-0492 records, to show them failing with `ofLog` swapped back to `newest`; pass or fail of the change is judged with `flai test`." (wip/archive/agents/S-0653.md, Decisions; D-0553). S-0649 did the same for I-0324, putting the old read back and running `go test -run` on the one test by hand (D-0543).

A later story that fixes a flaky test can do the same, and record the exception in its Decisions. A hand `go test` of `flai/internal/serve` must pass `-skip '^TestGateway'`, or it runs the paid gateway checks (D-0451). Reproduce before lengthening a wait (D-0413).

S-0651, which remediated I-0309 (D-0557), confined the load further, to a single CPU: "each test binary, built with -race, ran pinned to one CPU with GOMAXPROCS=1 beside 40 busy loops on that CPU, so other agents on the host were not loaded. main's versions came from git archive main into a scratch copy." Main's channel test then failed 5 of 10 runs and its MCP test 1 of 5. A sub-agent cannot run `git archive`, which flai guard refuses it, and copies the worktree instead (D-0527).

S-0777, which fixed eight kluster tests that failed on macOS's symlinked temp dir (D-0616), made the exception for a flag `flai test` cannot pass either, an environment variable: "The suite runs by hand with go test, against the rule that tests run through flai test, because flai test cannot set TMPDIR; it is the reproduction the story needs (D-0492). They ran with FLAI_STORY and FLAI_STARTED_BY unset, so that this session's variables did not reach the tests; flai guard refuses unsetting FLAI_AGENT, which was left as set. Pass or fail of each change was judged with flai test." (wip/archive/agents/S-0777.md, Decisions). So a hand run under `flai serve` unsets `FLAI_STORY` and `FLAI_STARTED_BY` for the tests, and leaves `FLAI_AGENT`, which the guard keeps; its TMPDIR lives under `/tmp` (D-0617).

S-0779, which remediated I-0372 (D-0618), made it a seventh time, citing this entry: "T-3333's stress run calls `go test` by hand, as D-0492 records for load-dependent remediations: `flai test` takes no `-count` or `-run`. Pass or fail of each change was judged with `flai test`." Its fifty-count `-run 'MCP|TheHost'` race run passed beside the whole `flai/cmd` package looping under race as the load: 2050 passes in 779 s (I-0372). Its first two load runs each caught a test no planned task covered, so a stress run is worth making before the fix is called done, not only after (wip/archive/agents/S-0779.md, Decisions).

Sources: S-0569 (wip/archive/agents/S-0569.md, Decisions); I-0247; S-0560 (wip/archive/agents/S-0560.md, Decisions and Log); I-0246; S-0543 (wip/archive/agents/S-0543.md, Decisions); S-0617 (wip/archive/agents/S-0617.md, Decisions); I-0252; S-0653 (wip/archive/agents/S-0653.md, Decisions); I-0314; S-0651 (wip/archive/agents/S-0651.md, Decisions); I-0309; S-0777 (wip/archive/agents/S-0777.md, Decisions); S-0779 (wip/archive/agents/S-0779.md, Decisions); I-0372.
