---
id: S-0435
type: story
nature: experiment
title: "Experiment: does priming with the digest cut the threads opened to the operator per completed story on this repository, two weeks before against two weeks after"
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T06:00:41Z
updated: 2026-10-16T22:44:54Z
transitions: []
tags: [cli]
topics: [analysis, planning]
touches: [design/experiments/README.md, design/experiments/S-0435-digest-priming-and-threads-to-the-operator.md, design/issues/summary.md]
after: [S-0434]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 25m
  delivery: 2026-10-17T22:13:00Z
  basis: "Its own forecast of 25m; 42nd in the pull order with an in-progress limit of 5, behind S-0436, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0428, S-0429, S-0430, S-0431, S-0432, S-0433 and S-0434."
  by: flai
  at: 2026-10-16T22:44:54Z
---
# S-0435 Experiment: does priming with the digest cut the threads opened to the operator per completed story on this repository, two weeks before against two weeks after

## Goal

Find out, on this repository, whether the digest does what E-0025 built it for. The hypothesis is that once the digest primes (S-0430's acceptance), fewer threads per completed story are opened to the operator, and fewer of them ask what an entry already answered. The measure is S-0434's, read from `flai stats --json` over the two weeks before S-0430 was accepted and the two weeks after, with the results document under `design/experiments/`.

## Acceptance criteria

- [ ] The results document, `design/experiments/S-0435-digest-priming-and-threads-to-the-operator.md` from `design/experiments/template.md`, states the hypothesis, the two windows by date, and the method: `flai stats --json` with each window, the `digest.weeks[]` figures, and the waits `waiting.weeks[].threads` reports, as flai gave them.
- [ ] It reports for each window the threads opened to the operator, the stories completed, the threads per completed story, the threads that asked what an entry already answered, and the thread waits' total and mean, in one table.
- [ ] It names what else changed between the windows that could explain a difference (other epics accepted, the orchestrator's `answer_threads` setting, the stories' mix by nature), and says whether the result supports the hypothesis, refutes it, or cannot tell, and why.
- [ ] It lists three threads of the after window that an entry answered and the entry each cited, or says that none did, and the entries the curator absorbed or compacted in the window.
- [ ] It ends with a recommendation: keep the digest as it is, change it (naming what, as an issue filed with `flai issue new` and linked), or drop it; and `design/experiments/README.md` lists the document.

## Tasks
- T-1847 The results document states the hypothesis and the two windows and gives the measured figures for each, as flai stats reports them
- T-1848 The results document weighs what else changed, cites threads an entry answered, says whether the hypothesis holds, and recommends keeping, changing, or dropping the digest

## Notes

- Drafted by the planner for E-0025: the epic's last story.
- It cannot start before 14 days have passed since S-0430's acceptance, since the after window must be whole. The operator promotes it no earlier than that; T-1847 blocks the story if it starts sooner.
- An experiment is accepted like any other story and releases nothing (ADR-0066).

### Planning

- **Touches.** The story declared none. Both tasks write the results document, named here so that no folder is claimed; T-1848 adds `design/experiments/README.md` and `design/issues/summary.md`, which `flai issue new` writes when the recommendation files a change.
  - **From design:** `design/experiments/template.md` and the folder's README, as the documentation standard asks of an experiment story.
  - **Co-change:** `flai touches suggest` found nothing changed often enough with `design/experiments/README.md` (one commit).
  - **No folder touch.** The issue file a recommendation may add is in `design/issues`, which is in `claims.shared`, and the close-out commits it.
- **Forecast.** Raised from flai's 13m to 25m. flai's figure is the median over all 236 done stories, since few experiments are done; reading two windows' stats, the manifest's history, and three threads into an argued document is more than the size, 5 criteria and 3 touches, counts. Delivery is no earlier than 14 days after S-0430's forecast delivery (2026-10-17T02:13Z), when the after window is whole, plus the duration.
- **Cost of delay.** Not set. Neither the story nor E-0025 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0476, question 1, and the value follows its answer.
