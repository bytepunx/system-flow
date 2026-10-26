---
id: T-2130
type: task
nature: improvement
title: flai serve notify keeps a reminder channel, a webhook URL and a host command, in the host configuration, and a sender posts a reminder's JSON to it
status: done
parent: S-0498
owner: alex
created: 2026-10-26T04:14:30Z
updated: 2026-10-26T04:49:32Z
transitions:
  - to: ready
    at: 2026-10-26T04:41:15Z
    by: agent-S-0498
  - to: in-progress
    at: 2026-10-26T04:41:16Z
    by: agent-S-0498
  - to: done
    at: 2026-10-26T04:49:32Z
    by: agent-S-0498
stream: S-0498
tags: [flai]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/serve/notify.go, flai/internal/serve/notify_test.go, docs/operators/settings.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/serve.go]
after: [T-2128]
usage:
  source: log
  seconds: 490
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 108
      output: 569
      cache_read: 5503416
      cache_write: 145359
      cost: 1.9112
---
# T-2130 flai serve notify keeps a reminder channel, a webhook URL and a host command, in the host configuration, and a sender posts a reminder's JSON to it

## Work

The channel the ADR of T-2127 names, kept where the host's providers and forges are:

- `flai/internal/config/config.go`: a `Notify` struct on `Config` (`json:"notify,omitempty"`) with `URL` and `Command`, both optional strings; `WithNotify` and `WithoutNotify` that copy and validate as `WithProvider` and `WithForge` do (the URL must be `http` or `https`; the command non-empty when given); `Parse` keeps rejecting unknown keys, so the field is declared. `config.Keys()` leaves it out as it leaves the other managed maps. Tests in `config_test.go`.
- `flai/cmd/serve_actions.go`: `flai serve notify set --url <url> --command <command>`, `flai serve notify show`, and `flai serve notify clear`, beside `flai serve agent provider` and `flai serve forge`; `--json` on `show`. Tests in `serve_actions_test.go`. The command runs only from a shell: `flai guard` already refuses a served agent `flai serve` writes, so no new guard rule.
- `flai/internal/serve/notify.go`: the sender. `Reminder` is the JSON body the ADR names (`project`, `thread`, `title`, `item`, `awaiting_seconds`, `since`, `recommendation` with `author`, `first_line`, and `source` when one is pending, `href`); `Send(ctx, cfg.Notify, Reminder)` POSTs the body to the URL (`content-type: application/json`, `user-agent: flai`, one attempt, a 5 s timeout, a non-2xx is an error) and runs the command through the host's shell (`sh -c` on Unix, as `fix_command` is run) with the body on standard input and a 30 s timeout, each when set, and returns both errors. It sends no token and no file content. Tests in `notify_test.go` with an `httptest` server and a command that writes its standard input to a file.
- `docs/operators/settings.md`: rows for `notify.url` and `notify.command` under `## flai configuration`, so `TestSettingsIndexIsComplete` passes.
- `docs/users/flai-reference.md`: regenerate with `make flai-reference` for the new sub-commands; `docs/users/flai.md`: `flai serve notify` under `### flai serve`, with an example for a push service (`curl` to an `ntfy` topic) and one for a desktop toast.

The webhook body is the shape `flaiover`'s `notify.ts` posts to `dashboard.notify_url`, extended, so one receiver may take both.

## Done when

- `flai serve notify set --url https://ntfy.example/flai` then `show` prints the channel, `clear` removes it, and `set --url ftp://x` is refused with the reason.
- `Send` posts the body once to the URL and pipes it to the command, each tested, and returns the error of a URL that answers 500 or a command that exits non-zero.
- `flai test` on the changed paths passes, `TestSettingsIndexIsComplete` and the reference's staleness test among them.

## Notes

Layer 3: waits for T-2128 because both add rows to `docs/operators/settings.md`. T-2131 (the reminder loop) waits for it and T-2128.
