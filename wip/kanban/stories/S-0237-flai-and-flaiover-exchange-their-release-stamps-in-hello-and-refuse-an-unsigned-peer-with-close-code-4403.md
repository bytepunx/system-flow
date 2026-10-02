---
id: S-0237
type: story
nature: feature
title: flai and flaiover exchange their release stamps in hello and refuse an unsigned peer with close code 4403
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:24Z
updated: 2026-10-02T12:37:24Z
transitions: []
tags: [cli, dashboard]
topics: [release, security]
touches: [flai/internal/channel, flai/internal/serve, flaiover/src/lib/server/agent.ts, flaiover/src/routes, docs/operators, docs/users/flaiover.md]
after: [S-0235]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0237 flai and flaiover exchange their release stamps in hello and refuse an unsigned peer with close code 4403

## Goal

The channel's handshake carries each side's release stamp; each verifies the other's with the public key and the component it names before the credential proof, and a missing, unverifiable, or wrong stamp closes the connection with `4403 unsigned peer` (ADR-0070). The refusal is visible on both sides.

## Acceptance criteria
- [ ] `hello` carries flai's `release` (statement and signature) and the answer carries the dashboard's; each side verifies the other's signature with the embedded key, checks the component (`flaiover` from the dashboard, `flai` from flai), and otherwise closes with `4403`; the protocol number is bumped and an older peer's hello without the field is treated as unsigned.
- [ ] flai backs off a minute on `4403` as on `4409`, logs one `error` with the peer's version and the reason, and `flai serve status` and `flai host status` say "the dashboard refused this flai: unsigned" or "flai refused the dashboard: unsigned".
- [ ] The dashboard shows a project whose flai it refused with the reason where it shows a flai that lacks a required method, and its connection list says per project whether the flai is signed.
- [ ] Tests on both sides cover a signed peer, an unsigned one, a wrong component, and an older peer without the field.
- [ ] The operator documentation states the promise and its limit as `release-signing.md § Verifying the peer` does, and `docs/users/flaiover.md` describes what the connection list shows.

## Tasks

## Notes
