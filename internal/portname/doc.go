// Package portname provides a lightweight resolver that maps port numbers to
// human-readable service names.
//
// A built-in table covers the most common well-known ports (ssh, http, https,
// mysql, redis, …). Callers may supply additional mappings via New which take
// precedence over the built-ins.
//
// Usage:
//
//	r := portname.New(nil)
//	fmt.Println(r.Lookup(443))          // "https"
//	fmt.Println(r.Annotate([]int{22, 80})) // ["22(ssh)" "80(http)"]
package portname
