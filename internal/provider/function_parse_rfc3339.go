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

var _ function.Function = &ParseRFC3339Function{}

type ParseRFC3339Function struct{}

func NewParseRFC3339Function() function.Function {
	return &ParseRFC3339Function{}
}

func (f *ParseRFC3339Function) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "parse_rfc3339"
	resp.Description = "Parse an RFC3339 timestamp and return components"
}

func (f *ParseRFC3339Function) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Parse RFC3339 timestamp into components",
		Description: "Parses an RFC3339 timestamp and returns a JSON string with year, month, day, hour, minute, second, and unix timestamp",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "timestamp",
				Description: "RFC3339 formatted timestamp string",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *ParseRFC3339Function) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
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

	// Return JSON string with components
	result := `{"year":"` + strconv.Itoa(t.Year()) +
		`","month":"` + strconv.Itoa(int(t.Month())) +
		`","day":"` + strconv.Itoa(t.Day()) +
		`","hour":"` + strconv.Itoa(t.Hour()) +
		`","minute":"` + strconv.Itoa(t.Minute()) +
		`","second":"` + strconv.Itoa(t.Second()) +
		`","unix":"` + strconv.FormatInt(t.Unix(), 10) +
		`","weekday":"` + strconv.Itoa(int(t.Weekday())) + `"}`

	resp.Result = function.NewResultData(types.StringValue(result))
}
