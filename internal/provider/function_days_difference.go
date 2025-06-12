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

var _ function.Function = &DaysDifferenceFunction{}

type DaysDifferenceFunction struct{}

func NewDaysDifferenceFunction() function.Function {
	return &DaysDifferenceFunction{}
}

func (f *DaysDifferenceFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "days_difference"
}

func (f *DaysDifferenceFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Calculate days between timestamps",
		Description: "Returns the number of complete days between two RFC3339 timestamps as an integer string. Positive if end is after start.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "start_timestamp",
				Description: "RFC3339 formatted start timestamp",
			},
			function.StringParameter{
				Name:        "end_timestamp",
				Description: "RFC3339 formatted end timestamp",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *DaysDifferenceFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var startTimestamp, endTimestamp string

	// Get both arguments at once
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &startTimestamp, &endTimestamp))
	if resp.Error != nil {
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimestamp)
	if err != nil {
		resp.Error = function.NewFuncError("Invalid start timestamp: " + err.Error())
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimestamp)
	if err != nil {
		resp.Error = function.NewFuncError("Invalid end timestamp: " + err.Error())
		return
	}

	duration := endTime.Sub(startTime)
	days := int(duration.Hours() / 24)

	resp.Result = function.NewResultData(types.StringValue(strconv.Itoa(days)))
}
