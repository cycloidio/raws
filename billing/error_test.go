package billing

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewBillingError(t *testing.T) {
	t.Run("NewBillingError", func(t *testing.T) {
		e := &BillingError{
			Err:  errors.New("test"),
			flag: CSVErrorType,
		}
		ce := NewBillingError(errors.New("test"), CSVErrorType)
		if !reflect.DeepEqual(e, ce) {
			t.Errorf("NewBillingError: received=%+v | expected=%+v",
				ce, e)
		}
		if ce.Error() != "test" {
			t.Errorf("NewBillingError Error(): received=%+v | expected=%+v",
				ce.Error(), "test")
		}
	})
}

func TestNewCSVError(t *testing.T) {
	t.Run("Valid new CSV Error", func(t *testing.T) {
		e := &BillingError{
			Err:  errors.New("test"),
			flag: CSVErrorType,
		}
		ce := NewCSVError(errors.New("test"))
		if !reflect.DeepEqual(e, ce) {
			t.Errorf("NewCSVError: received=%+v | expected=%+v",
				ce, e)
		}
	})
}

func TestIsCSVError(t *testing.T) {
	t.Run("Is a CSV Error", func(t *testing.T) {
		e := NewCSVError(errors.New("test"))
		if !IsCSVError(e) {
			t.Errorf("IsCSVError should be a CSV error: sent=%+v", e)
		}
	})
	t.Run("Is not a CSV Error", func(t *testing.T) {
		e := errors.New("test")
		if IsCSVError(e) {
			t.Errorf("IsCSVError should not be a CSV error: sent=%+v", e)
		}
	})
}

func TestNewConvertError(t *testing.T) {
	t.Run("Valid new Convert Error", func(t *testing.T) {
		e := &BillingError{
			Err:  errors.New("test"),
			flag: ConvertErrorType,
		}
		ce := NewConvertError(errors.New("test"))
		if !reflect.DeepEqual(e, ce) {
			t.Errorf("NewConvertError: received=%+v | expected=%+v",
				ce, e)
		}
	})
}

func TestIsConvertError(t *testing.T) {
	t.Run("Is a Convert Error", func(t *testing.T) {
		e := NewConvertError(errors.New("test"))
		if !IsConvertError(e) {
			t.Errorf("IsConvertError should be a Convert error: sent=%+v", e)
		}
	})
	t.Run("Is not a Convert Error", func(t *testing.T) {
		e := errors.New("test")
		if IsConvertError(e) {
			t.Errorf("IsConvertError should be a Convert error: sent=%+v", e)
		}
	})
}

func TestNewDynamoDBError(t *testing.T) {
	t.Run("Valid new DynamoDB Error", func(t *testing.T) {
		e := &BillingError{
			Err:  errors.New("test"),
			flag: DynamoDBErrorType,
		}
		ce := NewDynamoDBError(errors.New("test"))
		if !reflect.DeepEqual(e, ce) {
			t.Errorf("NewDynamoDBError: received=%+v | expected=%+v",
				ce, e)
		}
	})
}

func TestIsDynamoDBError(t *testing.T) {
	t.Run("Is a DynamoDB Error", func(t *testing.T) {
		e := NewDynamoDBError(errors.New("test"))
		if !IsDynamoDBError(e) {
			t.Errorf("IsDynamoDBError should be a DynamoDB error: sent=%+v", e)
		}
	})
	t.Run("Is not a DynamoDB Error", func(t *testing.T) {
		e := errors.New("test")
		if IsDynamoDBError(e) {
			t.Errorf("IsDynamoDBError should be a DynamoDB error: sent=%+v", e)
		}
	})
}

func TestNewS3Error(t *testing.T) {
	t.Run("Valid new S3 Error", func(t *testing.T) {
		e := &BillingError{
			Err:  errors.New("test"),
			flag: S3ErrorType,
		}
		ce := NewS3Error(errors.New("test"))
		if !reflect.DeepEqual(e, ce) {
			t.Errorf("NewS3Error: received=%+v | expected=%+v",
				ce, e)
		}
	})
}

func TestIsS3Error(t *testing.T) {
	t.Run("Is a S3 Error", func(t *testing.T) {
		e := NewS3Error(errors.New("test"))
		if !IsS3Error(e) {
			t.Errorf("IsS3Error should be an S3 error: sent=%+v", e)
		}
	})
	t.Run("Is not a S3 Error", func(t *testing.T) {
		e := errors.New("test")
		if IsS3Error(e) {
			t.Errorf("IsS3Error should be an S3 error: sent=%+v", e)
		}
	})
}
