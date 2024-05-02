FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/free-chat-to-api .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/free-chat-to-api /app/free-chat-to-api

EXPOSE 8484

CMD [ "/app/free-chat-to-api" ]