---
id: T-0684
type: task
nature: feature
title: flai check validates design/experiments documents, and the folder, its README, and its document template ship here and in the template
status: done
parent: S-0194
owner: arobson
created: 2026-10-02T10:08:08Z
updated: 2026-10-02T10:13:37Z
transitions:
  - to: ready
    at: 2026-10-02T10:08:27Z
    by: agent-S-0194
  - to: in-progress
    at: 2026-10-02T10:11:14Z
    by: agent-S-0194
  - to: done
    at: 2026-10-02T10:13:37Z
    by: agent-S-0194
stream: S-0194
tags: []
usage:
  source: log
  seconds: 143
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 43
      output: 10384
      cache_read: 3004642
      cache_write: 34761
      cost: 1.0869
---

# T-0684 flai check validates design/experiments documents, and the folder, its README, and its document template ship here and in the template

## Work

`flai check` validates every `design/experiments/<S-nnnn>-<slug>.md`: front matter `title`, `updated`, `status`, and `story`, which names an existing story, and the sections Hypothesis, Success measure, What was done, Results, and Recommendation, whose text says adopt, adapt, or drop. `design/experiments/README.md` and a template for the document ship here and in `template/root/`. `CLAUDE.md`, the template's `CLAUDE.md`, and `design/system/repository-layout.md` list the folder.

## Done when

- Check tests cover a valid document and each thing missing; the template render test passes with the new files.

## Notes
