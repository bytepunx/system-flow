---
id: I-0350
title: The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-11-11T16:04:40Z
last_reported: 2026-11-11T16:04:40Z
updated: 2026-11-12T15:56:43Z
---

# I-0350 The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container

## Description
The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container

## Instances

### 2026-11-11T16:04:40Z
Story: S-0705.
S-0705's docs pass (T-2979) found two statements about the container harness that later stories made untrue and left: design/system/agent-adapters.md § A container harness of flai's own, What is built, still lists the close-out in the image (done since S-0628, ADR-0274) and the image's permission_prompt check (done since S-0618) as not done; docs/operators/index.md § Where a story's code runs and § Strategic agents in a container, What it cannot do yet, say the orchestrator's flai verify and flai accept run tiers inside its container, which ADR-0278 makes refuse in a role's container. A reader of either is told the wrong thing. S-0705 left them, as outside its scope.

## Remediation

Story S-0725 remediates this issue, created from it at 2026-11-11T18:16:51Z.
Closed 2026-11-12T15:56:43Z: S-0725 corrected both passages. design/system/agent-adapters.md § A container harness of flai's own now says the close-out in the image (S-0628, ADR-0274) and the image's permission_prompt check (S-0618) are built. docs/operators/index.md § Strategic agents in a container now says the orchestrator verifies and accepts through the MCP tools verify and accept on the host, and that its shell's flai accept and tier-running flai verify refuse in its container (ADR-0278).
