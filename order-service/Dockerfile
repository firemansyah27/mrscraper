FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 4000

CMD ["./main"]
