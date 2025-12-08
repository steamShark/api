# --- Build stage ---
FROM golang:1.24-alpine AS build
WORKDIR /src

# Install build deps if needed (e.g., git)
RUN apk add --no-cache git

# Leverage Go modules cache
COPY go.mod go.sum ./
RUN go mod download
# Copy the rest of the source
COPY . .
# Set reproducible build flags
ENV CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64 \
GOFLAGS=-buildvcs=false
# Adjust the build path to your main package
# If your main is at ./cmd/api/main.go, build like this:
RUN go build -ldflags "-s -w" -o /out/app ./cmd/api

# --- Runtime stage ---
FROM alpine:3.20
# Create non-root user
RUN addgroup -S app && adduser -S app -G app
# Add a tiny tool for healthcheck
RUN apk add --no-cache wget ca-certificates && update-ca-certificates
# Copy the binary
COPY --from=build /out/app /bin/app
# Runtime env
ENV PORT=8800
# Listen on PORT
EXPOSE 8800
# Basic healthcheck hitting /health (customize to your route)
HEALTHCHECK --interval=30s --timeout=3s --retries=3 CMD wget -qO- http://localhost:${PORT}/health || exit 1

USER app

ENTRYPOINT ["/bin/app"]