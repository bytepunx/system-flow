---
id: D-0661
title: "Accept ADR-0312 as written: Recommended: accept ADR-0312 as written"
kind: answer
sources: [TH-1258]
topics: [verify, workflow, cli, template, flai, go, conventions, code, all]
paths: [design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md, design/adrs/README.md, flai/internal/manifest/manifest_tests.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/tests_test.go, flai/internal/verify/verify.go, flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/cmd/test.go, flai/cmd/test_test.go, flai/internal/mcpserver/test.go, flai/internal/mcpserver/test_test.go, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, system-flow.yaml, design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
date: 2026-11-17
status: absorbed
absorbed_by: "design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md#decision"
asked_again: []
cited: [S-0806]
---

# D-0661 Accept ADR-0312 as written: Recommended: accept ADR-0312 as written

The operator, alex, answered TH-1258 at 2026-11-17T07:46:25Z.

What was asked:

```text
ADR-0312 (design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) is `proposed`. S-0805's criterion 1 needs it accepted, or superseded by an ADR with your changes.
```

The answer:

```text
**Recommended:** accept ADR-0312 as written. S-0805 then goes on with §1 and the template's four labels, as agent-S-0805 proposes.

**Why:**

- Its Context records your decisions of 2026-11-16: steps name suites, not tiers; a cluster step is skipped where no cluster is configured; this repository changes nothing until its own file says so.
- Its Decision is what E-0034's ten stories S-0805 to S-0815 are planned on. Their plans are approved and finalized, all but S-0815.
- One tension: §8 puts the new keys in the template's file now, and a Consequence says the template should not use them until the release that reads them is the oldest a project runs. S-0806's T-3473 settles it. Listing the keys raises `flai.minimum` on publish, so an older flai stops with a message naming the release instead of a bare parse error.

**Why I am asking you:** accepting an ADR is your decision.

**Alternative:** name what to change. agent-S-0805 then writes a superseding ADR under T-3444, and the stories after it are planned again where they read the change.

Source: design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md § Decision
```

Source: TH-1258.
