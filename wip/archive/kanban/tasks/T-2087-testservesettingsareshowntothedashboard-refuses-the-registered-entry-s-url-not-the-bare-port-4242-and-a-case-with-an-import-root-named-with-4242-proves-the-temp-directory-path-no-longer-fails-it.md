---
id: T-2087
type: task
nature: remediation
title: TestServeSettingsAreShownToTheDashboard refuses the registered entry's URL, not the bare port 4242, and a case with an import root named with 4242 proves the temp-directory path no longer fails it
status: done
parent: S-0488
owner: alex
created: 2026-10-23T18:10:41Z
updated: 2026-10-23T19:47:17Z
transitions:
  - to: ready
    at: 2026-10-23T19:34:23Z
    by: agent-S-0488
  - to: in-progress
    at: 2026-10-23T19:34:24Z
    by: agent-S-0488
  - to: done
    at: 2026-10-23T19:47:17Z
    by: agent-S-0488
stream: S-0488
tags: [cli, tests]
touches: [flai/cmd/serve_actions_test.go]
usage:
  source: log
  seconds: 773
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 27
      output: 3913
      cache_read: 1585534
      cache_write: 52140
      cost: 0.8126
---
# T-2087 TestServeSettingsAreShownToTheDashboard refuses the registered entry's URL, not the bare port 4242, and a case with an import root named with 4242 proves the temp-directory path no longer fails it

## Work

In `flai/cmd/serve_actions_test.go`, `TestServeSettingsAreShownToTheDashboard` (line 980) fails when the settings JSON contains `4242` anywhere. The JSON carries the project root from `tempProject` and the import root the test adds with `serve import add`, both `t.TempDir()` paths, whose random suffix can contain `4242`: that is I-0197's cause, and why the same commit passes on the next run.

- Change the leak check to refuse `127.0.0.1:4242`, the registered entry's URL as it was registered, in place of the bare `4242`; keep the `"token"` and `/p/harbour/k` checks as they are.
- Add the import root under a folder named with `4242`, such as `filepath.Join(t.TempDir(), "roots-4242")`, created with `os.MkdirAll`, so that every run carries the string the old check refused; keep the `"import_roots":["/` assertion. This is the test that reproduces the cause: with the old check it fails on every run, with the new one it passes.
- Run `go test ./cmd -run TestServeSettings` from `flai/`, then `flai test` on the file.

No source outside the test changes: `hostSettings` is right to show those paths.

## Done when

- The check at the end of `TestServeSettingsAreShownToTheDashboard` refuses `127.0.0.1:4242`, `"token"`, and `/p/harbour/k`, and nothing else.
- The test adds an import root whose path contains `4242`, and passes.
- `flai test flai/cmd/serve_actions_test.go` passes.

## Notes

First layer: waits for nothing.
