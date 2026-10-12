---
id: TH-0417
title: "S-0237: refuse unsigned by default before S-0239's allowance exists?"
anchor:
  path: wip/kanban/stories/S-0237-flai-and-flaiover-exchange-their-release-stamps-in-hello-and-refuse-an-unsigned-peer-with-close-code-4403.md
  item: S-0237
status: resolved
participants: [agent-S-0237, alex]
created: 2026-10-12T17:05:46Z
updated: 2026-10-12T17:56:44Z
---

# TH-0417 S-0237: refuse unsigned by default before S-0239's allowance exists?

On wip/kanban/stories/S-0237-flai-and-flaiover-exchange-their-release-stamps-in-hello-and-refuse-an-unsigned-peer-with-close-code-4403.md.

## Entries

### 2026-10-12T17:05:46Z agent-S-0237
**Question.** S-0237 makes both sides close with 4403 on a peer without a valid stamp. The allowance for development builds, `dashboard.allow_unsigned`, is S-0239, which is in the backlog. Your host runs `ghcr.io/bytepunx/flaiover:latest` (`~/.flai/config.json`), and `latest` from main carries no stamp (`release-signing.md § The stamp as built`). So once flai and the dashboard both carry S-0237, your signed flai refuses that dashboard. A flai older than S-0237 (protocol 1) is refused by the new dashboard too, as criterion 1 asks.

**Recommendation: A.** Build it as ADR-0070 and the criteria say: refuse by default. The operator documentation says to run a released dashboard image (pick one on the Updates page, ADR-0118, or set `dashboard.tag`) and a released flai, both from the release that carries S-0237, and upgrade them together. Development builds stay refused until S-0239.

**B.** Also take the off switch only into S-0237: flai reads `dashboard.allow_unsigned` from the host's configuration, and the dashboard reads `FLAIOVER_ALLOW_UNSIGNED=1`. With it on, an unsigned peer is accepted and recorded as unsigned. S-0239 keeps passing the variable from `flai dashboard`, the warnings, the banner, and the statuses.

I am working the handshake meanwhile; the answer changes only one gate on each side.

### 2026-10-12T17:56:44Z alex
Resolved: S-0237 was accepted
