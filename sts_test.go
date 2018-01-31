package raws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type mockSTS struct {
	stsiface.STSAPI

	// Mocking of GetCallerIdentity
	gcio   *sts.GetCallerIdentityOutput
	gcierr error
}

func (m mockSTS) GetCallerIdentityWithContext(_ aws.Context, _ *sts.GetCallerIdentityInput, _ ...request.Option) (*sts.GetCallerIdentityOutput, error) {
	return m.gcio, m.gcierr
}
