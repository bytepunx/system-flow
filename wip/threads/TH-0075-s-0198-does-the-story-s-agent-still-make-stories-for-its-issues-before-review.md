---
id: TH-0075
title: "S-0198: does the story's agent still make stories for its issues before review?"
anchor:
  path: wip/kanban/stories/S-0198-give-the-operator-the-option-to-have-all-issues-turned-into-stories.md
  item: S-0198
status: resolved
participants: [agent-S-0198, alex]
created: 2026-10-03T02:09:55Z
updated: 2026-10-03T02:12:10Z
---

# TH-0075 S-0198: does the story's agent still make stories for its issues before review?

On wip/kanban/stories/S-0198-give-the-operator-the-option-to-have-all-issues-turned-into-stories.md.

## Entries

### 2026-10-03T02:09:55Z agent-S-0198
The continuous-improvement baseline tells a story's agent to make a backlog story, at review, for each issue it recorded or bumped. Your answer on TH-0074 has the review page list those issues first, checked, so that you choose at Accept. If the agent had already made the stories, those issues would be linked and never shown.

Recommendation: change the baseline so the story's agent records each issue with `flai issue new` or `bump`, which now name the story, and makes no story for it. The operator chooses at Accept on the review page. The agent makes a story with `flai issue story` or the MCP tool `issue_story` only when the operator asks, or when `flai check` warns `issues.no-story` about an issue a week old. `harness.Prompt` would say the same.

The alternative is to keep the agent making the stories before review, so the review page shows only the issues the agent did not cover. I am writing the recommendation into T-0735 now, and will change it if you choose otherwise.

### 2026-10-03T02:12:01Z alex
take the recommended approach

### 2026-10-03T02:12:10Z agent-S-0198
Resolved: Confirmed; T-0735 (6dd9a71) says it in harness.Prompt, the continuous-improvement baseline and its copy, and the design.
