-- Database initialization script
-- This script runs when PostgreSQL container starts for the first time
-- It ensures the database exists and is ready for migrations

-- Create the main database if it doesn't exist
-- Note: This is mainly for safety, as POSTGRES_DB env var should handle this
SELECT 'CREATE DATABASE go_clean_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'go_clean_db')\gexec

-- Connect to the database
\c go_clean_db;

-- Enable UUID extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create a schema_migrations table for golang-migrate
-- This will be used by golang-migrate to track migration versions
CREATE TABLE IF NOT EXISTS schema_migrations (
    version bigint NOT NULL PRIMARY KEY,
    dirty boolean NOT NULL
);

-- Add some helpful comments
COMMENT ON TABLE schema_migrations IS 'Tracks database migration versions for golang-migrate';
COMMENT ON COLUMN schema_migrations.version IS 'Migration version number';
COMMENT ON COLUMN schema_migrations.dirty IS 'Indicates if migration failed and needs manual intervention';

-- Log successful initialization
\echo 'Database initialization completed successfully';