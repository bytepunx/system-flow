---
id: I-0279
title: "scripts/agent-container.sh always runs flai-agent-claude-code:local and ignores the manifest's agent.image.layer, so a close-out run by hand in it has none of the project's test tools"
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-07T06:58:47Z
last_reported: 2026-11-07T06:58:47Z
updated: 2026-11-10T08:34:46Z
---

# I-0279 scripts/agent-container.sh always runs flai-agent-claude-code:local and ignores the manifest's agent.image.layer, so a close-out run by hand in it has none of the project's test tools

## Description
scripts/agent-container.sh always runs flai-agent-claude-code:local and ignores the manifest's agent.image.layer, so a close-out run by hand in it has none of the project's test tools

## Instances

### 2026-11-07T06:58:47Z
Story: S-0565.
Met on S-0565's T-2452, 2026-11-07. TH-0821's option (a) runs S-0556's close-out once by hand inside the layered image with `scripts/agent-container.sh S-0556 -- scripts/close-out.sh S-0556`. But the script sets `IMAGE="flai-agent-claude-code:local"` (line 32) and never reads `agent.image.layer` from `system-flow.yaml`. A run there would lack Go, golangci-lint, and ShellCheck, which `flai-agent-sf-claude-code:local` carries. `flai serve`'s container adapter does pick the layer image (S-0510). The by-hand run used `docker run` directly on `flai-agent-sf-claude-code:local`, with the eight mounts and the environment flai serve gave S-0556's agent. Remedy: the script takes the layer's tag as `scripts/agent-image.sh` names it (`flai-agent-<key>-claude-code:local`) when the manifest names a layer, or an image variable.

## Remediation

Story S-0670 remediates this issue, created from it at 2026-11-10T08:34:46Z.
