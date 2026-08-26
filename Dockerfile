FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM gcr.io/distroless/static-debian12

WORKDIR /srv

COPY --from=builder /app/server /srv/server
COPY --from=builder /app/web/public /srv/web/public

ENV STATIC_DIR=/srv/web/public

ENTRYPOINT ["/srv/server"]
