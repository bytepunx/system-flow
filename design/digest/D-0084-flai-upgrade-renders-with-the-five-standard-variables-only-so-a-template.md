---
id: D-0084
title: flai upgrade renders with the five standard variables only, so a template variable a fork adds fails every upgrade
kind: cause
sources: [I-0040]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/cmd/new.go, flai/cmd/upgrade.go, flai/cmd/import.go, flai/cmd/new_test.go, flai/cmd/upgrade_test.go, flai/cmd/import_test.go, flai/internal/template, flai/internal/manifest, flai/internal/lock, flai/internal/upgrade, docs/contributors/template.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/system/template.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/repository-layout.md, design/issues]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0084 flai upgrade renders with the five standard variables only, so a template variable a fork adds fails every upgrade

I-0040 was closed on 2026-10-01.

The cause:

```text
flai upgrade renders with the five standard variables only, so a template variable a fork adds fails every upgrade
```

The fix:

```text
S-0185: system-flow.lock.yaml records every template variable's value under vars; flai upgrade renders with --var, the manifest's five fields, the recorded values, then defaults, and refuses a required variable with no value before changing anything (TestUpgradeRendersAForksOwnVariables)
```

Source: I-0040.
