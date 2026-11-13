---
id: S-0757
type: story
nature: feature
title: The dashboard shows each agent image job as it runs, its output, and a failure's cause, and offers the operator a retry
status: backlog
parent: E-0036
owner: alex
created: 2026-11-12T16:03:25Z
updated: 2026-11-13T07:28:14Z
transitions: []
tags: [cli, dashboard]
topics: [docker, dashboard]
touches: [flai/internal/hostapi/images.go, flai/internal/hostapi/images_test.go, flai/internal/hostapi/contract_test.go, flai/internal/serve/imagejob.go, flai/internal/serve/imagejob_test.go, flaiover/src/routes/api/images/+server.ts, flaiover/src/routes/api/images/retry/+server.ts, flaiover/src/routes/ops/images/+page.svelte, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/ImageJob.svelte, flaiover/src/lib/components/ImageJob.svelte.test.ts, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, design/components/flaiover/src/lib/components/ImageJob.examples.json, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0751]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 8
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 7
          cache_read: 61934
          cache_write: 11101
          cost: 0.0208
forecast:
  duration: 50m
  delivery: 2026-11-13T20:37:00Z
  basis: "Its own forecast of 50m; 31st in the pull order with an in-progress limit of 3, behind S-0659, S-0241, S-0660, S-0661, S-0662, S-0663, S-0730, S-0733, S-0734, S-0735, S-0736, S-0737, S-0738, S-0739, S-0740, S-0741, S-0742, S-0743, S-0744, S-0745, S-0746, S-0747, S-0748, S-0749, S-0750, S-0751, S-0752, S-0753, S-0754, S-0755 and S-0756."
  by: flai
  at: 2026-11-13T07:23:41Z
finalized:
  by: alex
  at: 2026-11-13T07:28:14Z
---
# S-0757 The dashboard shows each agent image job as it runs, its output, and a failure's cause, and offers the operator a retry

## Goal

When an agent start fails today, the dashboard shows the refusal's text on the story's agent (`StoryAgent.svelte`), and nothing more. After S-0751, a missing image is brought back by a host job that may take minutes. The operator should see it while it runs, read its output when it fails, and start it again without a shell.

After this story the dashboard shows, under Ops, the project's agent images and their state from `flai image status`, and each image job: running or ended, its way, its step, and its output as it is written. A failed job shows its cause, the failing step, and its output, with a retry. A run held on the job shows that it waits for its image, linked to the job, in place of a failure.

## Acceptance criteria

- [ ] A host API read gives the project's agent images with their state, and the image jobs, running and the last ended, each with its id, its images, its way, its current step, its start and end, and its outcome. A log read gives a job's output from an offset, as `kluster.log` does for a deploy's step
- [ ] A host API write starts a new image job for the project, behind the `agent` host action. It is refused while a job for the same images runs, and it says why
- [ ] An Ops page lists the images and the jobs, follows a running job's output live, and shows a failed job's cause, step, and output with a Retry button that calls the write
- [ ] A story's agent held on an image job shows "waiting for its image" with a link to the job, and a run failed by a job shows the cause and a link, in place of the bare refusal
- [ ] flaiover's tests cover the page, the component, and the held run, and the component's examples file passes `components.examples`
- [ ] `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe the page, the retry, and the host action that gates it

## Tasks
- T-3181 The host API reads the project's agent images and image jobs, reads a job's log from an offset, and starts a new job behind the agent host action
- T-3182 flaiover's API routes read the agent images, the jobs, and a job's log, and post a retry
- T-3183 An ImageJob component shows one image job's state, step, live output, and a failure's cause with a Retry button
- T-3184 A story's agent held on an image job shows that it waits for its image, and one failed by a job shows the cause, each linked to the job
- T-3185 An Ops › Agent images page lists the project's images and jobs, follows a running job, and retries a failed one
- T-3186 The dashboard's design and user guide describe the Agent images page, a held run, and the retry behind the agent host action

## Notes

Proposed slice 4b of E-0036. The log read follows `kluster.log` (ADR-0263), and the page's live log reuses `LogFeed.svelte`.

### Planning

- Touches:
  - `flai/internal/hostapi/images.go`, `images_test.go`: new methods, registered as `kluster.go` registers its own (line 1497). `contract_test.go` lists every method.
  - `flai/internal/serve/imagejob.go`, `imagejob_test.go`: S-0751's runner gives what the read needs.
  - `flaiover/src/routes/api/images/+server.ts`, `retry/+server.ts`: new routes, modelled on `api/kluster/log` and `api/kluster/deploy`.
  - `flaiover/src/routes/ops/images/+page.svelte`, `flaiover/src/lib/sitemenu.ts`, `sitemenu.test.ts`: a new Ops page beside Klusters and Deployments, which the menu lists.
  - `flaiover/src/lib/components/ImageJob.svelte`, its test, and `design/components/flaiover/src/lib/components/ImageJob.examples.json`: a new component, whose examples `flai check` requires (ADR-0164).
  - `flaiover/src/lib/components/StoryAgent.svelte`, its test: layout. The refusal is shown there today (`act`, lines 115 to 128).
  - `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`: named in the last criterion.
  - Left out from `flai touches suggest`: `design/system/flai-cli.md` (40%) and `docs/users/flai.md` (39%). No command changes, and the host methods are the dashboard's design.
- `after` S-0751: this story reads the image job, its record on the run in `agents.json`, and its log file, which that story writes.
- Tasks, in four layers: T-3181 (the host API) first; T-3182, T-3183, and T-3184 share no file and run together; T-3185 uses T-3182's route and T-3183's component; T-3186 describes the page and the held run.
- Forecast 50m: `flai forecast` gives 40m (109 s per unit over 118 done large feature stories, size 22). Raised for six tasks across Go and flaiover, each with its own vitest and svelte-check cycle. Delivery is as flai gives it.
- Cost of delay: not set yet. E-0036 has no inputs; TH-1094 asks the operator for them.
