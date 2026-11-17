---
id: S-0836
type: story
nature: remediation
title: flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry
status: done
owner: alex
created: 2026-11-17T08:17:31Z
updated: 2026-11-17T12:04:42Z
transitions:
  - to: ready
    at: 2026-11-17T10:05:44Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-17T11:15:21Z
    by: agent-S-0836
  - to: review
    at: 2026-11-17T12:03:13Z
    by: agent-S-0836
  - to: done
    at: 2026-11-17T12:04:42Z
    by: orchestrator
tags: [flai, mdlint]
touches: [flai/internal/mdlint/testdata/cases/code-spans.md, flai/internal/mdlint/testdata/cases/expected.txt, design/adrs/0318-flai-digest-new-and-edit-refuse-an-entry-whose-markdown-the-project-s-lint.md, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/digest.go, flai/internal/digest/write.go, flai/internal/digest/write_test.go, design/issues/I-0363-flai-s-markdown-lint-reports-md038-on-a-code-span-holding-escaped-backticks.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2886
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 3
      test_runs: 2
      work: 91
  models:
    - model: claude-opus-5-5
      input: 200
      output: 39703
      cache_read: 25868381
      cache_write: 323689
      cost: 8.558
  strategic:
    - kind: planner
      seconds: 2
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 2
          output: 3
          cache_read: 176625
          cache_write: 3207
          cost: 0.0123
        - model: claude-opus-5-5
          input: 10
          output: 1682
          cache_read: 711691
          cache_write: 19703
          cost: 0.3238
    - kind: orchestrator
      seconds: 93
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 25
          output: 231
          cache_read: 1371340
          cache_write: 114075
          cost: 0.4631
        - model: claude-sonnet-5-5
          input: 6
          output: 56
          cache_read: 48764
          cache_write: 36075
          cost: 0.0684
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-17T08:17:31Z
  value: 12.5
  by: planner-S-0836
  at: 2026-11-17T08:29:08Z
forecast:
  duration: 45m
  delivery: 2026-11-17T12:39:00Z
  basis: "Its own forecast of 45m; 3rd in the pull order with an in-progress limit of 3, behind S-0756, S-0806, S-0832, S-0811 and S-0809."
  by: flai
  at: 2026-11-17T11:04:27Z
finalized:
  by: orchestrator
  at: 2026-11-17T08:33:55Z
---
# S-0836 flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry

## Goal

This story remediates [I-0363](../../../design/issues/I-0363-flai-s-markdown-lint-reports-md038-on-a-code-span-holding-escaped-backticks.md), "flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0363 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0363 is closed with `flai issue close I-0363 --reason` saying what fixed it

## Tasks
- T-3529 mdlint and markdownlint-cli2 agree on D-0592's line 17 in both its wordings, pinned by the code span fixture
- T-3531 Close I-0363 with the cause and the fix T-3529 found
- T-3558 flai digest new and edit refuse an entry whose markdown the project's lint rejects

## Notes

Cost of delay inputs set by flai from I-0363. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T16:00:01Z, 3.7 days before this story; under one cycle counts as one).

### Outcome

The planned fix, in mdlint's parser, was not needed. markdownlint-cli2 0.20.0 reports the same MD038 as mdlint on D-0592's line 17 as the curator wrote it, and neither reports one on main's wording, which was hand-edited after the issue was filed. So mdlint's finding was true, and the issue's account of a false positive was wrong. The cause was that `flai digest new` and `edit` never linted an entry's markdown, so an entry markdownlint rejects could reach main and fail `TestRepositoryLintsClean`. Criterion 1 is met against that cause: digest writes now refuse such an entry ([ADR-0318](../../../design/adrs/0318-flai-digest-new-and-edit-refuse-an-entry-whose-markdown-the-project-s-lint.md)), with `TestWritesRefuseWhatTheMarkdownLintRejects` in `flai/internal/digest/write_test.go` using the original line, and the code span fixture pins the two tools' agreement on both wordings. The story's touches were narrowed from the planned mdlint files to what the branch changes.

### Planning

Proposed solution, from I-0363's one instance: make flai's markdown lint (`flai/internal/mdlint`) parse the code spans of line 17 of `design/digest/D-0592-flai-stream-sync-json-prints-a-map-that-flai-cmd-stream-go-builds-by-hand-not.md` as markdownlint-cli2 0.20.0 does, so it reports no MD038 there. The fix belongs in the lint, not in the digest entry: markdownlint passes the line, so mdlint's finding is a false positive, and rewriting the entry would leave the next such line to fail again. The issue's account of the cause, that mdlint ends the span at an escaped backtick, is to be checked first: CommonMark gives a backslash no escaping power inside a code span, while one before a backtick outside a span escapes it, so the divergence may lie in which backtick opens the next span, or in MD038's padding rule, rather than in the closing itself. The story's agent establishes it by running markdownlint-cli2 0.20.0 (`npx`, as `scripts/mdlint-fixtures.sh` does) and mdlint on the line, and records the cause in the narrative's Decisions before changing the parser.

Touches, none declared before this plan:

- `flai/internal/mdlint/inline.go`: layout. It holds `parseRange` and `codeClose`, which open and close the code spans MD038 reads.
- `flai/internal/mdlint/rules.go`: layout. It holds `md038`, should the divergence be in its padding rule rather than in the parse.
- `flai/internal/mdlint/mdlint_test.go`: co-change, 4 of the 5 commits that changed `inline.go` or the code span fixture (80%). Its MD038 table (I-0072) takes the reproducing case.
- `flai/internal/mdlint/testdata/cases/code-spans.md`: layout. The fixture that `TestFixturesMatchMarkdownlint` lints against markdownlint's own findings; the D-0592 line is added here.
- `flai/internal/mdlint/testdata/cases/expected.txt`: co-change, 4 of 5 (80%). Regenerated with `scripts/mdlint-fixtures.sh` after the fixture changes.

Left out from `flai touches suggest`: `doc.go` and `mdlint.go` (2 of 5, 40%): the block parse and the rule table need no change for an inline code span's parse. The task widens to them if the cause lies there. The issue file under `design/issues/` is flai's to write with `flai issue close`, and the touches drift leaves that folder out (ADR-0284).

No folder touches. No `after`: the story reads no other open story's output.

Forecast: 45m, adjusted from flai forecast's 18m (147 s per unit of size over 68 done medium remediation stories, size 7). The cause must be found against markdownlint-cli2 before the fix, the fixtures regenerated, and the close-out runs the integration tier, `TestRepositoryLintsClean` among it, since the story changes mdlint's code.

Cost of delay value: 12.50 USD a week, as flai cod works it out from the operator's input of 5m lost per 168h cycle at 150 USD an hour. It stands: the one occurrence fails only the integration and race tiers outside a close-out, which is what the input already prices.

### Accepted by the orchestrator

- Verified: 4d3eef401daab04573f90f029bec8e24f8e641b8
- At: 2026-11-17T12:04:42Z

Verdict: meets both criteria (1 and 2), with no convention breach, at commit 4d3eef401daab04573f90f029bec8e24f8e641b8.
- 1: flai/internal/digest/write.go, flai/internal/digest/write_test.go, flai/internal/mdlint/testdata/cases/code-spans.md, flai/internal/mdlint/testdata/cases/expected.txt, flai/cmd/digest.go, design/adrs/0318-flai-digest-new-and-edit-refuse-an-entry-whose-markdown-the-project-s-lint.md
- 2: design/issues/I-0363-flai-s-markdown-lint-reports-md038-on-a-code-span-holding-escaped-backticks.md, design/issues/summary.md
