---
id: I-0254
title: An agent container on a native Linux Docker engine cannot reach flai's HTTP MCP server on its default 127.0.0.1 bind through host.docker.internal
class: impression
status: open
count: 1
cost: 15m
first_reported: 2026-11-02T07:30:21Z
last_reported: 2026-11-02T07:30:21Z
updated: 2026-11-07T23:39:38Z
---

# I-0254 An agent container on a native Linux Docker engine cannot reach flai's HTTP MCP server on its default 127.0.0.1 bind through host.docker.internal

## Description
An agent container on a native Linux Docker engine cannot reach flai's HTTP MCP server on its default 127.0.0.1 bind through host.docker.internal

## Instances

### 2026-11-02T07:30:21Z
Story: S-0508.
S-0508: the container harness maps host.docker.internal to Docker's host gateway and rewrites FLAI_MCP_URL's host to it. On this host, Docker Desktop on WSL2 forwards the gateway to the host's loopback, where the server listens (S-0505, S-0506). A native Linux engine's host gateway is the docker0 address, which a server bound to 127.0.0.1 does not answer, so the agent's first MCP call fails unless the operator sets mcp.bind. Not measured here: no native engine on this host. The docs say to set mcp.bind; a remedy would have flai serve refuse a container start, saying so, when the bind cannot be reached from the bridge.

## Remediation

Story S-0619 remediates this issue, created from it at 2026-11-07T23:39:38Z.
