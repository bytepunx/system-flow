# flaiover: the system-flow dashboard. Build from the repository root:
#   docker build -f flaiover/Dockerfile -t flaiover:local .
# Runs as any UID (flai dashboard passes --user, so the two secrets it mounts can be read); listens on 3000. The image
# holds no flai, no git, and no ssh (S-0075, ADR-0029): everything it shows
# and everything it changes is asked of flai serve on the host, which dials
# the /agent endpoint. It needs no writable home, and nothing of the project is
# mounted into it (S-0077, ADR-0031).

ARG NODE_VERSION=24

# ---- flaiover build --------------------------------------------------------
# Runs on the build platform for every target (S-0119): the build output and the
# production node_modules are plain JavaScript, so only the runtime stage runs
# emulated when CI builds linux/arm64 on an amd64 runner.
FROM --platform=$BUILDPLATFORM node:${NODE_VERSION}-alpine AS build
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
# The flai of the commit this image was built from: the one it was tested against, not one it holds.
ARG FLAI_VERSION=dev
LABEL dev.system-flow.flai="${FLAI_VERSION}" \
      org.opencontainers.image.source="https://github.com/bytepunx/system-flow" \
      org.opencontainers.image.description="flaiover: the system-flow dashboard; it reaches a project through flai serve on the host" \
      org.opencontainers.image.licenses="LicenseRef-system-flow-Shield-1.0"
RUN apk add --no-cache tini
WORKDIR /app
COPY --from=build /src/flaiover/build ./build
COPY --from=build /src/flaiover/node_modules ./node_modules
COPY --from=build /src/flaiover/package.json ./package.json
COPY --from=build /src/flaiover/server.js ./server.js
COPY --from=build /src/flaiover/scheme.js ./scheme.js
# The license travels with every copy of the software (S-0231): the Host › License page reads
# this file (FLAIOVER_LICENSE names another), and it is the same text flai license prints.
COPY LICENSE.md ./LICENSE.md
# The release stamp (S-0235, release-signing.md § A release stamp in hello): on a flaiover/v* tag,
# release-flaiover.yml signs `flaiover <version> <commit>` with the release key and passes the
# statement and its base64 signature, which this keeps as the two lines of /app/release-stamp for
# the server to check at start (src/lib/server/release.ts). It must name the version and commit this
# image reports. Given neither, as flai dashboard --build and scripts/flaiover-image.sh build it, the
# image has no stamp; given one alone, the build fails. Declared here, after the copies, so that a
# new stamp rebuilds only this layer.
ARG FLAIOVER_STAMP_STATEMENT=
ARG FLAIOVER_STAMP_SIGNATURE=
RUN if [ -n "$FLAIOVER_STAMP_STATEMENT" ] && [ -n "$FLAIOVER_STAMP_SIGNATURE" ]; then \
        if [ "$FLAIOVER_STAMP_STATEMENT" != "flaiover $FLAIOVER_VERSION $FLAI_COMMIT" ]; then \
            echo "the stamp names '$FLAIOVER_STAMP_STATEMENT', not this build's 'flaiover $FLAIOVER_VERSION $FLAI_COMMIT'" >&2; exit 1; \
        fi; \
        printf '%s\n%s\n' "$FLAIOVER_STAMP_STATEMENT" "$FLAIOVER_STAMP_SIGNATURE" > /app/release-stamp && \
        chmod 0644 /app/release-stamp; \
    elif [ -n "$FLAIOVER_STAMP_STATEMENT$FLAIOVER_STAMP_SIGNATURE" ]; then \
        echo "FLAIOVER_STAMP_STATEMENT and FLAIOVER_STAMP_SIGNATURE go together: give both or neither" >&2; exit 1; \
    fi
ENV NODE_ENV=production \
    PORT=3000 \
    HOST=0.0.0.0 \
    HOME=/tmp \
    LOG_FORMAT=json \
    FLAIOVER_COMMIT=${FLAI_COMMIT} \
    FLAIOVER_VERSION=${FLAIOVER_VERSION}
EXPOSE 3000
# Docker marks the container unhealthy when /_health stops answering (S-0184); flai host restarts a
# dashboard that does not answer, and flai dashboard status reports it. busybox wget is in the base image.
HEALTHCHECK --interval=15s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -q -O /dev/null "http://127.0.0.1:${PORT}/_health" || exit 1
# Unprivileged by default; any --user works because nothing here needs a home.
USER 65532:65532
ENTRYPOINT ["/sbin/tini", "--"]
# server.js is adapter-node's handler plus the WebSocket upgrade for /agent (ADR-0029), and tells
# the handler which scheme each request came over (scheme.js, S-0083).
CMD ["node", "server.js"]
