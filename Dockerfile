FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg ./pkg
COPY cmd ./cmd

RUN CGO_ENABLED=0 go build -o main cmd/main/main.go

FROM chromedp/headless-shell:145.0.7587.5

RUN apt-get update && apt-get install -y tzdata && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/main /app/main

ENTRYPOINT ["/app/main"]
