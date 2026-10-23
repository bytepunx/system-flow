---
id: T-1812
type: task
nature: feature
title: The digest folder's README, the layout, the CLI design, and the user docs say what the digest is, where it lives, and how flai writes it, and the template ships the folder
status: done
parent: S-0428
owner: alex
created: 2026-10-16T06:03:37Z
updated: 2026-10-23T17:11:55Z
transitions:
  - to: ready
    at: 2026-10-23T17:06:51Z
    by: agent-S-0428
  - to: in-progress
    at: 2026-10-23T17:06:52Z
    by: agent-S-0428
  - to: done
    at: 2026-10-23T17:11:55Z
    by: agent-S-0428
stream: S-0428
tags: [cli, template]
touches: [design/digest/README.md, template/root/design/digest/README.md, design/README.md, template/root/design/README.md.tmpl, design/system/repository-layout.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/conventions.md, CLAUDE.md, template/root/CLAUDE.md.tmpl, template/CHANGELOG.md, template/template.yaml]
after: [T-1809, T-1810, T-1811]
usage:
  source: log
  seconds: 303
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 93
      output: 38105
      cache_read: 5131885
      cache_write: 195119
      cost: 2.9377
---
# T-1812 The digest folder's README, the layout, the CLI design, and the user docs say what the digest is, where it lives, and how flai writes it, and the template ships the folder

## Work

- Write `design/digest/README.md` and its template copy `template/root/design/digest/README.md`: what the folder is for, one fact per entry, that only flai writes it (the seed, the curator, the operator's answers), and that an entry never holds a secret's value.
- Add the folder to `design/README.md`, `template/root/design/README.md.tmpl`, the tree in `design/system/repository-layout.md`, the layout and artifact table in `docs/users/conventions.md`, and the layout table of `CLAUDE.md` and `template/root/CLAUDE.md.tmpl`, beside `design/analysis/`.
- Describe `flai digest` and the two check findings in `design/system/flai-cli.md` and `docs/users/flai.md`.
- Add a `template/CHANGELOG.md` entry and bump the version in `template/template.yaml`.

Waits for the check, commands, and search tasks, so the docs describe what was built.

## Done when

The documents say what the code does, `flai check --strict` reports nothing on them, and the markdown lint passes.

## Notes

`docs/operators/settings.md` and `design/system/project-manifest.md` are left out: the folder follows `layout.design` and adds no setting.
