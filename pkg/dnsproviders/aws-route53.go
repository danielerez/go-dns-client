package dnsproviders

import (
	"fmt"
	"strings"

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

func (r Route53) getService() (*route53.Route53, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", "route53"),
	})
	if err != nil {
		return nil, err
	}
	return route53.New(sess), nil
}

// CreateRecordSet creates a record set
func (r Route53) CreateRecordSet(recordSetName, recordSetValue string) (string, error) {
	return r.ChangeRecordSet("CREATE", recordSetName, recordSetValue)
}

// DeleteRecordSet deletes a record set
func (r Route53) DeleteRecordSet(recordSetName, recordSetValue string) (string, error) {
	return r.ChangeRecordSet("DELETE", recordSetName, recordSetValue)
}

// UpdateRecordSet updates a record set
func (r Route53) UpdateRecordSet(recordSetName, recordSetValue string) (string, error) {
	return r.ChangeRecordSet("UPSERT", recordSetName, recordSetValue)
}

// ChangeRecordSet change record set according to the specified action
func (r Route53) ChangeRecordSet(action, recordSetName, recordSetValue string) (string, error) {
	svc, err := r.getService()
	if err != nil {
		return "", err
	}

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

// GetRecordSet returns a record set according to the specified name
func (r Route53) GetRecordSet(recordSetName string) (string, error) {
	svc, err := r.getService()
	if err != nil {
		return "", err
	}

	listParams := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(r.RecordSet.HostedZoneID),
		MaxItems:        aws.String("1"),
		StartRecordName: aws.String(recordSetName),
	}
	respList, err := svc.ListResourceRecordSets(listParams)
	if err != nil {
		return "", err
	}

	if len(respList.ResourceRecordSets) == 0 {
		// RecordSet not found
		return "", nil
	}

	recordSetNameAWSFormat := strings.Replace(recordSetName, "*", "\\052", 1) + "."
	if recordSetNameAWSFormat != *respList.ResourceRecordSets[0].Name {
		// RecordSet not found
		return "", nil
	}

	return respList.ResourceRecordSets[0].String(), nil
}
