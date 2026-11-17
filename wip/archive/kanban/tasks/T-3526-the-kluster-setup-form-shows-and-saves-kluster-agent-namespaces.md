---
id: T-3526
type: task
nature: remediation
title: The kluster setup form shows and saves kluster.agent_namespaces
status: done
parent: S-0830
owner: alex
created: 2026-11-17T08:28:58Z
updated: 2026-11-17T10:08:21Z
transitions:
  - to: ready
    at: 2026-11-17T10:06:23Z
    by: agent-S-0830
  - to: in-progress
    at: 2026-11-17T10:06:23Z
    by: agent-S-0830
  - to: done
    at: 2026-11-17T10:08:21Z
    by: agent-S-0830
stream: S-0830
tags: [flaiover, kluster]
touches: [flaiover/src/lib/components/DeploySettings.svelte, flaiover/src/lib/components/DeploySettings.svelte.test.ts, flaiover/src/lib/server/repo.ts, design/components/flaiover/src/lib/components/DeploySettings.examples.json]
usage:
  source: log
  seconds: 118
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 3256
      cache_read: 1748588
      cache_write: 53774
      cost: 0.8451
---
# T-3526 The kluster setup form shows and saves kluster.agent_namespaces

## Work

flai already lists `kluster.agent_namespaces` in `kluster.status`'s `settings`, as kind `number`, whole, default `2`, its value a JSON number (`manifest.Settings` in `flai/internal/manifest/settings.go`, `klusterSettings` in `flai/internal/hostapi/kluster.go`). The dashboard drops it in three places:

- `PLAIN_KEYS` in `flaiover/src/lib/components/DeploySettings.svelte` does not name it, so neither the form nor `change()` reaches it.
- `typed()` turns any value that is not a string or a command into an empty string, so a number would show as empty.
- `KlusterSetting` in `flaiover/src/lib/server/repo.ts` has no `number` kind and no number value.

Add `kluster.agent_namespaces` to `PLAIN_KEYS` after `kluster.timeout`. Give `KlusterSetting` the `number` kind and a number value. Make `typed()` show a number as its digits. Render the field as a numeric input (`inputmode="numeric"`) with the default as its placeholder. Save it as the other plain keys are saved: the trimmed text when changed, unset when emptied. flai's `settings.manifest` parses the text and refuses what is not a whole number of 0 or more; the form shows that refusal beside the field, as it does for the others.

Write a test in `DeploySettings.svelte.test.ts` that fails before the change. It gives the form a `kluster.agent_namespaces` setting of value `3`, finds `field-kluster.agent_namespaces` showing `3`, changes it to `1`, saves, and asserts the `POST /api/settings` body sets `kluster.agent_namespaces` to `"1"` and nothing else. A second case empties it and asserts it is unset, and a third shows flai's refusal beside the field. Add the setting to the examples in `DeploySettings.examples.json`, so `components.examples` stays in step with the props.

## Done when

- The form on the Klusters page's step 4 and the Deployments page's Setup pane shows `agent_namespaces` with its value, its meaning, and its default.
- A change to it is saved through `settings.manifest`, and an emptied field unsets it.
- The new tests pass, and fail with `kluster.agent_namespaces` taken out of `PLAIN_KEYS`.
- `flai test` passes on the changed paths, svelte-check and the examples check included.

## Notes
