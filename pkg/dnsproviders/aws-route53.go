package dnsproviders

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

// Route53 represnets a Route53 client
type Route53 struct {
	RecordSet RecordSet
}

// CreateRecordSet creates a record set
func (r Route53) CreateRecordSet(recordSetName, recordSetValue string) (string, error) {
	return r.ChangeRecordSet("UPSERT", recordSetName, recordSetValue)
}

// DeleteRecordSet deletes a record set
func (r Route53) DeleteRecordSet(recordSetName, recordSetValue string) (string, error) {
	return r.ChangeRecordSet("DELETE", recordSetName, recordSetValue)
}

// ChangeRecordSet change record set according to specified action
func (r Route53) ChangeRecordSet(action, recordSetName, recordSetValue string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", "route53"),
	})
	if err != nil {
		return "", err
	}
	svc := route53.New(sess)

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String(action),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(recordSetName),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(recordSetValue),
							},
						},
						TTL:  aws.Int64(r.RecordSet.TTL),
						Type: aws.String(r.RecordSet.RecordSetType),
					},
				},
			},
		},
		HostedZoneId: aws.String(r.RecordSet.HostedZoneID),
	}

	result, err := svc.ChangeResourceRecordSets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case route53.ErrCodeNoSuchHostedZone:
				fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
			case route53.ErrCodeNoSuchHealthCheck:
				fmt.Println(route53.ErrCodeNoSuchHealthCheck, aerr.Error())
			case route53.ErrCodeInvalidChangeBatch:
				fmt.Println(route53.ErrCodeInvalidChangeBatch, aerr.Error())
			case route53.ErrCodeInvalidInput:
				fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
			case route53.ErrCodePriorRequestNotComplete:
				fmt.Println(route53.ErrCodePriorRequestNotComplete, aerr.Error())
			default:
				return "", aerr
			}
		} else {
			return "", err
		}
	}

	return result.String(), nil
}