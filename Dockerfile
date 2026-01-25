# -------- Build stage --------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for go modules sometimes)
RUN apk add --no-cache git

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server ./cmd/server

# -------- Runtime stage --------
FROM gcr.io/distroless/base-debian12:debug

WORKDIR /app

COPY --from=builder /app/server /server

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/server"]

# FROM golang:1.25-alpine AS builder

# WORKDIR /app

# COPY . .

# RUN go mod download

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o server ./cmd/server

# FROM gcr.io/distroless/base-debian12:debug

# WORKDIR /app

# COPY --from=builder /app/server .

# EXPOSE 8080

# CMD ["/app/server"]
