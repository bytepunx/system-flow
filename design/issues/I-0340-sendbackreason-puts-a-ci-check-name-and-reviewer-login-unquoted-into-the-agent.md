---
id: I-0340
title: sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt
class: impression
status: closed
count: 1
first_reported: 2026-11-10T09:16:23Z
last_reported: 2026-11-10T09:16:23Z
updated: 2026-11-11T22:03:05Z
---

# I-0340 sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt

## Description
sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt

## Instances

### 2026-11-10T09:16:23Z
First occurrence.

## Remediation

Story S-0711 remediates this issue, created from it at 2026-11-10T09:56:05Z.
Closed 2026-11-11T22:03:05Z: S-0711 fixed it: sendBackReason now quotes each check name and reviewer login through forgeName (flai/internal/serve/forgestatus.go). Each one is on a single line in a code span, with line ends and tabs made a space and other control and format characters dropped. So no forge name reaches a sent-back agent's start prompt as a line of its own. TestSendBackReasonQuotesWhatTheForgeNamed reproduces the case.
