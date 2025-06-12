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

func TestDaysDifferenceFunction(t *testing.T) {
	testCases := []struct {
		name      string
		start     string
		end       string
		expected  string
		expectErr bool
	}{
		{
			name:     "same day",
			start:    "2024-01-15T10:30:00Z",
			end:      "2024-01-15T15:45:00Z",
			expected: "0",
		},
		{
			name:     "one day difference",
			start:    "2024-01-15T10:30:00Z",
			end:      "2024-01-16T10:30:00Z",
			expected: "1",
		},
		{
			name:     "multiple days",
			start:    "2024-01-15T10:30:00Z",
			end:      "2024-01-20T10:30:00Z",
			expected: "5",
		},
		{
			name:     "negative difference",
			start:    "2024-01-20T10:30:00Z",
			end:      "2024-01-15T10:30:00Z",
			expected: "-5",
		},
		{
			name:     "cross month boundary",
			start:    "2024-01-30T10:30:00Z",
			end:      "2024-02-05T10:30:00Z",
			expected: "6",
		},
		{
			name:     "leap year",
			start:    "2024-02-28T10:30:00Z",
			end:      "2024-03-01T10:30:00Z",
			expected: "2", // 2024 is a leap year
		},
		{
			name:      "invalid start timestamp",
			start:     "invalid",
			end:       "2024-01-15T10:30:00Z",
			expectErr: true,
		},
		{
			name:      "invalid end timestamp",
			start:     "2024-01-15T10:30:00Z",
			end:       "invalid",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := NewDaysDifferenceFunction()

			// Create request with proper arguments
			req := function.RunRequest{}
			var argValues []attr.Value
			argValues = append(argValues, types.StringValue(tc.start))
			argValues = append(argValues, types.StringValue(tc.end))
			req.Arguments = function.NewArgumentsData(argValues)

			// Create response
			resp := &function.RunResponse{}

			// Run function
			f.Run(context.Background(), req, resp)

			// Check for expected error
			if tc.expectErr {
				if resp.Error == nil {
					t.Errorf("Expected error for inputs %q, %q, but got none", tc.start, tc.end)
				}
				return
			}

			// Check for unexpected error
			if resp.Error != nil {
				t.Errorf("Unexpected error for inputs %q, %q: %v", tc.start, tc.end, resp.Error)
				return
			}

			// Check result - call the Value method
			if resp.Result == (function.ResultData{}) {
				t.Errorf("Expected result, got empty ResultData")
				return
			}

			resultValue := resp.Result.Value()
			result := resultValue.(types.String)
			if result.ValueString() != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result.ValueString())
			}
		})
	}
}
