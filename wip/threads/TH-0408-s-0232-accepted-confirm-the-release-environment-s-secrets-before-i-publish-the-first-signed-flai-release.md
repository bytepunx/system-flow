---
id: TH-0408
title: "S-0232 accepted: confirm the release environment's secrets before I publish the first signed flai release"
anchor:
  path: wip/archive/kanban/stories/S-0232-the-release-key-signs-flai-s-checksums-txt-in-ci-and-both-components-carry-the-public-key.md
  item: S-0232
status: resolved
participants: [orchestrator, alex]
created: 2026-10-12T15:28:13Z
updated: 2026-10-12T15:36:17Z
---

# TH-0408 S-0232 accepted: confirm the release environment's secrets before I publish the first signed flai release

On wip/archive/kanban/stories/S-0232-the-release-key-signs-flai-s-checksums-txt-in-ci-and-both-components-carry-the-public-key.md.

## Entries

### 2026-10-12T15:28:13Z orchestrator
**Recommendation:** confirm that the GitHub environment `release` holds both `COSIGN_PRIVATE_KEY` and `COSIGN_PASSWORD`, as `docs/operators/runbooks/release-key.md` § Store the secrets says. Then reply "publish", and I publish S-0232 with the next release.

## Why I am holding it

I accepted S-0232 at 339ac9aa. From now on, `release-flai.yml` runs in the `release` environment, and its "Release key present" step fails the job when either secret is empty. The next `flai/v*` tag is the first release that runs it.

I cannot see the environment's secrets. If one is missing, the tag is pushed and the release workflow fails with no assets. That is hard to undo.

`release_evaluate` gives 1 pending story, S-0232, at 2.63 USD a week, so nothing is lost by waiting.

## Also on main now

- Until S-0340 is accepted, a plain `make flai-snapshot` on main needs cosign and the key. Pass `--skip=sign` by hand meanwhile (S-0232's Notes, § Snapshot builds).
- The attestation step is skipped while the repository is private, as you took on TH-0324.

If you would rather publish yourself, run `flai release --pending`, and I will see it on the next acceptance.

### 2026-10-12T15:36:17Z alex
Resolved.
