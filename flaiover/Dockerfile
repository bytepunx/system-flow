# flaiover: the system-flow dashboard, with a flai binary from the same
# commit (ADR-0016). Build from the repository root:
#   docker build -f flaiover/Dockerfile -t flaiover:local .
# Runs as any UID (flai dashboard passes --user); listens on 3000. The
# project is whatever PROJECT_DIR names: flai dashboard mounts the repository
# at its host path and sets PROJECT_DIR to it, so the links git keeps for
# story worktrees resolve here (ADR-0022); /project is the default for running
# the image by hand. flai's config and cache live under PROJECT_DIR/.flai-cache,
# so no writable home is needed.

ARG GO_VERSION=1.26
ARG NODE_VERSION=24

# ---- flai ------------------------------------------------------------------
FROM golang:${GO_VERSION}-alpine AS flai
ARG FLAI_VERSION=dev
ARG FLAI_COMMIT=unknown
ARG FLAI_DATE=unknown
WORKDIR /src/flai
COPY flai/go.mod flai/go.sum ./
RUN go mod download
COPY flai/ ./
RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags "-s -w -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Version=${FLAI_VERSION} -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Commit=${FLAI_COMMIT} -X github.com/bytepunx/system-flow/flai/internal/buildinfo.Date=${FLAI_DATE}" \
    -o /out/flai .

# ---- flaiover build --------------------------------------------------------
FROM node:${NODE_VERSION}-alpine AS build
ARG FLAI_COMMIT=unknown
WORKDIR /src/flaiover
ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0
RUN corepack enable
COPY flaiover/package.json flaiover/pnpm-lock.yaml flaiover/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY flaiover/ ./
RUN pnpm build && pnpm prune --prod

# ---- runtime ---------------------------------------------------------------
FROM node:${NODE_VERSION}-alpine
ARG FLAI_COMMIT=unknown
ARG FLAIOVER_VERSION=0.0.0
LABEL org.opencontainers.image.source="https://github.com/bytepunx/system-flow" \
      org.opencontainers.image.description="flaiover: the system-flow dashboard with a bundled flai" \
      org.opencontainers.image.licenses="UNLICENSED"
# openssh-client: git pushes an acceptance over SSH when the operator has
# given the container a key (ADR-0026). Unused, and harmless, otherwise.
RUN apk add --no-cache git openssh-client-default tini
WORKDIR /app
COPY --from=build /src/flaiover/build ./build
COPY --from=build /src/flaiover/node_modules ./node_modules
COPY --from=build /src/flaiover/package.json ./package.json
COPY --from=build /src/flaiover/server.js ./server.js
COPY --from=flai /out/flai /usr/local/bin/flai
ENV NODE_ENV=production \
    PORT=3000 \
    HOST=0.0.0.0 \
    PROJECT_DIR=/project \
    FLAI_BIN=/usr/local/bin/flai \
    HOME=/tmp \
    LOG_FORMAT=json \
    FLAIOVER_COMMIT=${FLAI_COMMIT} \
    FLAIOVER_VERSION=${FLAIOVER_VERSION}
EXPOSE 3000
# Unprivileged by default; any --user works because nothing here needs a home.
USER 65532:65532
ENTRYPOINT ["/sbin/tini", "--"]
# server.js is adapter-node's handler plus the WebSocket upgrade for /agent (ADR-0029).
CMD ["node", "server.js"]
