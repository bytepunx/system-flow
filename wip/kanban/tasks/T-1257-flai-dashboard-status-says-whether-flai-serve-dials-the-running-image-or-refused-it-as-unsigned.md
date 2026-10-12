---
id: T-1257
type: task
nature: feature
title: flai dashboard status says whether flai serve dials the running image or refused it as unsigned
status: done
parent: S-0238
owner: alex
created: 2026-10-07T23:05:59Z
updated: 2026-10-12T20:33:29Z
transitions:
  - to: ready
    at: 2026-10-12T20:23:23Z
    by: system-flow
  - to: in-progress
    at: 2026-10-12T20:23:24Z
    by: system-flow
  - to: done
    at: 2026-10-12T20:33:29Z
    by: agent-S-0238
stream: S-0238
tags: [cli]
touches: [flai/cmd/dashboard.go, flai/cmd/dashboard_test.go]
after: [T-1255]
usage:
  source: log
  seconds: 601
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 70
      output: 409
      cache_read: 2705010
      cache_write: 199610
      cost: 1.3201
---
# T-1257 flai dashboard status says whether flai serve dials the running image or refused it as unsigned

## Work

Show T-1255's verdict in `flai dashboard status`, in `newDashboardStatusCmd` in `flai/cmd/dashboard.go`. This task waits for T-1255, which records the verdict. It shares no path with T-1256, so the two run together.

- Read the verdict from `flai serve`'s status, as `a.serveDir()` already does there.
- Beside the image, say whether `flai serve` dials it (signed), refused it ("the dashboard runs an unsigned image", with the digest), or did not measure it (with the reason).
- Add the verdict to `--json`.
- Update the command's `Long` help to describe the line.
- S-0236's digest line in this command says whether the digest is on a signed list. Keep the two lines consistent.

## Done when

- `dashboard_test.go` covers a listed image, an unlisted one, and an unmeasured one, in text and JSON, with `fakeRunner`.
- `flai test flai/cmd` passes.

## Notes
