---
id: T-2180
type: task
nature: improvement
title: "scripts/shellcheck.sh lints the scripts given, or every scripts/*.sh, and the existing scripts pass it"
status: done
parent: S-0503
owner: alex
created: 2026-10-26T05:22:59Z
updated: 2026-10-26T05:34:38Z
transitions:
  - to: ready
    at: 2026-10-26T05:31:27Z
    by: agent-S-0503
  - to: in-progress
    at: 2026-10-26T05:31:28Z
    by: agent-S-0503
  - to: done
    at: 2026-10-26T05:34:37Z
    by: agent-S-0503
stream: S-0503
tags: [scripts, tooling]
touches: [scripts/shellcheck.sh, scripts/flai-test.sh, scripts/README.md, Makefile, scripts/flaiover-install-test.sh, scripts/install-published-test.sh, scripts/litellm.sh, scripts/release-key-setup.sh, scripts/env.sh]
after: [T-2179]
usage:
  source: log
  seconds: 189
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 35
      output: 10965
      cache_read: 1997854
      cache_write: 59769
      cost: 1.0131
---
# T-2180 scripts/shellcheck.sh lints the scripts given, or every scripts/*.sh, and the existing scripts pass it

## Work

- Add `scripts/shellcheck.sh`: POSIX `sh` under `set -eu`, sourcing `env.sh` so `bin/shellcheck` is on PATH. It runs `shellcheck -x -P SCRIPTDIR` over the paths given, or over every `scripts/*.sh` when none is given, and exits as shellcheck exits. `-x -P SCRIPTDIR` lets shellcheck follow `. "$(dirname "$0")/env.sh"`; where it still reports SC1090 or SC1091, put a `# shellcheck source=env.sh` directive on the line rather than excluding the rule. When shellcheck is missing, say so on stderr, name `scripts/install-tools.sh`, and exit non-zero.
- Fix the seven findings ShellCheck 0.10.0 reports today, in `scripts/flaiover-install-test.sh`, `scripts/install-published-test.sh`, `scripts/litellm.sh`, and `scripts/release-key-setup.sh` (SC2015, SC2041, SC2016, SC2001). Where the code is right as it stands, disable the rule on that line by directive with the reason beside it, never file-wide.
- Add a `shellcheck` target to the `Makefile` that calls the script, a row for the script in `scripts/README.md`, and run it from `scripts/flai-test.sh` ahead of gofmt, which keeps that script in step with the tiers as its header asks.

Waits for T-2179: it runs the pinned ShellCheck that task installs into `bin/`.

## Done when

- `scripts/shellcheck.sh` with no arguments exits 0 on this tree, and exits non-zero when given a scratch script with a finding, naming the file and line.
- `make shellcheck` runs it, and `scripts/flai-test.sh` runs it before gofmt.
- `flai check --strict` and the markdown lint are clean on the files changed.

## Notes

Drafted by planner-S-0503. The template's `template/root/scripts/*.sh` are left out of the default set: they are rendered into other projects and the `template` tier checks them; widening to them is a choice for the story's agent to raise on the plan thread if it looks worth it.
