FROM golang:1.23.2 as go

RUN mkdir /app
COPY . /app

WORKDIR /app/

RUN go build src/server/main.go

EXPOSE 8081

CMD ["./main"]