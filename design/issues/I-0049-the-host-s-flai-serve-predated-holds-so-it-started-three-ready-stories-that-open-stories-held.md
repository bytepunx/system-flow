---
id: I-0049
title: The host's flai serve predated holds, so it started three ready stories that open stories held
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-09-26T18:17:31Z
last_reported: 2026-09-26T18:17:31Z
updated: 2026-10-01T10:55:14Z
---

# I-0049 The host's flai serve predated holds, so it started three ready stories that open stories held

## Description
The host's flai serve predated holds, so it started three ready stories that open stories held

## Instances

### 2026-09-26T18:17:31Z
2026-09-26: the host ran flai 1.18.5 after 1.19.0 (S-0128, the hold) was published. It started S-0130 while S-0129 and S-0133 claimed flaiover, then S-0131 and S-0132 while S-0130 claimed flai/cmd, flai/internal/mcpserver, docs, and design/system. S-0130 kept to Go until S-0129 was accepted (TH-0022); smoke's strict repository check fails on the wip.overlap warnings between S-0130 and S-0132 until one is accepted. Remedy: upgrade the host flai when a release changes what flai serve starts.

## Remediation
Closed 2026-10-01T10:55:14Z: S-0181: flai serve logs, the MCP inbox reports (flai_outdated), and the dashboard's host badge shows when the host's flai is older than the newest flai/v* tag reachable from the project's HEAD, with the upgrade command; system-flow.yaml's flai.minimum stops an older flai before it reads any item, and publishing a flai release that changes the front-matter fields raises it.
