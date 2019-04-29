package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateName(t *testing.T) {
	tests := []struct {
		name string
		tmp  Function
		opt  string
	}{
		{
			name: "Basic",
			tmp: Function{
				Entity: "Entity",
			},
			opt: "GetEntity",
		},
		{
			name: "FnName",
			tmp: Function{
				FnName: "FnEntity",
			},
			opt: "FnEntity",
		},
		{
			name: "FilterByOwner",
			tmp: Function{
				Entity:        "Entity",
				FilterByOwner: "not-relevant",
			},
			opt: "GetOwnEntity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.opt, tt.tmp.Name())
		})
	}
}

func TestTemplateOutput(t *testing.T) {
	tests := []struct {
		name string
		tmp  Function
		opt  string
	}{
		{
			name: "Basic",
			tmp: Function{
				Service: "Service",
				Entity:  "Entity",
				Prefix:  "Prefix",
			},
			opt: "Service.PrefixEntityOutput",
		},
		{
			name: "FnOutput",
			tmp: Function{
				Service:  "Service",
				FnOutput: "FnOutput",
			},
			opt: "Service.FnOutput",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.opt, tt.tmp.Output())
		})
	}
}

func TestTemplateInput(t *testing.T) {
	tests := []struct {
		name string
		tmp  Function
		opt  string
	}{
		{
			name: "Basic",
			tmp: Function{
				Service: "Service",
				Entity:  "Entity",
				Prefix:  "Prefix",
			},
			opt: "Service.PrefixEntityInput",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.opt, tt.tmp.Input())
		})
	}
}

func TestTemplateSignature(t *testing.T) {
	tests := []struct {
		name string
		tmp  Function
		opt  string
	}{
		{
			name: "Basic",
			tmp: Function{
				Service: "Service",
				Entity:  "Entity",
				Prefix:  "Prefix",
			},
			opt: "GetEntity (ctx context.Context, input *Service.PrefixEntityInput) (map[string]Service.PrefixEntityOutput, error)",
		},
		{
			name: "FnSignature",
			tmp: Function{
				FnSignature: "SomeSignature",
			},
			opt: "SomeSignature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.opt, tt.tmp.Signature())
		})
	}
}

func TestTemplateExecute(t *testing.T) {
	tests := []struct {
		name string
		tmp  Function
		opt  string
	}{
		{
			name: "Basic",
			tmp: Function{
				FnSignature: "Signature",
				Service:     "Service",
				Entity:      "Entity",
				Prefix:      "Prefix",
			},
			opt: `
			func (c *connector) Signature {
				var errs Errors
				var regionsOpts = map[string]Service.PrefixEntityOutput{}

				for _, svc := range c.svcs {
					if svc.Service == nil {
						svc.Service = Service.New(svc.session)
					}

					opt, err := svc.Service.PrefixEntityWithContext(ctx, input)
					if err != nil {
						errs = append(errs, NewError(svc.region, Service.ServiceName, err))
					} else {
						regionsOpts[svc.region] = *opt
					}
				}

				if errs != nil {
					return regionsOpts, errs
				}

				return regionsOpts, nil
			}`,
		},
		{
			name: "FilterByOwner",
			tmp: Function{
				FilterByOwner: "OwnerField",
				FnSignature:   "Signature",
				Service:       "Service",
				Entity:        "Entity",
				Prefix:        "Prefix",
			},
			opt: `
			func (c *connector) Signature {
				var errs Errors
				var regionsOpts = map[string]Service.PrefixEntityOutput{}

				if input == nil {
					input = &Service.PrefixEntityInput{}
				}
				input.OwnerField = append(input.OwnerField, c.accountID)

				for _, svc := range c.svcs {
					if svc.Service == nil {
						svc.Service = Service.New(svc.session)
					}

					opt, err := svc.Service.PrefixEntityWithContext(ctx, input)
					if err != nil {
						errs = append(errs, NewError(svc.region, Service.ServiceName, err))
					} else {
						regionsOpts[svc.region] = *opt
					}
				}

				if errs != nil {
					return regionsOpts, errs
				}

				return regionsOpts, nil
			}`,
		},
		{
			name: "NoGenerateFn",
			tmp: Function{
				NoGenerateFn: true,
			},
			opt: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buff := bytes.Buffer{}
			err := tt.tmp.Execute(&buff)
			require.NoError(t, err)
			ttopt := strings.Join(strings.Fields(tt.opt), " ")
			buffs := strings.Join(strings.Fields(buff.String()), " ")
			assert.Equal(t, ttopt, buffs)
		})
	}
}
