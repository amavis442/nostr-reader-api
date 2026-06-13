CREATE INDEX IF NOT EXISTS idx_notes_event_created_at
    ON notes (event_created_at DESC)
    WHERE kind = 1 AND garbage = false AND root = true;
