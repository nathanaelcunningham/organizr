-- Add category column to downloads table
ALTER TABLE downloads ADD COLUMN category TEXT DEFAULT '';
