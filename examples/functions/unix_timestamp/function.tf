
locals {
  start_date   = "2024-01-15T10:30:00Z"
  end_date     = "2024-06-11T15:45:30Z"
  current_date = timestamp()

  # Unix Timestamp
  unix_time = provider::timeutils::unix_timestamp(local.end_date)
}

output "unix_timestamp" {
  description = "Straight Unix timestamp"
  value       = parseint(local.unix_time, 10)
}
