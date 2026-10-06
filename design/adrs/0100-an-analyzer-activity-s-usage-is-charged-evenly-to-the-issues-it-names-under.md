---
id: ADR-0100
title: "An analyzer activity's usage is charged evenly to the issues it names, under each issue's usage.strategic, and the story made from an issue carries it, counted once"
status: accepted
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0083, ADR-0095]
topics: [cli, analysis, planning]
---

# ADR-0100 An analyzer activity's usage is charged evenly to the issues it names, under each issue's usage.strategic, and the story made from an issue carries it, counted once

## Context

ADR-0083 charges a planner activity's usage to the item it planned, and ADR-0095 an orchestrator activity's to the work items it named, each under `usage.strategic`, and leaves the analyzer's in its project strategic total. The analyzer's work ends in a report and in the issues it files or bumps for it (S-0224), and issues carried no usage. The designer asked (S-0227) that the analyzer's cost land on the issues it found, and travel with an issue to the draft story the issue step makes from it, so that the cost of a remediation includes the cost of finding it. The analyzer runs once and ends without calling `activity_log`, so on TH-0200 the designer chose that its run's end name the issues that name its report.

## Decision

An issue may carry `usage` with `strategic` entries alone, an analyzer activity's usage is split evenly between the issues it names, and the story the issue step makes from an issue carries the issue's entries, counted once in `flai stats`.

1. **An issue's usage.** An issue's front matter may carry `usage` with `strategic` only: one entry per strategic agent kind, in the shape ADR-0083 gives items (`kind`, `seconds`, `estimated: true`, `models`). It has no agents' figures. An issue without usage has no `usage` key.
2. **Charged to the issues it names.** When flai serve logs an analyzer activity, the usage apportioned to its span, as ADR-0083 apportions a planner's, is split evenly between the issues the entry names, as ADR-0095 splits an orchestrator's between work items, and each share is added under the issue's `analyzer` entry. An ID that names no issue takes no share, and neither does a work item the entry names. An `activity_log` call names its issues itself; an analyzer run's end names the issues, open or closed, that name the report the run wrote, by an instance's `Report:` line or a Remediation link (TH-0200). An activity that names no issue is left in the analyzer's project strategic total.
3. **Carried to the story made from it.** The story `flai issue story` or the MCP tool `issue_story` makes from an issue is charged each of the issue's strategic entries, under its kind, and so is every item above it, up to its epic, as ADR-0083 sums a planner's charge. Its Notes say what was carried. The issue keeps its own entries.
4. **Counted once.** A story is made from an issue when the issue's `## Remediation` names it in the line the issue step writes, `Story S-nnnn remediates this issue, created from it at <ts>.`, and the story exists. `flai stats` lists each issue's strategic usage (`strategic_issues`), with the story made from it and whether it is counted. In each kind's totals it counts an issue's entries under `issues` while no story was made from it, and the story's, under `items`, after. The kind's project strategic total is its document's accrued totals less `items` less `issues`. The agents' figures never include an issue's usage.

## Consequences

- The cost of finding a problem is on the issue that records it, and on the remediation story and its epic once the operator chooses to make one.
- An issue's file changes when an analyzer activity is charged to it. The analyzer works in the main checkout, where its issues are written and left uncommitted, so the charge is committed with them.
- An issue bumped by an analysis after a story was made from it is charged, but not counted: its story's entry was copied when the story was made. That share stays out of `items` and `issues` and so is counted in the kind's project total.
- A story merely linking an issue, by naming it, carries nothing and does not stop the issue's usage being counted.
- `usage` is listed for issues in `front-matter-fields.txt`, so publishing the flai release that carries it raises `flai.minimum`.
- `design/system/continuous-improvement.md`, `design/system/strategic-agents.md`, and `design/system/metrics.md` describe the issue's usage, the charge, the carry-over, and the figures.

## Alternatives considered

- **Move the usage off the issue to the story when one is made.** The issue would lose the record of what finding it cost, and a story cancelled or archived would take it with it. Keeping both and counting once keeps each document true.
- **Charge the issues' parent story or epic.** An issue has no parent until a story is made from it, and the operator may never make one.
- **Charge the whole of an activity to each issue it names.** The issues would add up to more than was spent.
- **Ask the analyzer to call `activity_log` naming the issues it filed.** The charge would depend on the agent remembering the call; the report already names the issues, so flai can find them (TH-0200).
- **Tell a story made from an issue by any link to it.** A story can name an issue without being made from it, and only the story the issue step made carries the usage.
