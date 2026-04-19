// Package env reads portwatch runtime configuration from environment variables.
//
// Supported variables:
//
//	PORTWATCH_PORT_RANGE    – port range to scan, e.g. "1-65535"
//	PORTWATCH_CONCURRENCY   – number of concurrent scan workers
//	PORTWATCH_INTERVAL      – polling interval as a Go duration, e.g. "60s"
//	PORTWATCH_STATE_PATH    – path to the state file
//	PORTWATCH_ALERT_WEBHOOK – webhook URL for change notifications
//
// Environment values take precedence over file-based configuration when
// merged by the config loader.
package env
