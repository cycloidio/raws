package billing

import (
	"testing"
)

func TestBillingManager_Import(t *testing.T) {

}

func TestBillingManager_getS3Filename(t *testing.T) {
	var accountId string = "111111111"
	var date string = "2000-01"
	var expectFilename string = accountId + "-aws-billing-detailed-line-items-with-resources-and-tags-" + date + ".csv.zip"

	mockedReader := mockAWSReader{
		accountID: accountId,
	}
	m := BillingManager{
		date:        date,
		s3Connector: mockedReader,
	}
	receivedFilename := m.getS3Filename()
	if m.getS3Filename() != expectFilename {
		t.Errorf("Invalid S3 filename, received: %q expected %q", receivedFilename, expectFilename)
	}
}
