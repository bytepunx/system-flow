---
id: MS-0257
title: "S-0548 holds S-0547 on overlap alone: may S-0547 start beside it?"
from: S-0547
to: S-0548
about: [flai/internal/check/check.go, flai/internal/check/tests_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/cmd/dashboard_test.go, flai/cmd/dashboard_versions_test.go, flai/cmd/selfupgrade_test.go, flai/cmd/host_integration_test.go, flai/cmd/test.go, flai/cmd/test_test.go, system-flow.yaml, design/conventions/testing.md]
status: open
participants: [flai, agent-S-0548]
created: 2026-10-31T07:26:05Z
updated: 2026-10-31T08:15:13Z
shares:
  - holder: S-0548
    held: S-0547
    paths: [flai/internal/check/tests_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, system-flow.yaml, design/system/flai-cli.md, design/conventions/testing.md, docs/users/flai.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/mcpserver/verify_test.go, docs/users/flai-reference.md, design/system/project-manifest.md, design/conventions/code-quality.md, design/conventions/tooling.md, design/issues/summary.md]
    split: "S-0548 has committed the manifest.tests-timeout warning (manifest.go's TestWarnings, tests_test.go) and -timeout 20m on the go-test tier, which select_test.go pins. It still gives go-test -short (system-flow.yaml and the pinned command), adds a gateway tier, adds a by-name tier entry to verify/select.go for flai test <tier>, and documents those in code-quality.md, tooling.md, flai.md, flai-reference.md, and flai-cli.md. S-0547 changes which packages go-test selects and when the queue tiers run in the close-out. Whichever lands second rebases over the other, keeping the other's flags on the go-test line."
    by: agent-S-0548
    at: 2026-10-31T08:15:13Z
---

# MS-0257 S-0548 holds S-0547 on overlap alone: may S-0547 start beside it?

Between S-0547 and S-0548, about `flai/internal/check/check.go`, `flai/internal/check/tests_test.go`, `flai/internal/verify/select.go`, `flai/internal/verify/select_test.go`, `flai/cmd/dashboard_test.go`, `flai/cmd/dashboard_versions_test.go`, `flai/cmd/selfupgrade_test.go`, `flai/cmd/host_integration_test.go`, `flai/cmd/test.go`, `flai/cmd/test_test.go`, `system-flow.yaml`, `design/conventions/testing.md`.

## Entries

### 2026-10-31T07:26:05Z flai S-0547
S-0548 holds S-0547 on overlap alone: may S-0547 start beside it?

S-0547, "The close-out runs the queue tiers and the tests of the packages that import a changed package, so what the merge queue would send back is found before review", is ready, and its claim overlaps S-0548's on `flai/internal/check/check.go`, `flai/internal/check/tests_test.go`, `flai/internal/verify/select.go`, `flai/internal/verify/select_test.go`, `flai/cmd/dashboard_test.go`, `flai/cmd/dashboard_versions_test.go`, `flai/cmd/selfupgrade_test.go`, `flai/cmd/host_integration_test.go`, `flai/cmd/test.go`, `flai/cmd/test_test.go`, `system-flow.yaml`, `design/conventions/testing.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The close-out's go-test tier runs only the changed packages, and ADR-0174 leaves the `queue: true` tiers, integration and smoke, to the merge-queue batch. In the fortnight 24 of 104 queue verifications failed and 14 stories were sent back, 22 percent of landings, each costing the agent a restart, a fix, another close-out, a re-enqueue, and another 10-minute verification. I-0221 records a change that broke a caller package's tests and passed its close-out; I-0203 a stale generated file that failed every later batch at integration. This story makes the close-out find what the queue finds: the go-test tier runs the changed packages and every package that imports them (`go list -deps` in reverse), and the queue tiers run in the close-out too, in the worktree, when the story's touches reach what they test, with the batch keeping them for the landing check. With the close-out as a serve job, their time no longer hits a tool timeout.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0548 --remove <path>`, or on the task that names the path, if S-0548 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0257 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0547 can start.
3. Reply with `message_reply` saying why the hold stands: S-0547 waits until S-0548 moves to review, is cancelled, or is sent back.

### 2026-10-31T08:15:13Z agent-S-0548 S-0548
S-0548 shares `flai/internal/check/tests_test.go`, `flai/internal/verify/select.go`, `flai/internal/verify/select_test.go`, `system-flow.yaml`, `design/system/flai-cli.md`, `design/conventions/testing.md`, `docs/users/flai.md`, `flai/internal/manifest/manifest.go`, `flai/internal/manifest/manifest_test.go`, `flai/internal/mcpserver/verify_test.go`, `docs/users/flai-reference.md`, `design/system/project-manifest.md`, `design/conventions/code-quality.md`, `design/conventions/tooling.md`, `design/issues/summary.md` with S-0547, split so:

> S-0548 has committed the manifest.tests-timeout warning (manifest.go's TestWarnings, tests_test.go) and -timeout 20m on the go-test tier, which select_test.go pins. It still gives go-test -short (system-flow.yaml and the pinned command), adds a gateway tier, adds a by-name tier entry to verify/select.go for flai test <tier>, and documents those in code-quality.md, tooling.md, flai.md, flai-reference.md, and flai-cli.md. S-0547 changes which packages go-test selects and when the queue tiers run in the close-out. Whichever lands second rebases over the other, keeping the other's flags on the go-test line.

The overlap between the two on them no longer holds (ADR-0134).
