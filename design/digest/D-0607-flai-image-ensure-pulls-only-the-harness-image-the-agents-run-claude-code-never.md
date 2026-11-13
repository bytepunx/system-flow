---
id: D-0607
title: "flai image ensure pulls only the harness image the agents run, claude-code, never bedrock, whose layers the harness image carries; the build builds bedrock, then claude-code over it"
kind: decision
sources: [S-0740]
topics: [docker, cli, release]
paths: [flai/internal/agentimage/pull.go, flai/internal/agentimage/ensure.go, flai/internal/agentimage/build.go, design/system/flai-cli.md]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0607 flai image ensure pulls only the harness image the agents run, claude-code, never bedrock, whose layers the harness image carries; the build builds bedrock, then claude-code over it

S-0740 built `flai image ensure` and `flai image build` (ADR-0295). It decided which images the pull fetches. From S-0740's Decisions (`wip/archive/agents/S-0740.md`):

```text
Only the harness image the agents run, claude-code, is pulled; bedrock is pulled by no one, since no agent runs it and the claude-code image carries its layers. The build builds bedrock, then claude-code over it.
```

So the pull asks the release for the harness image's list, `flai-agent-claude-code_<version>.digests`, and pulls `ghcr.io/bytepunx/flai-agent-claude-code` by its digest. Bedrock is still published and listed by each release (S-0734), and the build still makes it locally, tagged `flai-agent-bedrock:<version>`, as the base of the harness image. A later story should not add a bedrock pull to `Puller` or `Ensurer`, nor read bedrock's absence after a pull as a missing image.

`design/system/flai-cli.md` § "Where the agent images come from (ADR-0295)" describes the pull, under "Pull first, build as the fallback", as fetching both images. Its "Built (S-0740)" paragraph says what the code does: `flai-agent-<inner>` is pulled.

Sources: S-0740.
