// Package digest provides hashing utilities for port snapshots,
// allowing quick comparison of scan results without full diff computation.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Digest is a compact hash representing a set of open ports.
type Digest string

// Empty is the digest of an empty port set.
const Empty Digest = "e3b0c44298fc1c149afb"

// FromPorts computes a stable SHA-256 digest from a slice of port numbers.
// The input slice does not need to be sorted.
func FromPorts(ports []int) Digest {
	if len(ports) == 0 {
		return Empty
	}

	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	strs := make([]string, len(sorted))
	for i, p := range sorted {
		strs[i] = strconv.Itoa(p)
	}

	input := strings.Join(strs, ",")
	sum := sha256.Sum256([]byte(input))
	return Digest(hex.EncodeToString(sum[:10]))
}

// Equal reports whether two digests are identical.
func Equal(a, b Digest) bool {
	return a == b
}

// String returns the string representation of the digest.
func (d Digest) String() string {
	return fmt.Sprintf("digest:%s", string(d))
}
