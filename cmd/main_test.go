package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_generate(t *testing.T) {
	buff := bytes.Buffer{}
	fns := []Function{
		Function{
			Entity:  "Instances",
			Prefix:  "Describe",
			Service: "ec2",
			Documentation: `
			// GetInstances returns all EC2 instances based on the input given.
			// Returned values are commented in the interface doc comment block.
			`,
		},
		Function{
			FnSignature:  "DownloadObject(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error)",
			NoGenerateFn: true,
			Documentation: `
			// DownloadObject downloads an object in a bucket based on the input given
			`,
		},
	}
	exopt, err := ioutil.ReadFile("./testdata/generated.go")
	err = generate(&buff, fns)
	require.NoError(t, err)
	assert.Equal(t, exopt, buff.Bytes())
}
