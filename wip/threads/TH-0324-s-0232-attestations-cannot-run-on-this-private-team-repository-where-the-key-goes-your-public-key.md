---
id: TH-0324
title: "S-0232: attestations cannot run on this private Team repository; where the key goes; your public key"
anchor:
  path: wip/kanban/stories/S-0232-the-release-key-signs-flai-s-checksums-txt-in-ci-and-both-components-carry-the-public-key.md
  item: S-0232
status: open
participants: [agent-S-0232, orchestrator]
created: 2026-10-07T22:23:53Z
updated: 2026-10-07T22:24:14Z
---

# TH-0324 S-0232: attestations cannot run on this private Team repository; where the key goes; your public key

On wip/kanban/stories/S-0232-the-release-key-signs-flai-s-checksums-txt-in-ci-and-both-components-carry-the-public-key.md.

## Entries

### 2026-10-07T22:23:53Z agent-S-0232
Two findings from checking the docs (T-1221, now in `design/system/release-signing.md`) change what S-0232 can deliver. I need your decision on both, and later your public key.

## 1. Criterion 3: attestations are not available here

GitHub: "If you are on a GitHub Free, GitHub Pro, or GitHub Team plan, artifact attestations are only available for public repositories." `bytepunx/system-flow` is private and `bytepunx` is on Team (`gh api`, 2026-10-07). An `actions/attest-build-provenance` step would fail every flai release.

Recommendation: add the step (v4.2.2, over `dist/checksums.txt` and the archives) with `if: ${{ !github.event.repository.private }}`, so it is skipped while the repository is private and runs once it is public. I tick criterion 3 on that basis and say so in the story's notes.

Alternatives:

- Leave the step out and move criterion 3 to a backlog story for when the repository is public or on Enterprise Cloud.
- Make the repository public, or move the organisation to Enterprise Cloud, and the step runs unconditionally.

## 2. Where the key's secrets live

Required reviewers on an Actions environment are also public-only on Team, so the design's reviewer gate is not available. Environment secrets and deployment tag rules are.

Recommendation: an environment `release` whose deployment rule admits only tags `flai/v*` and `flaiover/v*`, holding `COSIGN_PRIVATE_KEY` and `COSIGN_PASSWORD`, plus a tag ruleset limiting who may create those tags. The pair is generated locally with `cosign generate-key-pair` and stored with `gh secret set --env release`, since `generate-key-pair github://` writes only repository secrets. `release-flai.yml` runs its job in that environment.

Alternative: plain repository secrets, written by `cosign generate-key-pair github://bytepunx/system-flow`. Simpler, but any workflow run on any branch can read the key.

## 3. Your public key (criteria 1 and 4)

The runbook, `docs/operators/runbooks/release-key.md` on `story/S-0232`, comes next. Once it is there I will reply here. Please run it, then reply with the whole `cosign.pub` (a PEM `PUBLIC KEY` block; it is public) and the fingerprint the runbook prints. Never paste `cosign.key` or the password.

Meanwhile I write the runbook and the workflow on recommendations 1 and 2, and change them if you choose otherwise.

### 2026-10-07T22:24:14Z orchestrator (recommendation)
Recommendation: take the agent's recommendations 1 and 2, and run the runbook for item 3 once the agent says it is on the branch.

1. **Attestations: the conditional step.** ADR-0070 adds `actions/attest-build-provenance` "beside, for anyone with cosign or `gh`; flai depends on neither." So a step skipped while the repository is private takes nothing from what flai verifies. The cosign signature over `checksums.txt` still carries the whole guarantee. Skipping it changes how criterion 3 is met, which is a scope change and yours to confirm. If you would rather not tick a criterion whose step does not run yet, take the agent's first alternative: move criterion 3 to a backlog story for when the repository is public.
2. **Secrets: the `release` environment limited to `flai/v*` and `flaiover/v*` tags.** ADR-0070 holds the key as two Actions secrets "used only by the two release workflows". Environment secrets behind a tag deployment rule keep to that. Plain repository secrets would let any workflow on any branch read the key. This does not supersede the ADR, which names no secret scope.
3. **Your public key.** Only you can do this. Generate the pair locally, store it with `gh secret set --env release`, and reply with `cosign.pub` and its fingerprint. Never paste `cosign.key` or its password.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/adrs/0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md § Decision
