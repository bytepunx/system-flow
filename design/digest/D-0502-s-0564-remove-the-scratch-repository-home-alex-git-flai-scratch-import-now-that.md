---
id: D-0502
title: "S-0564: remove the scratch repository /home/alex/git/flai-scratch-import now that S-0001 is accepted: done"
kind: answer
sources: [TH-0863]
topics: [orchestration, cli, planning, template, conventions, all]
paths: [wip/kanban/stories/S-0564-the-import-and-analysis-cadence-formulas-are-each-run-once-in-this-repository.md, design/issues/summary.md, design/issues/I-0282-a-formula-s-instance-declares-no-touches-so-the-story-that-started-it-holds-it.md, design/issues/I-0300-a-story-already-in-ready-when-the-agent-host-action-is-turned-on-gets-no-agent.md, design/issues/I-0302-flai-import-writes-no-system-flow-lock-yaml-so-the-import-formula-s-flai-check.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0502 S-0564: remove the scratch repository /home/alex/git/flai-scratch-import now that S-0001 is accepted: done

The scratch repository S-0564 made for its import run, `/home/alex/git/flai-scratch-import`, is gone: the operator removed it on their own word once both runs were accepted, S-0001 (import) at `049cf36` there and S-0584 (cadence) at `07fa45be7` here. A later import run, such as a check of I-0300 or I-0302 once fixed, makes a scratch repository of its own below `/home/alex/git`, a folder named for import, rather than reusing this one.

The operator, alex, answered TH-0863 at 2026-11-07T08:53:53Z.

What was asked:

```text
**Question:** may the scratch repository `/home/alex/git/flai-scratch-import` be removed? Its Notes say to remove it only on your word.

**Recommended:** remove it yourself, since it is outside this repository and flai serve serves it as a project. First turn off its host actions, then delete the folder.

**Alternative:** keep it as a fixture for the next import run, such as a check of I-0300 or I-0302 once fixed.
```

The answer:

```text
done
```

Source: TH-0863, on S-0564.
