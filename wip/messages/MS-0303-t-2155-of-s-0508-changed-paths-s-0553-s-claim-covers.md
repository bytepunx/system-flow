---
id: MS-0303
title: T-2155 of S-0508 changed paths S-0553's claim covers
from: S-0508
to: S-0553
about: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, docs/operators/settings.md, docs/users/flai-reference.md, design/system/agent-adapters.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md]
status: closed
participants: [agent-S-0508, flai, orchestrator]
created: 2026-11-02T07:09:44Z
updated: 2026-11-02T18:11:43Z
---

# MS-0303 T-2155 of S-0508 changed paths S-0553's claim covers

Between S-0508 and S-0553, about `flai/internal/harness/harness.go`, `flai/internal/harness/harness_test.go`, `docs/operators/settings.md`, `docs/users/flai-reference.md`, `design/system/agent-adapters.md`, `design/system/flai-cli.md`, `docs/operators/index.md`, `docs/users/flai.md`.

## Entries

### 2026-11-02T07:09:44Z agent-S-0508 S-0508
T-2155 of S-0508 changed paths S-0553's claim covers.

T-2155, The container adapter: harness.For selects it, config.inner names the inner harness, Capabilities and Reader are the inner harness's, and Start gives the docker run argv with three mounts and -e names only, committed c73bdf2 on story/S-0508, `feat: [S-0508] the container harness adapter runs a story's agent in docker with three mounts, -e names only, and pilot's request on stdin`, changing `flai/internal/harness/harness.go`, `flai/internal/harness/harness_test.go`. It reaches the main branch when S-0508 is accepted; `git show c73bdf2` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-02T07:10:00Z agent-S-0508 S-0508
T-2156 of S-0508 changed paths S-0553's claim covers.

T-2156, flai serve agent harness container sets the image per inner harness, the network, and the login path on the host, shown by flai serve settings and indexed in settings.md, committed 1951454 on story/S-0508, `feat: [S-0508] flai serve agent harness container sets the image per inner harness, the network, and the login on the host`, changing `docs/operators/settings.md`, `docs/users/flai-reference.md`. It reaches the main branch when S-0508 is accepted; `git show 1951454` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-02T07:30:33Z agent-S-0508 S-0508
T-2169 of S-0508 changed paths S-0553's claim covers.

T-2169, One story of this repository is worked end to end by flai serve on harness container, and the docs describe the adapter, its settings, its mounts, and what it never mounts, committed 237a1d1 on story/S-0508, `docs: [S-0508] describe the container harness, its settings, its three mounts, and what it never mounts`, changing `design/system/agent-adapters.md`, `design/system/flai-cli.md`, `docs/operators/index.md`, `docs/users/flai.md`. It reaches the main branch when S-0508 is accepted; `git show 237a1d1` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-02T07:30:33Z flai S-0508
story/S-0508 and story/S-0553 conflict when merged.

A trial merge of the two at flai stream sync conflicts in:

- `design/system/flai-cli.md`

Whichever of S-0508 and S-0553 is accepted second will stop on these paths when it rebases. Agree here who changes what: one narrows its change, or names the other in `after:` and waits for it. The next sync that finds the two merging cleanly closes this conversation. When you do not agree, either of you asks the operator with `flai message escalate` on this conversation, or the MCP tool `message_escalate`, saying what you could not agree.

### 2026-11-02T07:30:55Z agent-S-0508 S-0508 (settled)
S-0508 takes the conflict. Both branches change the package map lines in `design/system/flai-cli.md` § Internal structure: S-0553 its roles line, S-0508 the `harness/` and `protected/` lines. S-0553 is in review and lands first, so S-0508 resolves it at its next sync and keeps both changes. Nothing for S-0553 to change. The harness.go and harness_test.go changes from T-2155 add an adapter and fields to Host and Start; they do not touch the role definitions.

### 2026-11-02T18:11:43Z orchestrator
Closed: S-0508 was accepted
