---
title: Two builds of authstar side by side in one kluster, a namespace per story
updated: 2026-11-06
status: active
story: S-0575
---

# Two builds of authstar side by side in one kluster, a namespace per story

E-0032 wants a story in review deployed to a local kluster for its integration or smoke tier. This experiment tried, by hand on authstar, the shape its outcome guesses at: one shared cluster, one namespace per story, and each story's own settings over the shared signet, RabbitMQ, Postgres, and ClickHouse.

## Hypothesis

One kluster brought up with authstar's `authstar` profile can host several builds of authstar at once, each in a namespace named for a story (`s-0901`, `s-0902`). The shared services need only per-namespace settings, not a cluster or an install each:

- a RabbitMQ vhost
- Postgres and ClickHouse databases with the namespace's prefix
- signet bundles under a signet namespace equal to the Kubernetes namespace, since SPIRE names each workload `spiffe://<trust domain>/ns/<namespace>/sa/<service>` and signet grants a workload its own namespace and service by convention

A story's deploy then costs a build, an image load, and a rollout, minutes rather than the cluster's bring-up. Deleting a story's namespace leaves the others and the shared services working.

## Success measure

- One `kluster up` with authstar's `integration/kluster.yaml` brings up signet, RabbitMQ, Postgres, and ClickHouse answering, and its time is recorded.
- Two builds from two commits of `~/git/bytepunx/authstar` run at once in `s-0901` and `s-0902`. Each has its own vhost, databases, and signet bundles, and authstar's golden-path scenario runs against each, passing or failing at a named step for a named reason.
- Deleting `s-0901` leaves `s-0902` and the shared services answering, and build A is then promoted to the default `authstar` namespace.
- The time of each build, load, deploy, and teardown is recorded, and the per-story cost is small beside the bring-up's.
- Nothing is committed in the authstar, signet, or kluster repositories.

## What was done

Host: `squatchship`, WSL2 on Linux 6.18, 24 cores, 62 GB. kluster 1.3.1, signet CLI 0.10.3. Another kluster cluster, `marketeer`, was running throughout and was left alone.

Scratch and copies:

- `/tmp/S-0575` holds the logs, the kubeconfig, and the two exported source trees.
- The story worktree's git-ignored `bin/s-0575/` holds what the experiment wrote: the scripts below, a copy of authstar's `integration/scripts`, and a copy of `kluster-lib`.
- Build A is authstar `e02e430`, main today. Build B is `af55c04`, main before S-0009 changed chronicle's token check. Each was exported with `git archive <commit> | tar -x -C /tmp/S-0575/src-<A|B>`, so the authstar checkout was only read.
- The host has no `nc`, which authstar's scripts call as `nc -z` to wait for a port-forward. A shim in `bin/s-0575/path/nc` opens `/dev/tcp/<host>/<port>` instead.
- The kubeconfig was written with `kluster kubeconfig --name authstar-integration --output /tmp/S-0575/kubeconfig` and used through `KUBECONFIG`, never merged into `~/.kube/config`.

### Bring-up

```sh
cd /tmp/S-0575
cp ~/git/bytepunx/authstar/integration/kluster.yaml .   # name authstar-integration, profile authstar
kluster --config ./kluster.yaml up
```

authstar's `integration/scripts/provision.sh` was not run. It adds a deploy key to the GitHub repository `bytepunx/authstar-integration` with `gh`, and commits and pushes secrets there for signet to sync. kluster 1.3.1's `authstar` profile no longer needs it: `kluster-lib/profile/authstar_services.go` seeds signet through its admin API, creates the databases, and deploys all eight services into the `authstar` namespace itself.

The shared services were checked like this:

```sh
kubectl -n postgres exec postgresql-0 -- env PGPASSWORD=postgres psql -U postgres -tAc "select datname from pg_database"
kubectl -n rabbitmq exec rabbitmq-0 -- rabbitmqctl list_vhosts -q
kubectl -n clickhouse exec clickhouse-shard0-0 -- clickhouse-client --password clickhouse -q "show databases"
kubectl -n signet port-forward deploy/signet 18444:8444 &
signet --server localhost:18444 status --token "$(kubectl create token signet-admin -n signet --audience signet --duration=1h)"
```

A baseline golden-path run against kluster's own `authstar` namespace followed. It used the operator tokens kluster writes to `~/.config/kluster/authstar-integration-authstar-tokens.yaml`, passed in as `TOWER_OPERATOR_TOKEN` and `HERALD_OPERATOR_TOKEN`.

## Results

### Bring-up

`kluster up` took 691 seconds, 11m31s, from 22:02:19Z to 22:13:50Z. The Go build of the seeding tool (below) ran on about five cores for two of those minutes, so an idle host may be a little faster.

| Step | Time |
|------|------|
| k3d cluster | 23s |
| cert-manager | 32s |
| SPIRE | 53s |
| traefik-tls | 5s |
| signet, unpinned, the latest chart | 2m50s |
| RabbitMQ | 2m01s |
| RabbitMQ topology operator | 25s |
| Dex | 12s |
| Postgres | 56s |
| ClickHouse | 1m04s |
| Configure SPIRE, signet | 25s, under 1s |
| Configure authstar: vhost, provisioner Job, databases, signet bundles, eight Deployments | 1m44s |

At the end every shared service answered:

| Service | Namespace | What answered |
|---------|-----------|---------------|
| Postgres | `postgres` | Databases `tower`, `keep`, `herald` |
| RabbitMQ | `rabbitmq` | Vhosts `/` and `authstar`; users `rmq-…` (administrator), `tower`, `keep` |
| ClickHouse | `clickhouse` | Database `chronicle` |
| signet | `signet` | `State: UNSEALED` |

All eight authstar services were running in `authstar`, from Docker Hub's `bytepunx/authstar-<service>:latest`.

The baseline golden path failed at its first step, in 8 seconds, before any story namespace existed:

```text
FAIL: tenant provisioning (HTTP 500): 14 UNAVAILABLE: No connection established.
Last error: SSL routines:tls_validate_record_header:wrong version number
```

Tower answered `POST /tenants` with a TLS call to a port that answers in plain text. So the `latest` images kluster pulls do not match the signet chart it installs, unpinned, today. kluster's own deployment of authstar is not a known-good baseline. Each story namespace therefore deploys a local build.

## Recommendation

Not settled yet. The choice between adopt, adapt, and drop rests on the story namespaces, which T-2467 and T-2468 try.
