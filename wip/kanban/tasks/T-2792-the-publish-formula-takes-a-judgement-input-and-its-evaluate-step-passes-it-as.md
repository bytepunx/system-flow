---
id: T-2792
type: task
nature: improvement
title: The publish formula takes a judgement input and its evaluate step passes it as --judged
status: done
parent: S-0625
owner: alex
created: 2026-11-09T22:02:27Z
updated: 2026-11-10T05:50:37Z
transitions:
  - to: ready
    at: 2026-11-10T05:45:00Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T05:45:00Z
    by: system-flow
  - to: done
    at: 2026-11-10T05:50:37Z
    by: agent-S-0625
stream: S-0625
tags: [flai, template]
touches: [design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, flai/cmd/publish_judgement_test.go]
after: [T-2790]
usage:
  source: log
  seconds: 337
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 42
      output: 184
      cache_read: 6318540
      cache_write: 35196
      cost: 2.3348
---
# T-2792 The publish formula takes a judgement input and its evaluate step passes it as --judged

## Work

Give the publish formula the input `judgement`, as T-2788's ADR decides. It waits for T-2790 because its step runs the `--judged` flag that T-2790 adds.

- In `design/workflows/formulas/publish.md` and `template/root/design/workflows/formulas/publish.md`, above the marker:
  - Declare the input `judgement`, with a description and `default: ""`.
  - Step 2's `run` becomes `release --evaluate --require-met --judged "{{judgement}}"`. Keep `stop: [2]`.
  - The body says that under `judgement` the orchestrator gives its reasoning as the input, and that with none step 2 stops as it does today.
- `flai/cmd/publish_judgement_test.go`, new: the test that reproduces I-0291 end to end. In a test project with `orchestration.release.policy: judgement` and an accepted story not yet released, it reads the real formula file:
  - `flai formula run publish --ready` with no input makes an instance whose step 2 command exits 2. This is the cause.
  - With `--input judgement=<reason>`, step 2's filled command exits 0, so step 3 runs.
  - The instance's step 2 task names the reason in its `## Work`.
- `TestRenderPrototypeTemplate` must still pass with the `{{judgement}}` placeholder in the template copy. S-0616 changes how that test treats `{{`. Check that it passes on the branch as it stands after a sync.

## Done when

- Both formula copies match above the marker. `flai formula show publish` lists the input and the filled step 2.
- `flai/cmd/publish_judgement_test.go` passes: with no input, step 2 exits 2; with a reason, it exits 0.
- `flai check --strict` is clean, and `TestRenderPrototypeTemplate` passes.

## Notes

Layer 3. It runs alone.
