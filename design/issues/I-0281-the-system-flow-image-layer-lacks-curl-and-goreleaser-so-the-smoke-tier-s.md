---
id: I-0281
title: The system-flow image layer lacks curl and GoReleaser, so the smoke tier's install-test fails inside the agent container and a Go story's close-out cannot pass there
class: blocker
status: closed
count: 1
cost: 30m
first_reported: 2026-11-07T07:12:34Z
last_reported: 2026-11-07T07:12:34Z
updated: 2026-11-10T10:10:21Z
---

# I-0281 The system-flow image layer lacks curl and GoReleaser, so the smoke tier's install-test fails inside the agent container and a Go story's close-out cannot pass there

## Description
The system-flow image layer lacks curl and GoReleaser, so the smoke tier's install-test fails inside the agent container and a Go story's close-out cannot pass there

## Instances

### 2026-11-07T07:12:34Z
Story: S-0565.
Met on S-0565's T-2452, 2026-11-07, on S-0556's accepted head `45b5126`, run inside `flai-agent-sf-claude-code:local` with flai serve's mounts. Every other tier passed inside: gofmt, vet, golangci-lint, go-test with -race (3m38s), markdown, template, and integration (4m10s). The smoke tier, scoped as a close-out scopes it (`CLOSE_OUT_STORY=S-0556`), passed its template render, its repository check (3152 items, 0 findings), and its markdown lint. It then failed in `scripts/install-test.sh`: `install.sh` stops with "curl is required", since the layer has wget but no curl. On a throwaway image of the layer plus curl, it failed next at `scripts/flai-snapshot.sh` line 65: "goreleaser: not found". GoReleaser is installed on the host into `bin/` by `scripts/install-tools.sh`, not in the layer. The host's close-out of the same commit passed smoke in 43.9s. So a story whose close-out selects smoke, as every story touching `flai/` or `template/` does (ADR-0233), cannot close out inside the container. Remedy: the layer adds curl and GoReleaser at the version `scripts/install-tools.sh` pins (S-0510's layer), and a check after the layer is built runs `flai test smoke` in it.

## Remediation

Story S-0672 remediates this issue, created from it at 2026-11-10T08:34:49Z.
Closed 2026-11-10T10:10:21Z: S-0672: the system-flow layer (flai/images/system-flow/Dockerfile) now carries curl and openssl, which install.sh needs, GoReleaser v2.18.1 at scripts/install-tools.sh's pin, which scripts/flai-snapshot.sh runs, and procps-ng, whose ps takes the -p that scripts/release-server.sh uses and BusyBox's ps refuses, a third gap found while fixing this. scripts/agent-image.sh --check runs the smoke tier inside the built layer with the close-out container's mounts. It failed on the layer built before (curl, openssl, goreleaser, ps -p) and passed on the new one, install test included.
