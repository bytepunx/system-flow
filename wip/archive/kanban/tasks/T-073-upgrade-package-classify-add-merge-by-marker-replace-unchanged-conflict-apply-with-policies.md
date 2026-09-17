---
id: T-073
type: task
nature: feature
title: "upgrade package: classify add, merge by marker, replace unchanged, conflict; apply with policies"
status: done
parent: S-020
owner: alex
created: 2026-09-17T03:19:34Z
updated: 2026-09-17T03:25:21Z
transitions:
  - to: ready
    at: 2026-09-17T03:25:21Z
    by: agent
  - to: in-progress
    at: 2026-09-17T03:25:21Z
    by: agent
  - to: done
    at: 2026-09-17T03:25:21Z
    by: agent
stream: S-020
tags: [cli, template]
---

# T-073 upgrade package: classify add, merge by marker, replace unchanged, conflict; apply with policies

## Work
internal/upgrade: render the new template in memory; classify each path: add (absent in project), merge (both sides carry the baseline marker: new above, project below), replace (project file hash equals the lock's applied hash), same (identical), conflict (otherwise); apply with a policy per conflict (keep, replace); update lock and manifest template fields only when no conflict is unresolved.

## Done when
Unit tests cover every class and both policies.

## Notes
