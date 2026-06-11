// ABOUTME: Content fingerprint for note deduplication: a sha256 over the
// ABOUTME: author pubkey and content, identifying reposts of identical text.
package db

import (
	"crypto/sha256"
	"fmt"
)

// contentFingerprint returns the sha256 hex digest of pubkey + content. The
// same author posting identical content yields the same fingerprint, even when
// the events have different ids, which is what lets us deduplicate reposts.
func contentFingerprint(pubkey, content string) string {
	hash := sha256.Sum256([]byte(pubkey + content))
	return fmt.Sprintf("%x", hash[:])
}
