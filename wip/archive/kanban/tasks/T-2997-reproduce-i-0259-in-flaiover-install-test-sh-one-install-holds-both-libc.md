---
id: T-2997
type: task
nature: remediation
title: "Reproduce I-0259 in flaiover-install-test.sh: one install holds both libc variants, and a newer configuration installs again"
status: done
parent: S-0706
owner: alex
created: 2026-11-10T10:30:52Z
updated: 2026-11-11T09:23:50Z
transitions:
  - to: ready
    at: 2026-11-11T09:20:58Z
    by: agent-S-0706
  - to: in-progress
    at: 2026-11-11T09:20:58Z
    by: agent-S-0706
  - to: done
    at: 2026-11-11T09:23:50Z
    by: agent-S-0706
stream: S-0706
tags: [scripts, flaiover]
touches: [scripts/flaiover-install-test.sh, scripts/README.md]
after: [T-2993, T-2994]
usage:
  source: log
  seconds: 172
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 7327
      cache_read: 2955644
      cache_write: 67465
      cost: 1.2775
---
# T-2997 Reproduce I-0259 in flaiover-install-test.sh: one install holds both libc variants, and a newer configuration installs again

## Work

It waits for T-2993, whose `supportedArchitectures` setting the test copies, and for T-2994, whose `--if-needed` rule it tests and whose `scripts/README.md` change it builds on.

`scripts/flaiover-install-test.sh` already reproduces I-0177 offline with local-folder packages. Add I-0259's case the same way, with no network. Use two local optional packages, one with `"libc": ["glibc"]` and one with `"libc": ["musl"]`, depended on by a fixture project. First install without `supportedArchitectures` and show that only the host's libc variant lands: that is the issue. Then install with the setting T-2993 chose, in the file it chose, and assert that both land. Also assert that `flaiover-install.sh --if-needed` installs again once the fixture's configuration file is newer than `node_modules/.modules.yaml`, and skips while it is not.

If pnpm cannot be made to resolve libc-specific optional packages offline, keep the `--if-needed` case. Then also assert that flaiover's own configuration declares both libcs. Say in this task's Notes why the fixture case could not run.

Add I-0259's case to the test's row in `scripts/README.md`.

## Done when

- `scripts/flaiover-install-test.sh` fails on the code before T-2993 and T-2994 and passes after them, offline.
- `scripts/README.md`'s row for the test names I-0259's case.
- `flai test scripts/flaiover-install-test.sh` passes, and so does `scripts/flaiover-test.sh`, which runs the test.

## Notes

The fixture case runs offline, with two local folder packages, one per libc. pnpm 10.34.5 honours `libc` for them only when it resolves: a `--frozen-lockfile` install links both folders, whatever the setting, and lists the other libc's as skipped in `.modules.yaml`. Tarball packages are not filtered at all. So the test resolves afresh (`--no-frozen-lockfile`) to show the setting's effect, and checks `--if-needed`'s rule apart, by whether it runs pnpm. It fails on main's `flaiover-install.sh` and on main's `pnpm-workspace.yaml`, each alone, and passes with both of this story's. The real lockfile's registry bindings are filtered by a frozen install too: that is I-0259's instance, and a fresh `--if-needed` install in the worktree put in both `@rolldown/binding-linux-x64-gnu` and `-musl`.
