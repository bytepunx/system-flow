---
id: T-2749
type: task
nature: remediation
title: An ADR and the design say a permission_prompt check that ends on the harness's usage limit records no outcome and is checked again after the limit resets
status: done
parent: S-0629
owner: alex
created: 2026-11-09T17:08:54Z
updated: 2026-11-09T17:27:13Z
transitions:
  - to: ready
    at: 2026-11-09T17:25:15Z
    by: agent-S-0629
  - to: in-progress
    at: 2026-11-09T17:25:16Z
    by: agent-S-0629
  - to: done
    at: 2026-11-09T17:27:13Z
    by: agent-S-0629
stream: S-0629
tags: [flai, serve, claude-code, adr]
touches: [design/adrs/README.md, design/system/agent-adapters.md, design/system/flai-cli.md, design/adrs/0270-a-permission-prompt-check-whose-claude-p-ends-on-the-harness-s-usage-limit.md]
usage:
  source: log
  seconds: 117
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 32
      output: 289
      cache_read: 4192324
      cache_write: 19037
      cost: 1.5372
---
# T-2749 An ADR and the design say a permission_prompt check that ends on the harness's usage limit records no outcome and is checked again after the limit resets

## Work

Write the remedy I-0296 lacks before the code, because it changes a rule of an accepted ADR. ADR-0106 § 2 says "A version with a record is not checked again, whatever its outcome", and ADR-0106 can't be edited. Write a new ADR with `flai adr new` that refines ADR-0106 § 2:

- A check whose `claude -p` ends on the harness's usage limit has no outcome. The limit is a 429 result or a rejecting `rate_limit_event`, as `flai/internal/usage/claudecode.go` reads it for ADR-0214. flai records nothing for that version and provider and opens no thread. It logs that the check was refused on the limit, with the reset time when Claude Code gave one.
- The next look after the reset checks that version again. A refusal that names no reset time is checked again an hour later, so a limit with no reset does not start a refused `claude -p` at every look. While flai serve runs, the in-process gate in `consider` and `considerImage` lets the version be checked again once the reset has passed.
- A `failed` record is still final for every other cause, as ADR-0106 decided.

Then update the sections that describe the record so they say the same and link the ADR:

- `design/system/agent-adapters.md` § Permissions
- `design/system/flai-cli.md` § Commands, where it says removing the entry has it checked again
- `design/system/flai-cli.md` § The container harness (S-0508)

The new ADR file's number isn't known until `flai adr new` writes it, so this task names `design/adrs/README.md`, the index, and `task_done` widens the touches to the new file.

## Done when

- A new accepted ADR refines ADR-0106 § 2 as above. It is indexed in `design/adrs/README.md`, and ADR-0106 is unchanged.
- `design/system/agent-adapters.md` § Permissions and both `design/system/flai-cli.md` sections say that a usage-limit refusal records no outcome and when the version is checked again. Each links the new ADR.
- `flai check --strict` and the markdown lint pass.

## Notes
