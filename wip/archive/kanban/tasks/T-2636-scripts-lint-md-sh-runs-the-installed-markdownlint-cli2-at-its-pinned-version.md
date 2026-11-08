---
id: T-2636
type: task
nature: improvement
title: scripts/lint-md.sh runs the installed markdownlint-cli2 at its pinned version, and install-tools.sh installs it
status: done
parent: S-0609
owner: alex
created: 2026-11-07T23:48:23Z
updated: 2026-11-08T02:26:31Z
transitions:
  - to: ready
    at: 2026-11-08T02:17:12Z
    by: agent-S-0609
  - to: in-progress
    at: 2026-11-08T02:17:13Z
    by: agent-S-0609
  - to: done
    at: 2026-11-08T02:26:31Z
    by: agent-S-0609
stream: S-0609
tags: [scripts]
touches: [scripts/lint-md.sh, scripts/install-tools.sh, design/tech/ci.md, design/tech/node.md, docs/contributors/index.md]
usage:
  source: log
  seconds: 558
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 10648
      cache_read: 4483940
      cache_write: 67523
      cost: 1.5441
---
# T-2636 scripts/lint-md.sh runs the installed markdownlint-cli2 at its pinned version, and install-tools.sh installs it

## Work

- In `scripts/lint-md.sh`, run `markdownlint-cli2` straight from `PATH` (env.sh puts the checkout's `bin/` first) when it answers the pinned `MDLINT_VERSION`; otherwise fall back to `npx --yes markdownlint-cli2@$MDLINT_VERSION` as today. Each `npx` call costs 1.0 to 1.2 s, of which markdownlint is 0.27 to 0.36 s (cause 19 in `design/system/server-performance.md`). Keep every glob, the close-out and merge-queue branches, and the exit codes as they are.
- Say on the first line which it ran, installed or through npx, as the script already names the version.
- In `scripts/install-tools.sh`, install `markdownlint-cli2` at the same pinned version into the main checkout's tools (an npm prefix under `.flai-cache`, linked into `bin/`), skipping it when that version is there already, as it does for ShellCheck.
- Record it in `design/tech/ci.md` and `design/tech/node.md`, and in what `docs/contributors/index.md` says `install-tools.sh` installs.
- Waits for no task: the Go tasks call the script through the manifest's `markdown` tier and do not depend on how it starts markdownlint.

## Done when

- With `markdownlint-cli2` installed by `scripts/install-tools.sh`, `scripts/lint-md.sh --fix <file>` and `scripts/lint-md.sh <file>` run it without `npx`, and give the same findings and exit code as the `npx` run.
- Without it installed, or at another version, the script runs through `npx` as before.
- `flai test scripts/lint-md.sh scripts/install-tools.sh` passes ShellCheck and the markdown lint.

## Notes
