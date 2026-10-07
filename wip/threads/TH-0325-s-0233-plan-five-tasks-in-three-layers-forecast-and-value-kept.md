---
id: TH-0325
title: "S-0233 plan: five tasks in three layers, forecast and value kept"
anchor:
  path: wip/kanban/stories/S-0233-flai-self-upgrade-flai-host-upgrade-and-install-sh-verify-the-release-s-signature-before-installing-it.md
  item: S-0233
status: resolved
participants: [planner-S-0233, orchestrator]
created: 2026-10-07T22:41:34Z
updated: 2026-10-07T22:41:44Z
---

# TH-0325 S-0233 plan: five tasks in three layers, forecast and value kept

On wip/kanban/stories/S-0233-flai-self-upgrade-flai-host-upgrade-and-install-sh-verify-the-release-s-signature-before-installing-it.md.

## Entries

### 2026-10-07T22:41:34Z planner-S-0233
S-0233 now has five tasks in the backlog, in three layers. Its touches gained three doc files. Its forecast (1h15m) and value (3.29 USD a week) stand. The reasons are under `### Planning` in its Notes. No finalized words changed.

## Tasks

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1226 selfupgrade verifies `checksums.txt.sig` with `crypto/ecdsa` and refuses missing, bad, or unknown-key signatures, with tests | none | `selfupgrade.go` and its test, new `signature.go` and its test, `testdata/` |
| 1 | T-1227 `install.sh` verifies with `openssl`; a test checks its key matches flai's | none | `install.sh`, `scripts/install-test.sh`, `installsh_test.go`, `system-flow-check.yml` |
| 2 | T-1228 `self-upgrade --list` and `host versions` mark unverifiable releases; `host upgrade` refuses them | T-1226 | `selfupgrade.go`, `signature.go` and its test, `cmd/selfupgrade.go` and its test, `cmd/host.go`, `host_versions_test.go` |
| 3 | T-1229 the Updates page shows such a release as not installable, with no install control | T-1228 | `HostProcesses.svelte` and its test, `docs/users/flaiover.md` |
| 3 | T-1230 user, operator, and design docs say what is verified and what each refusal means | T-1227, T-1228 | `docs/users/flai.md`, `flai-reference.md`, `runbooks/install.md`, `runbooks/update.md`, `design/system/flai-cli.md` |

T-1226 and T-1227 share no path, so they run together. T-1229 and T-1230 do the same.

## Assumptions

- **The key.** Both checks use the PEM constant that S-0232's T-1225 puts in `flai/internal/buildinfo/releasekey.go`. `install.sh` carries its own copy of the PEM, and a test fails when the two differ.
- **`flai host upgrade`** runs `flai self-upgrade` as a child. It is verified through T-1226 without code of its own, and T-1228 checks the refusal.
- **CI step.** `scripts/install-test.sh` installs the latest published release. It passes only once a release signed by S-0232's workflow is published. If S-0233 merges before that, the smoke step fails until the release is out.
- **The mark.** A release with no `checksums.txt.sig` asset, or one signed by an unknown key, counts as unverifiable. The agent decides whether listing downloads each signature or checks only for the asset.
- **Release notes.** The story's Notes ask that the release notes say the first verifying release is installed without a check. T-1230 puts this in the docs. The story's agent decides where the release notes take it.
- **Left out.** I left `flai/internal/hostapi/writes.go` and the MCP `versions` tool unchanged. The host API's upgrade ends in self-upgrade's refusal, and no criterion names MCP. Say if MCP's list should carry the mark too.

## Proposals

No split, merge, or drop.

### 2026-10-07T22:41:43Z orchestrator
Approved. The tasks cover the six criteria:

- T-1226 covers criteria 1, 2, and 4 for flai.
- T-1227 covers criteria 3 and 4 for `install.sh`.
- T-1228 covers criterion 1 for `host upgrade` and the CLI half of criterion 6.
- T-1229 covers the Updates page half of criterion 6.
- T-1230 covers criterion 5.

The layers follow the `after`. No two tasks of a layer share a path. The two folder touches are explained: `docs/operators` was declared, and `testdata` holds files not named yet. The 1h15m forecast and the 3.29 USD a week value stand.

On MCP: leave the `versions` tool and `hostapi/writes.go` unchanged, as you planned. Criterion 6 names only `self-upgrade --list`, `host versions`, and the Updates page. Widening it is the operator's call, and it can be its own story.

Source: wip/kanban/stories/S-0233-flai-self-upgrade-flai-host-upgrade-and-install-sh-verify-the-release-s-signature-before-installing-it.md

### 2026-10-07T22:41:44Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories; MCP left out as the criteria name it not
