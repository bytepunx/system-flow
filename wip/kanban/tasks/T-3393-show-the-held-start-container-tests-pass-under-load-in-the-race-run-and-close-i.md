---
id: T-3393
type: task
nature: remediation
title: Show the held-start container tests pass under load in the race run, and close I-0369
status: backlog
parent: S-0795
owner: alex
created: 2026-11-16T00:41:11Z
updated: 2026-11-16T00:41:11Z
transitions: []
stream: S-0795
tags: [flai, serve, test]
touches: [flai/internal/serve/container_test.go, design/issues/I-0369-testcontainerstrategicstartunderalayerwaitsfortheharnessimagealone-fails-under.md]
after: [T-3392]
---
# T-3393 Show the held-start container tests pass under load in the race run, and close I-0369

## Work

Waits for T-3392: it checks that fix and closes the issue on it.

Run the image-held container tests in `flai/internal/serve/container_test.go` under `-race` and under load: `go test -race -count=50 -run 'TestContainer.*(Held|Image|Layer)' ./internal/serve/` with the package's other tests running beside it, or with `-cpu` raised. Then run `flai test race` once. I-0369's two instances failed only under load, so a quiet run proves nothing: record the commands, counts, and results in this task's Notes. If the window still shows, or the check at line 1616 should now read the state after the planner ends, change the test only so that it checks the same thing without racing, and say why in the narrative.

Then close the issue: `flai issue close I-0369 --reason "<what fixed it>"`, naming T-3392's change and the runs that show it.

## Done when

- The image-held container tests pass in every run of the loaded `-race` loop, and `flai test race` passes, each recorded in Notes.
- I-0369 is closed with a reason that names the fix.
- Both of S-0795's criteria can be ticked.

## Notes
