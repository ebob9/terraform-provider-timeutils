// Copyright (C) Aaron Edwards
// SPDX-License-Identifier: MPL-2.0
//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	providerName = "terraform-provider-timeutils"
	namespace    = "ebob9"
	name         = "timeutils"
	version      = "0.1.0"
)

var (
	hostname = "registry.terraform.io"
	osArch   = runtime.GOOS + "_" + runtime.GOARCH
)

// Default target to run when none is specified
var Default = Install

// Build builds the provider binary
func Build() error {
	fmt.Println("Building provider...")
	return sh.Run("go", "build", "-o", providerName)
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	os.RemoveAll("bin")
	os.Remove(providerName)
	return nil
}

// Test runs the test suite
func Test() error {
	fmt.Println("Running tests...")
	return sh.Run("go", "test", "-count=1", "-parallel=4", "./...")
}

// TestAcc runs acceptance tests
func TestAcc() error {
	fmt.Println("Running acceptance tests...")
	env := map[string]string{"TF_ACC": "1"}
	return sh.RunWith(env, "go", "test", "-count=1", "-parallel=4", "-timeout", "10m", "-v", "./...")
}

// Install builds and installs the provider to Go bin directory
func Install() error {
	mg.Deps(Build)

	fmt.Println("Installing provider to Go bin directory...")

	// Get Go bin directory
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// Default GOPATH if not set
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		gopath = filepath.Join(homeDir, "go")
	}

	binDir := filepath.Join(gopath, "bin")

	// Create bin directory if it doesn't exist
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Copy the binary to the bin directory
	srcPath := providerName
	dstPath := filepath.Join(binDir, providerName)

	if err := sh.Copy(dstPath, srcPath); err != nil {
		return fmt.Errorf("failed to copy provider binary: %w", err)
	}

	// Make it executable
	if err := os.Chmod(dstPath, 0755); err != nil {
		return fmt.Errorf("failed to make provider executable: %w", err)
	}

	fmt.Printf("Provider installed to: %s\n", dstPath)
	fmt.Printf("Make sure %s is in your PATH\n", binDir)
	return nil
}

// InstallLocal builds and installs the provider locally for Terraform development
func InstallLocal() error {
	mg.Deps(Build)

	fmt.Println("Installing provider locally for Terraform...")

	// Determine the local plugin directory for Terraform
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	pluginDir := filepath.Join(
		homeDir,
		".terraform.d",
		"plugins",
		hostname,
		namespace,
		name,
		version,
		osArch,
	)

	// Create the plugin directory
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Copy the binary to the plugin directory
	srcPath := providerName
	dstPath := filepath.Join(pluginDir, providerName)

	if err := sh.Copy(dstPath, srcPath); err != nil {
		return fmt.Errorf("failed to copy provider binary: %w", err)
	}

	// Make it executable
	if err := os.Chmod(dstPath, 0755); err != nil {
		return fmt.Errorf("failed to make provider executable: %w", err)
	}

	fmt.Printf("Provider installed to: %s\n", dstPath)
	fmt.Println("You can now use this provider in Terraform configurations with:")
	fmt.Printf("  source = \"%s/%s\"\n", namespace, name)
	return nil
}

// Release builds binaries for multiple platforms
func Release() error {
	fmt.Println("Building release binaries...")

	if err := os.MkdirAll("bin", 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	platforms := []struct {
		goos, goarch string
	}{
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"linux", "386"},
		{"linux", "amd64"},
		{"linux", "arm"},
		{"linux", "arm64"},
		{"windows", "386"},
		{"windows", "amd64"},
		{"freebsd", "386"},
		{"freebsd", "amd64"},
		{"freebsd", "arm"},
		{"openbsd", "386"},
		{"openbsd", "amd64"},
		{"solaris", "amd64"},
	}

	for _, platform := range platforms {
		binary := fmt.Sprintf("%s_%s_%s_%s", providerName, version, platform.goos, platform.goarch)
		if platform.goos == "windows" {
			binary += ".exe"
		}

		outputPath := filepath.Join("bin", binary)
		fmt.Printf("Building %s...\n", binary)

		env := map[string]string{
			"GOOS":   platform.goos,
			"GOARCH": platform.goarch,
		}

		if err := sh.RunWith(env, "go", "build", "-o", outputPath); err != nil {
			return fmt.Errorf("failed to build %s: %w", binary, err)
		}
	}

	fmt.Println("Release binaries built successfully in ./bin/")
	return nil
}

// Tidy runs go mod tidy
func Tidy() error {
	fmt.Println("Running go mod tidy...")
	return sh.Run("go", "mod", "tidy")
}

// Fmt formats the Go code
func Fmt() error {
	fmt.Println("Formatting Go code...")
	return sh.Run("go", "fmt", "./...")
}

// Vet runs go vet
func Vet() error {
	fmt.Println("Running go vet...")
	return sh.Run("go", "vet", "./...")
}

// Check runs formatting, vetting, and testing
func Check() error {
	mg.Deps(Fmt, Vet, Test)
	return nil
}

// Dev does a full development cycle: tidy, check, install to Go bin
func Dev() error {
	mg.Deps(Tidy, Check, Install)
	return nil
}

// DevLocal does a full development cycle for Terraform: tidy, check, install locally
func DevLocal() error {
	mg.Deps(Tidy, Check, InstallLocal)
	return nil
}
