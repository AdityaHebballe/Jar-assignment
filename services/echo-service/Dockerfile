# Build stage
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /echo-service

# Final stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /echo-service /echo-service
EXPOSE 8081
ENTRYPOINT ["/echo-service"]
