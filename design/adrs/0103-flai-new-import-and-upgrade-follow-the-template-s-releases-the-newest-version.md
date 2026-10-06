---
id: ADR-0103
title: "flai new, import, and upgrade follow the template's releases: the newest version tag unless --ref or an edit to system-flow.yaml says otherwise, and upgrade asks when they disagree"
status: proposed
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0015]
---

# ADR-0103 flai new, import, and upgrade follow the template's releases: the newest version tag unless --ref or an edit to system-flow.yaml says otherwise, and upgrade asks when they disagree

## Context

`flai upgrade` kept rendering an old template: 1.0.18 while `v1.0.60` was tagged. Three things caused it.

- The template cache cloned each repo and ref once and never fetched again. A clone of `main` stayed at the commit it was first cloned at.
- `flai upgrade` replaced `--ref` with the manifest's `template.ref` whenever `--template` was not given, so the operator could not name a version.
- A `template.version` edited by hand in `system-flow.yaml` made upgrade report "already at" that version and write it back, though the files were older.

Every config and every manifest so far holds `template.ref: main`. `flai template push --tag` tags each template release `v<version>`. [ADR-0015](0015-template-lock-file.md) made `system-flow.lock.yaml` the record of what was rendered; it did not say which version a project is at or which version an upgrade applies. The question is which template version `flai new`, `flai import`, and `flai upgrade` apply, and what happens when the operator's edit to the manifest disagrees with the newest release.

## Decision

**With no `--ref`, `flai new`, `flai import`, and `flai upgrade` apply the newest version tag of a git template when the ref they would use follows releases. `--ref` always wins. The lock records what was applied and is never a target. `flai upgrade` asks which version to apply when an edit to `system-flow.yaml` disagrees with the newest tag.**

- **Which refs follow releases.** The ref is the config's `template.ref` for `flai new` and `flai import`, and the manifest's `template.ref` for `flai upgrade`. It follows releases when it is empty, the template's default branch (`main`), or a version tag. Another branch or a commit is used as given. The config's default `template.ref` stays `main`.
- **The newest version tag.** A tag `vX.Y.Z`, the newest in semantic version order, pre-releases excluded. A template with no version tags falls back to the ref. A local template directory has no tags and is used as it is.
- **`--ref`.** It names what to apply and wins over the config and the manifest. `--ref 1.0.60` matches the tag `v1.0.60`. It is written to `template.ref` in `system-flow.yaml` and to the lock.
- **The version a project is at** is the lock's, or the manifest's when there is no lock. The lock records what was applied; flai never upgrades towards it.
- **When upgrade asks.** With no `--ref`, `flai upgrade` asks the operator which version to apply when `system-flow.yaml` names a different ref or version than the lock recorded and that differs from the newest tag. A changed `template.version` names the tag `v<version>`. The choices are the version the manifest names, the newest tag, and changing nothing. An edit that equals the newest tag is applied without asking. A recorded tag equal to the lock's is not a pin: the next upgrade takes the newest tag without asking.
- **What is recorded.** The version applied, chosen or not, is written to `template.ref` and `template.version` in `system-flow.yaml` and to the lock.
- **Without a terminal, or with `--yes`,** a run that would ask changes nothing and exits non-zero, naming the choices and the `flai upgrade --ref <tag>` to run. `--dry-run` says which version it would apply and whether it would ask.
- **The cache.** A cached clone at a branch is fetched again before it is used. A clone at a tag is reused. When the fetch fails, as offline, flai warns and uses the clone it has.

This refines ADR-0015 and does not supersede it: the lock still records hashes, topics, and variables as that ADR says.

## Consequences

- Existing projects and installs follow releases with no edit, because `main` follows releases.
- `flai upgrade` brings a project to the newest release. A stale cache or a hand-edited version no longer holds it back.
- To stay on an older version, the operator runs `flai upgrade --ref <tag>` or edits the manifest and chooses it when asked. Holding a project at a version needs a branch, a commit, or answering the prompt; a tag that equals the lock's does not hold it.
- An upgrade that would ask cannot run unattended. Scripts and CI pass `--ref`.
- With no `--ref`, `main` means the newest release, not the head of the branch. A project that wants unreleased template work names another branch or a commit.
- A run against a git template at a branch fetches once per use, which costs a network round trip. Offline, a cached clone that cannot be fetched is used with a warning, and `flai new` and `flai import` use the configured ref when the tags cannot be listed; `flai upgrade` with no `--ref` stops instead and names `--ref`, since without the tags it cannot tell the newest release from an older clone.

## Alternatives considered

- **Changing the config's default ref to empty.** Every existing config says `main`, so existing installs would stay on the branch.
- **Treating any recorded tag older than the newest as a pin that asks.** Every upgrade after the first would ask, against "upgrade takes the newest".
- **Always offering "stay at the lock's version".** The lock is never a target, so the choice became "change nothing".
