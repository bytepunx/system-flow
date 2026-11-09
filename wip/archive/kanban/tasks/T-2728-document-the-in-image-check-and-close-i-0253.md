---
id: T-2728
type: task
nature: remediation
title: Document the in-image check and close I-0253
status: done
parent: S-0618
owner: alex
created: 2026-11-09T10:48:25Z
updated: 2026-11-09T11:18:39Z
transitions:
  - to: ready
    at: 2026-11-09T11:16:56Z
    by: agent-S-0618
  - to: in-progress
    at: 2026-11-09T11:16:56Z
    by: agent-S-0618
  - to: done
    at: 2026-11-09T11:18:39Z
    by: agent-S-0618
stream: S-0618
tags: [docs, issues]
touches: [design/system/flai-cli.md, design/system/agent-adapters.md, docs/users/flai.md, docs/operators/index.md, design/issues/I-0253-flai-serve-reads-the-claude-code-version-in-an-agent-image-but-never-runs-adr.md]
after: [T-2727]
usage:
  source: log
  seconds: 103
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 35
      output: 12456
      cache_read: 4624504
      cache_write: 75626
      cost: 1.7174
---
# T-2728 Document the in-image check and close I-0253

## Work

Rewrite what the docs say of a container start's Claude Code version, which today says the check is not run inside the image:

- `design/system/flai-cli.md` § The container harness (S-0508), the `containerReady` paragraph and step 4 of § A strategic agent's container (S-0509);
- `design/system/agent-adapters.md`, the paragraph on the check (`claudecheck.go`): add that a container's image is checked inside it, through the image's own `flai` on stdio, and recorded with the image;
- `docs/users/flai.md`, the container harness paragraph that names `container claude code version unchecked`;
- `docs/operators/index.md`, the `claude-checks.json` entry, which gains the image a record was checked in.

Then close the issue with `flai issue close I-0253 --reason`, saying that a container start whose image's Claude Code version has no record now runs ADR-0106's check inside the image (T-2726, T-2727). It writes the issue's file and the generated `design/issues/summary.md`.

This waits for T-2727, since the docs describe the behaviour it and T-2726 ship.

## Done when

- No document says the check is not run inside an image, and each names how it is run there.
- I-0253 is closed with a reason naming what fixed it, and `flai check --strict` is clean.

## Notes
