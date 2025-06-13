# GNUmakefile for compatibility with existing workflows
# Delegates to Mage for actual build tasks

.PHONY: default install build test testacc generate fmt vet clean release

# Check if mage is installed
MAGE := $(shell which mage)
ifeq ($(MAGE),)
$(error "mage is not installed. Install with: go install github.com/magefile/mage@latest")
endif

# Default target
default: install

# Installation targets
install:
	@$(MAGE) install

install-local:
	@$(MAGE) installlocal

# Build targets
build:
	@$(MAGE) build

# Test targets
test:
	@$(MAGE) test

testacc:
	@$(MAGE) testacc

testunit:
	@$(MAGE) testunit

testcoverage:
	@$(MAGE) testcoverage

# Code quality targets
fmt:
	@$(MAGE) fmt

vet:
	@$(MAGE) vet

check:
	@$(MAGE) check

# Generate target (for compatibility with workflows)
generate:
	@echo "Running go generate..."
	@go generate ./...
	@echo "Formatting generated code..."
	@$(MAGE) fmt

# Dependency management
tidy:
	@$(MAGE) tidy

# Clean up
clean:
	@$(MAGE) clean

# Release builds
release:
	@$(MAGE) release

# Development workflows
dev:
	@$(MAGE) dev

devlocal:
	@$(MAGE) devlocal

# Help target
help:
	@echo "Available targets:"
	@echo "  install      - Build and install to Go bin directory"
	@echo "  install-local- Install to Terraform plugins directory"
	@echo "  build        - Build the provider binary"
	@echo "  test         - Run unit tests"
	@echo "  testacc      - Run acceptance tests (requires TF_ACC=1)"
	@echo "  testunit     - Run only unit tests"
	@echo "  testcoverage - Run tests with coverage"
	@echo "  fmt          - Format Go code"
	@echo "  vet          - Run go vet"
	@echo "  check        - Run fmt, vet, and test"
	@echo "  generate     - Run go generate and format code"
	@echo "  tidy         - Run go mod tidy"
	@echo "  clean        - Remove build artifacts"
	@echo "  release      - Build release binaries for all platforms"
	@echo "  dev          - Full development cycle (tidy, check, install)"
	@echo "  devlocal     - Full development cycle for Terraform"
	@echo ""
	@echo "Note: This GNUmakefile delegates to Mage. Install with:"
	@echo "  go install github.com/magefile/mage@latest"
