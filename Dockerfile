FROM golang:1.14.3

ENV GO111MODULE=on

WORKDIR /app
COPY . /app

RUN cd cmd/dns-client; go build -o /app/dns-client .

ENTRYPOINT ["./dns-client"]
CMD ["--help"]