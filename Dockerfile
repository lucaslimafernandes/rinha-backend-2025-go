# Golang image
FROM golang:1.24-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build main.go
# RUN go build -o main .

FROM alpine:3.19

COPY --from=builder /app/main .

EXPOSE 9999

CMD ["./main"]