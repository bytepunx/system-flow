---
id: TH-0738
title: S-0515 is verified and meets both criteria, but the queue refuses it while TH-0683 awaits your answer
anchor:
  path: wip/kanban/stories/S-0515-flai-s-http-mcp-server-refuses-an-agent-in-a-container-with-403-invalid-host-header-since-the-mcp-sdk-s-localhost-protection-admits-only-a-loopback-host-on-a-loopback-listener.md
  item: S-0515
status: resolved
participants: [orchestrator, alex]
created: 2026-10-29T08:05:44Z
updated: 2026-10-29T13:55:19Z
---

# TH-0738 S-0515 is verified and meets both criteria, but the queue refuses it while TH-0683 awaits your answer

On wip/kanban/stories/S-0515-flai-s-http-mcp-server-refuses-an-agent-in-a-container-with-403-invalid-host-header-since-the-mcp-sdk-s-localhost-protection-admits-only-a-loopback-host-on-a-loopback-listener.md.

## Entries

### 2026-10-29T08:05:44Z orchestrator
**Recommended:** confirm my recommendation on TH-0683, which approves the plan with its ADR and turns off the SDK's localhost Host check behind flai's two gates. Or answer it otherwise and resolve it. Then I enqueue S-0515 at once.

## Where it stands

- flai verify passed at the head, d0e52ac5. Integration and smoke are left to the queue.
- The verifier met both criteria:
  - `mcphttp.go` builds both SDK handlers with `DisableLocalhostProtection`.
  - `TestAForeignHostWithTheTokenIsServed` reproduces I-0213's 403 on `host.docker.internal:4243`.
  - ADR-0216 refines ADR-0030, and the docs say the token is the gate.
  - I-0213 is closed with its reason.
- The `Origin` refusal and the token check still run before the SDK handlers, and the test checks both on a foreign Host.
- `queue_add` refused it with `thread_open: thread TH-0683 is not resolved on S-0515`.

## Why I do not resolve TH-0683 myself

TH-0683 asks you to approve a change to the HTTP server's security posture. No recorded decision covered that until ADR-0216, which this story writes. Moving S-0515 to ready suggests you agree, but I do not take that as the answer. Under `confirm_after` (48h), flai serve confirms my recommendation on 2026-10-30 at about 13:05Z if you do not answer first.

### 2026-10-29T13:55:19Z alex
Resolved: S-0515 was accepted
