---
id: I-0218
title: scripts/close-out.sh stops at its first step inside the agent container, since scripts/flai.sh builds flai with go, which the image lacks
class: blocker
status: open
count: 1
cost: 10m
first_reported: 2026-10-28T12:31:00Z
last_reported: 2026-10-28T12:31:00Z
updated: 2026-10-28T12:31:00Z
---

# I-0218 scripts/close-out.sh stops at its first step inside the agent container, since scripts/flai.sh builds flai with go, which the image lacks

## Description
scripts/close-out.sh stops at its first step inside the agent container, since scripts/flai.sh builds flai with go, which the image lacks

## Instances

### 2026-10-28T12:31:00Z
Story: S-0505.
S-0505 T-2166: agent-S-0513 ran `scripts/close-out.sh S-0513` in flai-agent-claude-code:local. Its first step, flai verify through scripts/flai.sh, stopped with `scripts/flai.sh: exec: line 8: .../bin/flai: not found` after `go: not found`; last line `close-out: S-0513 stopped at flai verify (exit 127)`. bin/ is the main checkout's and not mounted, and the image has no go. No tier was reached, so which tiers lack a tool could not be seen from inside: this repository's tiers need go and golangci-lint (gofmt, vet, golangci-lint, go-test), node and pnpm with flaiover's node_modules (flaiover-lint, vitest, svelte-check), markdownlint-cli2 (markdown), and go again for integration and smoke. The MCP verify tool, run on the host, passed S-0513 at every step. Remediation for S-0510: the project's image layer carries the tiers' tools, and the close-out in a container runs the image's flai rather than building one.

## Remediation
