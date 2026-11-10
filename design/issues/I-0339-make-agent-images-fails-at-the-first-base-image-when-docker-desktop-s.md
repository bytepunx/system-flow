---
id: I-0339
title: make agent-images fails at the first base image when Docker Desktop's credential helper desktop.exe cannot start from WSL
class: efficiency
status: closed
count: 1
cost: 5m
first_reported: 2026-11-10T09:13:13Z
last_reported: 2026-11-10T09:13:13Z
updated: 2026-11-10T10:56:45Z
---

# I-0339 make agent-images fails at the first base image when Docker Desktop's credential helper desktop.exe cannot start from WSL

## Description
make agent-images fails at the first base image when Docker Desktop's credential helper desktop.exe cannot start from WSL

## Instances

### 2026-11-10T09:13:13Z
Story: S-0682.
S-0682's T-2900 ran make agent-images in its worktree on 2026-11-10T09:09Z. The bedrock build stopped at loading golang's metadata with "error getting credentials - err: exit status 1", after WSL logged "UtilAcceptVsock:273: accept4 failed 110": ~/.docker/config.json names credsStore desktop.exe, and Windows interop was not answering. The base images are public, so DOCKER_CONFIG pointing at a scratch config.json of {} built them anonymously, leaving the operator's config untouched. scripts/agent-image.sh could say so, or fall back to an anonymous pull, when the helper fails.

## Remediation

Story S-0713 remediates this issue, created from it at 2026-11-10T10:11:03Z.
Closed 2026-11-10T10:56:45Z: S-0713: before its first docker build, scripts/agent-image.sh asks the credential helper credsStore names to list its logins, ten seconds at most; when it is not on PATH, fails, or does not answer, as desktop.exe does from WSL while Windows interop is down, the script says so and builds with no login, with DOCKER_CONFIG at a scratch configuration that keeps the current context and links the original folder's other entries, leaving the operator's config.json unchanged. scripts/agent-image-test.sh reproduces it with no Docker, and smoke.sh runs it.
