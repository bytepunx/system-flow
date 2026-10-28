---
id: D-0085
title: An explicitly empty --var for a required template variable passes the required check, and flai new writes the files before failing
kind: cause
sources: [I-0041]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/cmd/new.go, flai/cmd/upgrade.go, flai/cmd/import.go, flai/cmd/new_test.go, flai/cmd/upgrade_test.go, flai/cmd/import_test.go, flai/internal/template, flai/internal/manifest, flai/internal/lock, flai/internal/upgrade, docs/contributors/template.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/system/template.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/repository-layout.md, design/issues]
date: 2026-10-01
status: absorbed
absorbed_by: "design/system/template.md#variables"
asked_again: []
cited: []
---

# D-0085 An explicitly empty --var for a required template variable passes the required check, and flai new writes the files before failing

I-0041 was closed on 2026-10-01.

The cause:

```text
An explicitly empty --var for a required template variable passes the required check, and flai new writes the files before failing
```

The fix:

```text
S-0185: flai new and flai import apply the required check to a value given with --var too, and stop before writing anything (TestNewRefusesAnEmptyRequiredVarBeforeWriting)
```

Source: I-0041.
