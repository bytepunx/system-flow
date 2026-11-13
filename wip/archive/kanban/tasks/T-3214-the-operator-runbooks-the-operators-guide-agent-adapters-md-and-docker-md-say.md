---
id: T-3214
type: task
nature: improvement
title: The operator runbooks, the operators guide, agent-adapters.md, and docker.md say how any project gets its agent images after E-0036, and what is this repository's own
status: done
parent: S-0760
owner: alex
created: 2026-11-12T16:07:00Z
updated: 2026-11-13T20:06:59Z
transitions:
  - to: ready
    at: 2026-11-13T20:00:59Z
    by: agent-S-0760
  - to: in-progress
    at: 2026-11-13T20:00:59Z
    by: agent-S-0760
  - to: done
    at: 2026-11-13T20:06:59Z
    by: agent-S-0760
stream: S-0760
tags: [cli]
touches: [docs/operators/runbooks/agent-container.md, docs/operators/runbooks/install.md, docs/operators/index.md, docs/users/flai.md, design/system/agent-adapters.md, design/tech/docker.md, design/system/flai-cli.md, design/system/project-manifest.md, flai/images/system-flow/Dockerfile]
after: [T-3207, T-3209]
usage:
  source: log
  seconds: 360
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 73
      output: 27087
      cache_read: 5334796
      cache_write: 135339
      cost: 2.4936
---
# T-3214 The operator runbooks, the operators guide, agent-adapters.md, and docker.md say how any project gets its agent images after E-0036, and what is this repository's own

## Work

Say what is true after E-0036, and bump each file's `updated`:

- `docs/operators/runbooks/agent-container.md`: the images come from flai (§ The images, § Build). `make agent-images` and the build command are this repository's, under a heading of their own. § Settings (line 139) names the build command and the switch. § Clean up (line 166) says flai brings back what is removed.
- `docs/operators/runbooks/install.md` § Agents in containers (line 85 on): `flai image ensure` in place of `make agent-images` and the manual layer build (lines 93 to 108).
- `docs/operators/index.md` § The image (line 1164): the label and where the images come from.
- `docs/users/flai.md` § In a container (line 501): `make agent-images` as this repository's alone.
- `design/system/agent-adapters.md` § A container harness of flai's own (line 466): where the images come from, as built.
- `design/tech/docker.md` § Agent images: the published images, the host build, the layer's arguments, the job, and this repository's build command.

It waits for T-3207 and T-3209, whose manifest and script it describes.

## Done when

- None of the six files says that only this repository can build the images, or tells a project other than this one to run `make agent-images`.
- `flai check --strict` and the markdown lint pass.

## Notes
