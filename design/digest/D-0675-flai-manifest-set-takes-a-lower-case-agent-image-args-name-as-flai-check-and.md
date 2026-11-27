---
id: D-0675
title: flai manifest set takes a lower-case agent.image.args name, as flai check and Docker do, so the command never refuses a build argument the check accepts
kind: decision
sources: [S-0839, T-3574]
topics: [flai, cli, go, docker]
paths: [flai/internal/manifest/settings.go, flai/internal/manifest/agent.go, flai/cmd/manifest.go, docs/operators/settings.md, design/system/project-manifest.md]
date: 2026-11-27
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0675 flai manifest set takes a lower-case agent.image.args name, as flai check and Docker do, so the command never refuses a build argument the check accepts

S-0839 made `flai manifest set` write `agent.image.args.<NAME>` and `agent.image.layer`, which before could only be written by hand ([D-0615](D-0615-s-0760-criterion-1-no-flai-command-writes-agent-image-args-so-the-manifest-s.md), I-0364).

Its narrative, `wip/archive/agents/S-0839.md`, records the decision under `## Decisions`:

```text
T-3574's sub-agent left lower-case ARG names allowed: agent.go's argPattern and Docker both allow them, and refusing them in manifest set alone would make it disagree with flai check. An unset prunes the args: and image: blocks it leaves empty, so a project can remove its whole image through the command.
```

So the `<NAME>` the docs write in upper case is a convention, not a rule. A later change that wants upper-case names enforces them in `argPattern`, where `flai check` and the command read it together, never in `flai manifest set` alone. `design/system/project-manifest.md` records the pruning on unset, but not the case of the names.
