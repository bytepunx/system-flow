---
id: T-0430
type: task
nature: feature
title: flai serve serves the unregistered projects under import_roots, and flai serve status says why any is not served
status: done
parent: S-0117
owner: alex
created: 2026-09-26T05:25:44Z
updated: 2026-09-26T05:32:25Z
transitions:
  - to: ready
    at: 2026-09-26T05:26:14Z
    by: agent-S-0117
  - to: in-progress
    at: 2026-09-26T05:27:39Z
    by: agent-S-0117
  - to: done
    at: 2026-09-26T05:32:25Z
    by: agent-S-0117
stream: S-0117
tags: []
---

# T-0430 flai serve serves the unregistered projects under import_roots, and flai serve status says why any is not served

## Work

`flai serve` serves every system-flow project below an `import_roots` folder that is not registered, as it serves the projects below the folder it was started in (S-0102), without writing them to the registry. One function finds the projects below a set of folders and splits them into served and not served, each with its reason (no key, a key another project has, no dashboard to serve it for, a manifest that does not load). `flai serve` writes both into its state; `flai serve status` lists the unserved with their reasons, computing them itself when `flai serve` is not running.

## Done when

Behavior tests in `internal/serve` and `cmd` cover served, each reason, and status output; `make test` passes.

## Notes
