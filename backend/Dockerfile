# ── Stage 1: Build ─────────────────────────────────────────────────────────────
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o server ./cmd/app/main.go

# ── Stage 2: Run ───────────────────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/server"]
