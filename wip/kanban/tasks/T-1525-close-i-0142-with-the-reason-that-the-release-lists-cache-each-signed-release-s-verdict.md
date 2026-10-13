---
id: T-1525
type: task
nature: remediation
title: Close I-0142 with the reason that the release lists cache each signed release's verdict
status: backlog
parent: S-0379
owner: alex
created: 2026-10-13T17:02:35Z
updated: 2026-10-13T17:02:35Z
transitions: []
stream: S-0379
tags: [flai, issues]
touches: [design/issues/I-0142-flai-self-upgrade-list-and-flai-host-versions-download-two-assets-for-every-signed-release-to-mark-the-ones-they-cannot-verify.md, design/issues/summary.md]
after: [T-1523, T-1524]
---
# T-1525 Close I-0142 with the reason that the release lists cache each signed release's verdict

## Work

Run `flai issue close I-0142 --reason "<reason>"` in the story's worktree. The reason names the verdict cache, its key, and the test that shows a second listing downloads no signed release's assets. It waits for T-1523 and T-1524 so that the reason names the fix as built and documented.

## Done when

- [ ] I-0142 is closed, with a reason that names the fix and its test.
- [ ] `design/issues/summary.md` no longer lists I-0142.
- [ ] `flai check --strict` is clean for the story.

## Notes

Run it from the worktree, not the main checkout: `flai issue` writes under the checkout it runs in.
