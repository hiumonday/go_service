FROM golang:1.24-alpine AS builder
RUN apk add --no-progress --no-cache gcc musl-dev
WORKDIR /app
COPY . .
RUN go mod download

RUN go build -tags musl -ldflags '-extldflags "-static"' -o server ./cmd/server/main.go

FROM alpine:3.19

# Chứng chỉ CA để gọi HTTPS
# RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
ENTRYPOINT ["/app/server"]