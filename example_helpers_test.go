package raws_test

import (
	"errors"
	"os"
)

func getAWSKeys() (accessKey string, secretKey string, err error) {
	var ok bool
	if accessKey, ok = os.LookupEnv("AWS_ACCESS_KEY"); !ok {
		return "", "", errors.New("AWS_ACCESS_KEY environment variable isn't set")
	}

	if secretKey, ok = os.LookupEnv("AWS_SECRET_KEY"); !ok {
		return "", "", errors.New("AWS_SECRET_KEY environment variable isn't set")
	}

	return accessKey, secretKey, nil
}
