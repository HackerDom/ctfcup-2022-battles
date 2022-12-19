FROM golang:1.16-alpine

WORKDIR /app

COPY main.go main.go

RUN go build -o checker main.go

EXPOSE 8080

CMD ["/app/checker"]