---
id: T-2361
type: task
nature: remediation
title: The mcpserver tests run in parallel where they share no process state, so the package finishes under two minutes under -race
status: done
parent: S-0548
owner: alex
created: 2026-10-31T07:24:38Z
updated: 2026-10-31T08:14:28Z
transitions:
  - to: ready
    at: 2026-10-31T07:32:06Z
    by: agent-S-0548
  - to: in-progress
    at: 2026-10-31T07:32:07Z
    by: agent-S-0548
  - to: done
    at: 2026-10-31T08:14:28Z
    by: agent-S-0548
stream: S-0548
tags: []
touches: [flai/internal/mcpserver/activity_test.go, flai/internal/mcpserver/adr_test.go, flai/internal/mcpserver/agents_test.go, flai/internal/mcpserver/analyze_test.go, flai/internal/mcpserver/cascade_test.go, flai/internal/mcpserver/criteria_test.go, flai/internal/mcpserver/curate_test.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/follow_test.go, flai/internal/mcpserver/inbox_test.go, flai/internal/mcpserver/installed_test.go, flai/internal/mcpserver/issues_test.go, flai/internal/mcpserver/items_write_test.go, flai/internal/mcpserver/messages_test.go, flai/internal/mcpserver/orchestrate_test.go, flai/internal/mcpserver/orchestrator_end_test.go, flai/internal/mcpserver/orchestrator_threads_test.go, flai/internal/mcpserver/outdated_test.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcpserver/plan_test.go, flai/internal/mcpserver/prime_test.go, flai/internal/mcpserver/queue_test.go, flai/internal/mcpserver/review_test.go, flai/internal/mcpserver/role_test.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/shared_test.go, flai/internal/mcpserver/story_start_test.go, flai/internal/mcpserver/stream_test.go, flai/internal/mcpserver/subproject_test.go, flai/internal/mcpserver/task_test.go, flai/internal/mcpserver/test_test.go, flai/internal/mcpserver/timing_test.go, flai/internal/mcpserver/unpublished_test.go, flai/internal/mcpserver/verify_test.go, flai/internal/mcpserver/versions_test.go, flai/internal/mcpserver/work_test.go]
usage:
  source: log
  seconds: 2496
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 243
      output: 1171
      cache_read: 20263774
      cache_write: 545050
      cost: 7.2404
---
# T-2361 The mcpserver tests run in parallel where they share no process state, so the package finishes under two minutes under -race

## Work

`go test -race ./internal/mcpserver` takes 8 to 10 minutes (I-0225, I-0227) and no test calls `t.Parallel`. Each test builds its own fixture in `t.TempDir()` with its own server and clock, so most can run together. Call `t.Parallel()` in every test, and every subtest, that sets no environment variable (`t.Setenv` panics under `t.Parallel`), changes no package variable (`permissionPoll`, `flaiExecutable`), and changes no shared file such as the real repository. Make such a package variable an option of the server, or leave its tests serial and name the variable in S-0548's notes as the reason. Story criterion 4 times the package without `-short`, so `TestPrimeOnThisRepository` (131 s under race) and `TestPrimeComesInPartsOnThisRepository` (111 s) count: find where their time goes and cut it, such as by priming the repository once for both, without weakening what they assert. T-2366 makes them skip under `-short` after this task. Where other slow tests spend their time on work they do not check, share it with `sync.Once` or cut it. Waits for no other task.

## Done when

- [ ] `go test -race -count=1 ./internal/mcpserver` passes and finishes under two minutes on the host, timed and recorded in the narrative with the five slowest tests
- [ ] No test asserts less than before
- [ ] Story criteria 3 and 4 for `mcpserver`

## Notes
