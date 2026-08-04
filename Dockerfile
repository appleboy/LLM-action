FROM golang:1.25-alpine AS builder

# Build arguments for version injection
ARG VERSION=dev
ARG COMMIT=unknown

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with version information
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT}" \
    -o llm-action .

# Final stage
FROM alpine:3.22

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/llm-action /app/

# GitHub docker actions must run as root: the runner's file-command files
# (GITHUB_OUTPUT etc.) are owned by the runner user and not writable by
# an unprivileged container user.
ENTRYPOINT ["/app/llm-action"]
