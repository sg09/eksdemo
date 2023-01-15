package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
)

type CloudtrailClient struct {
	*cloudtrail.Client
}

func NewCloudtrailClient() *CloudtrailClient {
	return &CloudtrailClient{cloudtrail.NewFromConfig(GetConfig())}
}

// func NewCloudtrailEventFilter(eventId string) types.LookupAttribute {
// 	return types.LookupAttribute{
// 		AttributeKey:   types.LookupAttributeKeyEventId,
// 		AttributeValue: aws.String(eventId),
// 	}
// }

func (c *CloudtrailClient) ListTrails() ([]types.TrailInfo, error) {
	trails := []types.TrailInfo{}
	pageNum := 0

	paginator := cloudtrail.NewListTrailsPaginator(c.Client, &cloudtrail.ListTrailsInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		trails = append(trails, out.Trails...)
		pageNum++
	}

	return trails, nil
}

// Looks up management events or CloudTrail Insights events that are captured by CloudTrail.
// You can look up events that occurred in a region within the last 90 days.
func (c *CloudtrailClient) LookupEvents(insights bool, filters []types.LookupAttribute) ([]types.Event, error) {
	events := []types.Event{}
	pageNum := 0

	input := cloudtrail.LookupEventsInput{
		LookupAttributes: filters,
	}

	if insights {
		input.EventCategory = types.EventCategoryInsight
	}

	paginator := cloudtrail.NewLookupEventsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		events = append(events, out.Events...)
		pageNum++
	}

	return events, nil
}
