---
id: T-2955
type: task
nature: remediation
title: scripts/agent-image.sh builds with an anonymous Docker configuration when the configured credential helper fails
status: done
parent: S-0713
owner: alex
created: 2026-11-10T10:23:47Z
updated: 2026-11-10T10:54:16Z
transitions:
  - to: ready
    at: 2026-11-10T10:53:50Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T10:53:51Z
    by: system-flow
  - to: done
    at: 2026-11-10T10:54:16Z
    by: agent-S-0713
stream: S-0713
tags: [scripts, docker]
touches: [scripts/agent-image.sh]
usage:
  source: log
  seconds: 24
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 8
      output: 35
      cache_read: 543188
      cache_write: 5841
      cost: 0.2021
---
# T-2955 scripts/agent-image.sh builds with an anonymous Docker configuration when the configured credential helper fails

## Work

Before its first `docker build`, in the default run and under `--layer`, `scripts/agent-image.sh` reads the configuration Docker would use: `$DOCKER_CONFIG/config.json`, else `$HOME/.docker/config.json`. When that file names a `credsStore`, the script runs `docker-credential-<store> list` under a short timeout. When the helper is not on `PATH`, exits non-zero, or does not answer in time, the script exports `DOCKER_CONFIG` as a scratch folder (`mktemp -d`, removed on exit) holding a `config.json` of `{}`, and prints one line to standard error. That line names the helper and says the public base images are pulled with no login. When the original configuration folder has a `cli-plugins` folder, link it into the scratch one, so that `docker build` still finds buildx there. Change nothing in the operator's configuration. A helper that answers, or no `credsStore`, leaves the build as it is. `--tag` and `--check` build nothing and probe nothing. This is the first layer: it waits for no task.

## Done when

- With a `credsStore` whose helper fails, the script's builds run with `DOCKER_CONFIG` naming a configuration of `{}`, and it says so in one line.
- With a helper that answers, or none named, `DOCKER_CONFIG` is left as it was.
- `~/.docker/config.json` is never written.
- ShellCheck passes on the script.

## Notes
