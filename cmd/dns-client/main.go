package main

import (
	"flag"
	"fmt"

	"github.com/danielerez/dns-lib-go/pkg/dnsproviders"
)

func main() {
	var action string
	var hostedZoneID, recordSetName, recordSetValue string
	var ttl int64

	flag.StringVar(&action, "action", "UPSERT", "Action to execute.")
	flag.StringVar(&hostedZoneID, "hosted-zone-id", "", "HostedZone ID.")
	flag.StringVar(&recordSetName, "record-set-name", "", "RecordSet name.")
	flag.StringVar(&recordSetValue, "record-set-value", "", "RecordSet value.")
	flag.Int64Var(&ttl, "ttl", 60, "TTL in seconds.")
	flag.Parse()

	var dnsProvider dnsproviders.Provider = dnsproviders.Route53{
		RecordSet: dnsproviders.RecordSet{
			HostedZoneID: hostedZoneID,
			RecordSetType: "A",
			TTL: ttl,
		},
	}

	var output string
	var err error
	switch action {
	case "CREATE":
		output, err = dnsProvider.CreateRecordSet(recordSetName, recordSetValue)
	case "UPSERT":
		output, err = dnsProvider.UpdateRecordSet(recordSetName, recordSetValue)
	case "DELETE":
		output, err = dnsProvider.DeleteRecordSet(recordSetName, recordSetValue)
	}

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(output)
	}
}
