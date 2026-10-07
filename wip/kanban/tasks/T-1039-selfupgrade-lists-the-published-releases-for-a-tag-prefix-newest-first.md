---
id: T-1039
type: task
nature: improvement
title: selfupgrade lists the published releases for a tag prefix, newest first
status: done
parent: S-0298
owner: alex
created: 2026-10-06T21:44:46Z
updated: 2026-10-07T14:32:33Z
transitions:
  - to: ready
    at: 2026-10-07T14:26:52Z
    by: agent-S-0298
  - to: in-progress
    at: 2026-10-07T14:26:52Z
    by: agent-S-0298
  - to: done
    at: 2026-10-07T14:32:33Z
    by: agent-S-0298
stream: S-0298
tags: [cli]
touches: [flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go, design/system/flai-cli.md, design/system/flaiover-dashboard.md]
usage:
  source: log
  seconds: 340
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 55
      output: 292
      cache_read: 1960000
      cache_write: 80170
      cost: 0.9243
---
# T-1039 selfupgrade lists the published releases for a tag prefix, newest first

## Work

`selfupgrade.Resolve` finds the newest `flai/v` release or one pinned version, reading only the first 50 releases, and nothing lists them. Add a `List` that returns every published release (no drafts, no prereleases) whose tag has a given prefix, `flai/v` or `flaiover/v`, newest first by semver, with its version, tag, and publish date, following the GitHub API's pages past the first. Make `Resolve` share it, so a pinned version that is not published is refused with the versions that are. Test it against a stand-in API (`FLAI_RELEASES_API`), as the existing tests do. Waits for nothing: the names the ADR task fixes do not reach this package.

## Done when

- `List` returns the flai and the flaiover releases newest first across more than one page, leaving out drafts, prereleases, and other prefixes.
- A pinned version that is not published is refused with an error naming the published ones.
- `scripts/flai-test.sh` passes for the package.

## Notes
