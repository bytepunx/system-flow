---
id: D-0296
title: "S-0380 decision: where the design system lives, how it is previewed, and how you and agents change it (10 choices): I would not like to force more things into a markdown document when it's not the…"
kind: recommendation
sources: [TH-0445]
topics: [design, flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [design/system/design-system-tool.md, design/system/README.md, design/adrs, design/adrs/README.md, design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md, design/issues/summary.md, design/issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md]
date: 2026-10-14
status: absorbed
absorbed_by: "design/system/design-system-tool.md#decision"
asked_again: []
cited: []
---

# D-0296 S-0380 decision: where the design system lives, how it is previewed, and how you and agents change it (10 choices): I would not like to force more things into a markdown document when it's not the…

The operator, alex, answered TH-0445 at 2026-10-14T20:30:34Z.

What was asked:

```text
**Recommendation:** take all ten choices below as recommended. Reply "all as recommended", or give the numbers you want otherwise.
```

The recommendation it overruled:

```text
**Recommendation:** reply "all as recommended". These are ten new design decisions, so they are yours, not mine.

I read the finding on `story/S-0380` (`design/system/design-system-tool.md`). Each choice follows from what flai and the dashboard do today:

- **1, 4, 8** rest on ADR-0023. `doc.save` takes only markdown, and it refuses any finding a save introduces. So tokens in markdown front matter work with no new host method. Contrast as a check finding would block every trial colour. A stylesheet left stale by a save would block the save too.
- **3 and 7:** the gallery is compiled with the dashboard, so it can only style role names it knows. A blank same-origin frame isolates it with nothing new in the container.
- **10:** the document's `topics` already decide which stories prime reaches. A second list of governed sub-projects in S-0386 could disagree with them.

**The one to weigh: choice 8.** A token save on main changes a code file, the UI's stylesheet, with no story and no close-out. ADR-0023 today commits only the one markdown path saved. If every code change must pass a story, pick the alternative: saves land only on a proposal's branch, and the page is for trying. I still recommend main. You review every change on this board anyway, and the stylesheet is generated, never hand-written.

**Choice 6:** a Google font is loaded only when the document names one. It sends a request from your browser to Google. Choose "system stacks only" if that matters to you.

Once you answer, S-0380's agent records the ADRs and names on TH-0438 what S-0381 to S-0386 must change. The planner then replans them, and I finalize those that come out consistent.

Source: design/adrs/0023-documents-are-saved-through-flai.md § Decision
```

The operator's answer:

```text
1. I would not like to force more things into a markdown document when it's not the obvious best choice simply because we currently only have API semantics for markdown documents. Introduce file.get and file.save as alternatives that we can leverage for different kinds of files.
```

Source: TH-0445.
