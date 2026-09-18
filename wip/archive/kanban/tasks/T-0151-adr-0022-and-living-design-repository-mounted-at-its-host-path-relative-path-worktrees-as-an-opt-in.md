---
id: T-0151
type: task
nature: remediation
title: "ADR-0022 and living design: repository mounted at its host path, relative-path worktrees as an opt-in"
status: done
parent: S-0050
owner: alex
created: 2026-09-18T19:47:23Z
updated: 2026-09-18T19:49:23Z
transitions:
  - to: ready
    at: 2026-09-18T19:48:22Z
    by: alex
  - to: in-progress
    at: 2026-09-18T19:48:23Z
    by: alex
  - to: done
    at: 2026-09-18T19:49:23Z
    by: alex
stream: S-0050
tags: []
touches: [design/adrs, design/system]
---

# T-0151 ADR-0022 and living design: repository mounted at its host path, relative-path worktrees as an opt-in

## Work
Write `design/adrs/0022-repository-mounted-at-its-host-path.md` refining ADR-0019 and add it to the index. Decision: `flai dashboard` mounts the repository at its absolute host path and sets `PROJECT_DIR` to it, at every git version; relative-path worktrees are an explicit, off-by-default setting and are never chosen by detecting a git version, because creating one sets `extensions.relativeWorktrees` on the whole clone. Record the alternatives the operator weighed, including the version switch. Update `design/system/flaiover-dashboard.md` (how the repository is mounted) and `design/system/flai-cli.md` (`flai dashboard`, `flai stream open`, the config key, the check rule).

## Done when
- ADR-0022 is accepted and indexed, and ADR-0019's status line in the index says it is refined
- The two design documents describe the mount and the setting
- `flai check --strict` is clean

## Notes
