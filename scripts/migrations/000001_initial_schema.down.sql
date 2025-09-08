-- Rollback initial schema migration
-- This removes the health check table created in the up migration

BEGIN;

-- Drop the health check table
DROP TABLE IF EXISTS health_check;

COMMIT;