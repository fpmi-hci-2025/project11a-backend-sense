-- Remove denormalized counter columns from publications and comments tables
-- These will be calculated dynamically using JOINs

BEGIN;

-- Drop counter columns from publications table
ALTER TABLE publications DROP COLUMN IF EXISTS likes_count;
ALTER TABLE publications DROP COLUMN IF EXISTS comments_count;
ALTER TABLE publications DROP COLUMN IF EXISTS saved_count;

-- Drop counter column from comments table
ALTER TABLE comments DROP COLUMN IF EXISTS likes_count;

COMMIT;

