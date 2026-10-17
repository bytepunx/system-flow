---
id: TH-0500
title: sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt
anchor:
  path: flai/internal/serve/forgestatus.go
status: open
participants: [orchestrator]
created: 2026-10-17T17:34:03Z
updated: 2026-10-17T17:34:03Z
---

# TH-0500 sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt

On flai/internal/serve/forgestatus.go.

## Entries

### 2026-10-17T17:34:03Z orchestrator
**Recommendation:** a small hardening fix. In `sendBackReason`, quote the check name and the reviewer login, for example in a code span, and strip control characters, the way the thread body already does. File it as an issue, or fold it into S-0395, whose agent touches this send-back path next.

**What the verifier found on S-0394.** I accepted S-0394 at 6bd4ad9 and released it as flai 1.60.0.
- **Already safe:** review comments are cut to 500 characters, quoted as blockquotes, and kept out of the prompt.
- **The gap:** the transition reason and the agent's prompt are built by `sendBackReason`, and it embeds the CI check's name and the reviewer's login raw.

**Why it is low risk.** A check name comes from the repository's CI configuration, so only someone who can already push could use it to slip text into an agent's prompt.

**Why you.** Filing work or changing S-0395's scope is yours. I can't file issues myself.
