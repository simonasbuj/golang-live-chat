FROM golang:1.25.3-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o chatApp ./cmd/main.go 

RUN chmod +x /app/chatApp

# build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/chatApp /app
COPY --from=builder /app/config ./config
COPY --from=builder /app/frontend ./frontend

CMD ["/app/chatApp"]
