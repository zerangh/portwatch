// Package checkpoint manages a lightweight on-disk record of the most
// recent successful port scan.
//
// After every scan the watcher saves a Checkpoint containing the scan
// timestamp and the number of open ports observed. On the next startup
// portwatch loads the checkpoint to determine how long ago the last scan
// ran, enabling it to warn the operator when an interval was missed (e.g.
// because the host was rebooted or the process was killed).
//
// Checkpoints are written atomically via a rename so a crash mid-write
// cannot corrupt the previous checkpoint.
package checkpoint
