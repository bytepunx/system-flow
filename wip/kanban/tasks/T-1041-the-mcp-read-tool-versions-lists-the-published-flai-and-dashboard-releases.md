---
id: T-1041
type: task
nature: improvement
title: The MCP read tool versions lists the published flai and dashboard releases
status: done
parent: S-0298
owner: alex
created: 2026-10-06T21:44:59Z
updated: 2026-10-07T14:51:41Z
transitions:
  - to: ready
    at: 2026-10-07T14:43:24Z
    by: agent-S-0298
  - to: in-progress
    at: 2026-10-07T14:43:24Z
    by: agent-S-0298
  - to: done
    at: 2026-10-07T14:51:27Z
    by: agent-S-0298
stream: S-0298
tags: [cli]
touches: [flai/internal/mcpserver/versions.go, flai/internal/mcpserver/versions_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/server_test.go, flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go, flai/cmd/selfupgrade.go, flai/cmd/dashboard_versions.go]
after: [T-1038, T-1039, T-1040]
usage:
  source: log
  seconds: 483
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 58
      output: 374
      cache_read: 2971654
      cache_write: 109959
      cost: 1.3956
---
# T-1041 The MCP read tool versions lists the published flai and dashboard releases

## Work

Add the MCP tool `versions` in `flai/internal/mcpserver/versions.go`, registered with the folder's tools in `folder.go`: it answers the published flai releases and the published dashboard releases, newest first, through T-1039's `List`, marking the flai this server runs, the newest of each, and any flai below the project's `flai.minimum`. It deploys nothing: as the operator answered on TH-0202, MCP lists only, and deploying stays with the operator's CLI and the dashboard's host actions. It is not added to `guard.MCPReads`: a story's sub-agents have no need of it. Waits for T-1038 for the tool's name and scope, and for T-1039 for `List`.

## Done when

- `versions` answers both lists, tested against a stand-in releases API.
- No MCP tool installs or deploys a release.
- `scripts/flai-test.sh` passes for `flai/internal/mcpserver`.

## Notes
