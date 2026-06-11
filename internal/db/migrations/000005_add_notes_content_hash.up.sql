-- Deduplicate notes by content fingerprint (same author + identical content).
-- The application stores content_hash = sha256(pubkey || content) and relies on
-- this unique index together with ON CONFLICT DO NOTHING to skip duplicates,
-- merging only the relay URLs onto the already stored note.
--
-- NOTE: if the notes table already holds duplicate content_hash values (reposts
-- stored before this feature), creating the unique index fails. De-duplicate
-- those rows first, for example:
--   DELETE FROM public.notes a USING public.notes b
--   WHERE a.content_hash = b.content_hash AND a.id > b.id;

ALTER TABLE public.notes ADD COLUMN IF NOT EXISTS content_hash character(64);

CREATE UNIQUE INDEX IF NOT EXISTS idx_notes_content_hash ON public.notes USING btree (content_hash);
