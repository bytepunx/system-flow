---
id: D-0652
title: "Go tests run inside a close-out job inherit FLAI_VERIFY_PROGRESS=1 from flai serve, so a test that expects no progress output clears the variable with t.Setenv(closeout.ProgressEnv, \"\") first (I-0395)"
kind: cause
sources: [I-0395, S-0797, T-3491, S-0804]
topics: [testing, go, code, close-out]
paths: [flai/cmd/verify_test.go, flai/internal/serve/verifyjob.go]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0652 Go tests run inside a close-out job inherit FLAI_VERIFY_PROGRESS=1 from flai serve, so a test that expects no progress output clears the variable with t.Setenv(closeout.ProgressEnv, "") first (I-0395)

Since S-0804, flai serve's close-out job runs `flai verify` with `FLAI_VERIFY_PROGRESS=1` (`flai/internal/serve/verifyjob.go`). The `go test` processes the integration tier starts inherit it. So a test that means to run without the variable still has it, unless it clears it itself. I-0395 records how this broke main for every story (`design/issues/I-0395-testverifywritesitsprogresstostderronlywhenasked-fails-in-every-close-out-job.md`, Instances):

```text
S-0797's close-out job S-0797-20261117T020356Z stopped at the integration tier on flai/cmd/verify_test.go:1502 TestVerifyWritesItsProgressToStderrOnlyWhenAsked, "without the variable: exit 0", with stderr full of progress lines. [...] The test fails in every close-out that runs integration, at a head where S-0797 changed nothing in flai/cmd.
```

The fix, from I-0395's Remediation (S-0797, T-3491):

```text
TestVerifyWritesItsProgressToStderrOnlyWhenAsked clears FLAI_VERIFY_PROGRESS with t.Setenv(closeout.ProgressEnv, "") before its run without the variable, so it no longer inherits the close-out job's FLAI_VERIFY_PROGRESS=1.
```

So any test that checks what flai does without `FLAI_VERIFY_PROGRESS`, or any other variable the close-out job sets, clears that variable with `t.Setenv` before it runs. A test that passes under `flai test` but fails only in the close-out job may have inherited the job's environment. S-0797 fixed the test itself rather than leaving it for a story of its own: the fix was one line, and no story's close-out could pass the integration tier without it (`wip/archive/agents/S-0797.md`, Decisions).

Sources: I-0395, S-0797 (T-3491), S-0804.
