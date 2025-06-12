// Copyright (C) Aaron Edwards
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestParseRFC3339Function(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expected  map[string]string
		expectErr bool
	}{
		{
			name:  "valid RFC3339 timestamp",
			input: "2024-01-15T10:30:45Z",
			expected: map[string]string{
				"year":    "2024",
				"month":   "1",
				"day":     "15",
				"hour":    "10",
				"minute":  "30",
				"second":  "45",
				"unix":    "1705314645",
				"weekday": "1", // Monday
			},
		},
		{
			name:  "different date",
			input: "2023-12-25T23:59:59Z",
			expected: map[string]string{
				"year":    "2023",
				"month":   "12",
				"day":     "25",
				"hour":    "23",
				"minute":  "59",
				"second":  "59",
				"unix":    "1703548799",
				"weekday": "1", // Monday
			},
		},
		{
			name:      "invalid timestamp",
			input:     "invalid",
			expectErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := NewParseRFC3339Function()

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
			result := resultValue.(types.String)

			// Parse the JSON result
			var actual map[string]string
			if err := json.Unmarshal([]byte(result.ValueString()), &actual); err != nil {
				t.Errorf("Failed to parse result JSON: %v", err)
				return
			}

			// Check each expected field
			for key, expectedValue := range tc.expected {
				if actual[key] != expectedValue {
					t.Errorf("Expected %s=%q, got %s=%q", key, expectedValue, key, actual[key])
				}
			}
		})
	}
}
