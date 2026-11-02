---
id: I-0253
title: flai serve reads the Claude Code version in an agent image but never runs ADR-0106's permission_prompt check inside it, so a new image's version is only warned of
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-02T07:30:18Z
last_reported: 2026-11-02T07:30:18Z
updated: 2026-11-02T07:30:18Z
---

# I-0253 flai serve reads the Claude Code version in an agent image but never runs ADR-0106's permission_prompt check inside it, so a new image's version is only warned of

## Description
flai serve reads the Claude Code version in an agent image but never runs ADR-0106's permission_prompt check inside it, so a new image's version is only warned of

## Instances

### 2026-11-02T07:30:18Z
Story: S-0508.
S-0508 (T-2162, TH-0666): on harness container, flai serve reads `claude --version` from the image onto the run as harness_version and logs `container claude code version unchecked` when claude-checks.json has no passed record of it. The headless check of ADR-0106 § Decision part 2, one Write through permission_prompt in a scratch project, runs only with the host's claude, so an image whose Claude Code changed what permission_prompt must answer is not caught before a story's agent hangs on a protected write. Remedy: run the check inside the image, through the container harness, when the image's version has no record.

## Remediation
