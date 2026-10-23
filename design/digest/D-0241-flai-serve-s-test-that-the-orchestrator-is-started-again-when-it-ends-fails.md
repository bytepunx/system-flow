---
id: D-0241
title: flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run
kind: cause
sources: [I-0106]
topics: [continuous-improvement, flai, cli, go, code, all]
paths: [flai/internal/serve/orchestrate_test.go, design/issues/I-0106-flai-serve-s-test-that-the-orchestrator-is-started-again-when-it-ends-fails-under-the-close-out-s-full-integration-run.md, design/issues/summary.md, design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md, design/issues/I-0113-a-failed-integration-tier-in-the-close-out-shows-only-the-last-lines-of-go-test-so-the-failing-test-is-not-named.md, design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0241 flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run

I-0106 was closed on 2026-10-08.

The cause:

```text
flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run
```

The fix:

```text
S-0320: TestTheOrchestratorIsStartedAgainWhenItEnds no longer reads the host's clock around the retry. `orchestrateLab` gained an `at` clock, and the test's looks right after the failure, one second before `Ended + orchestrateRetry`, and at it take their time from the failed run's recorded `Ended`. The cause was a host-relative shift: `Ended` is written to the second, so the truncated fraction plus the host's delay could make the run due early. The boundary is now exact and does not depend on the host (commit 5a37eb2f).
```

Source: I-0106.
