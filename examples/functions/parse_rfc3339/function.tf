
locals {
  start_date   = "2024-01-15T10:30:00Z"
  end_date     = "2024-06-11T15:45:30Z"
  current_date = timestamp()

  # Timestamp parsing
  parsed = jsondecode(provider::timeutils::parse_rfc3339(local.end_date))

}

output "timestamp_components" {
  description = "Map of timestamp components"
  value = {
    year    = parseint(local.parsed.year, 10)
    month   = parseint(local.parsed.month, 10)
    day     = parseint(local.parsed.day, 10)
    hour    = parseint(local.parsed.hour, 10)
    minute  = parseint(local.parsed.minute, 10)
    second  = parseint(local.parsed.second, 10)
    unix    = parseint(local.parsed.unix, 10)
    weekday = parseint(local.parsed.weekday, 10) # 0=Sunday, 1=Monday, etc.
  }
}