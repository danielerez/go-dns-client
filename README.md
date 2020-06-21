# go-dns-client

## Run

    $ go run cmd/dns-client/main.go --help
    -action string
            Action to execute (CREATE/UPSERT/DELETE). (default "CREATE")
    -hosted-zone-id string
            HostedZone ID.
    -record-set-name string
            RecordSet name.
    -record-set-value string
            RecordSet value.
    -ttl int
            TTL in seconds. (default 60)

## Docker

### Build
    $ docker image build -t dns-client .

### Run
    $ docker run -v $HOME/.aws/credentials:/root/.aws/credentials:ro dns-client \
    -hosted-zone-id <hosted-zone-id> \
    -record-set-name <record-set-name> \
    -record-set-value <record-set-value>