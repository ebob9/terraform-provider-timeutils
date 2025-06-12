# Terraform Time Utils Provider

A custom Terraform provider that provides advanced time manipulation functions missing from Terraform's built-in capabilities, including RFC3339 parsing, Unix timestamp conversion, strftime formatting, and precise date calculations.

## Features

This provider adds the following functions to Terraform:

- `unix_timestamp(rfc3339_string)` - Convert RFC3339 timestamp to Unix timestamp
- `strftime(format, rfc3339_string)` - Format timestamps using strftime format specifiers
- `days_difference(start_rfc3339, end_rfc3339)` - Calculate exact days between timestamps
- `parse_rfc3339(rfc3339_string)` - Parse timestamp into components (year, month, day, etc.)

## Installation

### Method 1: Local Development Install

1. **Clone and build the provider:**
   ```bash
   git clone https://github.com/ebob9/terraform-provider-timeutils.git
   cd terraform-provider-timeutils
   mage devlocal 
   ```

2. **The provider will be installed to:**
   ```
   ~/.terraform.d/plugins/registry.terraform.io/ebob9/timeutils/0.1.0/<OS>_<ARCH>/
   ```

### Method 2: Manual Install

1. **Build the binary:**
   ```bash
   go build -o terraform-provider-timeutils
   ```

2. **Create the plugin directory:**
   ```bash
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/ebob9/timeutils/0.1.0/$(go env GOOS)_$(go env GOARCH)/
   ```

3. **Copy the binary:**
   ```bash
   cp terraform-provider-timeutils ~/.terraform.d/plugins/registry.terraform.io/your-org/timeutils/0.1.0/$(go env GOOS)_$(go env GOARCH)/
   ```

### Method 3: Cross-platform Release Build

```bash
make release
```

This creates binaries for multiple platforms in the `./bin/` directory.

## Usage

### Provider Configuration

```hcl
terraform {
   required_providers {
      timeutils = {
         source = "ebob9/timeutils"
         version = "~> 0.1"
      }
   }
}

provider "timeutils" {}
```

### Function Examples

#### Calculate Days Between Timestamps

```hcl
locals {
  start_date = "2024-01-15T10:30:00Z"
  end_date = "2024-06-11T15:45:30Z"
  
  days_diff = provider::timeutils::days_difference(local.start_date, local.end_date)
}

output "days_between" {
  value = parseint(local.days_diff, 10)  # Convert string to number
}
```

#### Unix Timestamp Conversion

```hcl
locals {
  timestamp = "2024-01-15T10:30:00Z"
  unix_time = provider::timeutils::unix_timestamp(local.timestamp)
}

output "unix_timestamp" {
  value = parseint(local.unix_time, 10)
}
```

#### strftime Formatting

```hcl
locals {
  timestamp = "2024-01-15T10:30:00Z"
  
  formatted_dates = {
    iso_date = provider::timeutils::strftime("%Y-%m-%d", local.timestamp)
    us_date = provider::timeutils::strftime("%m/%d/%Y", local.timestamp)
    full_date = provider::timeutils::strftime("%A, %B %d, %Y", local.timestamp)
    time_12h = provider::timeutils::strftime("%I:%M %p", local.timestamp)
    time_24h = provider::timeutils::strftime("%H:%M:%S", local.timestamp)
  }
}
```

#### Parse Timestamp Components

```hcl
locals {
  timestamp = "2024-01-15T10:30:00Z"
  parsed = jsondecode(provider::timeutils::parse_rfc3339(local.timestamp))
}

output "timestamp_components" {
  value = {
    year = parseint(local.parsed.year, 10)
    month = parseint(local.parsed.month, 10)
    day = parseint(local.parsed.day, 10)
    hour = parseint(local.parsed.hour, 10)
    minute = parseint(local.parsed.minute, 10)
    second = parseint(local.parsed.second, 10)
    unix = parseint(local.parsed.unix, 10)
    weekday = parseint(local.parsed.weekday, 10)  # 0=Sunday, 1=Monday, etc.
  }
}
```

### Days Between Timestamp and Now

```hcl
variable "start_timestamp" {
  description = "Start timestamp in RFC3339 format"
  type        = string
}

locals {
  # Get current time and calculate difference
  current_time = timestamp()
  days_since_start = provider::timeutils::days_difference(var.start_timestamp, local.current_time)
}

output "days_between" {
  description = "Number of days between start timestamp and now"
  value       = parseint(local.days_since_start, 10)
}
```

## Supported strftime Format Specifiers

| Format | Description | Example |
|--------|-------------|---------|
| %Y | 4-digit year | 2024 |
| %y | 2-digit year | 24 |
| %m | Month (01-12) | 01 |
| %B | Full month name | January |
| %b | Abbreviated month | Jan |
| %d | Day of month (01-31) | 15 |
| %A | Full weekday name | Monday |
| %a | Abbreviated weekday | Mon |
| %H | Hour (00-23) | 10 |
| %I | Hour (01-12) | 10 |
| %M | Minute (00-59) | 30 |
| %S | Second (00-59) | 00 |
| %p | AM/PM | AM |
| %Z | Timezone name | UTC |
| %z | Timezone offset | +0000 |

## Development

### Prerequisites

- Go 1.21+
- [Mage](https://magefile.org/) - Install with: `go install github.com/magefile/mage@latest`
- Terraform 1.0+ or OpenTofu 1.6+

### Available Mage Targets

```bash
# Installation
mage install      # Build and install to Go bin directory (default)
mage installlocal # Install to Terraform plugins directory

# Development
mage dev          # Full development cycle (tidy, check, install to Go bin)
mage devlocal     # Full development cycle for Terraform (install locally)
mage build        # Build the provider binary
mage clean        # Remove build artifacts

# Code Quality
mage fmt          # Format Go code
mage vet          # Run go vet
mage test         # Run unit tests
mage testacc      # Run acceptance tests
mage check        # Run fmt, vet, and test

# Dependencies
mage tidy         # Run go mod tidy

# Release
mage release      # Build binaries for all platforms
```

### Quick Start Development

1. **Set up the development environment:**
   ```bash
   git clone https://github.com/ebob9/terraform-provider-timeutils
   cd terraform-provider-timeutils
   mage dev
   ```

   This will automatically detect your OS and architecture and install to the appropriate directory.

2. **Create a test Terraform configuration using the provider**

3. **Test your changes:**
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

### Platform Detection

The build system automatically detects your platform:

```bash
mage install 
```

### Testing

```bash
# Run unit tests
mage test

# Run acceptance tests
mage testacc
```

## Why This Provider?

Terraform's built-in time functions are limited:
- No Unix timestamp conversion
- No strftime-style formatting
- No precise date arithmetic
- Limited date manipulation capabilities

This provider fills those gaps by providing:
- ✅ Precise date calculations
- ✅ Unix timestamp support
- ✅ Full strftime formatting
- ✅ RFC3339 parsing capabilities
- ✅ Cross-platform compatibility

## License

MPL-2.0

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request