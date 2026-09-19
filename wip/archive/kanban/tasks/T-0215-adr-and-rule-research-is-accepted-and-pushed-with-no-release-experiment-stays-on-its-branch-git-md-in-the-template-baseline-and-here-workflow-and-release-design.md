---
id: T-0215
type: task
nature: feature
title: "ADR and rule: research is accepted and pushed with no release, experiment stays on its branch; git.md in the template baseline and here, workflow and release design"
status: done
parent: S-0053
owner: alex
created: 2026-09-19T07:15:37Z
updated: 2026-09-19T07:17:00Z
transitions:
  - to: ready
    at: 2026-09-19T07:15:38Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T07:15:38Z
    by: system-flow
  - to: done
    at: 2026-09-19T07:17:00Z
    by: system-flow
stream: S-0053
tags: []
---

# T-0215 ADR and rule: research is accepted and pushed with no release, experiment stays on its branch; git.md in the template baseline and here, workflow and release design

## Work
Write the ADR (the tooling enforces the rule, so the change is recorded): a research story's findings (ADRs, design documents, new stories) are deliverables of the repository, not of a component, so acceptance merges, archives, commits, and pushes them and cuts no release, whatever files the commits touched; an experiment stays on its branch and acceptance keeps refusing it, the operator's decision on 2026-09-19. Replace the sentence in `design/conventions/git.md`, in `template/root/design/conventions/git.md` first and then here above the end-of-baseline marker, the same text in both; the story's criterion, approved by the operator, is the authority for the baseline edit. Bring `design/system/workflow.md` and the release design into agreement. The pre-release idea in `git.md` stays as it is.

## Done when
- The ADR is accepted-format and indexed
- `git.md` reads the same in the template and here, and the design documents agree with it

## Notes
