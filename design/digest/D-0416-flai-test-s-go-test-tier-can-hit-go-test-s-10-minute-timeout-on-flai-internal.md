---
id: D-0416
title: flai test's go-test tier can hit go test's 10-minute timeout on flai/internal/mcpserver under -race on a loaded host with no test hung, naming only the package (I-0227)
kind: cause
sources: [S-0535, I-0227]
topics: [flai, go, mcp, code]
paths: [flai/internal/mcpserver, system-flow.yaml]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0416 flai test's go-test tier can hit go test's 10-minute timeout on flai/internal/mcpserver under -race on a loaded host with no test hung, naming only the package (I-0227)

S-0535 changed `flai/internal/mcpserver`. `flai test` on its paths failed the go-test tier twice, at 10m21s and 10m1.6s, with a finding that named only the package. The tier in `system-flow.yaml` passes no `-timeout`, so go test's 10-minute default applies to the `-race` run, and the panic names no test.

The agent recorded a deviation in its narrative: it ran `go test` by hand three times, "one test with `-run`; the `mcpserver` package with `-timeout 560s`; it and main's copy without `-race`". "That showed no test hangs or slowed." The host's load average was about 13 on 24 cores. Without `-race` the package took 130s on the branch against 76s for main's copy; the gap is `TestPrimeOnThisRepository` and `TestPrimeComesInPartsOnThisRepository`, which prime the real repository.

So a later story that changes `mcpserver` and sees this timeout should suspect load, not a hang, and check that before chasing a test. [I-0227](../issues/I-0227-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md) is open and lists the remedies: a longer `-timeout` on the tier, naming the tests still running when the timeout fires, or keeping the two prime tests out of `-race`. Running the tests by hand is still a deviation from `tooling.md`, and has to be recorded as one.

Source: S-0535's narrative (`wip/archive/agents/S-0535.md`), § Decisions, and I-0227.
