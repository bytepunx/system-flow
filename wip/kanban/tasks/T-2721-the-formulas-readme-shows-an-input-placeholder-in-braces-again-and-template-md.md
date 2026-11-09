---
id: T-2721
type: task
nature: remediation
title: The formulas README shows an input placeholder in braces again, and template.md says copied files may hold them
status: done
parent: S-0616
owner: alex
created: 2026-11-09T10:21:01Z
updated: 2026-11-09T10:31:51Z
transitions:
  - to: ready
    at: 2026-11-09T10:31:26Z
    by: agent-S-0616
  - to: in-progress
    at: 2026-11-09T10:31:27Z
    by: agent-S-0616
  - to: done
    at: 2026-11-09T10:31:51Z
    by: agent-S-0616
stream: S-0616
tags: [template, docs]
touches: [template/root/design/workflows/formulas/README.md, design/workflows/formulas/README.md, design/system/template.md, template/CHANGELOG.md]
after: [T-2720]
usage:
  source: log
  seconds: 24
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 1041
      cache_read: 1303508
      cache_write: 15501
      cost: 0.4056
---
# T-2721 The formulas README shows an input placeholder in braces again, and template.md says copied files may hold them

## Work

S-0421 had to reword the formulas README so that it said "between double curly braces" instead of showing the placeholder. This task puts the placeholder back, on the real template, to show that I-0249's instance no longer fails.

- In `template/root/design/workflows/formulas/README.md`, rewrite the sentence "An input is used by its name between double curly braces, as ADR-0238 §4 shows" so that it shows the form in code, `{{name}}`.
- Make the same edit in `design/workflows/formulas/README.md`, this project's copy. The two files are identical today and stay identical.
- In `design/system/template.md`, the `render.suffix` comment says that files with the suffix are rendered. Add that every other file is copied byte for byte, so a `{{` in it is content: for example, a formula's input placeholders or a role file's start prompt. Bump `updated`.
- Add a patch line for S-0616 under `## Unreleased` in `template/CHANGELOG.md`.

This task waits for T-2720. Before T-2720's change, `TestRenderPrototypeTemplate` fails on this README. This task's check is that test passing over the README.

## Done when

- Both formulas READMEs show `{{name}}` and match byte for byte.
- `go test -run TestRenderPrototypeTemplate ./cmd/` passes without `-short` in `flai/`.
- `design/system/template.md` says which files are copied verbatim.
- `template/CHANGELOG.md` has the S-0616 patch line.
- `flai test` on the four paths passes, the markdown lint included.

## Notes
