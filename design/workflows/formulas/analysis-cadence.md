---
title: Analyze the project and list the issues the report names
updated: 2026-11-02
status: active
name: analysis-cadence
description: "Run the analyzer on its focus, or on all three, and list the open issues its report names, for the operator to choose stories from at acceptance."
nature: improvement
inputs:
  - name: focus
    description: "bottlenecks, intent, or risk; empty for all three"
    default: ""
steps:
  - id: analyze
    title: Run the analyzer and write its report
    run: {role: analyzer, prompt: "Analyze this project with the focus \"{{focus}}\": bottlenecks, intent, or risk, and all three when the focus is empty. Write one report under design/analysis, and file or bump an issue for each actionable finding, linked from the report."}
  - id: issues
    title: List the open issues the analyzer's report names
    run: {flai: issue list --report latest --json}
    after: [analyze]
    in: main
---

# Analysis cadence

`flai serve` starts this formula when `analysis.formula` in `system-flow.yaml` names it, each time `analysis.schedule` comes round with the `analyze` host action on, in place of starting the analyzer, and with the focus empty. To start it with a focus yourself, run `flai formula run analysis-cadence --input focus=<focus> --ready`, the focus `bottlenecks`, `intent`, or `risk`. Step 2 runs from the project's main checkout, where the analyzer writes its report and files its issues; `latest` is the newest report under `design/analysis`.

- The `analyze` host action must be on for the role step: the analyzer runs under it.
- Record in the story's Notes the report step 2 printed and the issues it names, each with its title and class. The operator chooses at acceptance which of them become stories.
- A failing step blocks its task: the analyzer's run failed or asked, or step 2 found no report. Block the task with the output and ask the operator on a thread.
- This story's branch changes nothing. The report and the issues are the analyzer's, in the main checkout.

<!-- system-flow:end-of-baseline -->

## Project additions
