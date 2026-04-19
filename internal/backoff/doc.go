// Package backoff implements exponential backoff with configurable policy
// for retrying transient failures in portwatch operations such as scanning
// and state persistence.
//
// Basic usage:
//
//	b := backoff.New(backoff.DefaultPolicy())
//	for {
//		err := doSomething()
//		if err == nil {
//			break
//		}
//		if !b.Wait(ctx) {
//			return errors.New("max retries exceeded")
//		}
//	}
package backoff
