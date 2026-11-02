---
id: I-0249
title: TestRenderPrototypeTemplate refuses any double curly brace in a rendered project, so a baseline formula the template ships cannot use an input
class: blocker
status: open
count: 1
cost: 10m
first_reported: 2026-11-02T02:50:10Z
last_reported: 2026-11-02T02:50:10Z
updated: 2026-11-02T02:50:10Z
---

# I-0249 TestRenderPrototypeTemplate refuses any double curly brace in a rendered project, so a baseline formula the template ships cannot use an input

## Description
TestRenderPrototypeTemplate refuses any double curly brace in a rendered project, so a baseline formula the template ships cannot use an input

## Instances

### 2026-11-02T02:50:10Z
Story: S-0421.
S-0421's close-out stopped at integration because template/root/design/workflows/formulas/README.md showed the input placeholder `{{name}}` (ADR-0238 §4). flai/cmd/new_test.go's TestRenderPrototypeTemplate fails on any file of the rendered project that holds `{{`, taking it for unrendered template syntax. S-0421 now describes the placeholder in words. S-0423 and S-0424 ship baseline formulas under template/root/design/workflows/formulas whose steps take inputs, so they will need the test to pass over formula files, or the renderer to say which files it renders.

## Remediation
