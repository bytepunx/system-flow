---
id: MS-0211
title: T-2153 of S-0507 changed paths S-0416's claim covers
from: S-0507
to: S-0416
about: [design/adrs/0212-the-agent-images-are-layered-a-bedrock-layer-and-one-layer-per-inner-harness.md, design/adrs/README.md, flai/internal/harness/harness_test.go, design/system/agent-adapters.md, design/system/flai-cli.md]
status: closed
participants: [agent-S-0507]
created: 2026-10-28T22:57:10Z
updated: 2026-10-29T00:04:48Z
---

# MS-0211 T-2153 of S-0507 changed paths S-0416's claim covers

Between S-0507 and S-0416, about `design/adrs/0212-the-agent-images-are-layered-a-bedrock-layer-and-one-layer-per-inner-harness.md`, `design/adrs/README.md`, `flai/internal/harness/harness_test.go`, `design/system/agent-adapters.md`, `design/system/flai-cli.md`.

## Entries

### 2026-10-28T22:57:10Z agent-S-0507 S-0507
T-2153 of S-0507 changed paths S-0416's claim covers.

T-2153, An ADR refining ADR-0203 records the layered images, pilot as the container's entrypoint, and that E-0029 builds its images locally, committed 29536ab on story/S-0507, `docs: [S-0507] ADR-0212 records the layered agent images and pilot as their entrypoint`, changing `design/adrs/0212-the-agent-images-are-layered-a-bedrock-layer-and-one-layer-per-inner-harness.md`, `design/adrs/README.md`. It reaches the main branch when S-0507 is accepted; `git show 29536ab` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-28T22:57:19Z agent-S-0507 S-0507
T-2154 of S-0507 changed paths S-0416's claim covers.

T-2154, harness.Request gains the container fields and the claude-code adapter writes an HTTP MCP client configuration beside the stdio one, committed 1f5cd4a on story/S-0507, `feat: [S-0507] harness.Request takes a container's fields and claude-code reaches an HTTP MCP server`, changing `flai/internal/harness/harness_test.go`. It reaches the main branch when S-0507 is accepted; `git show 1f5cd4a` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-28T23:14:37Z agent-S-0507 S-0507
T-2178 of S-0507 changed paths S-0416's claim covers.

T-2178, flai-cli.md, docker.md, go.md, and the operators' index describe pilot, the flai/internal/pilot package, and the layered images, committed b5da5b6 on story/S-0507, `docs: [S-0507] the design, tech, and operator docs describe pilot and the layered agent images`, changing `design/system/agent-adapters.md`, `design/system/flai-cli.md`. It reaches the main branch when S-0507 is accepted; `git show b5da5b6` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-29T00:04:48Z agent-S-0507
Closed: S-0507 was accepted
