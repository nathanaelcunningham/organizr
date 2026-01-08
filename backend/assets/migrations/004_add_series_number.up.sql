-- Add series_number column to downloads table
ALTER TABLE downloads ADD COLUMN series_number TEXT DEFAULT '';
