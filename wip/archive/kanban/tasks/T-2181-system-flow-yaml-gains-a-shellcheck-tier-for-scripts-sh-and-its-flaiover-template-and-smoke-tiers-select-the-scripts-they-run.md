---
id: T-2181
type: task
nature: improvement
title: "system-flow.yaml gains a shellcheck tier for scripts/**/*.sh, and its flaiover, template, and smoke tiers select the scripts they run"
status: done
parent: S-0503
owner: alex
created: 2026-10-26T05:23:22Z
updated: 2026-10-26T05:37:38Z
transitions:
  - to: ready
    at: 2026-10-26T05:34:46Z
    by: agent-S-0503
  - to: in-progress
    at: 2026-10-26T05:34:47Z
    by: agent-S-0503
  - to: done
    at: 2026-10-26T05:37:38Z
    by: agent-S-0503
stream: S-0503
tags: [scripts, tests, manifest]
touches: [system-flow.yaml, design/conventions/tooling.md, scripts/README.md, design/system/project-manifest.md]
after: [T-2180]
usage:
  source: log
  seconds: 171
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 35
      output: 10823
      cache_read: 1971939
      cache_write: 58994
      cost: 1
---
# T-2181 system-flow.yaml gains a shellcheck tier for scripts/**/*.sh, and its flaiover, template, and smoke tiers select the scripts they run

## Work

- In `system-flow.yaml`'s `tests`, add a tier `shellcheck` ahead of `markdown`: `command: [scripts/shellcheck.sh, "{files}"]`, `paths: ["scripts/**/*.sh"]`, `format: plain`, `all_command: [scripts/shellcheck.sh]`. Not `all_only`: a change to a script runs it in `flai test` and in `flai verify`. Check that `scripts/**/*.sh` selects `scripts/flai.sh` as `flai/**/*.go` selects `flai/main.go`; if the pattern needs `scripts/*.sh` too, give both.
- Widen the `paths` of the `all_only` tiers that run scripts, so a change to a script runs the tests that exercise it in `flai verify` and in the merge queue's batch, which runs the tiers its paths select (ADR-0174): `flaiover` adds `scripts/flaiover-*.sh`, since `scripts/flaiover-test.sh` runs `flaiover-install-test.sh`, `flaiover-install.sh`, `flaiover-lint.sh`, and `flaiover-unit.sh`; `template` adds `scripts/template-test.sh`; `smoke` adds `scripts/**`, since `scripts/smoke.sh` runs `template-test.sh`, `check.sh`, `lint-md.sh`, `install-test.sh`, and `gateway-smoke.sh`. This is the gap I-0199 names: S-0481 changed `scripts/flaiover-install.sh` and added `scripts/flaiover-install-test.sh`, and no tier ran it.
- Document the tiers where this repository lists them: the `flai test` bullet under Project additions in `design/conventions/tooling.md`, the first paragraph of `scripts/README.md`, and the sentence in `design/system/project-manifest.md`'s `tests` bullet that says which `paths` this repository gives its `all_only` tiers.

Waits for T-2180: the tier's command must exist and pass on the tree before the manifest names it.

## Done when

- `flai check --strict` reports no `manifest.tests` finding.
- `scripts/flai.sh test scripts/flai.sh` runs the `shellcheck` tier and passes, and `scripts/flai.sh test --all` lists it.
- The three documents name the `shellcheck` tier and the widened paths; the markdown lint is clean.

## Notes

Drafted by planner-S-0503. The template's manifest (`template/root/system-flow.yaml.tmpl`) is left alone: its `test` tier selects `**`, so a project made from it has no such gap.
