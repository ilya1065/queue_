# ---------- build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/app ./cmd
RUN go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate


# ---------- runtime stage ----------
FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

RUN addgroup -S app && adduser -S app -G app \
  && mkdir -p /app/data /app/internal/config /app/db/migrations \
  && chown -R app:app /app

COPY --from=builder /out/app /app/app
COPY --from=builder /out/migrate /app/migrate

# конфиг по пути, который ищет код
COPY internal/config/config.yml /app/internal/config/config.yml
COPY db/migrations /app/db/migrations

USER app
ENTRYPOINT ["/app/app"]
