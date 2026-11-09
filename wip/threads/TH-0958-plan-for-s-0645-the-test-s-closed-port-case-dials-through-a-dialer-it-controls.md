---
id: TH-0958
title: "Plan for S-0645: the test's closed-port case dials through a dialer it controls, two tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0645-testtcpanswersalistenerandnotaclosedone-fails-on-every-run-on-the-operator-s.md
  item: S-0645
status: resolved
participants: [system-flow, agent-S-0645]
created: 2026-11-09T22:11:25Z
updated: 2026-11-09T22:54:48Z
---

# TH-0958 Plan for S-0645: the test's closed-port case dials through a dialer it controls, two tasks in two layers

On wip/kanban/stories/S-0645-testtcpanswersalistenerandnotaclosedone-fails-on-every-run-on-the-operator-s.md.

## Entries

### 2026-11-09T22:11:25Z system-flow
Plan for S-0645: two tasks in two layers, 20m, 30.00 USD a week.

Proposed solution: fix the test, not the probe. D-0444 records that this host runs WSL2 with mirrored networking. I-0321's second instance says the failure comes and goes. A run of `go test -count=10` on main passed 10 of 10 at about 22:08 UTC. So a loopback port that was free in Linux is answered some of the time, likely by the relay to a Windows listener. A real closed-port dial tests the host, not flai's `tcpAnswers`.

Tasks:

1. T-2816 (layer 1): confirm the cause on the host with repeated runs and `ss -ltnp`. Then give `tcpAnswers` and `publishAnswers` a dial seam, so the closed-port cases go through a dialer that refuses and only a listener the test holds open is dialled for real. Production keeps `net.DialTimeout`.
2. T-2817 (layer 2, after T-2816): close I-0321 with a reason that quotes T-2816's finding and fix.

Figures: forecast 20m, flai's 15m plus 5m for confirming an intermittent cause. Cost of delay 30.00 USD a week from the 12m input, unchanged.

Assumptions:

- Criterion 1's reproducing test does not fit: no test can make the host's relay. T-2816 records that in the narrative, and the test it writes asserts through a dialer that accepts every address, as the relay does.
- The production probe stays as it is. The dashboard's port-held state (ADR-0267) may also be answered by the relay on this host and hold a restart wrongly. That is a separate defect if it is seen; I would not add it to this story.
- D-0444's advice still holds for tests that need only a failed HTTP request. The curator may want to narrow it: it does not hold for a test that asserts a bare TCP connect is refused.

### 2026-11-09T22:54:38Z agent-S-0645
Taken as drafted. T-2816 is done at 2b6821420: `tcpAnswersWith` and `publishAnswersWith` take the dialer, production keeps `net.DialTimeout`, and the test dials for real only a listener it holds open. `.wslconfig` confirms `networkingMode=mirrored`; 300 bind-close-connect probes found no closed port answering at 22:45 UTC, so the failure depends on which ephemeral port the test draws. T-2817 closes I-0321 next. The production probe stays out of scope, as you propose.

### 2026-11-09T22:54:48Z agent-S-0645
Resolved: Plan taken as drafted; T-2816 done at 2b6821420, T-2817 closes I-0321.
