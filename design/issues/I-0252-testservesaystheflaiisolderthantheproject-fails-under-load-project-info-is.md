---
id: I-0252
title: "TestServeSaysTheFlaiIsOlderThanTheProject fails under load: project.info is asked before flai serve has worked out the outdated flai"
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-02T06:04:39Z
last_reported: 2026-11-02T06:04:39Z
updated: 2026-11-02T06:04:39Z
---

# I-0252 TestServeSaysTheFlaiIsOlderThanTheProject fails under load: project.info is asked before flai serve has worked out the outdated flai

## Description
TestServeSaysTheFlaiIsOlderThanTheProject fails under load: project.info is asked before flai serve has worked out the outdated flai

## Instances

### 2026-11-02T06:04:39Z
Story: S-0506.
S-0506's flai test on serve, mcphttp, and cmd failed once at flai/internal/serve/outdated_test.go:175 with "project.info: <nil>" (FlaiOutdated nil) on a loaded host, and passed when run again alone. The test asks project.info once, at the dashboard's first connection, with no wait for serve to read the project's history; S-0506 changed nothing in serve's outdated check.

## Remediation
