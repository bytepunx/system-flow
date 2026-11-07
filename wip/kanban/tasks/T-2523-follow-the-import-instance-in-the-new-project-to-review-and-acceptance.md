---
id: T-2523
type: task
nature: improvement
title: Follow the import instance in the new project to review and acceptance
status: backlog
parent: S-0564
owner: alex
created: 2026-11-07T03:47:29Z
updated: 2026-11-07T03:47:29Z
transitions: []
stream: S-0564
tags: [template]
after: [T-2514]
---
# T-2523 Follow the import instance in the new project to review and acceptance

## Work

`flai serve` starts the instance's agent in the new project. That agent works the formula's two steps from the instance's worktree: `flai check --strict`, then `flai import test`. It records in the instance's Notes what each reported, and which tests `flai import test` found and how each went. Follow it with `flai board` and `flai show <instance>` in the new project, or with the MCP tools and the new project's key.

When a step fails, the instance's agent blocks its task and asks the operator on a thread in the new project. Answering it is the operator's. Name that thread on a thread on S-0564, so the operator sees it beside this plan. Fixing the scratch repository is not this story's work.

The instance is accepted by the operator, or by the orchestrator where `orchestration.permissions.accept_reviews` is on for the new project. Write nothing in S-0564's Notes yet: the record of both runs is a later task.

This task waits for T-2514, which makes the instance.

## Done when

- The instance is done in the new project, and its Notes name each step's exit and the tests `flai import test` ran.

## Notes
