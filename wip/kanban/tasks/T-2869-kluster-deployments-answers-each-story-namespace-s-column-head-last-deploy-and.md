---
id: T-2869
type: task
nature: feature
title: kluster.deployments answers each story namespace's column, head, last deploy and checks run, and who asked, the stories in progress not deployed, and the count against the limit
status: backlog
parent: S-0663
owner: alex
created: 2026-11-10T08:24:53Z
updated: 2026-11-10T08:24:53Z
transitions: []
stream: S-0663
tags: [cli, dashboard]
touches: [flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go]
after: [T-2868]
---
# T-2869 kluster.deployments answers each story namespace's column, head, last deploy and checks run, and who asked, the stories in progress not deployed, and the count against the limit

## Work

Give the dashboard what T-2868 put in the list, through the host read `kluster.deployments` in `flai/internal/hostapi/kluster.go` (`klusterDeployments`, `reviewStories`, `reviewRows`).

- Each row of `deployments` carries the story's column, whether its commit is the head, its last deploy and its last checks run with outcome, log, `log_tail`, commit, and who asked, as T-2868's `Deployments` gives them, with the project's root left out of the paths as today.
- The stories with no record widen from those in review to those in progress too, each with its column, so the page can offer **Deploy** for one in progress. Keep `review` as it is for an older dashboard, and add the stories in progress beside it.
- The answer gains the count of story namespaces that live, the limit, and the requests refused for it.
- `kluster.status`'s and `kluster.deployments`' other fields stay as they are.

It waits for T-2868, whose `Deployments` fields it reads.

## Done when

- Tests in `flai/internal/hostapi/kluster_test.go` show a story in progress's namespace with its column, head, last deploy and checks run, and who asked; a story in progress with no record; and the count against the limit with a refused request.
- `flai test` passes on the changed paths, the `flaiover-flai` tier among them.

## Notes
