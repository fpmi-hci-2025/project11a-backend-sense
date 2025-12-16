-- Add title column to publications table
ALTER TABLE publications ADD COLUMN IF NOT EXISTS title text NOT NULL DEFAULT '';

