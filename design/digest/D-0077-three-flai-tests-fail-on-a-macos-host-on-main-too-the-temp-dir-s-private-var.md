---
id: D-0077
title: "Three flai tests fail on a macOS host, on main too: the temp dir's /private/var path and a git identity the host has"
kind: cause
sources: [I-0045]
topics: [flai, cli, go, code, all]
paths: [flai/internal/serve/checks_test.go, flai/internal/serve/stop_test.go, flai/cmd/accept_resume_test.go, flai/cmd/doc_test.go, flai/internal/gittest, flai/internal/preview, scripts/install-test.sh, ".github/workflows/flai.yml", design/tech/ci.md, design/issues]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0077 Three flai tests fail on a macOS host, on main too: the temp dir's /private/var path and a git identity the host has

I-0045 was closed on 2026-10-01.

The cause:

```text
Three flai tests fail on a macOS host, on main too: the temp dir's /private/var path and a git identity the host has
```

The fix:

```text
S-0180: checks_test reads pwd -P and install-test.sh resolves its mktemp path, so both sides of each path comparison are resolved; the identity test isolates git with internal/gittest and sets user.useConfigOnly, so git no longer works out an identity from the hostname; the stop test waits until the stub has set its TERM trap, which macOS's slower-starting bash had not when the signal came; flai.yml now runs the tests on macos-latest too.
```

Source: I-0045.
