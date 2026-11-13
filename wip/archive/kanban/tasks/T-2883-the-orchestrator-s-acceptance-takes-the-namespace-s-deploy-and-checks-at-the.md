---
id: T-2883
type: task
nature: feature
title: The orchestrator's acceptance takes the namespace's deploy and checks at the verified commit as evidence for a criterion
status: done
parent: S-0662
owner: alex
created: 2026-11-10T08:27:33Z
updated: 2026-11-13T14:27:31Z
transitions:
  - to: ready
    at: 2026-11-13T13:56:32Z
    by: agent-S-0662
  - to: in-progress
    at: 2026-11-13T13:56:32Z
    by: agent-S-0662
  - to: done
    at: 2026-11-13T14:27:31Z
    by: agent-S-0662
stream: S-0662
tags: [cli, template]
touches: [flai/internal/preview/orchestrator.go, flai/internal/preview/evidence_test.go, flai/cmd/accept.go, flai/cmd/accept_orchestrator_test.go, design/system/strategic-agents.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, template/CHANGELOG.md]
after: [T-2875, T-2877]
usage:
  source: log
  seconds: 1859
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 64
      output: 24693
      cache_read: 4744590
      cache_write: 203179
      cost: 2.5523
---
# T-2883 The orchestrator's acceptance takes the namespace's deploy and checks at the verified commit as evidence for a criterion

## Work

Let a criterion's evidence name the kluster record in place of changed files, in the form S-0659's ADR gives, such as `- <n>: kluster checks passed at <commit>`. `ParseEvidence` in `flai/internal/preview/orchestrator.go` reads it, and `unevidenced` accepts it only when T-2875's evidence says the deploy and the checks passed at the verified commit; otherwise it names the criterion and what the record says. Change the message in `flai/cmd/accept.go` that names the evidence forms.

Say it in `design/system/strategic-agents.md` § Accepting a story (S-0221), in `design/conventions/strategic-agents.md` § Accepting a story and its template copy, and in the orchestrator's role file, its template copy, and `flai/internal/roles/baseline/orchestrator.md`, which `baseline_test.go` holds to the template's. Add a line under `## Unreleased` in `template/CHANGELOG.md`.

It waits for T-2875, whose evidence it reads, and for T-2877, which writes `template/CHANGELOG.md` too.

## Done when

- `evidence_test.go` covers the new line parsed, accepted at the verified commit, and refused at another commit, after a failed checks run, and with the action off.
- `accept_orchestrator_test.go` lands a story with such evidence and refuses one whose record does not match.
- The three role copies are the same, and `TestBaselineRolesAreTheTemplates` or its equal passes.
- `flai test` passes on the changed files.

## Notes

- Criterion 5 of S-0662. Today only changed files, or a formula step's task, are evidence.
