FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.12
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]