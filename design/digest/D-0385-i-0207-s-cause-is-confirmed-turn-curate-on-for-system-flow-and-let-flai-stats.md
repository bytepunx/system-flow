---
id: D-0385
title: "I-0207's cause is confirmed: turn curate on for system-flow, and let flai stats tell an unmeasured week from zero repeats: Recommended for the operator: take agent-S-0504's recommendation and fix…"
kind: answer
sources: [TH-0649, S-0504]
topics: [digest, metrics, flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [design/adrs, design/adrs/README.md, design/system/metrics.md, flai/internal/metrics/digest.go, flai/internal/metrics/digest_test.go, flai/internal/metrics/metrics.go, flai/cmd/stats.go, flai/cmd/stats_test.go, flai/cmd/check_stats_test.go, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/explain.ts, flaiover/src/lib/viz/explain.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, design/issues/I-0207-the-curator-has-never-run-on-system-flow-so-no-thread-is-in-any-entry-s-asked-again-and-the-digest-s-repeat-measure-reads-0-whatever-was-asked.md, design/issues/summary.md, template/template.yaml]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0385 I-0207's cause is confirmed: turn curate on for system-flow, and let flai stats tell an unmeasured week from zero repeats: Recommended for the operator: take agent-S-0504's recommendation and fix…

The operator, alex, answered TH-0649 at 2026-10-28T12:02:08Z.

What was asked:

```text
**Cause confirmed on the host.** `~/.flai/config.json` has `host_actions.curate` set for `/home/alex/git/bytepunx/signet` only. `wip/agents/curator.md` does not exist on main. So no curator has run here. No entry lists a thread in `asked_again`, and every week reads `asked_again 0`.
```

The answer:

````text
**Recommended for the operator:** take agent-S-0504's recommendation and fix both causes. Run these on the host, in the system-flow checkout:

```sh
flai serve enable curate
flai curate
```

Then reply here once the curator has run.

- **Why it needs you.** A host action is the operator's to turn on, and an agent that `flai serve` started cannot start the curator. The question is a host change, so I escalate it rather than answer it.
- **Why `flai curate` now.** ADR-0199 § Decision makes `curate` start a run only on an acceptance seen while the action is on, or on `digest.schedule`. An acceptance made while `curate` is off starts nothing, then or later. Without a manual run, the first run waits for the next acceptance after you enable it.
- **Why the code half too.** Without it, every project whose curator has not yet run reads `asked_again 0` as if measured. The ADR the story writes keeps that 0 from being read as a measurement.

The story builds item 2 meanwhile. Only T-2122 and T-2126 wait on your reply.

Source: design/adrs/0199-the-curator-runs-after-acceptances-coalesced-and-on-digest-schedule-marks.md § Decision
````

The operator, alex, answered TH-0649 again at 2026-10-28T12:13:40Z.

The answer:

```text
let the next acceptance trigger the curator
```

Source: TH-0649.

The outcome: S-0504 closed I-0207 on the code half, ADR-0204's measure, and left the curator's first run to I-0212, which records the run of 2026-10-26 that failed on the usage limit.

Source: S-0504.
