---
id: D-0305
title: "S-0381 stays a draft until S-0380's decision: I would prefer we not store front-matter in tokens"
kind: answer
sources: [TH-0439]
topics: [design, flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [flai/internal/designsystem/tokens.go, flai/internal/designsystem/tokens_test.go, flai/internal/designsystem/resolver.go, flai/internal/designsystem/resolver_test.go, flai/internal/designsystem/color.go, flai/internal/designsystem/color_test.go, flai/internal/designsystem/css.go, flai/internal/designsystem/css_test.go, flai/internal/designsystem/starter.md, flai/internal/designsystem/starter.tokens.json, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/check/design.go, flai/internal/check/design_test.go, flai/internal/check/check.go, flai/cmd/design.go, flai/cmd/design_test.go, flai/cmd/root.go, flai/internal/hostapi/design.go, flai/internal/hostapi/design_test.go, flai/internal/hostapi/hostapi.go, design/system/flai-cli.md, design/system/project-manifest.md, design/system/documentation-standard.md, design/conventions/documentation.md, template/root/design/conventions/documentation.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, flai/internal/designsystem/testdata/full.tokens.json, flai/internal/designsystem/testdata/full.css.golden, flai/internal/designsystem/testdata/full.tailwind.css.golden, flai/internal/designsystem/testdata/show-missing.json.golden, flai/internal/designsystem/testdata/show-problems.json.golden, flai/internal/designsystem/testdata/show.json.golden, flai/internal/designsystem/testdata/themes.css.golden, flai/internal/designsystem/testdata/themes.tailwind.css.golden, flai/internal/designsystem/testdata/themes.tokens.json, flai/internal/designsystem/starter.go, flaiover/src/lib/server/agent.ts]
date: 2026-10-14
status: absorbed
absorbed_by: "design/adrs/0147-the-dashboard-reads-and-saves-a-text-file-that-is-not-markdown-through-file-get.md#decision"
asked_again: []
cited: []
---

# D-0305 S-0381 stays a draft until S-0380's decision: I would prefer we not store front-matter in tokens

The operator, alex, answered TH-0439 at 2026-10-14T23:11:49Z.

What was asked:

```text
**Recommendation:** keep S-0381 a draft until S-0380 is decided, then finalize it as its planner replans it.
```

The answer:

```text
I would prefer we not store front-matter in tokens. In the precursor research, my request was that we instead design file.get and file.save (analogs to doc.get and doc.save) that are capable of working outside of markdown to store the tokens in a different format and also allow for expansion into features that give the operator design tool surfaces, component viewers/testers, and the ability to preview UI/UX without the need to complete a story to see the outcome visually
```

Source: TH-0439.
