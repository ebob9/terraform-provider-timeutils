// Copyright (C) Aaron Edwards
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStrftimeFunction(t *testing.T) {
	testCases := []struct {
		name      string
		format    string
		timestamp string
		expected  string
		expectErr bool
	}{
		{
			name:      "ISO date format",
			format:    "%Y-%m-%d",
			timestamp: "2024-01-15T10:30:00Z",
			expected:  "2024-01-15",
		},
		{
			name:      "US date format",
			format:    "%m/%d/%Y",
			timestamp: "2024-01-15T10:30:00Z",
			expected:  "01/15/2024",
		},
		{
			name:      "Full weekday and month",
			format:    "%A, %B %d, %Y",
			timestamp: "2024-01-15T10:30:00Z",
			expected:  "Monday, January 15, 2024",
		},
		{
			name:      "Time format 12-hour",
			format:    "%I:%M %p",
			timestamp: "2024-01-15T14:30:00Z",
			expected:  "02:30 PM",
		},
		{
			name:      "Time format 24-hour",
			format:    "%H:%M:%S",
			timestamp: "2024-01-15T14:30:45Z",
			expected:  "14:30:45",
		},
		{
			name:      "Custom format",
			format:    "Today is %A, the %d day of %B, %Y",
			timestamp: "2024-01-15T10:30:00Z",
			expected:  "Today is Monday, the 15 day of January, 2024",
		},
		{
			name:      "invalid timestamp",
			format:    "%Y-%m-%d",
			timestamp: "invalid",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := NewStrftimeFunction()

			// Create request with proper arguments
			req := function.RunRequest{}
			var argValues []attr.Value
			argValues = append(argValues, types.StringValue(tc.format))
			argValues = append(argValues, types.StringValue(tc.timestamp))
			req.Arguments = function.NewArgumentsData(argValues)

			// Create response
			resp := &function.RunResponse{}

			// Run function
			f.Run(context.Background(), req, resp)

			// Check for expected error
			if tc.expectErr {
				if resp.Error == nil {
					t.Errorf("Expected error for format %q, timestamp %q, but got none", tc.format, tc.timestamp)
				}
				return
			}

			// Check for unexpected error
			if resp.Error != nil {
				t.Errorf("Unexpected error for format %q, timestamp %q: %v", tc.format, tc.timestamp, resp.Error)
				return
			}

			// Check result - call the Value method
			if resp.Result == (function.ResultData{}) {
				t.Errorf("Expected result, got empty ResultData")
				return
			}

			resultValue := resp.Result.Value()
			result, ok := resultValue.(types.String)
			if !ok {
				t.Errorf("Expected types.String, got %T", resultValue)
				return
			}

			if result.ValueString() != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result.ValueString())
			}
		})
	}
}
