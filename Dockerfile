# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /usr/local/bin/go-memtest .

FROM gcr.io/distroless/static

COPY --from=builder /usr/local/bin/go-memtest /usr/local/bin/go-memtest

ENTRYPOINT ["/usr/local/bin/go-memtest"]