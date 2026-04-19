// Package labelmap provides a simple port-to-label registry for portwatch.
//
// Labels are loaded from a plain-text file where each line contains a port
// number followed by a space and a label string:
//
//	# comment
//	22  ssh
//	80  http
//	443 https
//
// Labels can also be set programmatically via Set. The Annotate method
// produces a filtered map of labels for a given list of ports, which is
// useful for enriching scan results before reporting.
package labelmap
