---
id: TH-0015
title: "S-0121 conflicts with S-0120: accept S-0120 first, then I rebase S-0121?"
anchor:
  path: wip/kanban/stories/S-0121-flai-serve-project-add-remove-and-list-manage-the-projects-the-host-serves-and-say-why-one-is-not-showing.md
  item: S-0121
status: resolved
participants: [agent-S-0118, alex, agent-S-0119]
created: 2026-09-26T05:39:31Z
updated: 2026-09-26T06:03:05Z
---

# TH-0015 S-0121 conflicts with S-0120: accept S-0120 first, then I rebase S-0121?

On wip/kanban/stories/S-0121-flai-serve-project-add-remove-and-list-manage-the-projects-the-host-serves-and-say-why-one-is-not-showing.md.

## Entries

### 2026-09-26T05:39:31Z agent-S-0118
S-0121 is done: every criterion is ticked and verified, and flai-test and flaiover-test pass. But `story/S-0121` and `story/S-0120` (now in review) conflict in five files: `flai/internal/serve/serve.go`, `flai/cmd/serve.go`, `flai/cmd/dashboard_agent.go`, `docs/users/flai.md`, and `design/system/flai-cli.md`. Whichever is accepted second will not rebase cleanly. S-0120 also added `Status.ImportProjects` and `Status.Unserved` for the projects below the import folders, which `flai serve project list` here lists as "not served" from its own scan. After S-0120 lands, list should take them from S-0120's state instead.

Recommended: accept S-0120 first. I keep S-0121 in progress until then, rebase it onto main, resolve the conflicts, make `list` use S-0120's `ImportProjects`/`Unserved`, rerun both test suites, and move it to review. Acceptance of S-0121 then needs no manual merge.

Alternative: I move S-0121 to review now as it is, and whoever accepts second resolves the conflicts by hand.

Which do you want? I will wait here.

### 2026-09-26T05:39:47Z agent-S-0118
Resolved: S-0120 was accepted into main (c1d3bf2) while this was open, so the recommended order happened on its own. I am rebasing S-0121 onto main, resolving the five conflicts, and switching `flai serve project list` to S-0120's `ImportProjects`/`Unserved` before I move it to review.

### 2026-09-26T05:40:15Z alex
117 is accepted, work through the conflicts and ask questions in threads as needed

### 2026-09-26T05:45:51Z agent-S-0118
Done. S-0121 is rebased onto main with S-0120. The five conflicts are resolved, and `flai serve project list` and `remove` now use S-0120's import-folder serving. flai-test and flaiover-test pass, and check --strict is clean. One consequence is open in the narrative: a project below an import folder is served whether registered or not, so `remove` cannot take it off the dashboard. It says so and names `flai serve import remove`. I propose a list of removed roots, perhaps with S-0122. Moving S-0121 to review now.

### 2026-09-26T05:58:33Z agent-S-0118
Resolved: S-0120 accepted first; S-0121 rebased onto it and the conflicts resolved.

### 2026-09-26T05:59:45Z alex
ok

### 2026-09-26T06:03:05Z agent-S-0119
Resolved: S-0120 (was S-0117) was accepted first and S-0121 (was S-0118) rebased onto it, as recommended
