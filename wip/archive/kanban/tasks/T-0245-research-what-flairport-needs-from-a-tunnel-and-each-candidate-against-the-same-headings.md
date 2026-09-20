---
id: T-0245
type: task
nature: research
title: "Research: what flairport needs from a tunnel, and each candidate against the same headings"
status: done
parent: S-0069
owner: alex
created: 2026-09-20T04:07:27Z
updated: 2026-09-20T04:26:48Z
transitions:
  - to: ready
    at: 2026-09-20T04:07:41Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T04:07:41Z
    by: system-flow
  - to: done
    at: 2026-09-20T04:26:48Z
    by: system-flow
stream: S-0069
tags: []
---
# T-0245 Research: what flairport needs from a tunnel, and each candidate against the same headings

## Work
Write down first what flairport and flaiover need from a tunnel, from the code and the design: long-lived server-sent events (the board and the ADRs page), streamed newline-delimited JSON (acceptance), MCP over Streamable HTTP with a bearer token (ADR-0024), one hostname or a wildcard (S-0068), and whether the identity the edge established reaches the origin as a header. Then, for cloudflared with Cloudflare Access, Tailscale Funnel and Serve, ngrok, one self-hosted option (frp, rathole, or a WireGuard VPS with a reverse proxy), and anything better found on the way, fill the same headings from the provider's documentation: how the tunnel is made, where TLS ends and who holds the certificate, the identity check at the edge and whether it is free or paid, custom domain or assigned name, cost, limits (body size, idle and streaming timeouts, WebSocket and SSE), install and run as a service on Linux, WSL2, macOS, Windows, how it is configured (file, CLI, API), licence and self-hosting, and what the provider sees. Cite each claim with the date read.

## Done when
- The needs are listed in the finding before any candidate
- Each candidate has every heading filled or marked unknown, with citations and the date read

## Notes
