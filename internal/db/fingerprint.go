// ABOUTME: Content fingerprints for note deduplication: a sha256 over the
// ABOUTME: author and content, branch-aware for replies, plain for root notes.
package db

import (
	"crypto/sha256"
	"fmt"
)

// contentFingerprint returns the sha256 hex digest of pubkey + content. The
// same author posting identical content yields the same fingerprint, even when
// the events have different ids, which is what lets us deduplicate reposts.
// This is the fingerprint for root notes; replies use replyFingerprint.
func contentFingerprint(pubkey, content string) string {
	hash := sha256.Sum256([]byte(pubkey + content))
	return fmt.Sprintf("%x", hash[:])
}

// replyFingerprint returns the sha256 hex digest identifying a reply for
// deduplication. It folds the root and reply event ids into the hash so the
// same author copying identical text into a different branch yields a distinct
// fingerprint, instead of being dropped as a false duplicate of the original
// reply. Identical text in the exact same branch still collides, so genuine
// reposts keep deduplicating.
//
// The '|' separator disambiguates the fields: pubkey and the event ids are
// fixed-length hex (or empty) and cannot contain '|', and content comes last,
// so free-form content can never be confused with a thread id.
func replyFingerprint(pubkey, content, rootTag, replyTag string) string {
	hash := sha256.Sum256([]byte(pubkey + "|" + rootTag + "|" + replyTag + "|" + content))
	return fmt.Sprintf("%x", hash[:])
}
