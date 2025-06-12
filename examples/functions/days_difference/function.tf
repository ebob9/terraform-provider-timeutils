locals {
  start_date = "2024-01-15T10:30:00Z"
  end_date   = "2024-06-11T15:45:30Z"
  current_date = timestamp()

  # Days difference
  days_diff = provider::timeutils::days_difference(local.start_date, local.end_date)
}

output "days_between" {
  description = "Days between two timestamps"
  value       = parseint(local.days_diff, 10)
}
