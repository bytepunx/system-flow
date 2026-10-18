---
id: T-1896
type: task
nature: feature
title: The components methods take their sub-project as sub_project, and Go tests send them through the channel with the project key and fail on any params type with a project field
status: done
parent: S-0448
owner: alex
created: 2026-10-18T02:28:33Z
updated: 2026-10-18T02:35:49Z
transitions:
  - to: ready
    at: 2026-10-18T02:28:56Z
    by: agent-S-0448
  - to: in-progress
    at: 2026-10-18T02:28:57Z
    by: agent-S-0448
  - to: done
    at: 2026-10-18T02:35:49Z
    by: agent-S-0448
stream: S-0448
tags: []
touches: [flai/internal/hostapi/components.go, flai/internal/hostapi/components_test.go, flai/internal/hostapi/params_test.go, flai/internal/channel/channel.go]
usage:
  source: log
  seconds: 412
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 47
      output: 14744
      cache_read: 2364224
      cache_write: 82587
      cost: 1.3192
---

# T-1896 The components methods take their sub-project as sub_project, and Go tests send them through the channel with the project key and fail on any params type with a project field

## Work

- Rename `componentsParams.Project` to `SubProject`, `json:"sub_project"`, in `flai/internal/hostapi/components.go`, so that none of the four components methods reads `project` from its params; the answers' `project` fields stay as they are.
- Update the existing handler tests in `components_test.go` to send `sub_project`.
- Add a test that runs a real `channel.Client` with `MethodsFor`'s table against `channeltest`, and sends `components.list`, `components.status`, `components.start`, and `components.stop` with the project key in their params, as the dashboard does: list and status answer every supported sub-project, start and stop with `sub_project` reach the fake supervisor.
- Add `params_test.go`: parse the package's non-test Go files and fail on any type decoded as a method's params (`decode[T]`, `params(raw, &v)`) with a field whose JSON name is `project`, case-insensitively, as `encoding/json` matches it.
- Say in `channel.go`'s `call` doc comment that `project` in params is the channel's alone.
- Waits for nothing: it shares no path with T-1897.

## Done when

- The through-channel test and `params_test.go` both fail on `main`'s `components.go` and pass on the change.
- `flai test` on the changed paths passes.

## Notes
