---
id: I-0252
title: "TestServeSaysTheFlaiIsOlderThanTheProject fails under load: project.info is asked before flai serve has worked out the outdated flai"
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-11-02T06:04:39Z
last_reported: 2026-11-02T06:04:39Z
updated: 2026-11-09T10:49:31Z
---

# I-0252 TestServeSaysTheFlaiIsOlderThanTheProject fails under load: project.info is asked before flai serve has worked out the outdated flai

## Description
TestServeSaysTheFlaiIsOlderThanTheProject fails under load: project.info is asked before flai serve has worked out the outdated flai

## Instances

### 2026-11-02T06:04:39Z
Story: S-0506.
S-0506's flai test on serve, mcphttp, and cmd failed once at flai/internal/serve/outdated_test.go:175 with "project.info: <nil>" (FlaiOutdated nil) on a loaded host, and passed when run again alone. The test asks project.info once, at the dashboard's first connection, with no wait for serve to read the project's history; S-0506 changed nothing in serve's outdated check.

## Remediation

Story S-0617 remediates this issue, created from it at 2026-11-07T23:39:25Z.
Closed 2026-11-09T10:49:31Z: S-0617: the one instance (S-0506, 2026-11-02, project.info: <nil>) was I-0247's symptom, a watcher notification read as the reply, which S-0569 fixed on 2026-11-06 (ea5280843); project.info works out flai_outdated at each call and does not wait on serve's own look. T-2724 ran the three outdated serve tests 100 times each on 2 loaded CPUs since that fix, and none failed. T-2723 closed the one way left for flai_outdated to read nil while a newer flai is tagged: release.FlaiOutdated no longer keeps a failed git look for a minute, with TestFlaiOutdatedKeepsNoAnswerFromAGitThatFailed reproducing it.
