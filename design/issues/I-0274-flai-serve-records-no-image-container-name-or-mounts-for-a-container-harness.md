---
id: I-0274
title: flai serve records no image, container name, or mounts for a container harness run, so a run cannot be shown afterwards to have used the project's layer image
class: efficiency
status: open
count: 1
cost: 15m
first_reported: 2026-11-07T05:41:15Z
last_reported: 2026-11-07T05:41:15Z
updated: 2026-11-10T08:34:37Z
---

# I-0274 flai serve records no image, container name, or mounts for a container harness run, so a run cannot be shown afterwards to have used the project's layer image

## Description
flai serve records no image, container name, or mounts for a container harness run, so a run cannot be shown afterwards to have used the project's layer image

## Instances

### 2026-11-07T05:41:15Z
Story: S-0567.
S-0567's live planner run on harness container (planner-S-0559, 2026-11-07T05:05:45Z): agents.json records harness container, command docker, and harness_version 2.1.312, and the journal says 'started docker (container)', but neither names the image the container ran nor its name, flai-agent-plan-<session>, nor its mounts. The container ran with --rm, and docker events kept nothing from the window, so the image had to be inferred from containerImageHost in flai/internal/serve/container.go and from Claude Code's version, 2.1.312 in the image against 2.1.319 on the host.

## Remediation

Story S-0665 remediates this issue, created from it at 2026-11-10T08:34:37Z.
