---
id: T-2164
type: task
nature: feature
title: agent.image.layer in the manifest names a project's Dockerfile, read by flai and documented in project-manifest.md, settings.md, and the template
status: done
parent: S-0510
owner: alex
created: 2026-10-26T04:49:39Z
updated: 2026-11-02T18:32:49Z
transitions:
  - to: ready
    at: 2026-11-02T18:18:28Z
    by: agent-S-0510
  - to: in-progress
    at: 2026-11-02T18:18:29Z
    by: agent-S-0510
  - to: done
    at: 2026-11-02T18:32:49Z
    by: agent-S-0510
stream: S-0510
tags: [cli, docker]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, design/system/project-manifest.md, docs/operators/settings.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, flai/internal/manifest/agent.go, flai/internal/manifest/agent_test.go, flai/internal/workitem/store_test.go, flai/internal/workitem/front-matter-fields.txt]
usage:
  source: log
  seconds: 860
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 56
      output: 14763
      cache_read: 6807925
      cache_write: 124340
      cost: 2.555
---
# T-2164 agent.image.layer in the manifest names a project's Dockerfile, read by flai and documented in project-manifest.md, settings.md, and the template

## Work

Add `agent.image.layer` to the manifest schema in `flai/internal/manifest/manifest.go`: an optional path, relative to the repository root, of the Dockerfile that builds the project's layer over flai's harness image; refuse an absolute path or one outside the repository with a message. Document it in `design/system/project-manifest.md` beside the other `agent` keys, give it a row in `docs/operators/settings.md` so `settings_doc_test` passes, and add it as a comment to `template/root/system-flow.yaml.tmpl` with a line in `template/CHANGELOG.md`. Test the read and the refusals in `manifest_test.go`. This task waits for nothing.

## Done when

A manifest naming `agent.image.layer` is read and one with an absolute path is refused in `manifest_test.go`, the three documents name the key, and `flai test flai/internal/manifest flai/cmd` passes.

## Notes

The image tag the layer builds to, `flai-agent-<project key>-<inner>:local`, is derived, not a setting.
