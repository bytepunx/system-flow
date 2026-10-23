---
id: D-0226
title: TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded
kind: cause
sources: [I-0102]
topics: [flai, cli, go, code, all]
paths: [flai/internal/mcpserver/folder_test.go, design/issues/I-0102-testwaitforworkacrossafolder-fails-under-the-close-out-s-full-go-run-when-the-host-is-loaded.md, design/issues/summary.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0226 TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded

I-0102 was closed on 2026-10-07.

The cause:

```text
TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded
```

The fix:

```text
S-0310. The failure was a tool error, not a timeout: readyStoryIn in flai/internal/mcpserver/folder_test.go added its criterion with os.WriteFile, which truncates the story file before filling it, and a wait_for_work held across it could poll in between, fail on the parse error, and answer an error the test dropped as map[]. readyStoryIn now writes through atomicfile.WriteFile, as flai writes items; every call in folder_test.go reports the tool's error; and TestAHeldWaitFailsOnATruncatedStory reproduces the truncated read.
```

Source: I-0102.
