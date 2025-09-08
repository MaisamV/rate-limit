# Database Migrations

This directory contains all database migration files for the project using [golang-migrate](https://github.com/golang-migrate/migrate).

## File Naming Convention

Migration files follow the format: `{version}_{description}.{direction}.sql`

- **version**: Sequential number with leading zeros (e.g., 000001, 000002)
- **description**: Brief description using underscores (e.g., initial_schema, add_user_roles)
- **direction**: Either `up` (apply migration) or `down` (rollback migration)

### Examples
```
000001_initial_schema.up.sql
000001_initial_schema.down.sql
000002_add_user_roles.up.sql
000002_add_user_roles.down.sql
```

## Migration Guidelines

### Writing Migrations

1. **Always create both up and down files** for every migration
2. **Use transactions** when possible for atomic operations
3. **Add comments** to explain complex changes
4. **Use IF EXISTS/IF NOT EXISTS** for idempotent operations
5. **Test migrations** on a copy of production data

### Best Practices

- **Backward Compatible**: Ensure migrations don't break existing code
- **Small Changes**: Keep migrations focused on single logical changes
- **No Data Loss**: Never write migrations that could lose data
- **Rollback Safe**: Ensure down migrations can safely reverse changes

## Running Migrations

### Development
```bash
# Apply all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check migration status
make migrate-status

# Force migration version (use with caution)
make migrate-force VERSION=1
```

### Production

**⚠️ IMPORTANT**: Never run migrations automatically in production!

1. **Review** all pending migrations
2. **Test** on staging environment first
3. **Backup** database before applying
4. **Run** migrations manually or via approved CI/CD pipeline
5. **Verify** application functionality after migration

## Environment-Specific Behavior

### Development
- Migrations run automatically via docker-compose
- Database is disposable, so aggressive migrations are acceptable

### CI/CD
- Migrations run as separate job before service deployment
- Must complete successfully before deploying application

### Production
- Migrations run manually or via controlled pipeline
- Requires explicit approval and monitoring
- Should be run during maintenance windows when possible

## Troubleshooting

### Common Issues

1. **Migration fails**: Check logs, fix issue, and force to previous version if needed
2. **Dirty state**: Use `migrate force` to set version manually (with caution)
3. **Schema drift**: Ensure all changes go through migration files

### Recovery

```bash
# Check current migration version
make migrate-version

# Force to specific version (emergency use only)
make migrate-force VERSION=X

# Create new migration to fix issues
make migrate-create NAME=fix_issue_description
```

## Creating New Migrations

```bash
# Create new migration files
make migrate-create NAME=your_migration_description
```

This will create both `.up.sql` and `.down.sql` files with the next sequential version number.