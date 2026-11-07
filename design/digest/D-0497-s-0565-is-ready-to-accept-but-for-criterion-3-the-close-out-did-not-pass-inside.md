---
id: D-0497
title: "S-0565 is ready to accept but for criterion 3: the close-out did not pass inside the container (smoke failed, I-0281): accept"
kind: answer
sources: [TH-0851, S-0565]
topics: [harness, docker, all]
paths: [design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md, design/experiments/README.md, design/issues, design/issues/summary.md, design/issues/I-0278-after-a-stop-from-the-dashboard-the-story-s-agent-page-shows-the-agent-stopped.md, design/issues/I-0279-scripts-agent-container-sh-always-runs-flai-agent-claude-code-local-and-ignores.md, design/issues/I-0280-scripts-close-out-sh-takes-no-fresh-so-at-a-commit-the-verify-record-holds-a.md, design/issues/I-0281-the-system-flow-image-layer-lacks-curl-and-goreleaser-so-the-smoke-tier-s.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [S-0565, TH-0851]
---

# D-0497 S-0565 is ready to accept but for criterion 3: the close-out did not pass inside the container (smoke failed, I-0281): accept

The operator, alex, answered TH-0851 at 2026-11-07T07:51:07Z: S-0565 was accepted with criterion 3 ticked, though the smoke tier failed inside the container.

What was asked:

```text
**Recommended:** accept S-0565 as it is. Write "accept" here and I accept it at `3a2bca9b4d5a`, or run `flai accept S-0565 --by alex` yourself. The experiment found what it set out to find, and it recorded the one gap as I-0281.
```

The orchestrator had held the acceptance because criterion 3 says `scripts/close-out.sh` "runs every tier the manifest's `tests` select and passes inside the container", and smoke failed there at `scripts/install-test.sh`, since the layer has no curl or GoReleaser (I-0281). The tick rested on the operator's entry on TH-0836 at 06:54:03Z, "we can mark our criteria here as completed", written before the in-container run.

The answer:

```text
accept
```

S-0565's narrative (`wip/archive/agents/S-0565.md`, Decisions) records the same choice from the story's side: "the operator wrote on TH-0836 at 06:54Z that the criteria could be marked complete, and ticked criterion 3 himself at 07:00:23Z. My in-container run then found that smoke fails inside the layer for want of curl and GoReleaser (I-0281); every other tier S-0556's close-out selected passed inside. Following the operator: criterion 3 stays ticked, and the results document says plainly what passed and what did not, so the operator sees both at acceptance."

So an experiment whose criterion the operator ticked is accepted on that tick when its results document states the gap and an issue records it.

Sources: TH-0851, S-0565.
