---
id: T-2367
type: task
nature: remediation
title: The cmd tests whose fixtures are per test run with t.Parallel, and the package globals that stop them are made per test or named
status: done
parent: S-0548
owner: alex
created: 2026-10-31T07:25:25Z
updated: 2026-10-31T09:02:52Z
transitions:
  - to: ready
    at: 2026-10-31T08:16:01Z
    by: agent-S-0548
  - to: in-progress
    at: 2026-10-31T08:16:02Z
    by: agent-S-0548
  - to: done
    at: 2026-10-31T09:02:52Z
    by: agent-S-0548
stream: S-0548
tags: [cli]
touches: [flai/cmd/dashboard_test.go, flai/cmd/dashboard.go, flai/cmd/root.go, flai/cmd/forge.go, flai/cmd/components.go, flai/cmd/serve_mcp.go, flai/cmd/selfupgrade.go, flai/cmd/issue.go, flai/cmd/accept_remote.go, flai/cmd/accept_batch_dirty_test.go, flai/cmd/accept_conflict_test.go, flai/cmd/accept_epic_test.go, flai/cmd/accept_fix_test.go, flai/cmd/accept_forge_test.go, flai/cmd/accept_install_test.go, flai/cmd/accept_lock_test.go, flai/cmd/accept_orchestrator_dirty_test.go, flai/cmd/accept_orchestrator_test.go, flai/cmd/accept_overlap_test.go, flai/cmd/accept_own_thread_test.go, flai/cmd/accept_protected_test.go, flai/cmd/accept_queue_test.go, flai/cmd/accept_queue_tiers_test.go, flai/cmd/accept_release_due_test.go, flai/cmd/accept_remote_test.go, flai/cmd/accept_research_test.go, flai/cmd/accept_resume_test.go, flai/cmd/accept_test.go, flai/cmd/accept_threads_test.go, flai/cmd/accept_uncommitted_test.go, flai/cmd/accept_undone_test.go, flai/cmd/accept_worktree_test.go, flai/cmd/activity_test.go, flai/cmd/adr_test.go, flai/cmd/agent_test.go, flai/cmd/analyze_test.go, flai/cmd/archive_test.go, flai/cmd/board_limit_test.go, flai/cmd/board_test.go, flai/cmd/branch_relative_test.go, flai/cmd/branch_test.go, flai/cmd/check_stats_test.go, flai/cmd/check_test.go, flai/cmd/checks_test.go, flai/cmd/cmd_test.go, flai/cmd/cod_test.go, flai/cmd/components_test.go, flai/cmd/criteria_test.go, flai/cmd/curate_test.go, flai/cmd/dashboard_agent_test.go, flai/cmd/dashboard_folder_test.go, flai/cmd/dashboard_versions_test.go, flai/cmd/dashboard_watch_test.go, flai/cmd/design_test.go, flai/cmd/digest_test.go, flai/cmd/doc_lint_test.go, flai/cmd/doc_test.go, flai/cmd/edit_test.go, flai/cmd/file_test.go, flai/cmd/forecast_test.go, flai/cmd/forge_test.go, flai/cmd/guard_test.go, flai/cmd/hold_test.go, flai/cmd/host_integration_test.go, flai/cmd/host_test.go, flai/cmd/host_versions_test.go, flai/cmd/import_commit_test.go, flai/cmd/import_serve_test.go, flai/cmd/import_test.go, flai/cmd/installsh_test.go, flai/cmd/issue_test.go, flai/cmd/items_body_test.go, flai/cmd/license_test.go, flai/cmd/logging_test.go, flai/cmd/manifest_test.go, flai/cmd/mcp_http_test.go, flai/cmd/mcp_test.go, flai/cmd/message_test.go, flai/cmd/migrate_test.go, flai/cmd/move_by_test.go, flai/cmd/move_cancel_test.go, flai/cmd/move_orchestrate_test.go, flai/cmd/move_refused_test.go, flai/cmd/new_test.go, flai/cmd/orchestrate_permissions_test.go, flai/cmd/order_test.go, flai/cmd/plan_candidates_test.go, flai/cmd/plan_orchestrate_test.go, flai/cmd/plan_test.go, flai/cmd/planning_test.go, flai/cmd/prime_test.go, flai/cmd/promote_drafts_test.go, flai/cmd/promote_test.go, flai/cmd/push_test.go, flai/cmd/queue_held_test.go, flai/cmd/queue_test.go, flai/cmd/reference_test.go, flai/cmd/release_evaluate_test.go, flai/cmd/release_moved_test.go, flai/cmd/release_pending_test.go, flai/cmd/repo_url_test.go, flai/cmd/review_uncommitted_test.go, flai/cmd/role_test.go, flai/cmd/selfupgrade_test.go, flai/cmd/serve.go, flai/cmd/serve_actions_test.go, flai/cmd/serve_agent_stop_test.go, flai/cmd/serve_import_test.go, flai/cmd/serve_mcp_test.go, flai/cmd/serve_project_test.go, flai/cmd/serve_status_below_test.go, flai/cmd/serve_test.go, flai/cmd/settings_doc_test.go, flai/cmd/shared_test.go, flai/cmd/show_test.go, flai/cmd/stats_test.go, flai/cmd/story_start_test.go, flai/cmd/storycommit_test.go, flai/cmd/stream_diff_test.go, flai/cmd/stream_open_test.go, flai/cmd/stream_reopen_test.go, flai/cmd/stream_restore_test.go, flai/cmd/stream_state_test.go, flai/cmd/stream_sync_test.go, flai/cmd/subproject_test.go, flai/cmd/task_done_test.go, flai/cmd/taskplan_test.go, flai/cmd/template_push_test.go, flai/cmd/thread_test.go, flai/cmd/topics_test.go, flai/cmd/touches_overlap_test.go, flai/cmd/touches_test.go, flai/cmd/upgrade_test.go, flai/cmd/verify_test.go, flai/cmd/version_test.go, flai/cmd/workitems_test.go]
after: [T-2365]
usage:
  source: log
  seconds: 2810
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 450
      output: 147489
      cache_read: 38823271
      cache_write: 784214
      cost: 15.1634
---
# T-2367 The cmd tests whose fixtures are per test run with t.Parallel, and the package globals that stop them are made per test or named

## Work

`flai/cmd` has 744 tests at about 0.6 s each and none calls `t.Parallel`, so the package's time is their sum, 450 s under `-race` in batch 20261030-190545, and a close-out that touches `cmd` runs go-test for 394 s at the median. Add `t.Parallel()` to every test, and subtest, whose fixture is its own (`t.TempDir`, its own `fakeRunner` and `app` through `runWith` and `runWithApp` in `dashboard_test.go`).

The package globals tests override stop that: `dashboardReleases` (`dashboard.go`, set by `setDashboardReleases`), `apps` (`root.go`), `githubAPIBase` (`forge.go`), `componentsRunner` and `componentsInstaller` (`components.go`), `exitWithProcessEvery` (`serve_mcp.go`), `releaseKeys` and `executablePath` (`selfupgrade.go`), and the vars in `issue.go` and `accept_remote.go`. Make each one a field of the `app` or its options so a test sets its own, or leave the tests that set it serial and name the global and the tests in the story's notes as the reason. `t.Setenv` and `t.Chdir` tests stay serial; Go refuses `t.Parallel` with them.

Leave `test.go` and `test_test.go` to T-2368, which changes them in the same layer; add `t.Parallel` there only once T-2368 is done. The task widens its touches to each test file it changes when it closes.

Second layer: it waits for T-2365, which changes `dashboard_test.go` and takes the 184 s dashboard test out of the package first, so this task's timing measures parallelism alone.

## Done when

`go test -race -count=1 ./cmd` passes three times in a row and finishes under two minutes on the host, its time and five slowest tests (`go test -json` elapsed) recorded in the narrative. Every test left serial is serial for a reason the story's notes name. Story criteria 3 and 4 for `cmd`.

## Notes

Drafted by the planner of S-0548 on 2026-10-31. It first drafted this task for `mcpserver`, which the story's agent's T-2361 covers, and rewrote it for `cmd`.
