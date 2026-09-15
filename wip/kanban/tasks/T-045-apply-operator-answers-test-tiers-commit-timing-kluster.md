---
id: T-045
type: task
nature: feature
title: "Apply operator answers: test tiers, commit timing, kluster"
status: done
parent: S-026
owner: alex
created: 2026-09-15T22:49:31Z
updated: 2026-09-15T22:52:23Z
transitions:
  - to: ready
    at: 2026-09-15T22:49:31Z
    by: agent
  - to: in-progress
    at: 2026-09-15T22:49:31Z
    by: agent
  - to: done
    at: 2026-09-15T22:52:23Z
    by: agent
stream: S-026
tags: [conventions]
---

# T-045 Apply operator answers: test tiers, commit timing, kluster

## Work
Add the three-tier test rules (behavior, integration, smoke) to code-quality.md and tooling.md in the template and here, with `make test`, `make integration`, `make smoke` and the folder layout; state that commits do not happen at task boundaries in git.md; leave kluster in the baseline; add the scripts and Makefile targets here and stubs in the template; mark the git-clone test as integration so `make test` stays fast.

## Done when
Conventions updated in both places, `make test`, `make integration`, `make smoke` run here, template renders and checks.

## Notes
