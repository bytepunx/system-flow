---
id: D-0674
title: "S-0810 left three limits of the state job's cluster steps: the kluster pass may deploy a story's pre-rebase head once, a checks step the pass never records waits out the job's two hours, and the review page labels a state job a close-out until S-0814"
kind: decision
sources: [S-0810]
topics: [kluster, workflow, verify, flai, go, code]
paths: [flai/internal/serve/statejob.go, flai/internal/serve/kluster.go]
date: 2026-11-21
status: active
absorbed_by: ""
asked_again: []
cited: [S-0835]
---

# D-0674 S-0810 left three limits of the state job's cluster steps: the kluster pass may deploy a story's pre-rebase head once, a checks step the pass never records waits out the job's two hours, and the review page labels a state job a close-out until S-0814

S-0810 made the cluster's deploy and checks the last steps of a state job, which rebases the story's branch first (ADR-0312 §4). The kluster pass still runs them on its own look, and the job waits for what the pass records at the head the job left. Its narrative's Decisions (`wip/archive/agents/S-0810.md`) record what it left as it is:

> Known limits, left as they are: with a cluster configured, the baseline's entry into review now rebases the branch before the deploy (ADR-0312 §4); the kluster pass may deploy the pre-rebase head once when its look precedes the state job's start; a checks step whose run the pass never records waits out the job's limit; the dashboard's review page labels a state job as a close-out until S-0814 shows the state job.

The first is the design now (`design/system/kluster.md` § Where the cluster runs). The other three are not:

- **A pre-rebase deploy.** When the pass's look comes before the job starts, it deploys the head before the rebase, then deploys again at the rebased head.
- **A checks run never recorded.** The job's checks step waits until the job's limit, two hours, the close-out's (`design/system/flai-cli.md`, the state job).
- **The review page's label.** It shows a state job as a close-out until S-0814 shows the state job.

A story that changes how the kluster pass and the state job meet, such as S-0814 or S-0835, starts from these.

Source: S-0810's narrative, § Decisions.
