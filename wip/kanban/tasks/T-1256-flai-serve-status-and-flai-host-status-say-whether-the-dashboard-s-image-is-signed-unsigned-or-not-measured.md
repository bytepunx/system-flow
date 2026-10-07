---
id: T-1256
type: task
nature: feature
title: flai serve status and flai host status say whether the dashboard's image is signed, unsigned, or not measured
status: backlog
parent: S-0238
owner: alex
created: 2026-10-07T23:05:55Z
updated: 2026-10-07T23:05:55Z
transitions: []
stream: S-0238
tags: [cli]
touches: [flai/cmd/serve.go, flai/cmd/serve_test.go, flai/cmd/host.go, flai/cmd/host_test.go]
after: [T-1255]
---
# T-1256 flai serve status and flai host status say whether the dashboard's image is signed, unsigned, or not measured

## Work

Show T-1255's verdict in two commands. This task waits for T-1255, which records the verdict.

- `flai serve status`, in `printServeStatus` in `flai/cmd/serve.go`, and in its `--json`. For each project, show one of three verdicts:
  - "the dashboard runs an unsigned image", with the image and the digest.
  - signed, with the digest.
  - "image not measured: <reason>; the release stamp alone is checked".
- `flai host status`, in `describeHost` in `flai/cmd/host.go`, and in its `--json`. Under the dashboard, show the same verdict, read from `flai serve`'s status.
- S-0237 changes both commands to say which side refused the other. Put the image's line beside its refusal line.

## Done when

- `serve_test.go` and `host_test.go` cover each verdict in text and JSON. Create `host_test.go` if S-0237 has not.
- `flai test flai/cmd` passes.

## Notes
