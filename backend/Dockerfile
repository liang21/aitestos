# ============================================================
# aittestos - Industrial-Grade Multi-Stage Dockerfile
# Target: < 20MB final image with Distroless security
# Platform: linux/amd64, linux/arm64
# ============================================================

# ---------- Stage 1: Builder ----------
FROM golang:1.24-bookworm AS builder

# Install build dependencies (for git describe)
RUN apt-get update && apt-get install -y --no-install-recommends git && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /src

# ========================================
# EXTREME CACHING: Go modules layer
# ========================================
# This layer only changes when go.mod/go.sum change
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# ========================================
# EXTREME CACHING: Source code layer
# ========================================
# Copy source code (invalidates above layer when any .go file changes)
COPY . .

# ========================================
# BUILD: Static binary with metadata
# ========================================
# Build args for version injection (optional)
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Compile with:
#   - CGO_ENABLED=0: pure Go binary (no libc dependency)
#   - -s -w: strip debug info (reduce size ~30%)
#   - -trimpath: remove file system paths from binary
#   - ldflags: inject version/build-time metadata
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=linux \
    go build \
        -trimpath \
        -ldflags="-s -w \
            -X main.version=${VERSION} \
            -X main.buildTime=${BUILD_TIME} \
            -X main.gitCommit=${GIT_COMMIT}" \
        -o /usr/local/bin/server \
        ./cmd/server

# Verify binary is statically linked (no libc dependency)
RUN ldd /usr/local/bin/server 2>&1 | grep -q "not a dynamic executable" || \
    (echo "ERROR: Binary is not static!" && exit 1)

# Display binary size for verification
RUN ls -lh /usr/local/bin/server

# ---------- Stage 2: Distroless Final ----------
FROM gcr.io/distroless/static-debian12:nonroot

# Metadata
LABEL maintainer="aitestos"
LABEL description="AI-powered test case management platform"
LABEL version="${VERSION}"
LABEL license="MIT"

# Copy binary from builder (preserves ownership: nonroot:nonroot)
COPY --from=builder --chown=nonroot:nonroot /usr/local/bin/server /server

# Expose application port
EXPOSE 8080

# Distroless has no shell; health checks must be external (orchestration-level)
# Kubernetes: use livenessProbe/readinessProbe with HTTPGet
# docker-compose: uses wget healthcheck defined in compose file

# Run the binary
ENTRYPOINT ["/server"]
