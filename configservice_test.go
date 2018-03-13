package raws

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/configservice/configserviceiface"
)

type mockServiceConfig struct {
	configserviceiface.ConfigServiceAPI

	gdrcctx    context.Context
	gdrcinput  *configservice.GetDiscoveredResourceCountsInput
	gdrcoutput *configservice.GetDiscoveredResourceCountsOutput
	gdrcerr    error
}

func (m *mockServiceConfig) GetDiscoveredResourceCountsWithContext(
	ctx aws.Context, input *configservice.GetDiscoveredResourceCountsInput, _ ...request.Option,
) (*configservice.GetDiscoveredResourceCountsOutput, error) {
	m.gdrcctx = ctx
	m.gdrcinput = input
	return m.gdrcoutput, m.gdrcerr
}

func TestGetRecordedResourceCounts(t *testing.T) {
	t.Run("input parameter is nil", testGetRecordedResourceCountsInputParameterNil)
	t.Run("input parameter isn't nil", testGetRecordedResourceCountsInputParameterNotNil)
	t.Run("one region with error", testGetRecordedResourceCountsOneRegionError)
	t.Run("one region no error", testGetRecordedResourceCountsOneRegionNoError)
	t.Run("multiple region no error", testGetRecordedResourceCountsMultipleRegionsNoError)
	t.Run("multiple region with error", testGetRecordedResourceCountsMultipleRegionsOneError)
}

// testGetRecordedResourceCountsInputParemeterNil checks that when the input
// parameters passed to the AWSReader instance GetRecordedResourceCounts method
// is nil, nil is passed down to the AWS SDK GetDiscoveredResourceCounts
func testGetRecordedResourceCountsInputParameterNil(t *testing.T) {
	var (
		ctx  = context.Background()
		mock = &mockServiceConfig{
			gdrcoutput: &configservice.GetDiscoveredResourceCountsOutput{},
			gdrcerr:    nil,
		}

		awsReader = &connector{
			svcs: []*serviceConnector{{region: "test", configservice: mock}},
		}
	)

	_, _ = awsReader.GetRecordedResourceCounts(ctx, nil)
	if !reflect.DeepEqual(mock.gdrcctx, ctx) {
		t.Errorf(
			"Passed context to AWS SDK isn't the one passed to the AWSReader method. received=%+v | expected=%+v",
			mock.gdrcctx,
			ctx,
		)
	}

	if mock.gdrcinput != nil {
		t.Errorf(
			"Passed input to AWS SDK isn't nil. received=%+v | expected=nil",
			mock.gdrcinput,
		)
	}
}

// testGetRecordedResourceCountsInputParemeterNotNil checks that when the input
// parameters passed to the AWSReader instance GetRecordedResourceCounts method
// is passed down to the AWS SDK GetDiscoveredResourceCounts
func testGetRecordedResourceCountsInputParameterNotNil(t *testing.T) {
	var (
		ctx   = context.Background()
		input = &configservice.GetDiscoveredResourceCountsInput{
			NextToken: aws.String("some-test-token"),
		}
		mock = &mockServiceConfig{
			gdrcoutput: &configservice.GetDiscoveredResourceCountsOutput{},
			gdrcerr:    nil,
		}

		awsReader = &connector{
			svcs: []*serviceConnector{{region: "test", configservice: mock}},
		}
	)

	_, _ = awsReader.GetRecordedResourceCounts(ctx, input)
	if !reflect.DeepEqual(mock.gdrcctx, ctx) {
		t.Errorf(
			"Passed context to AWS SDK isn't the one passed to the AWSReader method. received=%+v | expected=%+v",
			mock.gdrcctx,
			ctx,
		)
	}

	if !reflect.DeepEqual(mock.gdrcinput, input) {
		t.Errorf(
			"Passed input to AWS SDK isn't the one passed to the AWSReader method. received=%+v | expected=%+v",
			mock.gdrcinput,
			input,
		)
	}
}

// testGetRecordedResourceCountsOneRegionError checks that when the AWS SDK
// GetDiscoveredResourceCounts return an error in the only one region, it's
// returned.
func testGetRecordedResourceCountsOneRegionError(t *testing.T) {
	var (
		ctx  = context.Background()
		mock = &mockServiceConfig{
			gdrcoutput: nil,
			gdrcerr: Errors{Error{
				err:     errors.New("error with test"),
				region:  "test",
				service: configservice.ServiceName,
			}},
		}

		awsReader = &connector{
			svcs: []*serviceConnector{{region: "test", configservice: mock}},
		}
	)

	var _, err = awsReader.GetRecordedResourceCounts(ctx, nil)

	var expectedErrors = Errors{
		Error{
			err:     mock.gdrcerr,
			region:  "test",
			service: configservice.ServiceName,
		},
	}
	checkError(t, err, expectedErrors)
}

// testGetRecordedResourceCountsOneRegionNoError checks that when the AWS SDK
// GetDiscoveredResourceCounts return a no error result in the only one region,
// the result is returned associated to the region.
func testGetRecordedResourceCountsOneRegionNoError(t *testing.T) {
	var (
		ctx  = context.Background()
		mock = &mockServiceConfig{
			gdrcoutput: &configservice.GetDiscoveredResourceCountsOutput{
				ResourceCounts: []*configservice.ResourceCount{
					{ResourceType: aws.String("some-resource-1")},
					{ResourceType: aws.String("some-resource-2")},
				},
			},
			gdrcerr: nil,
		}
		awsReader = &connector{
			svcs: []*serviceConnector{{region: "test", configservice: mock}},
		}
	)

	var counts, err = awsReader.GetRecordedResourceCounts(ctx, nil)
	if err != nil {
		t.Errorf("Unexpected returned error. received= %+v | expected=nil", err)
	}

	var expectedCounts = map[string]configservice.GetDiscoveredResourceCountsOutput{
		"test": *mock.gdrcoutput,
	}
	if !reflect.DeepEqual(counts, expectedCounts) {
		t.Errorf("resource counts: received=%+v | expected=%+v", counts, expectedCounts)
	}
}

// testGetRecordedResourceCountsMultipleRegionsNoError checks that when the AWS
// SDK GetDiscoveredResourceCounts return a no error result in all the regions,
// and the results are associated to the corresponding regions.
func testGetRecordedResourceCountsMultipleRegionsNoError(t *testing.T) {
	var (
		ctx      = context.Background()
		mockReg1 = &mockServiceConfig{
			gdrcoutput: &configservice.GetDiscoveredResourceCountsOutput{
				ResourceCounts: []*configservice.ResourceCount{
					{ResourceType: aws.String("some-resource-reg1-1")},
					{ResourceType: aws.String("some-resource-reg2-2")},
				},
			},
			gdrcerr: nil,
		}
		mockReg2 = &mockServiceConfig{
			gdrcoutput: &configservice.GetDiscoveredResourceCountsOutput{
				ResourceCounts: []*configservice.ResourceCount{
					{ResourceType: aws.String("some-resource-reg2-1")},
					{ResourceType: aws.String("some-resource-reg2-2")},
				},
			},
			gdrcerr: nil,
		}
		awsReader = &connector{
			svcs: []*serviceConnector{
				{region: "region-1", configservice: mockReg1},
				{region: "region-2", configservice: mockReg2},
			},
		}
	)

	var counts, err = awsReader.GetRecordedResourceCounts(ctx, nil)
	if err != nil {
		t.Errorf("Unexpected returned error. received= %+v | expected=nil", err)
	}

	var expectedCounts = map[string]configservice.GetDiscoveredResourceCountsOutput{
		"region-1": *mockReg1.gdrcoutput,
		"region-2": *mockReg2.gdrcoutput,
	}
	if !reflect.DeepEqual(counts, expectedCounts) {
		t.Errorf("resource counts: received=%+v | expected=%+v", counts, expectedCounts)
	}
}

// testGetRecordedResourceCountsMultipleRegionsOneError checks that when the AWS
// SDK GetDiscoveredResourceCounts return the error result in the correct region,
// and it returns the results (partial) to the regions which didn't respond with
// an error
func testGetRecordedResourceCountsMultipleRegionsOneError(t *testing.T) {
	var (
		ctx      = context.Background()
		mockReg1 = &mockServiceConfig{
			gdrcoutput: nil,
			gdrcerr:    errors.New("error with region-1"),
		}
		mockReg2 = &mockServiceConfig{
			gdrcoutput: &configservice.GetDiscoveredResourceCountsOutput{
				ResourceCounts: []*configservice.ResourceCount{
					{ResourceType: aws.String("some-resource-reg2-1")},
					{ResourceType: aws.String("some-resource-reg2-2")},
				},
			},
			gdrcerr: nil,
		}
		awsReader = &connector{
			svcs: []*serviceConnector{
				{region: "region-1", configservice: mockReg1},
				{region: "region-2", configservice: mockReg2},
			},
		}
	)

	var counts, err = awsReader.GetRecordedResourceCounts(ctx, nil)

	var expectedErrors = Errors{
		Error{
			err:     mockReg1.gdrcerr,
			region:  "region-1",
			service: configservice.ServiceName,
		},
	}
	checkError(t, err, expectedErrors)

	var expectedCounts = map[string]configservice.GetDiscoveredResourceCountsOutput{
		"region-2": *mockReg2.gdrcoutput,
	}
	if !reflect.DeepEqual(counts, expectedCounts) {
		t.Errorf("resource counts: received=%+v | expected=%+v", counts, expectedCounts)
	}
}
