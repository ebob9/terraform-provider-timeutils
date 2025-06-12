// Copyright (C) Aaron Edwards
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/lestrrat-go/strftime"
)

var _ function.Function = &StrftimeFunction{}

type StrftimeFunction struct{}

func NewStrftimeFunction() function.Function {
	return &StrftimeFunction{}
}

func (f *StrftimeFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "strftime"
}

func (f *StrftimeFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Format timestamp using strftime",
		Description: "Takes an RFC3339 timestamp and formats it using strftime format specifiers (e.g., '%Y-%m-%d %H:%M:%S')",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "format",
				Description: "strftime format string (e.g., '%Y-%m-%d %H:%M:%S')",
			},
			function.StringParameter{
				Name:        "timestamp",
				Description: "RFC3339 formatted timestamp string",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *StrftimeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var format, timestamp string

	// Get both arguments at once
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &format, &timestamp))
	if resp.Error != nil {
		return
	}

	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		resp.Error = function.NewFuncError("Invalid RFC3339 timestamp: " + err.Error())
		return
	}

	formatter, err := strftime.New(format)
	if err != nil {
		resp.Error = function.NewFuncError("Invalid strftime format: " + err.Error())
		return
	}

	formatted := formatter.FormatString(t)
	resp.Result = function.NewResultData(types.StringValue(formatted))
}
