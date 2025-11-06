FROM golang:latest AS builder
WORKDIR /app

# Copy source code
COPY . .

# Install build tools
RUN apt-get update && apt-get install -y make git && rm -rf /var/lib/apt/lists/*

# Copy go.mod & go.sum first to cache dependencies
COPY go.mod go.sum ./

# Download all external packages
RUN make force-download
RUN make prepare

# Cache dependencies using Docker BuildKit
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

CMD ["tail", "-f", "/dev/null"]

