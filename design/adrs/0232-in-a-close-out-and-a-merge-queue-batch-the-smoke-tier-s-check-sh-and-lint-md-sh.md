---
id: ADR-0232
title: "In a close-out and a merge queue batch the smoke tier's check.sh and lint-md.sh are the only whole-repository check and markdown lint, and the integration tier's TestMonorepoIsClean and TestRepositoryLintsClean skip there"
status: accepted
date: 2026-10-31
supersedes: []
superseded_by: []
refines: [ADR-0183]
---

# ADR-0232 In a close-out and a merge queue batch the smoke tier's check.sh and lint-md.sh are the only whole-repository check and markdown lint, and the integration tier's TestMonorepoIsClean and TestRepositoryLintsClean skip there

## Context

A batch ran `flai check --strict` and the markdown lint of the whole repository twice: the integration tier's `TestMonorepoIsClean` (54 to 64 s under race) and `TestRepositoryLintsClean` (40 to 45 s), and the smoke tier's `scripts/check.sh` (3 s) and `scripts/lint-md.sh` (21 s). ADR-0183 scoped all four to the batch's stories. S-0548 brings a verification under five minutes and asks that each run once.

## Decision

Under `flai verify`, that is in a close-out (`CLOSE_OUT_STORY` set) or a merge queue batch (`MERGE_QUEUE_STORIES` set), the smoke tier's `scripts/check.sh` and `scripts/lint-md.sh` are the repository check and the markdown lint, scoped as ADR-0183 says. `TestMonorepoIsClean` skips there. `TestRepositoryLintsClean`, which checks that flai's own markdown rules report nothing markdownlint-cli2 passes, skips there unless the stories change `flai/internal/mdlint`, a file `mdlint.ConfigFiles` names, or `scripts/lint-md.sh`. Both still run in full where no story scopes them: `make integration`, CI's `go test -race ./...`, and `flai test --all`.

## Consequences

A batch spends about a minute and a half less on the integration tier. `check.sh` runs `flai check --strict`, which fails on the warnings `TestMonorepoIsClean` passes over, so nothing the batch caught is lost. A false positive of flai's markdown rules on a file outside `wip/`, brought by markdown alone, is caught by CI and `make integration` rather than in the batch; under `wip/`, `flai check --strict` lints with flai's rules, so `check.sh` still catches it. ADR-0183's scoping code in `TestMonorepoIsClean` is removed.

## Alternatives considered

- The integration tests run the check and the lint, and smoke drops its scripts: about 95 s of race CPU against 24 s, and the scripts are the binary and markdownlint-cli2, the authorities.
- A second `-skip` on the integration tier's command: it would replace `scripts/integration.sh`'s `-skip ^TestGateway`.
