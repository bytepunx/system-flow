---
id: T-3460
type: task
nature: feature
title: flai verify and the close-out job run the suites the transitions out of the story's state need, and --all every suite
status: done
parent: S-0807
owner: alex
created: 2026-11-16T23:49:44Z
updated: 2026-11-18T10:34:37Z
transitions:
  - to: ready
    at: 2026-11-18T09:55:59Z
    by: agent-S-0807
  - to: in-progress
    at: 2026-11-18T09:55:59Z
    by: agent-S-0807
  - to: done
    at: 2026-11-18T10:34:37Z
    by: agent-S-0807
stream: S-0807
tags: [flai, template]
touches: [flai/internal/verify/suites.go, flai/internal/verify/suites_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/mcpserver/verify.go, scripts/close-out.sh, template/root/scripts/close-out.sh, docs/users/flai-reference.md, flai/internal/verify/record.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go]
after: [T-3454]
usage:
  source: log
  seconds: 2318
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 205
      output: 84212
      cache_read: 16642661
      cache_write: 766858
      cost: 9.1974
---
# T-3460 flai verify and the close-out job run the suites the transitions out of the story's state need, and --all every suite

## Work

ADR-0312 §2's last sentence and criterion 2. It waits for T-3454 because the job's answer from the record must match on the suites the record holds.

- Add `flai/internal/verify/suites.go`, which gives the suites a story needs: the union of the `verify:` suites of the workflow's transitions out of the story's state. `all` on any of them means every suite. With none, it gives every suite, so a workflow without `verify:` runs what it runs today.
- Filter `Select` and `selectStoryTiers` in `select.go` by those suites. `--all` runs every suite's tiers, as today. A tier's `queue` and `when` keep their meaning until S-0813.
- `flai verify <S-nnnn>` (`flai/cmd/verify.go`) uses the suites from the story's state. Its human and `--json` output name the suites it ran.
- The close-out job (`flai/internal/serve/verifyjob.go`, `verifyjob_container.go`) runs the same suites. `fromRecord` answers from the record only when it holds a passing verdict for every suite the job would run.
- Make the MCP tool `verify`'s description and the close-out scripts' usage text say what is run. Keep `scripts/close-out.sh` and `template/root/scripts/close-out.sh` byte for byte equal.

## Done when

- `flai/internal/verify` tests show the suites worked out from a workflow whose out-transitions list different suites.
- They show `all`, a workflow with no `verify:`, and `--all`.
- `flai/cmd/verify_test.go` and `flai/internal/serve/verifyjob_test.go` show the job running only the needed suites. They also show it answering from the record only when every needed suite passed there.
- `flai test` passes on the changed paths.

## Notes
