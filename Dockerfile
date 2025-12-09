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

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /home/appuser

# Copy the binary from builder
COPY --from=builder /app/llm-action /home/appuser/

# Change ownership to non-root user
RUN chown -R appuser:appuser /home/appuser

# Switch to non-root user
USER appuser

# Run the application
ENTRYPOINT ["/home/appuser/llm-action"]
