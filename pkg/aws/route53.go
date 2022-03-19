package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

func Route53ChangeResourceRecordSets(changeBatch *route53.ChangeBatch, zoneId string) error {
	sess := GetSession()
	svc := route53.New(sess)

	_, err := svc.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  changeBatch,
		HostedZoneId: aws.String(zoneId),
	})

	return err
}

func Route53GetHostedZone(zoneId string) (*route53.GetHostedZoneOutput, error) {
	sess := GetSession()
	svc := route53.New(sess)

	zone, err := svc.GetHostedZone(&route53.GetHostedZoneInput{
		Id: aws.String(zoneId),
	})

	if err != nil {
		return nil, err
	}

	return zone, nil
}

func Route53ListHostedZones() ([]*route53.HostedZone, error) {
	sess := GetSession()
	svc := route53.New(sess)

	zones := []*route53.HostedZone{}
	pageNum := 0

	err := svc.ListHostedZonesPages(&route53.ListHostedZonesInput{},
		func(page *route53.ListHostedZonesOutput, lastPage bool) bool {
			pageNum++
			zones = append(zones, page.HostedZones...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return zones, nil
}

func Route53ListHostedZonesByName(name string) ([]*route53.HostedZone, error) {
	sess := GetSession()
	svc := route53.New(sess)

	zones, err := svc.ListHostedZonesByName(&route53.ListHostedZonesByNameInput{
		DNSName: aws.String(name),
	})

	if err != nil {
		return nil, err
	}

	return zones.HostedZones, nil
}

func Route53ListResourceRecordSets(zoneId string) ([]*route53.ResourceRecordSet, error) {
	sess := GetSession()
	svc := route53.New(sess)

	recordSets := []*route53.ResourceRecordSet{}
	pageNum := 0

	err := svc.ListResourceRecordSetsPages(&route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneId),
	},
		func(page *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			pageNum++
			recordSets = append(recordSets, page.ResourceRecordSets...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return recordSets, nil
}
