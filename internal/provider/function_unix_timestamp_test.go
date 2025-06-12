// Copyright (C) Aaron Edwards
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUnixTimestampFunction(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expected  string
		expectErr bool
	}{
		{
			name:     "valid RFC3339 timestamp",
			input:    "2024-01-15T10:30:00Z",
			expected: "1705314600", // Unix timestamp for 2024-01-15T10:30:00Z
		},
		{
			name:     "valid RFC3339 with timezone",
			input:    "2024-01-15T10:30:00-08:00",
			expected: "1705343400", // Unix timestamp adjusted for PST
		},
		{
			name:      "invalid timestamp format",
			input:     "2024-01-15 10:30:00",
			expectErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			expectErr: true,
		},
		{
			name:      "invalid date",
			input:     "2024-13-45T25:70:80Z",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := NewUnixTimestampFunction()

			// Create request with proper arguments
			req := function.RunRequest{}
			var argValues []attr.Value
			argValues = append(argValues, types.StringValue(tc.input))
			req.Arguments = function.NewArgumentsData(argValues)

			// Create response
			resp := &function.RunResponse{}

			// Run function
			f.Run(context.Background(), req, resp)

			// Check for expected error
			if tc.expectErr {
				if resp.Error == nil {
					t.Errorf("Expected error for input %q, but got none", tc.input)
				}
				return
			}

			// Check for unexpected error
			if resp.Error != nil {
				t.Errorf("Unexpected error for input %q: %v", tc.input, resp.Error)
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
