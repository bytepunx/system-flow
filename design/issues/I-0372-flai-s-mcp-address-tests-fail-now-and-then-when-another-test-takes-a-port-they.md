---
id: I-0372
title: flai's MCP address tests fail now and then when another test takes a port they closed or never held
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-11-13T20:39:59Z
last_reported: 2026-11-13T20:39:59Z
updated: 2026-11-13T20:39:59Z
---

# I-0372 flai's MCP address tests fail now and then when another test takes a port they closed or never held

## Description
flai's MCP address tests fail now and then when another test takes a port they closed or never held

## Instances

### 2026-11-13T20:39:59Z
Story: S-0779.
Run 43455794271 on main (flai 1.93.1), test (ubuntu-latest, -race): TestMCPListenFailureClosesWhatItOpened failed at mcp_http_test.go:805 with `listen tcp 127.0.0.2:42545: bind: address already in use`. The test holds 127.0.0.1:<port> and counts on 127.0.0.2:<port> being free before listenMCP runs and again after it closes what it opened; nothing ensures either, and another parallel test's connection to a 127.0.0.2 server can take that port as its source. flai behaved as it should; the test's assumption broke. Cost: reading the failing job's log and running it again.

## Remediation
