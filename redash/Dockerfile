# Build stage
FROM harbor.razorpay.com/razorpay/rzp-docker-image-inventory-multi-arch:rzp-golden-image-base-golang-1.24-alpine3.21 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o redash-mcp ./cmd/main/

# Final stage
FROM harbor.razorpay.com/razorpay/rzp-docker-image-inventory-multi-arch:rzp-golden-image-base-alpine3.21

RUN apk --no-cache add ca-certificates

# Create a non-root user to run the application
RUN addgroup -S rzpgroup && adduser -S rzp -G rzpgroup

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/redash-mcp .

# Change ownership of the application to the non-root user
RUN chown -R rzp:rzpgroup /app

ENV CONFIG="" \
    REDASH_URL="" \
    REDASH_API_KEY="" \
    LOG_FILE=""

# Switch to the non-root user
USER rzp

# Fixed entrypoint - removed echo statement that interferes with MCP stdio
ENTRYPOINT ["./redash-mcp", "stdio", "--redash-url", "$REDASH_URL", "--redash-api-key", "$REDASH_API_KEY"]
