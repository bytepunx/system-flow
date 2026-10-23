---
id: D-0278
title: A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named
kind: cause
sources: [I-0113]
topics: [testing, flai, cli, go, code, all]
paths: [flai/internal/verify/parse.go, flai/internal/verify/parse_test.go, scripts/integration.sh, system-flow.yaml, scripts/README.md, design/system/project-manifest.md, docs/operators/settings.md, docs/users/flai.md, design/issues/I-0113-a-failed-integration-tier-in-the-close-out-shows-only-the-last-lines-of-go-test-so-the-failing-test-is-not-named.md, design/issues/summary.md, design/issues/I-0114-flai-verify-s-integration-tier-keeps-only-the-last-lines-of-go-test-s-output-so-the-failing-test-is-not-named.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0278 A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named

I-0113 was closed on 2026-10-12.

The cause:

```text
A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named
```

The fix:

```text
fixed by S-0313: a failing plain tier's finding keeps go test's --- FAIL, FAIL<tab><package>, and panic lines ahead of its last 20 lines (TestPlainNamesTheFailuresAheadOfItsLastLines), and this repository's integration tier runs scripts/integration.sh -json in flai as go-test-json, so its findings name each failing test by file and line
```

Source: I-0113.
