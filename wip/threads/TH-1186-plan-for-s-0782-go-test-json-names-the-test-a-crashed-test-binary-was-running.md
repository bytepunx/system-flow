---
id: TH-1186
title: "Plan for S-0782: go-test-json names the test a crashed test binary was running"
anchor:
  path: wip/kanban/stories/S-0782-flai-test-reports-a-go-test-failure-in-flai-cmd-that-names-no-test-and-the-same.md
  item: S-0782
status: resolved
participants: [planner-S-0782, orchestrator]
created: 2026-11-13T22:53:31Z
updated: 2026-11-13T22:54:05Z
---

# TH-1186 Plan for S-0782: go-test-json names the test a crashed test binary was running

On wip/kanban/stories/S-0782-flai-test-reports-a-go-test-failure-in-flai-cmd-that-names-no-test-and-the-same.md.

## Entries

### 2026-11-13T22:53:31Z planner-S-0782
S-0782 is planned: two tasks in two layers, touches, a forecast of 20m, and a cost of delay of 15 USD a week.

## Proposed solution

The `go-test-json` format (`flai/internal/verify/gotest.go`) reports a package that fails with no failed test and no timeout panic using only the package's own output lines. When a test binary crashes mid-run (a goroutine panic, `fatal error:`, a test's `os.Exit`, a kill), `go test -json` writes the crash as output of the running test and sends that test no `fail`. So the finding names the package with no test, line, or message: the I-0365 shape. S-0540 fixed only the timeout panic. The fix names the test that started and never ended, places it by the crash's stack frame, and gives the crash and exit lines. When no such test exists, it falls back to the package's last output lines, never an empty message.

## Tasks and layers

| Layer | Task | after | Touches |
|-------|------|-------|---------|
| 1 | T-3344 go-test-json names the test a crashed test binary was running, with its crash and exit | none | `flai/internal/verify/gotest.go`, `flai/internal/verify/parse_test.go`, `flai/internal/verify/testdata/go-test-crash.json`, `design/system/project-manifest.md` |
| 2 | T-3346 Close I-0365 with flai issue close | T-3344, since its reason names what T-3344 changed | the I-0365 file |

## Figures

- Forecast 20m, delivery 2026-11-14T05:28:00Z. I raised it from flai's 6m because S-0540 took 908 s for the same parser's timeout case.
- Cost of delay 15 USD a week, as `flai cod` gives it from flai's 6m input. Unchanged.

## Assumptions

- I read the criterion "the cause no longer occurs" as the unreadable finding. The flake behind the instance kept no output, so nothing can name it now. After this fix a recurrence names it, gets recorded as a new instance, and becomes its own story.
- The crash shapes are recorded from a real `go test -json` run on go 1.26 as a fixture, as S-0540 recorded its timeout fixture.
- No `after` on another story, and no folder touch.

## Proposals

- Not drafted: `flai test` and the MCP tool `test` could keep the failing tier's whole output in a log and name it, as the merge queue's batch already does (`RunOptions.LogDir`). That would cover crash shapes the parser still misses. I left it out to keep the story small. Say if you want it as a follow-up story.

### 2026-11-13T22:53:53Z orchestrator
Approved: the plan fits S-0782.

- **Tasks:** T-3344 meets criterion 1 with the crash fixture and its test. T-3346 meets criterion 2. The `after` from T-3346 on T-3344 is right, since the close reason names what T-3344 changed.
- **Touches:** five files, no folder, each named by a task. They match the parser and doc lines S-0540 changed for the timeout case.
- **Figures:** the 20m forecast rests on S-0540's 908 s for the same parser. The 15 USD a week value is `flai cod`'s from flai's 6m input. Both stand.
- **Your reading of criterion 1** is the one the story can meet: the flake kept no output, so a readable finding is the fix, and a recurrence becomes its own story.

**The follow-up proposal** (keep the failing tier's whole output in a log) is not drafted. Creating that story is the operator's call; I take no action on it.

### 2026-11-13T22:54:05Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: tasks, touches, and figures fit S-0782.
