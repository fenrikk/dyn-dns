package dnsprovider

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

const (
	DefaultTTL = 60
)

// Route53Provider implements the DNSProvider interface for AWS Route53
type Route53Provider struct {
	client       *route53.Client
	hostedZoneID string
	recordName   string
	ttl          int64
}

type Route53ProviderOptions struct {
	HostedZoneID string
	RecordName   string
	TTL          int64
}

func NewRoute53Provider(opts Route53ProviderOptions) (*Route53Provider, error) {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("aws config load error: %w", err)
	}

	client := route53.NewFromConfig(cfg)

	recordName := opts.RecordName
	if !strings.HasSuffix(recordName, ".") {
		recordName = recordName + "."
	}

	ttl := opts.TTL
	if ttl == 0 {
		ttl = DefaultTTL
	}

	return &Route53Provider{
		client:       client,
		hostedZoneID: opts.HostedZoneID,
		recordName:   recordName,
		ttl:          ttl,
	}, nil
}

func (r *Route53Provider) GetCurrentDNSRecord() (string, error) {
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(r.hostedZoneID),
		StartRecordName: aws.String(r.recordName),
		StartRecordType: types.RRTypeA,
		MaxItems:        aws.Int32(1),
	}

	resp, err := r.client.ListResourceRecordSets(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("error while getting DNS records: %w", err)
	}

	if len(resp.ResourceRecordSets) == 0 {
		return "", fmt.Errorf("record '%s' not found", r.recordName)
	}

	recordSet := resp.ResourceRecordSets[0]

	if *recordSet.Name != r.recordName || recordSet.Type != types.RRTypeA {
		return "", fmt.Errorf("record with A type for '%s' not found", r.recordName)
	}

	if len(recordSet.ResourceRecords) == 0 {
		return "", fmt.Errorf("record '%s' does not contain IP address", r.recordName)
	}

	return *recordSet.ResourceRecords[0].Value, nil
}

func (r *Route53Provider) UpdateDNSRecord(ip string) (bool, error) {
	currentIP, err := r.GetCurrentDNSRecord()

	if err != nil {
		return false, fmt.Errorf("error getting current DNS record: %w", err)
	}

	// Skip if no changes
	if currentIP == ip {
		return false, nil
	}

	// Trying to update record
	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(r.hostedZoneID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(r.recordName),
						Type: types.RRTypeA,
						TTL:  aws.Int64(r.ttl),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
					},
				},
			},
		},
	}

	_, err = r.client.ChangeResourceRecordSets(context.TODO(), input)
	if err != nil {
		return false, fmt.Errorf("error while updating DNS record: %w", err)
	}

	return true, nil
}
