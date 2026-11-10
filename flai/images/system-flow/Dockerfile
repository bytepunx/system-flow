# flai-agent-sf-claude-code: this repository's layer over the claude-code harness image (S-0510,
# ADR-0203), named by system-flow.yaml's agent.image.layer. It carries the tools the manifest's test
# tiers run, so that flai test and scripts/close-out.sh run inside the agent's container: Go at
# flai/go.mod's version with gcc and musl-dev for -race, golangci-lint v2 built with that Go,
# ShellCheck, and pnpm on the base image's Node 24. Build it after the harness image:
#   scripts/agent-image.sh --layer      (make agent-images builds all three)
# which passes BASE, GO_VERSION from flai/go.mod, and the versions scripts/install-tools.sh pins.
# No login state and no key, as the base has none; it keeps the base's entrypoint, pilot.

ARG BASE=flai-agent-claude-code:local
ARG GO_VERSION=1.26.0
# The Alpine release of the golang image, exact (S-0630): golang publishes 1.26.0 on 3.23 at the newest.
ARG GO_ALPINE_VERSION=3.23
ARG SHELLCHECK_VERSION=v0.10.0

FROM golang:${GO_VERSION}-alpine${GO_ALPINE_VERSION} AS go
ARG GOLANGCI_LINT_VERSION=v2.5.0
# golangci-lint built with the module's own Go, as scripts/install-tools.sh builds it, since a
# release binary built with an older Go refuses to lint a module that targets a newer one.
RUN --mount=type=cache,target=/root/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    GOBIN=/out go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}"
# flai's modules, fetched once into a cache every uid can read, so that a test run in the
# container needs no download.
COPY flai/go.mod flai/go.sum /src/flai/
RUN cd /src/flai && GOMODCACHE=/opt/go/mod go mod download && chmod -R a+rX /opt/go/mod

FROM koalaman/shellcheck-alpine:${SHELLCHECK_VERSION} AS shellcheck

FROM ${BASE}
ARG GO_VERSION
ARG GOLANGCI_LINT_VERSION=v2.5.0
ARG SHELLCHECK_VERSION
ARG PNPM_VERSION=10.34.5
LABEL dev.system-flow.go="${GO_VERSION}" \
      dev.system-flow.golangci-lint="${GOLANGCI_LINT_VERSION}" \
      dev.system-flow.shellcheck="${SHELLCHECK_VERSION}" \
      dev.system-flow.pnpm="${PNPM_VERSION}" \
      org.opencontainers.image.description="flai-agent-sf-claude-code: system-flow's test tools on flai-agent-claude-code"
USER root
# gcc and musl-dev for go test -race, which needs cgo; the rest is the base's.
RUN apk add --no-cache build-base && \
    npm install -g "pnpm@${PNPM_VERSION}" && \
    npm cache clean --force && rm -rf /root/.npm /tmp/node-compile-cache
COPY --from=go /usr/local/go /usr/local/go
COPY --from=go /out/golangci-lint /usr/local/bin/golangci-lint
COPY --from=go /opt/go/mod /opt/go/mod
COPY --from=shellcheck /bin/shellcheck /usr/local/bin/shellcheck
# The toolchain the image holds, never one go downloads; the module cache read-only, filled above;
# the build cache in the home any uid can write. No VCS stamp: in the container go build's stamping
# fails with git's exit status 128 though git itself answers there (I-0255), and a test's
# build needs none.
ENV PATH=/usr/local/go/bin:$PATH \
    GOTOOLCHAIN=local \
    GOMODCACHE=/opt/go/mod \
    GOFLAGS=-buildvcs=false
USER 65532:65532
