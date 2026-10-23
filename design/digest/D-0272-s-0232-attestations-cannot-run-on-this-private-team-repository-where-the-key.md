---
id: D-0272
title: "S-0232: attestations cannot run on this private Team repository; where the key goes; your public key: -----BEGIN PUBLIC KEY----- MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEXGMlBGpRPBSMMOcQrk81afzRbx/M…"
kind: recommendation
sources: [TH-0324]
topics: [release, security, releases, flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [".github/workflows/release-flai.yml", flai/.goreleaser.yaml, flaiover/src/lib/server/release.ts, design/tech/ci.md, design/system/release-signing.md, flaiover/src/lib/server/release.test.ts, flai/internal/buildinfo/releasekey.go, flai/internal/buildinfo/releasekey_test.go, docs/operators/runbooks/release-key.md, docs/operators/runbooks/README.md, docs/operators/settings.md, design/issues/I-0121-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md, design/issues/summary.md, design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md, design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md, design/issues/I-0129-an-empty-narrative-file-under-wip-agents-makes-every-flai-stream-log-in-the-project-fail.md, flai/internal/buildinfo/cosign.pub, design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md, design/issues/I-0134-flai-check-finds-issues-no-story-outside-the-story-at-close-out.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0272 S-0232: attestations cannot run on this private Team repository; where the key goes; your public key: -----BEGIN PUBLIC KEY----- MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEXGMlBGpRPBSMMOcQrk81afzRbx/M…

The operator, alex, answered TH-0324 at 2026-10-12T15:04:07Z.

What was asked:

```text
Two findings from checking the docs (T-1221, now in `design/system/release-signing.md`) change what S-0232 can deliver. I need your decision on both, and later your public key.
```

The recommendation it overruled:

```text
Recommendation: take the agent's recommendations 1 and 2, and run the runbook for item 3 once the agent says it is on the branch.

1. **Attestations: the conditional step.** ADR-0070 adds `actions/attest-build-provenance` "beside, for anyone with cosign or `gh`; flai depends on neither." So a step skipped while the repository is private takes nothing from what flai verifies. The cosign signature over `checksums.txt` still carries the whole guarantee. Skipping it changes how criterion 3 is met, which is a scope change and yours to confirm. If you would rather not tick a criterion whose step does not run yet, take the agent's first alternative: move criterion 3 to a backlog story for when the repository is public.
2. **Secrets: the `release` environment limited to `flai/v*` and `flaiover/v*` tags.** ADR-0070 holds the key as two Actions secrets "used only by the two release workflows". Environment secrets behind a tag deployment rule keep to that. Plain repository secrets would let any workflow on any branch read the key. This does not supersede the ADR, which names no secret scope.
3. **Your public key.** Only you can do this. Generate the pair locally, store it with `gh secret set --env release`, and reply with `cosign.pub` and its fingerprint. Never paste `cosign.key` or its password.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/adrs/0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md § Decision
```

The operator's answer:

```text
1. -----BEGIN PUBLIC KEY----- MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEXGMlBGpRPBSMMOcQrk81afzRbx/M izI2Hxp0X5OBfTx1U0q1kqKv2oMeL9d0FRsHo7kNz/pNi/gEs82LtNPbJA== -----END PUBLIC KEY----- 2. ee6df72313c9dbac1cf2270ed7f01cc523786d37bfce2dec12210215d6c5da40 3. take recommendations
```

Source: TH-0324.
