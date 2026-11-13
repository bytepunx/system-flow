---
id: T-3212
type: task
nature: improvement
title: The template's operators guide and manifest tell a new project that flai supplies its agent images, and how it names a layer and its arguments
status: done
parent: S-0760
owner: alex
created: 2026-11-12T16:06:49Z
updated: 2026-11-13T20:00:45Z
transitions:
  - to: ready
    at: 2026-11-13T19:48:13Z
    by: agent-S-0760
  - to: in-progress
    at: 2026-11-13T19:48:13Z
    by: agent-S-0760
  - to: done
    at: 2026-11-13T20:00:45Z
    by: agent-S-0760
stream: S-0760
tags: [template]
touches: [template/root/docs/operators/index.md.tmpl, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md]
usage:
  source: log
  seconds: 752
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 37
      output: 13677
      cache_read: 2693664
      cache_write: 68336
      cost: 1.2591
---
# T-3212 The template's operators guide and manifest tell a new project that flai supplies its agent images, and how it names a layer and its arguments

## Work

- `template/root/docs/operators/index.md.tmpl`: a section on agents in containers. Choose the `container` harness, and flai supplies the images at the first start. `flai image status` and `flai image ensure` show and bring them back, and the dashboard's Agent images page shows the job. A layer is named in the manifest with its arguments, and an operator may name their own build command or turn the automatic job off.
- `template/root/system-flow.yaml.tmpl` (lines 46 to 47): the commented `image:` block shows a layer with its arguments, in S-0745's shape, and says flai builds it.
- `template/CHANGELOG.md`: a line under `## Unreleased`.

It waits for no task: the shapes it describes are on main once S-0745 and S-0759 are accepted.

## Done when

- `scripts/template-test.sh`, `flai check --strict`, and the markdown lint pass.
- A project made from the template with `flai new` in a scratch folder renders the section and the block.

## Notes
