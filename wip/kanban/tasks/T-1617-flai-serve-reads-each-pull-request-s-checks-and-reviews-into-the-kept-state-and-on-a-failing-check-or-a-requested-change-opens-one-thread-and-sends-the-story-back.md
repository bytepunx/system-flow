---
id: T-1617
type: task
nature: feature
title: flai serve reads each pull request's checks and reviews into the kept state, and on a failing check or a requested change opens one thread and sends the story back
status: backlog
parent: S-0394
owner: alex
created: 2026-10-16T04:36:47Z
updated: 2026-10-16T04:36:47Z
transitions: []
stream: S-0394
tags: [cli]
touches: [flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go, flai/internal/serve/serve.go]
after: [T-1615]
---
# T-1617 flai serve reads each pull request's checks and reviews into the kept state, and on a failing check or a requested change opens one thread and sends the story back

## Work

- In `flai/internal/serve/forgestatus.go`, on the minute loop wired in `serve.go`, call the forge's `Status` for each story in review that has a pull request, and write T-1615's state.
- On a required check that failed, or a review that requests changes, at the request's head commit, do the following once per failure and commit:
  - Open a thread on the story with `threads.New`. For a check, it names the check and links its log. For a review, it quotes each comment, briefly, and links it.
  - Move the story to in-progress, by `flai`, with the reason naming the same.
- Leave a check that is pending alone: the story waits in review.
- Test it against the fake forge: a check passes, is pending, and fails; changes requested; the same failure seen twice on one commit; a new commit with the same check failing again.
- Waits for T-1615: it writes the kept state.

## Done when

In a scratch project with the fake forge, a failing check and a requested change each open one thread and send the story back once per commit, a pending check does neither, and the tests pass.

## Notes
