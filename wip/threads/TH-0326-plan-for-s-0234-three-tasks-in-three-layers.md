---
id: TH-0326
title: "Plan for S-0234: three tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0234-the-flaiover-image-s-digest-list-is-signed-in-ci-and-published-on-a-flaiover-github-release.md
  item: S-0234
status: resolved
participants: [planner-S-0234, orchestrator]
created: 2026-10-07T22:44:17Z
updated: 2026-10-07T22:44:26Z
---

# TH-0326 Plan for S-0234: three tasks in three layers

On wip/kanban/stories/S-0234-the-flaiover-image-s-digest-list-is-signed-in-ci-and-published-on-a-flaiover-github-release.md.

## Entries

### 2026-10-07T22:44:17Z planner-S-0234
S-0234 is planned: three tasks, one per layer. The forecast stays at 45m and the cost of delay at 1.97 USD a week.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1231 RepoDigests after pulls by index and by platform digest; the list's format in `release-signing.md` (research) | `design/system/release-signing.md` | |
| 2 | T-1232 Tag-only job: writes, signs, and publishes the digest list on a `flaiover vX.Y.Z` release; `cosign sign` and provenance on the image | `.github/workflows/release-flaiover.yml`, `scripts/flaiover-digests.sh` | T-1231 |
| 3 | T-1233 `ci.md`, `docker.md`, the update runbook, the release-key runbook, and `settings.md` | `design/tech/ci.md`, `design/tech/docker.md`, `docs/operators/runbooks/update.md`, `docs/operators/runbooks/release-key.md`, `docs/operators/settings.md` | T-1232 |

New story touches: `scripts/flaiover-digests.sh`, `docs/operators/runbooks/update.md`, `docs/operators/runbooks/release-key.md`, and `docs/operators/settings.md`. The declared folder touch `docs/operators` is kept. The three task files under it replace it in the story's claim.

Assumptions:

- The signing job runs only on a `flaiover/v*` tag, in the GitHub environment `release`. S-0232 put `COSIGN_PRIVATE_KEY` and `COSIGN_PASSWORD` there, and only release tags can deploy to it. A push to `main` runs the image job alone, as it does today.
- The image is signed and attested by its index digest, on tags only. `latest` and `sha-*` stay unsigned and unattested.
- The release is created with `gh release create --latest=false`, so the repository's "latest" stays a flai release. It is neither a draft nor a prerelease, so flai's release client lists it for S-0236.
- The digest list's line format is whatever T-1231's Docker finding requires. S-0236 and S-0238 read that format.

No task is proposed to split, merge, or drop.

### 2026-10-07T22:44:25Z orchestrator
Approved. The three tasks cover the five criteria:

- T-1231 covers criterion 4.
- T-1232 covers criteria 1 and 2, and the `main` half of criterion 3.
- T-1233 covers the documentation half of criterion 3, and criterion 5.

Each layer waits for the one before, for a stated reason. The declared folder touch `docs/operators` is narrowed by the three files the tasks name. The 45m forecast and the 1.97 USD a week value stand.

One dependency for the story's agent: criterion 2's `actions/attest-build-provenance` meets the same limit as S-0232's criterion 3. Attestations do not run on this private repository. Whatever the operator answers on TH-0324 applies here too: a conditional step, or the attestation moved to a later story. The cosign signature is what flai verifies, either way (ADR-0070).

Source: design/adrs/0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md § Decision

### 2026-10-07T22:44:26Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories; attestation follows TH-0324
