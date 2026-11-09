---
id: T-2720
type: task
nature: remediation
title: TestRenderPrototypeTemplate looks for template syntax only in files the renderer rendered or changed
status: done
parent: S-0616
owner: alex
created: 2026-11-09T10:20:52Z
updated: 2026-11-09T10:30:46Z
transitions:
  - to: ready
    at: 2026-11-09T10:24:26Z
    by: agent-S-0616
  - to: in-progress
    at: 2026-11-09T10:24:28Z
    by: agent-S-0616
  - to: done
    at: 2026-11-09T10:30:46Z
    by: agent-S-0616
stream: S-0616
tags: [flai, test]
touches: [flai/cmd/new_test.go]
usage:
  source: log
  seconds: 378
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 1387
      cache_read: 1737411
      cache_write: 20661
      cost: 0.5406
---
# T-2720 TestRenderPrototypeTemplate looks for template syntax only in files the renderer rendered or changed

## Work

`flai/internal/template/render.go` renders only files whose name ends in `render.suffix` (`.tmpl`). It copies every other file byte for byte. `TestRenderPrototypeTemplate` in `flai/cmd/new_test.go` still refuses `{{` in every file of the rendered project. It now works around that with `withoutFormulaInputs` (S-0424) and `withoutRoleStart` (S-0554).

- Move the walk's syntax check into a helper. It takes the template's `root/` folder, the render suffix, and the rendered project.
- For each rendered file, find its source: the same relative path, or that path with the suffix added.
- Report `{{` in a file rendered from a suffixed source, and in a file with no source in `root/` (one `flai new` wrote itself, such as the lock). Report it too in a file whose source has no suffix but whose bytes differ from that source.
- Pass over a file that matches its unsuffixed source byte for byte: its `{{` is the template's content.
- Keep the `.tmpl` leftover check and the front matter parse as they are.
- Remove `withoutFormulaInputs`, `withoutRoleStart`, `TestWithoutFormulaInputs`, and `TestWithoutRoleStart`, and any import only they used.
- Add a short-mode test of the helper over a small template under `t.TempDir()` that reproduces I-0249:
  - a verbatim `design/workflows/formulas/README.md` holding `{{name}}` passes;
  - a `.tmpl` source whose output still holds `{{` is reported;
  - a verbatim file changed after the render is reported.

This task waits for none: it is the first layer.

## Done when

- `TestRenderPrototypeTemplate` uses the helper, and the two workaround helpers and their tests are gone.
- The new short test passes and covers the three cases above.
- `flai test flai/cmd/new_test.go` passes.
- `go test -run TestRenderPrototypeTemplate ./cmd/` passes without `-short` in `flai/`.

## Notes
