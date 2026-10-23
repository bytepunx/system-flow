---
id: D-0276
title: flai test vets and tests only the changed Go packages, so a change that breaks a package importing them is found only by the integration tier
kind: cause
sources: [I-0128]
topics: [testing, flai, cli, go, code, all]
paths: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/verify/importers.go, flai/internal/verify/importers_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/verify/verify.go, system-flow.yaml, design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md, design/issues/I-0128-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md, design/issues/summary.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0276 flai test vets and tests only the changed Go packages, so a change that breaks a package importing them is found only by the integration tier

I-0128 was closed on 2026-10-12.

The cause:

```text
flai test vets and tests only the changed Go packages, so a change that breaks a package importing them is found only by the integration tier
```

The fix:

```text
Fixed by S-0364: a test tier's command may hold {importers}, the changed Go packages and every package of the module that imports one of them, in its own files or its test files. This repository's vet tier runs go vet on {importers}, so a change that breaks a package importing the changed one, as S-0338's change to internal/workitem broke cmd/story_start_test.go, fails flai test and the close-out's vet tier rather than only the integration tier. flai test flai/internal/workitem now vets ./cmd and every other importer, in 2.5s.
```

Source: I-0128.
