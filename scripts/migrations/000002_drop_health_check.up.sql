-- Drop health_check table migration
-- This migration removes the health_check table that was created in 000001

BEGIN;

-- Drop the health check table
DROP TABLE IF EXISTS health_check;

COMMIT;