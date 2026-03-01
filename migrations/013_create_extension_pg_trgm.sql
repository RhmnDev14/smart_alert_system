-- Enable pg_trgm extension for fuzzy text matching (trigram similarity)
-- This allows matching text even with typos, e.g. "joging" matches "Jogging"
CREATE EXTENSION IF NOT EXISTS pg_trgm;
