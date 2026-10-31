---
id: S-0548
type: story
nature: remediation
title: "A queue verification takes under five minutes: the slow and network-bound tests are fixed, tests run in parallel, the gateway tests leave smoke, and every go test tier sets its timeout"
status: done
parent: E-0031
owner: alex
created: 2026-10-30T01:52:08Z
updated: 2026-10-31T21:02:24Z
transitions:
  - to: ready
    at: 2026-10-31T07:21:29Z
    by: alex
  - to: in-progress
    at: 2026-10-31T07:21:39Z
    by: agent-S-0548
  - to: review
    at: 2026-10-31T09:59:12Z
    by: agent-S-0548
  - to: done
    at: 2026-10-31T21:02:24Z
    by: alex
tags: [flai]
topics: [testing]
touches: [flai/internal/mcpserver, flai/internal/serve, flai/internal/selfupgrade, flai/cmd, flai/internal/check, flai/internal/mdlint, scripts/smoke.sh, scripts/gateway-smoke.sh, scripts/integration.sh, system-flow.yaml, ".github/workflows/flai.yml", design/conventions/testing.md, design/tech/go-libraries.md, docs/operators, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, design/system/project-manifest.md, design/system/flai-cli.md, design/system/agent-adapters.md, design/conventions/code-quality.md, design/conventions/tooling.md, docs/users/flai.md, docs/users/flai-reference.md, docs/contributors/index.md, template/CHANGELOG.md, design/issues/summary.md, design/issues/I-0225-flai-internal-mcpserver-s-tests-take-8-to-10-minutes-under-race-at-go-s-10-minute-default-timeout-so-the-integration-tier-fails-a-merge-queue-batch-on-time-alone.md, design/issues/I-0226-two-flai-serve-tests-check-before-what-they-wait-for-has-happened-and-fail-a-loaded-merge-queue-batch-on-timing.md, design/issues/I-0227-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md, design/issues/I-0233-testachangedfilereachesthedashboardasanotification-fails-its-temp-folder-s-cleanup-with-directory-not-empty-as-i-0048-s-serve-tests-did.md, design/adrs/0232-in-a-close-out-and-a-merge-queue-batch-the-smoke-tier-s-check-sh-and-lint-md-sh.md, design/adrs/README.md, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/verify/manifest.go, flai/internal/verify/verify.go, scripts/README.md, design/issues/I-0234-flai-internal-mcpserver-s-tests-still-take-136-to-162-s-under-race-since-the.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 9469
  priced_by: harness
  turns:
    - day: 2026-10-31
      ceremony: 6
      test_runs: 7
      hand_edits: 6
      work: 125
  models:
    - model: claude-opus-5-5
      input: 1552
      output: 508883
      cache_read: 133952517
      cache_write: 2705785
      cost: 52.3186
  strategic:
    - kind: planner
      seconds: 538
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 100
          output: 46177
          cache_read: 5457420
          cache_write: 180629
          cost: 3.4605
    - kind: orchestrator
      seconds: 75
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 36
          output: 235
          cache_read: 1083315
          cache_write: 210627
          cost: 0.2938
forecast:
  duration: 6h
  delivery: 2026-10-31T13:30:00Z
  basis: flai forecast's 39m from size 24 raised to 6h for nine tasks, two of which make 990 tests parallel, and repeated timed race runs of whole packages and a batch.
  by: planner-S-0548
  at: 2026-10-31T07:28:17Z
---
# S-0548 A queue verification takes under five minutes: the slow and network-bound tests are fixed, tests run in parallel, the gateway tests leave smoke, and every go test tier sets its timeout

## Goal

A merge-queue verification took 10.4 minutes at the median, 15.2 hours in the fortnight 2026-10-15 to 2026-10-29, and a close-out's go-test tier 349 s, while every other tier together takes under ten seconds. Timing the tiers on 2026-10-31 found where the minutes go, and most of them test nothing about the story:

- Go runs packages in parallel, so a tier lasts as long as its slowest package: `mcpserver` 500 s and `cmd` 450 s in batch 20261030-190545, `serve` next at 86 s, the other 70 packages under 80 s. A close-out that touches `cmd` runs go-test for 394 s at the median (35 of 52 runs) and 23 s without `cmd` or `mcpserver`.
- No test runs in parallel inside a package: 3,480 test functions and no `t.Parallel`, so a package's time is the sum of its tests, `cmd` 744 tests at 0.6 s and `mcpserver` 246 at 2 s, on a 24-core host.
- `TestDashboardRunsTheCachedListWhenTheReleasesCannotBeReached` points the release client at `127.0.0.1:1` for a refused connection; on WSL2 with mirrored networking a closed local port hangs to the 60 s connect timeout, and `selfupgrade` tries three times, so the test takes 184 s, 40 percent of `cmd`.
- `TestPrimeOnThisRepository` and `TestPrimeComesInPartsOnThisRepository` prime every open story of the monorepo under race, 131 s and 111 s, half of `mcpserver`; `TestMonorepoIsClean` (54 s) and `TestRepositoryLintsClean` (40 s) run `flai check` and the markdown lint over the whole repository under race, which the smoke tier runs again as scripts in 3 s and 21 s.
- The race detector costs 3.4 times (`hostapi` 7 s plain, 25 s under `-race`) and `-count=1` caches nothing; compiling and linking every test binary takes 8 s and is not the cost. 47 mcpserver calls exercise `wait_for_events` and `wait_for_work` with a 3 s MaxWait, and 92 `time.Sleep` calls sit in tests across the tree.
- The smoke tier ends with `scripts/gateway-smoke.sh`, which runs `TestGatewayOpenRouter` and `TestGatewayLiteLLM` through live `claude -p` sessions whenever both keys are set in the host's environment: 70 s and 165 to 194 s in the last four batches, about 0.32 USD of gateway spend each, and an outcome the provider decides. On 2026-10-31 the OpenRouter test failed after 129 s on a 402 answer, so the next batch fails on the operator's OpenRouter balance.
- Integration and smoke run side by side and slow each other: the same steps take 170 s alone and 333 s and 291 s in a batch.

Other failures were time alone: the mcpserver package against Go's 10-minute default (I-0225, I-0227), two serve tests that check before what they wait for has happened (I-0226), the dashboard notification test's temp folder cleanup (I-0233), and the macOS CI run that exceeded the default on 2026-10-25.

This story brings a verification under five minutes and makes it fail only on a defect in the story: the dashboard test closes a listener of its own or gives the client a short timeout instead of a closed port; the gateway tests leave the smoke tier for one the operator runs on demand or on a schedule, so no batch spends gateway money or fails on a provider; the fixture-based tests of `cmd` and `mcpserver`, whose fixtures use `t.TempDir`, run with `t.Parallel`; the four whole-repository tests are the integration tier's own, skipped under `-short`, and the close-out's go-test runs with `-short` again; the check and the markdown lint run once per batch, in the tests or in the smoke script, not both; the serve and dashboard tests wait for the condition they assert; and every go test tier in the manifest and in CI names its timeout, so a slow host fails with the package's name rather than a panic.

## Acceptance criteria
- [x] `TestDashboardRunsTheCachedListWhenTheReleasesCannotBeReached` takes under two seconds: it reaches a listener it opened and closed, or a client whose connect timeout is under a second, and no test in the tree dials a port nothing listens on
- [x] `scripts/smoke.sh` runs no gateway test; `scripts/gateway-smoke.sh` is a tier of its own in system-flow.yaml with `queue: false` and `all_only: true` that runs only when the operator asks (`flai test gateway`) or on a schedule the host sets, with the testing convention and the operator docs saying so and what it spends
- [x] The `cmd` and `mcpserver` tests whose fixtures are per test run with `t.Parallel`, and the shared state that stops any of them from doing so (globals set by `setDashboardReleases` and the like) is made per test or named in the story's notes as the reason it stays serial
- [x] `go test -race ./internal/cmd` finishes under two minutes and `go test -race ./internal/mcpserver` under three on the host under its usual load, from 332 s and 9m14s, and the five slowest tests of each are named in the close-out record with what keeps the rest serial
- [x] `TestPrimeOnThisRepository`, `TestPrimeComesInPartsOnThisRepository`, `TestMonorepoIsClean`, and `TestRepositoryLintsClean` skip under `-short`; the close-out's go-test tier runs with `-short` and the integration tier without, and the testing convention says which tests are the integration tier's own
- [x] `flai check --strict` and the markdown lint of the whole repository run once in a batch: either the integration tests or the smoke script runs them, and the other does not
- [x] The tests I-0226 and I-0233 name pass ten times in a row under `-race -count=10` with `GOMAXPROCS=2`, and the issues are closed with the reason
- [x] Every go test tier in system-flow.yaml, `scripts/integration.sh`, and both `go test` lines in the CI workflow carry an explicit `-timeout`, and `flai check` warns of a go test tier without one
- [x] A batch's verification, integration and smoke together, takes under five minutes on the host, as `flai verify --all` at the close-out commit records it, from 10.4 minutes at the median before; I-0225 and I-0227 are closed

## Tasks

The planner drafts them.
- T-2361 The mcpserver tests run in parallel where they share no process state, so the package finishes under two minutes under -race
- T-2362 The serve tests I-0226 and I-0233 name wait for what they assert, and pass ten times in a row under -race -count=10 with GOMAXPROCS=2
- T-2363 The go-test tier names its timeout, and flai check warns of a tier that runs go test without -timeout
- T-2364 Close I-0225, I-0226, I-0227, and I-0233 with their reasons, and record the timings in the story's notes
- T-2365 The dashboard releases test reaches a listener it opened and closed, and no cmd test dials a port nothing listens on
- T-2366 The four whole-repository tests skip under -short, the close-out's go-test tier runs with -short, and the check and the markdown lint run once in a batch
- T-2367 The cmd tests whose fixtures are per test run with t.Parallel, and the package globals that stop them are made per test or named
- T-2368 flai test runs a tier the operator names, all_only or not, so flai test gateway runs the gateway checks on demand
- T-2369 The gateway checks leave the smoke tier for a gateway tier of their own that runs only when the operator asks or on the host's schedule
- T-2420 A pack's search for the paths its texts name reads only their slash-bearing runs, so priming the repository takes less time

## Notes

Measured after landing, not a criterion (the merge queue takes only a story whose criteria can be ticked at review): the fortnight after landing shows verification p50 under five minutes in `flai stats`, and a batch verification under five minutes with nothing else running. Criteria 4 and 9 were reworded by the operator on 2026-10-31 to what the close-out record measures under the host's usual load, which is where every verification runs.

Evidence: `flai stats` merge-queue figures for the fortnight; the per-tier durations in `.flai-cache/verify/*.json` since 2026-10-15; the per-test elapsed times in `.flai-cache/merge-queue/batches/20261030-190545/integration.log` and the gateway test times in each batch's `smoke.log`; timings taken by hand on 2026-10-31 (closed-port `curl` 70 s to time out, compile and link 8 s, `hostapi` 7 s plain and 25 s under race, the dashboard test 187 s alone, `gateway-smoke.sh` 130 s and failing on a 402); the CI run 40921428484 fixed by hand on 2026-10-29, whose race fixes are in the main checkout uncommitted. The gateway tests stay valuable: they are the only check that the guard and `permission_prompt` hold through a gateway (S-0351, S-0356); they move out of every batch, they are not dropped.

### Planning

Planned by planner-S-0548 on 2026-10-31 while S-0548's agent started: the agent wrote T-2361 to T-2364 against the story's first four criteria, and the planner wrote T-2365 to T-2369 for criteria 1, 2, 3 and 4 for `cmd`, 5, and 6, and set T-2364's `after` and touches so that it measures last.

Touches:

- Declared, all kept: the operator's 15, among them the folders `flai/internal/mcpserver`, `flai/internal/serve`, `flai/internal/dashboard`, `flai/internal/selfupgrade`, `flai/cmd`, `flai/internal/check`, `flai/internal/mdlint`, and `docs/operators`. The folders stay because T-2361 and T-2367 add `t.Parallel` to test files across `mcpserver` and `cmd` that no task can list in full yet, and T-2367 makes package globals per test in files it finds as it goes; flai replaces each folder in the claim with the task files it holds (ADR-0096). `flai/internal/dashboard` holds no task's file, since the dashboard test is in `flai/cmd`; the agent may narrow it with `flai touches`. `design/conventions/testing.md` does not exist: the testing tiers are in `design/conventions/code-quality.md` and `tooling.md`, which T-2366 and T-2369 change; TH-0771 named this assumption and the operator resolved it.
- Design and layout, from the tasks: `flai/internal/verify/select.go` and its test (T-2368), `flai/internal/manifest/manifest.go` and its test and `design/system/project-manifest.md` (T-2363), `design/system/flai-cli.md` and `docs/users/flai.md` (T-2363, T-2368), `docs/users/flai-reference.md` (T-2368), `template/CHANGELOG.md` (T-2363), `design/conventions/code-quality.md` and `tooling.md` (T-2366, T-2369), `docs/contributors/index.md` and `design/system/agent-adapters.md`, which describe `gateway-smoke.sh` in smoke (T-2369), and `design/issues/summary.md` with the four issue files (T-2364).
- Co-change: `flai touches suggest` ranks `docs/users/flai.md` (25%), `design/system/flai-cli.md` (24%), and `docs/users/flai-reference.md` (22%) first, which the tasks name; the rest (`hostapi`, `flaiover`, `guard`, `harness`) no criterion reaches, so they are left out.

Forecast: `flai forecast` gave 39m, 97 s per unit of size over 46 remediation stories at size 24. Raised to 6h: the size counts criteria and touches, not the 990 tests T-2361 and T-2367 make parallel or the globals they make per test, and four tasks need timed race runs of whole packages, three to ten runs of several minutes each, and a batch. Delivery counts from the start at 07:22.

Cost of delay: none written. Neither S-0548 nor E-0031 has inputs, so `flai cod` refuses. TH-0771 asked for them with a recommendation (penalty 14 USD a week, 12m lost per cycle), and the operator resolved it without setting any, so no value is written until inputs are given.

After: none. S-0548 needs no other story's accepted change. It overlaps S-0552, in progress, on `flai/cmd`, `flai/internal/serve`, `system-flow.yaml`, and `docs/operators`, which the two agents settle in their conversation.

### Close-out record

Timings by hand with `go test -race -count=1 -json`, each package alone, since `flai test` reports a tier's time and not each test's (a deviation from `tooling.md`, recorded in the narrative). The host ran four other stories' agents throughout; load is the one-minute average at start and end.

| Package | Before | After | Load | Five slowest tests after |
|---------|--------|-------|------|--------------------------|
| `mcpserver` | 9m14s (load 10) | 136 s, 141 s, 162 s | 40 to 57, 14 to 31, 28 to 55 | `TestAFirstLookIsBoundedAndLeavesTasksOut` 11.2 s, `TestInboxSaysTheFlaiIsOlderThanTheProject` 10.9 s, `TestTheMessageToolsAndInstructions` 10.5 s, `TestStoryStartLeavesOutOldChangesToKeepPartOne` 8.7 s, `TestALaterLookReportsTasksAndIsStillCapped` 8.2 s; the two prime-on-this-repository tests, parallel, end the run 55 to 70 s after the serial phase |
| `cmd` | 332 s (load 45 to 71) | 75 s, 80 s, 84 s, 87 s | 26 to 36, 31 to 33 | `TestHostRunsServeAndTheMCPServersAndTakesThemWithIt` 14.1 s, `TestSettingsIndexIsComplete` 10.2 s, `TestVerifyLeavesTheQueueTiersToTheMergeQueueWhileItCanBatch` 6.8 s, `TestALandingIsRefusedUnlessTheBatchVerifiedWhatLands` 6.0 s, `TestEveryAcceptanceWaitsForThePullRequestsChecksAndReviews` 4.8 s |

`TestDashboardRunsTheCachedListWhenTheReleasesCannotBeReached`: 184 s to 0.10 s. The three tests I-0226 and I-0233 name: 30 of 30 under `GOMAXPROCS=2 go test -race -count=10`, at load about 40.

Tests left serial, and why (criterion 3):

- `mcpserver`, 41 top-level tests: each sets an environment variable mid-test or one no other test shares (`FLAI_STORY`, `FLAI_ROLE`, `FLAI_REASONS`, `FLAI_STARTED_BY`, `GH_TOKEN`, `GITHUB_TOKEN`, `FLAI_RELEASES_API`, `FLAI_CONFIG`), or a package variable (`flaiExecutable`, `buildinfo.Version`). 54 more run as parallel subtests of seven serial parents that set one variable each (`together`): `TestInAGitSession`, `TestPermissionPromptPollingFast` (`permissionPoll`), `TestAsTheStorysAgent`, `TestAsTheOrchestrator`, `TestWithAGitIdentity`, `TestInStreamSession`, and the release-publish tests. The server reads the session's story, role, session, reasons, and started-by from the environment when it starts (`server.go`), which is I-0234.
- `cmd`, 110 tests: package variables `componentsRunner`, `componentsInstaller`, `githubAPIBase`, `releaseKeys`, `executablePath`, `exitWithProcessEvery`, `hostapi.Commands`, `release.RemoteTTL`, `buildinfo.Version`; and environment variables set mid-test or in combinations no group shares (`FLAI_ROLE`, `FLAI_AGENT`, `FLAI_STORY`, `FLAI_SESSION`, `FLAI_STARTED_BY`, `LOG_FORMAT`, `GITHUB_TOKEN`, `FLAI_CACHE_DIR`, `PATH`, git identity). `dashboardReleases`, `apps`, `acceptNewForge`, and `acceptMergeWait` became the app's own, so their tests run in parallel; eight groups run 137 tests as parallel subtests.

Criteria left unticked:

- 4, for `mcpserver`: 136 to 162 s on this host at load 14 to 57, not under two minutes. About 70 s is the serial phase above, and the prime tests cannot start before it. I-0234 records it with the remedy: the server takes the session as options, so those tests run in parallel. A faster search in `Selection.refs` measured no gain under race and was reverted (T-2420, cancelled). `cmd` meets it.
- 9: a batch's verification with nothing else running cannot be timed here, where four other agents keep the load between 14 and 60; and the fortnight of `flai stats` lies after landing. After the close-out passed at 797b165, `flai verify S-0548 --all` ran the two queue tiers one after the other at load 64 to 74: integration 3m33s and smoke 1m08s, 4m40s together, against 10.4 minutes at the median before. The fortnight's p50 in `flai stats` is for the operator to read after landing. I-0225 and I-0227 are closed.
