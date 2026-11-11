---
id: D-0416
title: flai test's go-test tier can hit go test's 10-minute timeout on flai/internal/mcpserver under -race on a loaded host with no test hung, naming only the package (I-0227)
kind: cause
sources: [S-0535, I-0227, S-0417, S-0548, S-0540, S-0722]
topics: [flai, go, mcp, code]
paths: [flai/internal/mcpserver, system-flow.yaml, flai/internal/verify/select.go, flai/internal/verify/gotest.go, design/system/project-manifest.md, scripts/integration.sh]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: [S-0540, T-2327, T-2328]
---

# D-0416 flai test's go-test tier can hit go test's 10-minute timeout on flai/internal/mcpserver under -race on a loaded host with no test hung, naming only the package (I-0227)

**No longer current in part: S-0722 changed the tiers on 2026-11-11 (D-0585).** The go-test tier and `scripts/integration.sh` now run on one core without `-race`, at `GOMAXPROCS=1` with `-p 1` (D-0584), and take `-timeout 10m`; only the manual race tier, `scripts/integration.sh --race`, runs under `-race` and keeps `-timeout 20m`. What follows is how the 20m came about.

S-0535 changed `flai/internal/mcpserver`. `flai test` on its paths failed the go-test tier twice, at 10m21s and 10m1.6s, with a finding that named only the package. The tier in `system-flow.yaml` passed no `-timeout`, so go test's 10-minute default applied to the `-race` run, and the panic named no test.

The agent recorded a deviation in its narrative: it ran `go test` by hand three times, "one test with `-run`; the `mcpserver` package with `-timeout 560s`; it and main's copy without `-race`". "That showed no test hangs or slowed." The host's load average was about 13 on 24 cores. S-0417 met it again and bumped I-0227: "the `mcpserver` tests ran into Go's 10-minute timeout when three packages ran together; alone they passed in 9m20s" (wip/archive/agents/S-0417.md, Decisions).

S-0548 closed I-0227 on 2026-10-31. The fix, as the issue records it: "the go-test tier passes -timeout 20m and flai check warns, as manifest.tests-timeout, of a tier that runs go test with none (T-2363); the close-out's go-test runs -short, so the two prime-on-this-repository tests are the integration tier's own (T-2366); and the package runs in parallel, 9m14s to 136 to 162 s under -race (T-2361)." I-0225 closed with it. `manifest.tests-timeout` is a plain warning, so `flai check --strict` fails a project whose go-test tier leaves the timeout out.

S-0540 fixed the second cause, the finding that names only the package. A scratch module run with `go test -json -timeout 2s` on go1.26.0 showed why: "the panic and its `running tests:` list come as `output` events of the running test, and only the package gets a `fail` event." The `go-test-json` format (`flai/internal/verify/gotest.go`, T-2326) now gives one finding per test still running, with how long it had run, placed at its own function's frame in the panic's stacks; a subtest at its top-level test's first closure, and a parent shown running only because its subtest runs is left to the subtest. `design/system/project-manifest.md` records the format's behaviour. S-0540 measured `flai test flai/internal/mcpserver` passing in 2m31.3s at a load average of 45.6, against 9m20s to 10m21s for the package alone before.

What remains is I-0234: about 70 s of tests run serially because the server reads the session from the environment (D-0443), and then the two prime tests, about 550 user-seconds under `-race` in `ctxpack.ForStory`. S-0548 tried a faster written-path search in `Selection.refs` (T-2420). It measured no gain under `-race`, "506 to 587 user-seconds on the prime tests against 570 to 589 without it", so it was reverted and T-2420 cancelled; "the profile that pointed there was taken without `-race`". So a later story that speeds up those tests profiles and measures them under `-race`, as the tiers run them.

Sources: S-0535 and S-0417's narratives, § Decisions; I-0227 and its `## Remediation`; S-0548's narrative (`wip/archive/agents/S-0548.md`), § Decisions; S-0540's narrative (`wip/archive/agents/S-0540.md`), § Decisions; I-0234; S-0722's narrative (`wip/archive/agents/S-0722.md`), § Decisions.
