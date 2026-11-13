---
id: D-0614
title: "A development flai with no build command still builds its agent images from the tree, as S-0740 left it; ADR-0295 §3's failing job with no command named is for S-0760, once scripts/agent-image.sh is this repository's command"
kind: decision
sources: [S-0759]
topics: [docker, cli, code, release]
paths: [flai/internal/agentimage/ensure.go, flai/internal/serve/imagejob.go, scripts/agent-image.sh, design/system/flai-cli.md]
date: 2026-11-13
status: absorbed
absorbed_by: "design/system/flai-cli.md#where-the-agent-images-come-from-adr-0295"
asked_again: []
cited: []
---

# D-0614 A development flai with no build command still builds its agent images from the tree, as S-0740 left it; ADR-0295 §3's failing job with no command named is for S-0760, once scripts/agent-image.sh is this repository's command

S-0759 built the operator's own image build command and the switch that turns the image job off (ADR-0295 §7). It chose not to build one part of ADR-0295 yet. From S-0759's Decisions (`wip/archive/agents/S-0759.md`):

```text
A development flai's build from the tree stays as S-0740 left it. ADR-0295 point 3's 'with none named, the job fails' is for S-0760, which makes `scripts/agent-image.sh` this repository's command, and changing it here would break this host's builds before that.
```

ADR-0295 §3 says a development flai is never pulled for, and its images come from the operator's build command; with none named, the job fails, naming the setting. Since S-0759, a development flai with no command named still builds from the tree, as `design/system/flai-cli.md` § "Where the agent images come from (ADR-0295)" says under **Built (S-0759)**: "With no command named, flai's own pull and build run as before, a development flai's build from the tree among them; S-0760 names this repository's command."

So S-0760, or whichever story makes `scripts/agent-image.sh` this repository's command, owns the change that makes the job fail for a development flai with no command named. It should name the command on this host first, so this host's builds do not break. A story before then should not treat the tree build as a defect.

Source: S-0759.
