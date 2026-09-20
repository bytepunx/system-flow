---
id: T-0279
type: task
nature: feature
title: ADR-0030, the design, the conventions, and the documentation say how an agent connects
status: done
parent: S-0076
owner: alex
created: 2026-09-20T12:16:19Z
updated: 2026-09-20T12:30:44Z
transitions:
  - to: ready
    at: 2026-09-20T12:27:43Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T12:27:43Z
    by: system-flow
  - to: done
    at: 2026-09-20T12:30:44Z
    by: system-flow
stream: S-0076
tags: []
---
# T-0279 ADR-0030, the design, the conventions, and the documentation say how an agent connects

## Work
ADR-0030 supersedes ADR-0024: MCP over HTTP is flai's on the host; project identity is carried forward unchanged. design/tech/go-libraries.md records the revision targeted and what changed in it. flai-cli.md, flaiover-dashboard.md, overview.md, workflow.md, agent-narrative.md, dashboard-host-channel.md; conventions session-start.md and work-management.md in template/ first, identical here; docs/users/flai.md, flaiover.md, docs/operators/index.md; template CLAUDE.md.tmpl.

## Done when
- flai check --strict is clean here and in a project made from the template
- .mcp.json unchanged and working over stdio

## Notes
