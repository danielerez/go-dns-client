package dnsproviders

type Provider interface {
	CreateRecordSet(recordSetName, recordSetValue string) (string, error)
	UpdateRecordSet(recordSetName, recordSetValue string) (string, error)
	DeleteRecordSet(recordSetName, recordSetValue string) (string, error)
	GetRecordSet(recordSetName string) (string, error)
}

type RecordSet struct {
	HostedZoneID  string
	RecordSetType string
	TTL           int64
}
