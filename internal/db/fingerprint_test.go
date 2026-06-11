// ABOUTME: Tests for the content fingerprint used to deduplicate notes from
// ABOUTME: the same author with identical content.
package db

import "testing"

func TestContentFingerprintEmpty(t *testing.T) {
	// sha256 of the empty string is a well-known vector.
	const emptySHA256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if got := contentFingerprint("", ""); got != emptySHA256 {
		t.Errorf("contentFingerprint(\"\",\"\") = %q, want %q", got, emptySHA256)
	}
}

func TestContentFingerprintHashesPubkeyPlusContent(t *testing.T) {
	// The fingerprint is the hash of pubkey + content concatenated, so these
	// two calls must produce the same value.
	if contentFingerprint("a", "b") != contentFingerprint("ab", "") {
		t.Error("fingerprint must hash pubkey+content concatenated")
	}
}

func TestContentFingerprintProperties(t *testing.T) {
	fp := contentFingerprint("npub1", "hello world")

	if len(fp) != 64 {
		t.Errorf("fingerprint length = %d, want 64 hex chars", len(fp))
	}
	// Deterministic.
	if fp != contentFingerprint("npub1", "hello world") {
		t.Error("fingerprint should be deterministic")
	}
	// Different content yields a different fingerprint.
	if fp == contentFingerprint("npub1", "hello worle") {
		t.Error("different content should yield a different fingerprint")
	}
	// Different author yields a different fingerprint.
	if fp == contentFingerprint("npub2", "hello world") {
		t.Error("different pubkey should yield a different fingerprint")
	}
}
