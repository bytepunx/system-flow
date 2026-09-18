---
id: T-0155
type: task
nature: remediation
title: flai check warns when the clone needs a newer git than the one on PATH
status: done
parent: S-0050
owner: alex
created: 2026-09-18T19:47:24Z
updated: 2026-09-18T19:55:19Z
transitions:
  - to: ready
    at: 2026-09-18T19:53:39Z
    by: alex
  - to: in-progress
    at: 2026-09-18T19:53:39Z
    by: alex
  - to: done
    at: 2026-09-18T19:55:19Z
    by: alex
stream: S-0050
tags: []
touches: [flai/internal/check]
---

# T-0155 flai check warns when the clone needs a newer git than the one on PATH

## Work
In `flai/internal/check`, add a warning `git.relative-worktrees` raised when the main repository's `.git/config` has `extensions.relativeWorktrees` true and the git on PATH is older than 2.48. Read the config file directly, since that git refuses the repository. The message names both versions and the ways out: upgrade git, or with a newer git run `git worktree repair` without relative paths and unset the extension. Try the recovery in a scratch repository before writing it into the message. Fixture or temp-dir tests for: extension absent, extension set with new git, extension set with old git.

## Done when
- The three cases are tested
- The recovery steps in the message were run in a scratch repository and worked, with the commands recorded in the narrative

## Notes
