# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN apk update && \
    apk add upx && \
    apk add --no-cache gcc && \
    rm -rf /var/lib/apt/lists/* && \
    apk add musl-dev
RUN go build -ldflags="-s -w" -o main /app/cmd/api/main.go && upx --brute main

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app/cmd/api/
COPY env env

RUN apk update
RUN apk add -U tzdata

ENV TZ=Asia/Taipei
EXPOSE 9600
CMD ["/app/cmd/api/main"]