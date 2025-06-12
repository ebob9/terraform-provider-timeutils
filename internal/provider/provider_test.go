// Copyright (C) Aaron Edwards
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"
)

func TestProvider(t *testing.T) {
	// Basic test to ensure provider can be instantiated
	provider := New("test")()

	if provider == nil {
		t.Error("Expected provider to be created, got nil")
	}
}

func TestProviderMetadata(t *testing.T) {
	provider := New("test")()

	// This is a basic smoke test - more comprehensive tests would require
	// setting up the full provider server infrastructure
	if provider == nil {
		t.Error("Provider should not be nil")
	}
}
