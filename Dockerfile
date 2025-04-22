FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o s3-aggregated-metrics-collector-collector .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/s3-aggregated-metrics-collector-collector .
ENTRYPOINT [ "/app/s3-aggregated-metrics-collector-collector" ]
