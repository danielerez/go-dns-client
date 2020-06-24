package dnsproviders

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockRoute53Client struct {
	route53iface.Route53API
}

func (m *mockRoute53Client) ChangeResourceRecordSets(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
	var output = route53.ChangeResourceRecordSetsOutput{
		ChangeInfo: &route53.ChangeInfo{
			Status: aws.String("PENDING"),
		},
	}
	return &output, nil
}

func (m *mockRoute53Client) ListResourceRecordSets(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
	var output = route53.ListResourceRecordSetsOutput{
		ResourceRecordSets: []*route53.ResourceRecordSet{
			{
				Name: aws.String("example.com."),
				Type: aws.String("A"),
			},
			{
				Name: aws.String("example2.com."),
				Type: aws.String("MX"),
			},
		},
	}
	return &output, nil
}

func TestRoute53(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "aws-route53_test")
}

var _ = Describe("RecordSet", func() {
	var dnsProvider Route53

	BeforeEach(func() {
		mockSvc := &mockRoute53Client{}
		dnsProvider = Route53{
			RecordSet: RecordSet{
				HostedZoneID:  "abc",
				RecordSetType: "A",
				TTL:           60,
			},
		}
		dnsProvider.SVC = mockSvc
	})

	It("change RecordSet", func() {
		output, err := dnsProvider.ChangeRecordSet("CREATE", "example.com", "1.2.3.4")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(output).ToNot(Equal(nil))
	})

	It("get RecordSet found", func() {
		output, err := dnsProvider.GetRecordSet("example.com")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(output).ToNot(Equal(""))
	})

	It("get RecordSet not found", func() {
		output, err := dnsProvider.GetRecordSet("test.com.")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(output).To(Equal(""))
	})

	It("get RecordSet wrong type", func() {
		output, err := dnsProvider.GetRecordSet("example2.com.")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(output).To(Equal(""))
	})
})
