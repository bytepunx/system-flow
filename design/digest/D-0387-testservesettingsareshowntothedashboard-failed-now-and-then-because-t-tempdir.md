---
id: D-0387
title: TestServeSettingsAreShownToTheDashboard failed now and then because t.TempDir paths in the settings JSON sometimes held 4242
kind: cause
sources: [I-0197, S-0488]
topics: [flai, cli, go, testing, code]
paths: [flai/cmd/serve_actions_test.go]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0387 TestServeSettingsAreShownToTheDashboard failed now and then because t.TempDir paths in the settings JSON sometimes held 4242

I-0197, closed 2026-10-23: `TestServeSettingsAreShownToTheDashboard` in `flai/cmd` failed now and then in the close-out's go-test tier on branches that did not touch serve.

The cause, as S-0488 found it: the test "refused any 4242 in the settings JSON, which also holds two t.TempDir() paths whose random numbers sometimes held 4242". The fix narrows the check to the registered entry's URL, `127.0.0.1:4242`, and names the import root `harbour-4242`, which failed the old check on every run and passes the new one.

The lesson for later tests: an assertion that a short number or string is absent from output that holds `t.TempDir()` paths is flaky; assert on the exact field instead.

Sources: I-0197, S-0488.
