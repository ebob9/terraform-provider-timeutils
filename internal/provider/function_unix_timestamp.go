// Copyright (C) Aaron Edwards
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &UnixTimestampFunction{}

type UnixTimestampFunction struct{}

func NewUnixTimestampFunction() function.Function {
	return &UnixTimestampFunction{}
}

func (f *UnixTimestampFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "unix_timestamp"
}

func (f *UnixTimestampFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Convert RFC3339 timestamp to Unix timestamp",
		Description: "Takes an RFC3339 formatted timestamp string and returns the Unix timestamp (seconds since epoch) as a string.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "timestamp",
				Description: "RFC3339 formatted timestamp string",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *UnixTimestampFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var timestamp string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &timestamp))
	if resp.Error != nil {
		return
	}

	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		resp.Error = function.NewFuncError("Invalid RFC3339 timestamp: " + err.Error())
		return
	}

	unixTime := strconv.FormatInt(t.Unix(), 10)
	resp.Result = function.NewResultData(types.StringValue(unixTime))
}
