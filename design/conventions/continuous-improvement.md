---
title: Continuous improvement
updated: 2026-10-07
audience: agent
order: 100
status: active
topics: [all]
roles: [story, analyze]
---

# Continuous improvement

How to record recurring friction, defects, blockers, and inefficiencies so they can be measured and remediated. The folder and schema are in `design/system/continuous-improvement.md`.

## Rules

- Keep a list of items that interfere with progress, quality, and efficacy filed under `design/issues`, one file per issue
- Classifications:
  - `defect`: something is wrong
  - `blocker`: progress stopped until someone acts
  - `efficiency`: it works but costs time
  - `impression`: a suspected problem without a measurement
- Provide a high-level description with more detailed instances of the identified issue or impression
- Front matter:
  - increment the `count` each time the issue comes up
  - keep `first_reported` and `last_reported` timestamps
  - if possible, record `cost`: average wall-clock time one occurrence costs, as a duration like `20m`
- `summary.md`
  - in the issues folder
  - a table of open issues with count, average cost, total cost, and last occurrence
  - order by issue cost
  - update it with every occurrence
- Record an occurrence with `flai issue new`, or `flai issue bump` when the issue exists: each instance names the story you work
- Working a story, record, bump, or close an issue in one call from the story's worktree: `flai issue new|bump|close --commit`, or the MCP tool `issue_new`, `issue_bump`, or `issue_close` with `commit`:
  - it numbers a new issue from the whole repository, commits exactly what it wrote on `story/S-nnnn` as `docs: [S-nnnn] ...`, with each `--trailer`, and adds it to the story's touches
  - follow it with no `git add`, `git commit`, or `flai touches`
  - off the story's branch it refuses and writes nothing; `--autocommit` is for the main checkout, not a story's agent
- Make no story for an issue yourself; the operator chooses, at acceptance, which issues become stories:
  - a story's review page lists the open issues it recorded or bumped, checked, and every other open issue no open story links, unchecked
  - accepting it makes a `remediation` or `improvement` backlog story for each checked issue
  - such a story links the issue document, recommends a solution, and its last criterion closes the issue
- Make that story with `flai issue story` (MCP `issue_story`) only when the operator asks, or when `flai check` warns `issues.no-story` about an issue open too long with no story

## When in doubt

- Record the impression as such rather than ignoring it or over-reporting.
- Increment similar existing issues and add the new instance to its body if it differs enough from earlier ones.

<!-- system-flow:end-of-baseline -->

## Project additions

- Record with `flai issue new`, `bump`, and `close`; never edit issue front matter or `summary.md` by hand. Instances and remediation text in the body are edited freely.
