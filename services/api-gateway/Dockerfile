FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app/services/api-gateway/cmd/api

RUN go build -o /bin/api

# === Stage 2 distroless ===

FROM gcr.io/distroless/base-debian12

COPY --from=builder /bin/api /api

EXPOSE 8080

ENTRYPOINT ["/api"]