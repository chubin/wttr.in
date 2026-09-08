# ========================
# Builder stage (Ubuntu/Debian-based)
# ========================
FROM golang:1.26-bookworm AS builder

WORKDIR /app

# Install build dependencies + fonts (Debian packages)
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    rsync \
    gcc \
    libc6-dev \
    fontconfig \
    fonts-dejavu-core \
    fonts-noto-core \
    fonts-noto-cjk \
    fonts-wqy-zenhei \
    fonts-symbola \
    fonts-motoya-l-cedar \
    fonts-lexi-gulim \
    && rm -rf /var/lib/apt/lists/*

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy source code
COPY . .

# Run the official build process
RUN bash build.sh build

# ========================
# Runtime stage (Debian-slim: the builder produces a glibc/cgo-linked
# binary via mattn/go-sqlite3, which will not run under Alpine's musl libc)
# ========================
FROM debian:bookworm-slim

WORKDIR /app

# Minimal runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    useradd -r -M -u 1000 -s /usr/sbin/nologin wttr && \
    mkdir -p /app/cache && \
    chown -R wttr:wttr /app

# Copy the built binary
COPY --from=builder /app/srv /app/bin/srv

USER wttr

EXPOSE 8002

# The binary's own CLI is `srv <subcommand> [args]`; the config file path
# is required, so it must be supplied on the command line (see README's
# "Installation" section for the config.yaml format).
CMD ["/app/bin/srv", "srv", "/app/config.yaml"]
