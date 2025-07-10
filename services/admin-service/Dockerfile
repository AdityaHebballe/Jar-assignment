# Build stage
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /admin-service

# Final stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /admin-service .
COPY --from=builder /app/admin.html .
EXPOSE 8082
ENTRYPOINT ["/admin-service"]
