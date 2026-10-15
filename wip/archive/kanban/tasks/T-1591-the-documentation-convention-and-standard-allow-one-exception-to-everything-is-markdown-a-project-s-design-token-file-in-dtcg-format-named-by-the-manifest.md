---
id: T-1591
type: task
nature: feature
title: "The documentation convention and standard allow one exception to Everything is markdown: a project's design token file, in DTCG format, named by the manifest"
status: done
parent: S-0381
owner: alex
created: 2026-10-14T23:53:00Z
updated: 2026-10-15T00:09:39Z
transitions:
  - to: ready
    at: 2026-10-14T23:57:53Z
    by: agent-S-0381
  - to: in-progress
    at: 2026-10-14T23:57:53Z
    by: agent-S-0381
  - to: done
    at: 2026-10-15T00:09:39Z
    by: agent-S-0381
stream: S-0381
tags: [template]
touches: [design/system/documentation-standard.md, design/conventions/documentation.md, template/root/design/conventions/documentation.md, template/CHANGELOG.md]
usage:
  source: log
  seconds: 706
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 17
      output: 8947
      cache_read: 840441
      cache_write: 40631
      cost: 0.5702
---
# T-1591 The documentation convention and standard allow one exception to Everything is markdown: a project's design token file, in DTCG format, named by the manifest

## Work

Carry ADR-0146 § Decision 5 into the standard. In `design/system/documentation-standard.md` first, then in the baseline rule "Everything is markdown" of `template/root/design/conventions/documentation.md` and its copy in `design/conventions/documentation.md`, add the one exception: a project's design token file, in DTCG format, at the path the manifest's `design.tokens` names, linking ADR-0146. Bump each file's `updated`. Add a patch entry for the template change to `template/CHANGELOG.md`, in the form its other entries take. The operator decided the exception on TH-0445, so this edits a baseline rule with the operator's consent, as `CLAUDE.md` asks. It waits for nothing: it changes no code, and shares no path with any other task.

## Done when

- The three documents state the same exception and link ADR-0146.
- The two copies of `documentation.md` agree, and `flai check --strict` and the markdown lint pass on them.

## Notes

- Layer 1, beside T-1532 and T-1590. Written by planner-S-0381.
