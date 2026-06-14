CREATE TABLE public.watermarks (
    context         VARCHAR(20) PRIMARY KEY,
    event_created_at BIGINT NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
