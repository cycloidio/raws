package billing

const (
	CSVErrorType = iota
	ConvertErrorType
	DynamoDBErrorType
	S3ErrorType
)

type BillingError struct {
	Err  error
	flag int
}

func (b BillingError) Error() string {
	return b.Err.Error()
}

func NewBillingError(err error, flag int) error {
	return &BillingError{
		Err:  err,
		flag: flag,
	}
}

func NewCSVError(err error) error {
	return &BillingError{
		Err:  err,
		flag: CSVErrorType,
	}
}

func NewConvertError(err error) error {
	return &BillingError{
		Err:  err,
		flag: ConvertErrorType,
	}
}

func NewDynamoDBError(err error) error {
	return &BillingError{
		Err:  err,
		flag: DynamoDBErrorType,
	}
}

func NewS3Error(err error) error {
	return &BillingError{
		Err:  err,
		flag: S3ErrorType,
	}
}
func IsConvertError(err error) bool {
	if v, ok := err.(*BillingError); ok {
		return v.flag == ConvertErrorType
	}
	return false
}

func IsCSVError(err error) bool {
	if v, ok := err.(*BillingError); ok {
		return v.flag == CSVErrorType
	}
	return false
}

func IsDynamoDBError(err error) bool {
	if v, ok := err.(*BillingError); ok {
		return v.flag == DynamoDBErrorType
	}
	return false
}

func IsS3Error(err error) bool {
	if v, ok := err.(*BillingError); ok {
		return v.flag == S3ErrorType
	}
	return false
}
