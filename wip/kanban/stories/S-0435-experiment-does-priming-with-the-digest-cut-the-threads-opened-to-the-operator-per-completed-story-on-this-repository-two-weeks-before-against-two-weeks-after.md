---
id: S-0435
type: story
nature: experiment
title: "Experiment: does priming with the digest cut the threads opened to the operator per completed story on this repository, two weeks before against two weeks after"
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T06:00:41Z
updated: 2026-10-16T06:00:41Z
transitions: []
tags: [cli]
topics: [analysis, planning]
after: [S-0434]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0435 Experiment: does priming with the digest cut the threads opened to the operator per completed story on this repository, two weeks before against two weeks after

## Goal

Find out, on this repository, whether the digest does what E-0025 built it for. The hypothesis is that once the digest primes (S-0430's acceptance), fewer threads per completed story are opened to the operator, and fewer of them ask what an entry already answered. The measure is S-0434's, read from `flai stats --json` over the two weeks before S-0430 was accepted and the two weeks after, with the results document under `design/experiments/`.

## Acceptance criteria

- [ ] The results document, `design/experiments/S-0435-<slug>.md` from `design/experiments/template.md`, states the hypothesis, the two windows by date, and the method: `flai stats --json` with each window, the `digest.weeks[]` figures, and the waits `waiting.weeks[].threads` reports, as flai gave them.
- [ ] It reports for each window the threads opened to the operator, the stories completed, the threads per completed story, the threads that asked what an entry already answered, and the thread waits' total and mean, in one table.
- [ ] It names what else changed between the windows that could explain a difference (other epics accepted, the orchestrator's `answer_threads` setting, the stories' mix by nature), and says whether the result supports the hypothesis, refutes it, or cannot tell, and why.
- [ ] It lists three threads of the after window that an entry answered and the entry each cited, or says that none did, and the entries the curator absorbed or compacted in the window.
- [ ] It ends with a recommendation: keep the digest as it is, change it (naming what, as an issue filed with `flai issue new` and linked), or drop it; and `design/experiments/README.md` lists the document.

## Tasks

## Notes

- Drafted by the planner for E-0025: the epic's last story.
- It cannot start before 14 days have passed since S-0430's acceptance, since the after window must be whole. The operator promotes it no earlier than that; the planner sets no date.
- An experiment is accepted like any other story and releases nothing (ADR-0066).
