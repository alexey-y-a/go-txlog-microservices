FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app/services/kv-service/cmd/kv

RUN go build -o /bin/kv

# === Stage 2 ===

FROM gcr.io/distroless/base-debian12

COPY --FROM=builder /bin/kv /kv

EXPOSE 8081

ENTRYPOINT ["/kv"]