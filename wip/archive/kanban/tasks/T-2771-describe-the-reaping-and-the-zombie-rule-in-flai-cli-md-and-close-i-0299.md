---
id: T-2771
type: task
nature: remediation
title: Describe the reaping and the zombie rule in flai-cli.md and close I-0299
status: done
parent: S-0632
owner: alex
created: 2026-11-09T21:57:43Z
updated: 2026-11-09T22:13:52Z
transitions:
  - to: ready
    at: 2026-11-09T22:13:01Z
    by: agent-S-0632
  - to: in-progress
    at: 2026-11-09T22:13:02Z
    by: agent-S-0632
  - to: done
    at: 2026-11-09T22:13:52Z
    by: agent-S-0632
stream: S-0632
tags: [flai, serve, docs]
touches: [design/system/flai-cli.md, design/issues/I-0299-a-role-run-the-stdio-mcp-server-starts-through-role-run-is-never-reaped-so-flai.md, design/issues/summary.md]
after: [T-2769, T-2770]
usage:
  source: log
  seconds: 50
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 745
      cache_read: 623442
      cache_write: 5288
      cost: 0.1819
---
# T-2771 Describe the reaping and the zombie rule in flai-cli.md and close I-0299

## Work

In `design/system/flai-cli.md`, the `flai serve agent` row says that `StartNow` hands the process over (`launcher.handOver`), that the command does not wait for it, and what `running()` tests (`Owns(pid, start)`). Say there that a handed-over process is now reaped in the background, so a long-lived parent such as `flai mcp` leaves no zombie. Say too that on Linux a process whose `/proc/<pid>/stat` state is `Z` or `X` is not running, so `settleOrphans` settles its run. Name S-0632 and I-0299. Bump the front matter's `updated`.

Then close the issue with `flai issue close I-0299 --reason "<what fixed it>"`, naming both changes. That command writes the issue file and `design/issues/summary.md`.

Waits for T-2769 and T-2770: it describes what they built, and the reason names it.

## Done when

- `design/system/flai-cli.md` describes the reaping and the zombie rule, and its `updated` is today's date.
- I-0299 is closed with a reason naming both fixes, and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` and the markdown lint pass on the touched files.

## Notes
