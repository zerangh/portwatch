// Package portgroup provides a registry for grouping ports into named
// categories such as "web", "db", or "monitoring". Groups can be defined
// programmatically and used to classify open ports discovered during a scan,
// enriching alerts and reports with human-readable context.
//
// Example:
//
//	r := portgroup.New()
//	r.Define("web", []int{80, 443, 8080})
//	names := r.Classify(443) // returns ["web"]
package portgroup
