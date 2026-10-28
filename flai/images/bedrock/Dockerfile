# flai-agent-bedrock: the base of every agent image (E-0029), flai built from this tree with git,
# ca-certificates, and tini. Build from the repository root:
#   docker build -f flai/images/bedrock/Dockerfile -t flai-agent-bedrock:local .
# scripts/agent-image.sh builds it, then the agent images on top of it. Not stamped, not signed.
#
# Alpine, not Debian slim: Claude Code's npm package ships a musl build of its native binary
# (@anthropic-ai/claude-code-linux-*-musl), which runs with libgcc, libstdc++, and the system
# ripgrep, and Alpine's own nodejs package is Node 24, so the claude-code layer needs no other source.

ARG GO_VERSION=1.26
ARG ALPINE_VERSION=3.24

# ---- flai build ------------------------------------------------------------
FROM golang:${GO_VERSION}-alpine AS build
ARG FLAI_VERSION=dev
ARG FLAI_COMMIT=unknown
ARG FLAI_DATE=
WORKDIR /src/flai
COPY flai/go.mod flai/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY flai/ ./
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -o /out/flai -ldflags "-s -w \
      -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Version=${FLAI_VERSION} \
      -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Commit=${FLAI_COMMIT} \
      -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Date=${FLAI_DATE}" .

# ---- runtime ---------------------------------------------------------------
FROM alpine:${ALPINE_VERSION}
ARG FLAI_VERSION=dev
LABEL dev.system-flow.flai="${FLAI_VERSION}" \
      org.opencontainers.image.source="https://github.com/bytepunx/system-flow" \
      org.opencontainers.image.description="flai-agent-bedrock: flai, git, and tini, the base of the flai agent images" \
      org.opencontainers.image.licenses="LicenseRef-system-flow-Shield-1.0"
RUN apk add --no-cache git ca-certificates tini
COPY --from=build /out/flai /usr/local/bin/flai
# The container runs as the operator's uid (--user), which has no entry here: a home any uid can
# write to holds what git and the agent keep there.
RUN addgroup -g 65532 agent && adduser -D -u 65532 -G agent -h /home/agent agent && \
    chmod 1777 /home/agent
ENV HOME=/home/agent
WORKDIR /home/agent
USER 65532:65532
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["sh"]
