package core

import (
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type mockSTS struct {
	stsiface.STSAPI

	// Mocking of GetCallerIdentity
	gcio   *sts.GetCallerIdentityOutput
	gcierr error
}

func (m mockSTS) GetCallerIdentity(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return m.gcio, m.gcierr
}
