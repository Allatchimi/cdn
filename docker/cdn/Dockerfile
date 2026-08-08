# ----------------------- Stage builder -----------------------
FROM golang:1.24-bookworm AS builder
ENV LANG=C.UTF-8

WORKDIR /app

# Install compilation packages and libvips
RUN apt-get update && apt-get install -y \
    bash \
    curl \
    build-essential \
    pkg-config \
    libvips-dev \
    && rm -rf /var/lib/apt/lists/*

# Install dependencies
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

# Build the app
COPY . .
RUN CGO_ENABLED=1 go build -a -installsuffix cgo -o ./.build/main ./cmd/main.go

# ----------------------- Stage runner -----------------------
FROM debian:bookworm-slim AS cdn
ENV LANG=C.UTF-8

WORKDIR /app

# Install runtime packages and libvips
RUN apt-get update && apt-get install -y \
    bash \
    ca-certificates \
    libvips42 \
    && rm -rf /var/lib/apt/lists/*

# Copy the binary and assets to the container
COPY --from=builder /app/.build/main ./main
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/app.env ./app.env
RUN chmod +x /app/main

EXPOSE 3100

ENTRYPOINT ["/app/main"]
