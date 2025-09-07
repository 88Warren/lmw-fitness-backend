-- Migration to add RoundRest column to workout_blocks table
-- This adds support for rest periods between rounds

ALTER TABLE workout_blocks ADD COLUMN round_rest VARCHAR(255) DEFAULT '';

-- Update existing records with round rest where mentioned in block notes
UPDATE workout_blocks SET round_rest = '60s' WHERE block_notes LIKE '%60 seconds rest between rounds%' OR block_notes LIKE '%60 Second rest between rounds%';
UPDATE workout_blocks SET round_rest = '90s' WHERE block_notes LIKE '%90 seconds rest between rounds%' OR block_notes LIKE '%Rest 90s between rounds%';
UPDATE workout_blocks SET round_rest = '120s' WHERE block_notes LIKE '%2 minutes between rounds%';
UPDATE workout_blocks SET round_rest = '75s' WHERE block_notes LIKE '%60-90s rest between rounds%';