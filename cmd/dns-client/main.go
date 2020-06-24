package main

import (
	"flag"
	"fmt"

	"github.com/danielerez/go-dns-client/pkg/dnsproviders"
)

func main() {
	var action string
	var hostedZoneID, recordSetName, recordSetValue string
	var ttl int64
	var sharedCreds bool
	var deafultEnvVars bool

	flag.StringVar(&action, "action", "CREATE", "Action to execute (CREATE/UPSERT/DELETE/GET).")
	flag.StringVar(&hostedZoneID, "hosted-zone-id", "", "HostedZone ID.")
	flag.StringVar(&recordSetName, "record-set-name", "", "RecordSet name.")
	flag.StringVar(&recordSetValue, "record-set-value", "", "RecordSet value.")
	flag.Int64Var(&ttl, "ttl", 60, "TTL in seconds.")
	flag.BoolVar(&sharedCreds, "shared-creds", false, "Use shared .aws/credentials file ('route53' profile).")
	flag.BoolVar(&deafultEnvVars, "default-env-vars", false, "Use default env vars (AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY).")
	flag.Parse()

	var dnsProvider dnsproviders.Provider = dnsproviders.Route53{
		RecordSet: dnsproviders.RecordSet{
			HostedZoneID:  hostedZoneID,
			RecordSetType: "A",
			TTL:           ttl,
		},
		SharedCreds: sharedCreds,
		DeafultEnvVars: deafultEnvVars,
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
	case "GET":
		output, err = dnsProvider.GetRecordSet(recordSetName)
	}

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(output)
	}
}
