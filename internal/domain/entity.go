// ABOUTME: Domain types shared across packages without external dependencies.
// ABOUTME: MaxNoteAge bounds how old an incoming note may be before it is ignored on ingest.
package domain

type MaxNoteAge struct {
	Days int `json:"days"`
}
