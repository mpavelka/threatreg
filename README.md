# threatreg

## Database Migrations Guide

```bash
# Initialize migrations (one-time setup)
python cli.py db init

# Create a new migration
python cli.py db migrate -m "Add user table"

# Apply migrations
python cli.py db upgrade

# Check migration status
python cli.py db status

# Rollback one migration
python cli.py db downgrade
```