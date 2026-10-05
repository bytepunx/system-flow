---
id: T-0840
type: task
nature: improvement
title: flai check --story reports findings outside the story as notes that --strict passes over
status: backlog
parent: S-0249
owner: alex
created: 2026-10-05T00:17:42Z
updated: 2026-10-05T00:17:53Z
transitions: []
stream: S-0249
tags: [flai, check]
touches: [flai/internal/check, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md]
---
# T-0840 flai check --story reports findings outside the story as notes that --strict passes over

## Work

Give `flai check` a `--story S-nnnn` flag. With it, each finding is either inside the story or outside it.

A finding is inside the story when its path is one of these:

- the story's item file, or the file of one of its tasks
- its narrative, `wip/agents/S-nnnn.md`
- a thread anchored on the story or one of its tasks
- a path in the story's diff against the main branch (`storygit.StoryDiff`, as `flai stream sync` reads it)

A `wip.overlap` with another open story is outside the story even when it is reported on the story's own file. Clearing it is the pull hold's business (ADR-0019, ADR-0046) and the other story's agent's, not this story's (I-0059, S-0244).

A finding outside the story keeps its level, rule, path, and message. It is marked as a note, for instance with an `outside` field on `Finding` and a count on `Result`, and the run passes over it, as `--strict` passes over the advisory warnings (ADR-0073). This applies to errors as well as warnings. I-0057's instances include an error on main that no story could clear: `issues.duplicate-id`, from two I-0062 files left by an acceptance. The unscoped check on main still fails on every error.

Say in the summary line how many findings were outside the story. Explain the flag in the help text, and regenerate the reference with `make flai-reference`.

Without `--story`, `flai check` behaves exactly as before, so CI on the main branch still gates every finding.

It waits for no task.

## Done when

- A test in `flai/internal/check` reproduces I-0057. The fixture has a story in progress with its narrative, and an archived story S-0002 whose branch was never merged, which gives `story.unaccepted`. `Run` scoped to the first story returns OK under strict. Unscoped, it fails as before. The test fails without the change.
- A test shows that a warning on the story's own item file, or on a path in its diff, still fails `--strict` with `--story`. It also shows that a `wip.overlap` with another in-progress story, and an error on another story's file, do not fail it.
- A test in `flai/cmd/check_test.go` runs `flai check --strict --story` and checks the exit code and the summary line, with `--json` carrying the outside marker.
- `docs/users/flai-reference.md` matches the help text, and `scripts/flai-test.sh` passes.

## Notes
