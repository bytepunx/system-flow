---
id: T-1262
type: task
nature: feature
title: With the allowance in force, flai serve dials an unlisted container, accepts a dashboard without a stamp, warns once per connection, and records the unsigned side
status: backlog
parent: S-0239
owner: alex
created: 2026-10-07T23:09:24Z
updated: 2026-10-07T23:10:05Z
transitions: []
stream: S-0239
tags: [cli]
touches: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/internal/channel/channel.go, flai/internal/channel/channel_test.go]
after: [T-1259]
---
# T-1262 With the allowance in force, flai serve dials an unlisted container, accepts a dashboard without a stamp, warns once per connection, and records the unsigned side

## Work

flai's half of criterion 2 and the `warn` of criterion 3. It waits for T-1259, which gives the allowance in force per project. It runs beside T-1261 and T-1263, with no path in common.

- `flai/internal/channel/channel.go`: the client takes the allowance. With it, a dashboard answer whose stamp is missing, unverifiable, or names another component goes on to the credential proof instead of closing with `4403`; the proof is never skipped. `channel.State` gains the unsigned side or sides of the open connection: `dashboard` when its stamp failed, `flai` when this build has no valid stamp of its own. A signed pair leaves it empty.
- `flai/internal/serve/serve.go`: pass each project's allowance to its client. S-0238's gate before each dial lets an unlisted image through with the allowance and refuses it without, as before. Log one `warn` per connection, when it opens, naming the project and the unsigned side, and nothing for a signed pair. The state reaches `serve.Status.Connections`, which the statuses read.
- Tests: `channel_test.go` with the test dashboard for an unsigned dashboard accepted with the allowance and closed with `4403` without it, and a signed pair recording nothing; `serve_test.go` for an unlisted image dialled with the allowance and refused without it, and exactly one `warn` per connection.

## Done when

- flai dials an unlisted container and accepts a dashboard without a stamp only when the allowance is in force for the project.
- Each connection with an unsigned side logs one `warn` naming it, and its state records the side; a signed pair logs and records nothing.
- The tests pass with `flai test` on the changed paths.

## Notes

Drafted by the planner. S-0238 judges the image in `flai/internal/serve/imagecheck.go`; if the allowance is best applied there rather than where `serve.go` acts on the verdict, the story's agent widens the touches.
