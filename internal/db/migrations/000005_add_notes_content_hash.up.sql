-- Deduplicate notes by content fingerprint (same author + identical content).
-- The application stores content_hash = sha256(pubkey || content) and relies on
-- this unique index together with ON CONFLICT DO NOTHING to skip duplicates,
-- merging only the relay URLs onto the already stored note.

BEGIN;

ALTER TABLE public.notes ADD COLUMN IF NOT EXISTS content_hash character(64);

-- Pre-existing duplicate notes (stored before this feature) would block the
-- unique index, so remove them first: keep the lowest id per content_hash and
-- delete the rest together with their dependent rows. reactions and
-- notifications have a foreign key to notes.id; seens and bookmarks reference
-- note_id without one. This is destructive, but only drops redundant copies.
CREATE TEMP TABLE dup_note_ids ON COMMIT DROP AS
SELECT id
FROM (
    SELECT id, row_number() OVER (PARTITION BY content_hash ORDER BY id) AS rn
    FROM public.notes
    WHERE content_hash IS NOT NULL
) ranked
WHERE rn > 1;

DELETE FROM public.reactions     WHERE note_id IN (SELECT id FROM dup_note_ids);
DELETE FROM public.notifications WHERE note_id IN (SELECT id FROM dup_note_ids);
DELETE FROM public.bookmarks     WHERE note_id IN (SELECT id FROM dup_note_ids);
DELETE FROM public.seens         WHERE note_id IN (SELECT id FROM dup_note_ids);
DELETE FROM public.notes         WHERE id      IN (SELECT id FROM dup_note_ids);

CREATE UNIQUE INDEX IF NOT EXISTS idx_notes_content_hash ON public.notes USING btree (content_hash);

COMMIT;
