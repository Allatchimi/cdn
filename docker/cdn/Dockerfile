# ----------------------- Stage builder -----------------------
FROM golang:1.25.0-alpine AS builder
ENV LANG=C.UTF-8
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

WORKDIR /app

# Outils nécessaires pour COMPILER l'application (avec les paquets -dev)
RUN for attempt in 1 2 3 4 5; do \
      if apk add --no-cache bash curl gcc g++ vips-dev poppler-dev musl-dev pkgconf libc6-compat; then exit 0; fi; \
      echo "apk install failed (attempt ${attempt}), retrying..."; \
      sleep 5; \
    done; \
    exit 1

# Gestion du cache des dépendances Go
COPY ./go.mod ./go.sum ./
RUN for attempt in 1 2 3 4 5; do \
      if go mod download; then exit 0; fi; \
      echo "Go module download failed (attempt ${attempt}), retrying..."; \
      sleep 10; \
    done; \
    exit 1

# Build de l'application (On active CGO pour Vips et on nettoie le binaire avec -ldflags)
COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /app/main ./cmd/main.go

# ----------------------- Stage runner -----------------------
FROM alpine:latest AS cdn
ENV LANG=C.UTF-8

WORKDIR /app

# CORRECTION : En production, on installe uniquement le RUNTIME (Pas de gcc, g++, ni de paquets -dev)
RUN for attempt in 1 2 3 4 5; do \
      if apk add --no-cache bash vips vips-poppler ca-certificates; then exit 0; fi; \
      echo "apk install failed (attempt ${attempt}), retrying..."; \
      sleep 5; \
    done; \
    exit 1

# Récupération du binaire léger et des assets
COPY --from=builder /app/main ./main
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/app.env ./app.env

RUN chmod +x /app/main

EXPOSE 3100

ENTRYPOINT ["/app/main"]
