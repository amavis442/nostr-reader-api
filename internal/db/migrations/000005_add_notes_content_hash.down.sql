DROP INDEX IF EXISTS public.idx_notes_content_hash;

-- The content_hash column is intentionally left in place: it may hold data and
-- can pre-date this migration. Uncomment to fully revert.
-- ALTER TABLE public.notes DROP COLUMN IF EXISTS content_hash;
