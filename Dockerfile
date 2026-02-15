# --- Build stage ---
FROM golang:1.24-alpine AS build
WORKDIR /src

# Leverage Go modules cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Set reproducible build flags
ENV CGO_ENABLED=1 \
GOOS=linux \
GOARCH=amd64 \
GOFLAGS=-buildvcs=false

# Build the actual api
RUN go build -o /out/app ./cmd/api

# --- Runtime stage ---
FROM alpine:3.20
# Set a working directory for the app
WORKDIR /app

# Add a tiny tool for healthcheck
RUN apk add --no-cache wget

# Create sqlite directory
RUN mkdir -p /app/data

# Copy the binary
COPY --from=build /out/app /app/app

# Listen on PORT
EXPOSE 8800

# Basic healthcheck hitting /health (customize to your route)
HEALTHCHECK --interval=30s --timeout=3s --retries=3 CMD wget -qO- http://localhost:${PORT}/health || exit 1

ENTRYPOINT ["/app/app"]