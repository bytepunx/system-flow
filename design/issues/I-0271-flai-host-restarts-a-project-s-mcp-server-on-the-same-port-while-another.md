---
id: I-0271
title: flai host restarts a project's MCP server on the same port while another process holds it, so the project has no MCP server until that process lets go
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-07T04:41:15Z
last_reported: 2026-11-07T04:41:15Z
updated: 2026-11-10T08:54:33Z
---

# I-0271 flai host restarts a project's MCP server on the same port while another process holds it, so the project has no MCP server until that process lets go

## Description
flai host restarts a project's MCP server on the same port while another process holds it, so the project has no MCP server until that process lets go

## Instances

### 2026-11-07T04:41:15Z
Story: S-0569.
Found by S-0569's T-2529 while reproducing I-0247's TestHostRunsServeAndTheMCPServersAndTakesThemWithIt failure. In flai/cmd/host.go, hostLauncher.mcpAddr checks a port with freePort and records it in l.taken[root], and the child listens on it later. When another process takes the port in between, the child exits with "cannot listen on 127.0.0.1:4256: bind: address already in use". The host then restarts it on the same port (Restarts:3, LastExit: exit status 1 in the reproduction), and mcpAddrFor never picks another. The live host already holds 13 of the default range 4243-4262. A live project whose remembered port another process holds gets no MCP server until that process lets go. The test now avoids the race by seeding its own port; flai's launcher is unchanged.

## Remediation

Story S-0697 remediates this issue, created from it at 2026-11-10T08:54:33Z.
