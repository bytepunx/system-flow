---
id: T-2000
type: task
nature: remediation
title: The manifest takes a sub-project's install command, flai check validates it, and flaiover's names scripts/flaiover-install.sh --if-needed
status: done
parent: S-0483
owner: alex
created: 2026-10-21T12:10:20Z
updated: 2026-10-21T13:16:42Z
transitions:
  - to: ready
    at: 2026-10-21T13:07:06Z
    by: agent-S-0483
  - to: in-progress
    at: 2026-10-21T13:07:08Z
    by: agent-S-0483
  - to: done
    at: 2026-10-21T13:16:42Z
    by: agent-S-0483
stream: S-0483
tags: [flai]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, docs/operators/settings.md, system-flow.yaml, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/manifest/settings.go, flai/internal/check/install_test.go, flai/internal/manifest/shared.go]
after: [T-1998]
usage:
  source: log
  seconds: 574
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 47
      output: 13176
      cache_read: 2366157
      cache_write: 119264
      cost: 1.4588
---
# T-2000 The manifest takes a sub-project's install command, flai check validates it, and flaiover's names scripts/flaiover-install.sh --if-needed

## Work

Build the manifest key T-1998's ADR decides. It waits for T-1998 because the key's name, shape, and rules are that ADR's.

- In `flai/internal/manifest/manifest.go`, give `Project` an `install` command, a list of strings like a test tier's `command`, with `omitempty`. Validate it where the manifest's other sub-project keys are checked: an empty first element, or a key set on a project of kind `template`, is a `flai check` error naming the fix.
- Test the load and each refusal in `flai/internal/manifest/manifest_test.go`.
- Add the `projects[].install` row to `docs/operators/settings.md`, which `flai/cmd/settings_doc_test.go` requires for every manifest key.
- In `system-flow.yaml`, give the `flaiover` project `install: [scripts/flaiover-install.sh, --if-needed]`.

## Done when

- `flai check --strict` passes on this repository with flaiover's `install` set.
- The new manifest tests and `settings_doc_test.go` pass.
- `flai test flai/internal/manifest docs/operators/settings.md` passes.

## Notes
