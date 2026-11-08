---
id: T-2653
type: task
nature: feature
title: Under flai verify, smoke skips the template test, the check, and the markdown lint a step of the same run already passed
status: backlog
parent: S-0610
owner: alex
created: 2026-11-07T23:49:21Z
updated: 2026-11-07T23:49:21Z
transitions: []
stream: S-0610
tags: [cli]
touches: [flai/internal/verify/run.go, flai/internal/verify/run_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, scripts/smoke.sh, design/system/flai-cli.md, docs/users/flai.md]
after: [T-2644]
---
# T-2653 Under flai verify, smoke skips the template test, the check, and the markdown lint a step of the same run already passed

## Work

For criterion 2: tell each tier what passed before it in the same run, and have `scripts/smoke.sh` skip what passed. It waits for T-2644 because both change `run.go` and its tests, and because it reads which tiers finished before `smoke` starts, which T-2644's ordering decides.

- In `verify.Run`, set on each tier's environment a variable, such as `FLAI_VERIFY_PASSED`, naming the tiers of the same run that passed before it started, and whether each ran whole (`all_command`) or for paths. In the close-out (`storyRun.run`, `flai/internal/verify/story.go`), add the close-out's own `check` step when it passed: it runs `flai check --strict --story` as `scripts/check.sh` does under `CLOSE_OUT_STORY`. A reused tier counts as passed. With T-2634's key, `smoke` must wait for the tiers whose verdict it reads; mark it so in the tiers the story's last task writes, or here if the key's shape makes that automatic.
- In `scripts/smoke.sh`, skip `template-test.sh` when `template` passed, `check.sh` when the close-out's check passed, and `lint-md.sh` when `markdown` passed whole; otherwise run it as now. Print one line per skipped step naming what covered it. Run as `make smoke` or `flai test smoke`, with the variable unset, it runs every step, as now. Keep it POSIX `sh` under `set -eu` and ShellCheck clean.
- Say how the close-out runs its tiers now, at once where the manifest allows and with smoke's skips, in `design/system/flai-cli.md` § The close-out job and in the `flai verify` section of `docs/users/flai.md`.

## Done when

- A test of `Run` shows each tier gets the passed tiers before it, and the close-out's check step, in its environment.
- `scripts/smoke.sh` with the variable naming `template`, the check, and `markdown` run whole runs only the install test; with it unset it runs every step; ShellCheck passes.
- The design and user docs say what the close-out runs at once and what smoke skips.

## Notes
