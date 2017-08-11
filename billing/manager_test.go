package billing

import (
	"errors"
	"testing"
)

func TestBillingManager_Import(t *testing.T) {
	const (
		date   = "2017-01"
		bucket = "test-bucket"
	)
	var awsReader mockAWSReader = mockAWSReader{
		accountID: "111111111111",
	}

	tests := []struct {
		name             string
		mockedChecker    mockChecker
		mockedDownloader mockDownloader
		mockedLoader     mockLoader
		mockedInjector   mockInjector
		expectedError    error
	}{
		{name: "no errors during any calls",
			mockedChecker: mockChecker{
				cb: true,
				ce: nil,
			},
			mockedDownloader: mockDownloader{
				ds: "",
				de: nil,
				us: "",
				ue: nil,
			},
			mockedInjector: mockInjector{
				crde: nil,
				crpe: nil,
			},
			mockedLoader:  mockLoader{},
			expectedError: nil,
		},
		{name: "errors during Check",
			mockedChecker: mockChecker{
				cb: true,
				ce: errors.New("error during check"),
			},
			expectedError: errors.New("error during check"),
		},
		{name: "errors during Download",
			mockedChecker: mockChecker{
				cb: true,
				ce: nil,
			},
			mockedDownloader: mockDownloader{
				ds: "",
				de: nil,
				us: "",
				ue: errors.New("error during download"),
			},
			expectedError: errors.New("error during download"),
		},
		{name: "errors during Unzip",
			mockedChecker: mockChecker{
				cb: true,
				ce: nil,
			},
			mockedDownloader: mockDownloader{
				ds: "",
				de: nil,
				us: "",
				ue: errors.New("error during unzip"),
			},
			expectedError: errors.New("error during unzip"),
		},
		{name: "errors during Inject",
			mockedChecker: mockChecker{
				cb: true,
				ce: nil,
			},
			mockedDownloader: mockDownloader{
				ds: "",
				de: nil,
				us: "",
				ue: nil,
			},
			mockedInjector: mockInjector{
				crde: nil,
				crpe: errors.New("error during CreateReport"),
			},
			mockedLoader:  mockLoader{},
			expectedError: errors.New("error during CreateReport"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &BillingManager{
				s3Connector: awsReader,
				checker:     tt.mockedChecker,
				downloader:  tt.mockedDownloader,
				loader:      tt.mockedLoader,
				injector:    tt.mockedInjector,
			}
			err := m.Import(date, bucket)
			checkErrors(t, tt.name, i, err, tt.expectedError)
		})
	}
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
