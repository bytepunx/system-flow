---
id: T-2204
type: task
nature: remediation
title: The install step in component-library.md names flai accept's install as the one exception to the main checkout never being installed into, with I-0179 as the reason and a link to ADR-0185
status: done
parent: S-0526
owner: alex
created: 2026-10-28T13:00:57Z
updated: 2026-10-28T13:06:40Z
transitions:
  - to: ready
    at: 2026-10-28T13:06:17Z
    by: agent-S-0526
  - to: in-progress
    at: 2026-10-28T13:06:18Z
    by: agent-S-0526
  - to: done
    at: 2026-10-28T13:06:40Z
    by: agent-S-0526
stream: S-0526
tags: [design, components]
touches: [design/system/component-library.md]
usage:
  source: log
  seconds: 18
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 12
      output: 61
      cache_read: 667385
      cache_write: 7128
      cost: 0.2282
---
# T-2204 The install step in component-library.md names flai accept's install as the one exception to the main checkout never being installed into, with I-0179 as the reason and a link to ADR-0185

## Work

Rewrite the last sentence of **The install step** under `design/system/component-library.md` § A story's worktree (S-0442) › Sources. Today it ends "The main checkout is never installed into." Since S-0483 that is false: `flai accept` runs a sub-project's `install` command in the main checkout after the acceptance commit when the accepted change touched a path the manifest's `projects[].install.paths` names ([ADR-0185](../../../design/adrs/0185-flai-accept-runs-a-sub-project-s-install-command-in-the-main-checkout-when-the.md)). The operator asked on TH-0578 that the rule be updated with the exception and its reason.

Write the sentence so that it says the dev server never installs into the main checkout, and that the one install there is `flai accept`'s, by the manifest's `install` command when the accepted change touched its `install.paths`, because a stale `node_modules` in the main checkout made every run that needed a new package fail (I-0179). Link ADR-0185 with a relative path, as the document's other ADR links are written. Say that this install runs with no time limit, unlike the dev server's ten minutes: `runInstalls` in `flai/cmd/accept_install.go` runs under `context.Background()`, and ADR-0185 states no timeout. Do not change the code.

Bump the document's `updated` date. Change nothing else in the document; `docs/users/flai.md` § Accept and release and `docs/operators/settings.md` already describe the install correctly and need no change.

No test fits: the defect is one sentence of the living design, which `flai check` and the markdown lint check for form, not for truth. Say so in the story's Notes when you tick criterion 1.

## Done when

- The install step in `design/system/component-library.md` no longer says the main checkout is never installed into, names `flai accept`'s install as the one exception, gives I-0179 as the reason, links ADR-0185, and says the install has no time limit.
- The document's `updated` date is bumped.
- `flai check --strict` and the markdown lint pass on the file.

## Notes

Drafted by planner-S-0526 from I-0192's instance and TH-0589.
