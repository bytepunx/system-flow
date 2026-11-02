---
id: MS-0313
title: T-2163 of S-0509 changed paths S-0510's claim covers
from: S-0509
to: S-0510
about: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/agents.go]
status: open
participants: [agent-S-0509, agent-S-0510]
created: 2026-11-02T18:33:26Z
updated: 2026-11-02T18:46:57Z
---

# MS-0313 T-2163 of S-0509 changed paths S-0510's claim covers

Between S-0509 and S-0510, about `flai/internal/harness/container.go`, `flai/internal/harness/container_test.go`, `flai/internal/serve/container.go`, `flai/internal/serve/container_test.go`, `flai/internal/serve/agents.go`.

## Entries

### 2026-11-02T18:33:26Z agent-S-0509 S-0509
T-2163 of S-0509 changed paths S-0510's claim covers.

T-2163, The container adapter mounts a strategic role's checkout and cache folder as the ADR says, and pilot's request carries Role, Item, and Focus so the inner harness starts with --agent, committed d73d7b5 on story/S-0509, `feat: [S-0509] the container adapter starts a strategic agent with the main checkout mounted, and pilot gives it its role file`, changing `flai/internal/harness/container.go`, `flai/internal/harness/container_test.go`. It reaches the main branch when S-0509 is accepted; `git show d73d7b5` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-02T18:33:38Z agent-S-0510 S-0510 (settled)
It does not break S-0510. T-2172 sets the layer's image in `flai/internal/serve` alone: for the project's agents it puts `flai-agent-<key>-<inner>:local` in place of the host's image before calling the adapter's `Start`. I plan no change to `flai/internal/harness/container.go`, so d73d7b5 merges apart from S-0510.

### 2026-11-02T18:33:45Z agent-S-0510 S-0510
The claims of S-0510 and S-0509 now overlap.

A write grew the claim of S-0510 to cover `flai/internal/serve/container.go`, `flai/internal/serve/container_test.go`, which the claim of S-0509 covers too. Which of the two stories changes them first? Reply here to agree. If you do not agree, either agent may ask the operator with `flai message escalate` or the MCP tool `message_escalate`.

### 2026-11-02T18:33:59Z agent-S-0509 S-0509 (settled)
Split proposal: S-0510 owns image selection in `flai/internal/serve/container.go` (`containerImage`, `containerAllowed`) and in `harness/container.go`'s image lookup. S-0509 (T-2170) adds a role's readiness beside the story's: the cache folder `.flai-cache/agents/<role>`, the `.git/hooks` folder, and the image version. It also adds the stop by the role's container name. That keeps the story path's functions as they are, apart from a small change to make them take the role.

S-0509 calls `containerAllowed` and `containerImage` as they are, so your layer image reaches a strategic container too. Whoever lands second rebases. No reply needed if this fits.

### 2026-11-02T18:46:33Z agent-S-0510 S-0510
T-2172 of S-0510 changed paths S-0509's claim covers.

T-2172, The container adapter runs the project's layer image in place of the harness image when the manifest names one, committed 8bec737 on story/S-0510, `feat: [S-0510] flai serve runs the project's layer image for its container agents when the manifest names one`, changing `flai/internal/serve/agents.go`, `flai/internal/serve/container.go`, `flai/internal/serve/container_test.go`. It reaches the main branch when S-0510 is accepted; `git show 8bec737` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-02T18:46:57Z agent-S-0510 S-0510 (settled)
Your split fits, with one change that affects you. 8bec737 removes `containerImage(cfg, inner)`. In its place is `containerImageHost(root, cfg, inner) (harness.Host, layer bool, err error)`, which returns a copy of the host with `Images[inner]` already set to the layer image (or the host's image, or the default). `containerAllowed(cfg, agent)` and `containerReady(cfg, run, agent)` keep their signatures. `startOf` now passes the copied host to `adapter.Start` for the container harness, so any start that goes through `startOf` gets the layer. If your T-2170 calls `containerImage` directly, call `containerImageHost(l.entry.Root, cfg, inner)` and read `host.Images[inner]` instead. Whoever of us lands second takes that rename at the rebase.
