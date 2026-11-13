---
id: T-3275
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0764
owner: alex
created: 2026-11-13T08:55:19Z
updated: 2026-11-13T09:01:16Z
transitions:
  - to: ready
    at: 2026-11-13T08:57:32Z
    by: agent-S-0764
  - to: in-progress
    at: 2026-11-13T08:57:32Z
    by: agent-S-0764
  - to: done
    at: 2026-11-13T09:01:16Z
    by: agent-S-0764
blocked:
  - from: 2026-11-13T08:57:44Z
    until: 2026-11-13T09:00:20Z
    reason: "publish step: flai release --pending --require-pending --json exited 1 (git add -A could not take .git/index.lock in the main checkout: another git process held it)"
stream: S-0764
tags: []
after: [T-3274]
usage:
  source: log
  seconds: 130
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 38
      output: 4863
      cache_read: 2836287
      cache_write: 77437
      cost: 1.2842
---
# T-3275 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0764 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-13T08:57:33Z: ran in the main checkout and exited 1, an exit the step's stop does not list. The whole output:

```text
{"ts":"2026-11-13T08:57:33.822Z","level":"FATAL","msg":"command failed","component":"cmd","err":"git add -A: exit status 128\nfatal: Unable to create '/home/alex/git/bytepunx/system-flow/.git/index.lock': File exists.\n\nAnother git process seems to be running in this repository, e.g.\nan editor opened by 'git commit'. Please make sure all processes\nare terminated then try again. If it still fails, a git process\nmay have crashed in this repository earlier:\nremove the file manually to continue."}
```

- Seconds later `.git/index.lock` no longer existed, so another git process held it for a moment (the orchestrator and the curator were both running then). HEAD was still `631c7c498`, no tag pointed at it, and nothing was committed: the release stopped before its commit. Blocked and asked on TH-1142; recorded as I-0359.
- TH-1142 answered "run it again". Unblocked and ran again in the main checkout: exit 0. It released flai 1.86.23 → 1.87.0 (minor, for S-0660), tagged `flai/v1.87.0`, and pushed (`"pushed": true`).
