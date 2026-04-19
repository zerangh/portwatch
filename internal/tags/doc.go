// Package tags provides port tagging support for portwatch.
//
// Tags associate human-readable labels with individual ports or port ranges.
// Labels are loaded from a JSON file and surfaced during alert reporting to
// give operators contextual information about which services changed.
//
// Example tags.json:
//
//	[
//	  {"label": "ssh",       "ports": [22]},
//	  {"label": "web",       "ports": [80, 443, 8080]},
//	  {"label": "ephemeral", "from": 49152, "to": 65535}
//	]
package tags
