# FROM golang:1.21.1-alpine AS builder
# WORKDIR /app
# COPY . .
# RUN go mod download
# RUN go build -o app cmd/app/main.go

# FROM alpine:latest
# COPY --from=builder /app/app /app
# WORKDIR /
# CMD ["./app"]

FROM golang:1.21.1-alpine
WORKDIR /app
COPY . .
RUN go mod download
CMD ["go", "run", "cmd/app/main.go"]