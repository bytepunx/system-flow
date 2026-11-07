---
id: T-2560
type: task
nature: improvement
title: flai serve enable kluster writes and commits a block and a config for a project with none, and flai manifest set changes the block's keys
status: done
parent: S-0586
owner: alex
created: 2026-11-07T08:29:53Z
updated: 2026-11-07T09:02:35Z
transitions:
  - to: ready
    at: 2026-11-07T08:38:57Z
    by: agent-S-0586
  - to: in-progress
    at: 2026-11-07T08:38:59Z
    by: agent-S-0586
  - to: done
    at: 2026-11-07T09:02:35Z
    by: agent-S-0586
stream: S-0586
tags: [kluster, cli]
touches: [flai/internal/kluster/init.go, flai/internal/kluster/init_test.go, flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/kluster_test.go, flai/internal/docedit/file.go, flai/internal/docedit/file_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, docs/users/flai-reference.md, flai/internal/hostapi/kluster_test.go, docs/operators/settings.md]
usage:
  source: log
  seconds: 1416
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 122
      output: 40276
      cache_read: 14645998
      cache_write: 246726
      cost: 5.1595
---
# T-2560 flai serve enable kluster writes and commits a block and a config for a project with none, and flai manifest set changes the block's keys

## Work

Waits for T-2559, whose ADR decides the block's shape and the config's path.

- **The scaffold.** Add `kluster.Init` in `flai/internal/kluster/init.go`. It writes the `kluster:` block with `config: kluster.yaml` into `system-flow.yaml`, editing lines rather than rewriting the file, as `flai/internal/manifest/projects.go` and `shared.go` do, since `flai manifest set tests=` requotes lines it did not change (D-0449). It writes `kluster.yaml` when that file is absent and keeps it when it is there. It commits those two files alone, with the rule `flai manifest set --autocommit` follows, so `dashboard.autocommit: false` leaves them uncommitted.
- **Enabling.** In `flai/cmd/serve_actions.go`, `flai serve enable kluster` on a project with no block calls `kluster.Init` before it turns the action on, and says what it wrote and the commit. The Features toggle runs this command through `settings.action`, so the dashboard needs no change for it. `--all-projects` writes nothing. Agents never reach it: `flai guard` refuses them `flai serve enable` (ADR-0220).
- **`deploy` optional.** `KlusterErrors` in `flai/internal/manifest/manifest.go` stops reporting a block with no `deploy`. A deploy or a promotion still refuses one, as `flai/internal/kluster/deploy.go` does today.
- **The sentence.** `kluster.ErrNoBlock` in `flai/internal/kluster/kluster.go` says to turn the action on, not to edit the file.
- **The keys.** Add `kluster.config`, `kluster.namespace`, `kluster.build`, `kluster.deploy`, `kluster.teardown`, and `kluster.timeout` to the setting catalog in `flai/internal/manifest/settings.go`, with a command kind for the argument lists, so that `flai manifest set` writes them, refuses a bad value by field, and commits. `settings.manifest` then reaches them from the dashboard under `settings`, with no new host method.
- **The config file.** Let `docedit.FileReach` in `flai/internal/docedit/file.go` reach the file `kluster.config` names, so that `file.get` and `file.save` read and save it (ADR-0147).

Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

## Done when

- A test on a repository with no block runs `flai serve enable kluster` and finds the block, `kluster.yaml`, one commit holding those two files, the action on, and `flai check --strict` clean. Run again, it writes nothing.
- A test finds an existing `kluster.yaml` kept, `dashboard.autocommit: false` honoured, and `--all-projects` writing nothing.
- `flai manifest set kluster.deploy=...` and each other key is tested, a refused value naming its field, and the rest of `system-flow.yaml` byte for byte unchanged.
- `file.get` and `file.save` reach the config file, and no path outside what the manifest names.
- `flai test` passes on the paths this task changes.

## Notes
