FROM golang:1.23.2 as go

RUN mkdir /app
COPY . /app

WORKDIR /app/src

RUN go build server.go

RUN chmod a+x server

EXPOSE 8081

CMD ["/app/src/server"]