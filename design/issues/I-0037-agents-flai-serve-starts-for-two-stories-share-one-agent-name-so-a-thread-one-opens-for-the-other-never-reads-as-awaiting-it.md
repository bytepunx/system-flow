---
id: I-0037
title: Agents flai serve starts for two stories share one agent name, so a thread one opens for the other never reads as awaiting it
class: defect
status: open
count: 2
cost: 4m
first_reported: 2026-09-24T01:44:46Z
last_reported: 2026-09-24T07:49:46Z
updated: 2026-09-24T07:49:46Z
---

# I-0037 Agents flai serve starts for two stories share one agent name, so a thread one opens for the other never reads as awaiting it

## Description
Agents flai serve starts for two stories share one agent name, so a thread one opens for the other never reads as awaiting it

## Instances

### 2026-09-24T01:44:46Z
S-0107: agent-S-0107 opened TH-0004 on S-0106 for agent-S-0106. flai serve started both with FLAI_AGENT=system-flow, so the thread is awaiting 'other' with the only participant being the same name, and S-0106's inbox does not list it as awaiting it. S-0107 fell back to adding host.start and host.stop itself after S-0106 lands.

### 2026-09-24T07:49:46Z
The reverse case: TH-0009, which S-0112's agent opened, read as awaiting S-0113's agent once the operator answered it, because both agents are named system-flow. S-0113's agent acted on it (T-0396) because the answer settled its own trailer decision. Otherwise the thread would have pulled it off its story.

## Remediation
