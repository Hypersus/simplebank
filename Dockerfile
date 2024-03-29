# build stage
FROM golang:1.20-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/wait-for .
COPY app.env .
COPY ./db/migration ./migration
COPY ./start.sh .
RUN chmod +x start.sh
RUN chmod +x migrate
RUN chmod +x wait-for

# configuration
EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]