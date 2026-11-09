---
id: I-0249
title: TestRenderPrototypeTemplate refuses any double curly brace in a rendered project, so a baseline formula the template ships cannot use an input
class: blocker
status: closed
count: 1
cost: 10m
first_reported: 2026-11-02T02:50:10Z
last_reported: 2026-11-02T02:50:10Z
updated: 2026-11-09T10:32:11Z
---

# I-0249 TestRenderPrototypeTemplate refuses any double curly brace in a rendered project, so a baseline formula the template ships cannot use an input

## Description
TestRenderPrototypeTemplate refuses any double curly brace in a rendered project, so a baseline formula the template ships cannot use an input

## Instances

### 2026-11-02T02:50:10Z
Story: S-0421.
S-0421's close-out stopped at integration because template/root/design/workflows/formulas/README.md showed the input placeholder `{{name}}` (ADR-0238 §4). flai/cmd/new_test.go's TestRenderPrototypeTemplate fails on any file of the rendered project that holds `{{`, taking it for unrendered template syntax. S-0421 now describes the placeholder in words. S-0423 and S-0424 ship baseline formulas under template/root/design/workflows/formulas whose steps take inputs, so they will need the test to pass over formula files, or the renderer to say which files it renders.

## Remediation

Story S-0616 remediates this issue, created from it at 2026-11-07T23:39:18Z.
Closed 2026-11-09T10:32:11Z: S-0616: TestRenderPrototypeTemplate (flai/cmd/new_test.go) now looks for template syntax only in a file rendered from a .tmpl source or one that does not match its source under template/root byte for byte, through leftTemplateSyntax and its behaviour test, so a {{ in a copied file, such as a formula's input placeholder or a role's start prompt, no longer fails it. The workarounds withoutFormulaInputs (S-0424) and withoutRoleStart (S-0554) are gone, and the formulas README shows {{name}} again, proving the fix on the real template.
