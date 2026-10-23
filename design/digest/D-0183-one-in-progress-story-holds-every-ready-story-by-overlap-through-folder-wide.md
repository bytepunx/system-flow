---
id: D-0183
title: One in-progress story holds every ready story by overlap, through folder-wide touches and docs/users/flai.md, so the board runs one story at a time under a limit of three
kind: cause
sources: [I-0087]
topics: [cli, mcp, hostapi, dashboard, conventions, planning, template, flai, go, flaiover, sveltekit, code, all]
paths: [design/adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md, design/adrs/README.md, design/conventions/strategic-agents.md, design/conventions/work-management.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/I-0080-a-story-worktree-has-no-flaiover-node-modules-so-the-close-out-stops-at-the-flaiover-step-until-flaiover-install-runs-there.md, design/issues/I-0085-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md, design/issues/I-0087-one-in-progress-story-holds-every-ready-story-by-overlap-through-folder-wide-touches-and-docs-users-flai-md-so-the-board-runs-one-story-at-a-time-under-a-limit-of-three.md, design/issues/summary.md, design/system/agent-coordination.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/system/project-manifest.md, design/system/workflow.md, design/tech/go-libraries.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai-reference.md, docs/users/flai.md, docs/users/flaiover.md, flai/cmd/board.go, flai/cmd/guard.go, flai/cmd/hold_test.go, flai/cmd/root.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/shared.go, flai/cmd/shared_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/writes_test.go, flai/internal/itemedit/claim.go, flai/internal/itemedit/claim_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/shared.go, flai/internal/manifest/shared_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/shared.go, flai/internal/mcpserver/shared_test.go, flai/internal/mcpserver/work_test.go, flai/internal/metrics/claims_test.go, flai/internal/serve/agents.go, flai/internal/serve/hold_test.go, flai/internal/serve/start_test.go, flai/internal/workitem/boardview.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_i0087_test.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/promote_test.go, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/settings.ts, flaiover/src/routes/api/settings/+server.ts, flaiover/src/routes/api/settings/settings.test.ts, system-flow.yaml, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, template/root/system-flow.yaml.tmpl, template/template.yaml, flai/internal/guard]
date: 2026-10-06
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0183 One in-progress story holds every ready story by overlap, through folder-wide touches and docs/users/flai.md, so the board runs one story at a time under a limit of three

I-0087 was closed on 2026-10-06.

The cause:

```text
A ready story is held while its `touches` overlap an in-progress story's, so that two agents do not edit the same file (ADR-0046). The in-progress limit is three. On 2026-10-06 the board still ran one story at a time, because every ready story overlapped the one in progress.

Two things made the overlap total. Stories claim whole folders, such as `flai/internal/mcpserver`, and a folder claim holds every file in it. And one document, `docs/users/flai.md`, is changed by almost every flai story. A story in review keeps its claim, so the hold lasts until the operator accepts it.
```

The fix:

```text
S-0295, ADR-0096: only a story in progress holds a ready story (one in review holds nothing); a story's folder touch is narrowed in its claim to the files its tasks name; an overlap wholly inside the manifest's claims.shared patterns holds nothing. flai/internal/workitem/hold_i0087_test.go rebuilds the board of 2026-10-06 and shows the ready stories clear, and each rule's revert holding one again.
```

Source: I-0087.
